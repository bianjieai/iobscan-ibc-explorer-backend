/* eslint-disable @typescript-eslint/camelcase */
import {Injectable} from '@nestjs/common';
import {Connection} from 'mongoose';
import {InjectConnection} from '@nestjs/mongoose';
import {IbcStatisticsSchema} from '../schema/ibc_statistics.schema';
import {IbcChainConfigSchema} from '../schema/ibc_chain_config.schema';
import {IbcTxSchema} from '../schema/ibc_tx.schema';
import {IbcDenomSchema} from '../schema/ibc_denom.schema';
import {TaskEnum, StatisticsNames} from '../constant';
import {AggregateBaseDenomCnt} from "../types/schemaTypes/ibc_denom.interface";
import {AggregateResult24hr} from "../types/schemaTypes/ibc_tx.interface";
import {IbcStatisticsType} from "../types/schemaTypes/ibc_statistics.interface";

@Injectable()
export class IbcStatisticsTaskService {
    private ibcStatisticsModel;
    private chainConfigModel;
    private ibcTxModel;
    private ibcDenomModel;

    constructor(@InjectConnection() private connection: Connection) {
        this.getModels();
        this.doTask = this.doTask.bind(this);
    }

    async doTask(taskName?: TaskEnum): Promise<void> {
        const dateNow = Math.floor(new Date().getTime() / 1000);
        this.parseIbcStatistics(dateNow);
    }

    // getModels
    async getModels(): Promise<void> {
        // ibcStatisticsModel
        this.ibcStatisticsModel = await this.connection.model(
            'ibcStatisticsModel',
            IbcStatisticsSchema,
            'ibc_statistics',
        );

        // chainConfigModel
        this.chainConfigModel = await this.connection.model(
            'chainConfigModel',
            IbcChainConfigSchema,
            'chain_config',
        );

        // ibcTxModel
        this.ibcTxModel = await this.connection.model(
            'ibcTxModel',
            IbcTxSchema,
            'ex_ibc_tx',
        );

        // ibcDenomModel
        this.ibcDenomModel = await this.connection.model(
            'ibcDenomModel',
            IbcDenomSchema,
            'ibc_denom',
        );

    }

    async aggregateFindSrcChannels(dateNow, chains: Array<string>): Promise<Array<AggregateResult24hr>> {
        return await this.ibcTxModel.aggregateFindSrcChannels24hr(dateNow, chains);
    }

    async aggregateFindDesChannels(dateNow, chains: Array<string>): Promise<Array<AggregateResult24hr>> {
        return await this.ibcTxModel.aggregateFindDesChannels24hr(dateNow, chains);
    }

    async updateStatisticsRecord(statisticsRecord: IbcStatisticsType) {
        await this.ibcStatisticsModel.updateStatisticsRecord(statisticsRecord);
    }

    async findStatisticsRecord(statistics_name: string): Promise<IbcStatisticsType> {
        return await this.ibcStatisticsModel.findStatisticsRecord(
            statistics_name,
        );
    }

    async aggregateBaseDenomCnt(): Promise<Array<AggregateBaseDenomCnt>> {
        return await this.ibcDenomModel.findBaseDenomCount()
    }

    // sync count
    async parseIbcStatistics(dateNow): Promise<void> {
        // tx_24hr_all
        const tx_24hr_all = await this.ibcTxModel.countActive();

        // chain_all
        const chain_all = await this.chainConfigModel.findCount();

        const chain_all_record = await this.chainConfigModel.findAll();
        const channels_all_record = [], chains = [];
        chain_all_record.forEach(chain => {
            chains.push(chain.chain_id)
            chain.ibc_info && chain.ibc_info.forEach(ibc_info_item => {
                ibc_info_item.paths.forEach(channel => {
                    channels_all_record.push({
                        channel_id: channel.channel_id,
                        state: channel.state,
                    });
                });
            });
        });

        //sc_chain_id,sc_channel
        const srcinfo_24hr = await this.aggregateFindSrcChannels(dateNow, chains);

        //dc_chain_id,dc_channel
        const desinfo_24hr = await this.aggregateFindDesChannels(dateNow, chains);


        const chainMap = new Map();
        for (const element of srcinfo_24hr) {
            if (chainMap.has(element._id.sc_chain_id) === false) {
                chainMap.set(element._id.sc_chain_id, '')
            }
        }

        for (const element of desinfo_24hr) {
            if (chainMap.has(element._id.dc_chain_id) === false) {
                chainMap.set(element._id.dc_chain_id, '')
            }
        }

        // chains_24hr_all
        const chains_24hr = chainMap.size;

        // channels_24hr
        const channels_24hr = srcinfo_24hr.length + desinfo_24hr.length;


        // channel_all
        const channel_all = channels_all_record.length;

        // channel_opened
        const channel_opened = channels_all_record.filter(channel => {
            return channel.state === 'STATE_OPEN';
        }).length;

        // channel_closed
        const channel_closed = channels_all_record.filter(channel => {
            return channel.state === 'STATE_CLOSED';
        }).length;

        // tx_all
        const tx_all = await this.ibcTxModel.countAll();

        // tx_success
        const tx_success = await this.ibcTxModel.countSuccess();

        // tx_failed
        const tx_failed = await this.ibcTxModel.countFaild();

        // denom_all
        const denom_all = await this.ibcDenomModel.findCount();

        // base_denom_all
        // const base_denom_all = await this.ibcBaseDenomModel.findCount();
        const base_denoms = await this.aggregateBaseDenomCnt();
        const base_denom_all = base_denoms.length;

        const parseCount = {
            tx_24hr_all,
            chains_24hr,
            channels_24hr,
            chain_all,
            channel_all,
            channel_opened,
            channel_closed,
            tx_all,
            tx_success,
            tx_failed,
            base_denom_all,
            denom_all,
        };

        for (const statistics_name of StatisticsNames) {
            const statisticsRecord = await this.findStatisticsRecord(
                statistics_name,
            );
            if (!statisticsRecord) {
                await this.ibcStatisticsModel.insertManyStatisticsRecord({
                    statistics_name,
                    count: parseCount[statistics_name],
                    create_at: dateNow,
                    update_at: dateNow,
                });
            } else {
                statisticsRecord.count = parseCount[statistics_name];
                statisticsRecord.update_at = dateNow;
                await this.updateStatisticsRecord(statisticsRecord);
            }
        }
    }
}
