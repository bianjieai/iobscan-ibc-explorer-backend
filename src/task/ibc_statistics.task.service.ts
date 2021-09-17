import { Injectable } from '@nestjs/common';
import { Connection } from 'mongoose';
import { InjectConnection } from '@nestjs/mongoose';
import { IbcStatisticsSchema } from '../schema/ibc_statistics.schema';
import { IbcChainSchema } from '../schema/ibc_chain.schema';
import { IbcTxSchema } from '../schema/ibc_tx.schema';
import { IbcDenomSchema } from '../schema/ibc_denom.schema';
import { IbcChannelSchema } from 'src/schema/ibc_channel.schema';
import { IbcBaseDenomSchema } from '../schema/ibc_base_denom.schema';
import { TaskEnum, IbcTxStatus, StatisticsNames } from '../constant';
import { flatten } from 'lodash';

@Injectable()
export class IbcStatisticsTaskService {
  private ibcStatisticsModel;
  private chainModel;
  private ibcTxModel;
  private ibcDenomModel;
  private ibcBaseDenomModel;
  private ibcChannelModel;
  constructor(@InjectConnection() private connection: Connection) {
    this.getModels();
    this.doTask = this.doTask.bind(this);
  }

  async doTask(taskName?: TaskEnum): Promise<void> {
    const dateNow = String(new Date().getTime());
    this.parseIbcStatistics(dateNow);
  }

  // 获取所有静态模型
  async getModels(): Promise<void> {
    // ibcStatisticsModel
    this.ibcStatisticsModel = await this.connection.model(
      'ibcStatisticsModel',
      IbcStatisticsSchema,
      'ibc_statistics',
    );

    // chainModel
    this.chainModel = await this.connection.model(
      'chainModel',
      IbcChainSchema,
      'chain_config',
    );

    // ibcChannelSchema
    this.ibcChannelModel = await this.connection.model(
      'ibcChannelModel',
      IbcChannelSchema,
      'ibc_channel',
    );

    // ibcTxModel
    this.ibcTxModel = await this.connection.model(
      'ibcTxModel',
      IbcTxSchema,
      'ibc_txs_test',
    );

    // ibcDenomModel
    this.ibcDenomModel = await this.connection.model(
      'ibcDenomModel',
      IbcDenomSchema,
      'ibc_denom',
    );

    // ibcBaseDenomModel
    this.ibcBaseDenomModel = await this.connection.model(
      'ibcBaseDenomModel',
      IbcBaseDenomSchema,
      'ibc_base_denom',
    );
  }

  // 同步首页数量
  async parseIbcStatistics(dateNow): Promise<void> {
    // tx_24hr_all
    const tx_24hr_all = await this.ibcTxModel.findCount({
      update_at: { $gte: dateNow - 24 * 60 * 60 * 1000 },
    });

    const sc_chains = await this.ibcTxModel.distinctChainList({
      type: 'sc_chain_id',
      dateNow,
      status: [
        IbcTxStatus.SUCCESS,
        IbcTxStatus.PROCESSING,
        IbcTxStatus.SETTING,
        IbcTxStatus.REFUNDED,
      ],
    });

    const dc_chains = await this.ibcTxModel.distinctChainList({
      type: 'dc_chain_id',
      dateNow,
      status: [IbcTxStatus.SUCCESS],
    });

    // chains_24hr_all
    const chains_24hr = Array.from(new Set([...sc_chains, ...dc_chains]))
      .length;

    // channels_24hr
    const channels_24hr = await this.ibcChannelModel.findCount({
      update_at: { $gte: dateNow - 24 * 60 * 60 * 1000 },
    });
    // chain_all
    const chain_all = await this.chainModel.findCount();

    // channel_all
    const channels_arr = await this.chainModel.aggregateFindChannels();
    let channel_all = 0;
    channels_arr.forEach(channels => {
      channel_all += flatten(channels['_id']).length;
    });

    // tx_all
    const tx_all = await this.ibcTxModel.findCount();

    // tx_success
    const tx_success = await this.ibcTxModel.findCount({
      status: IbcTxStatus.SUCCESS,
    });

    // tx_failed
    const tx_failed = await this.ibcTxModel.findCount({
      status: { $in: [IbcTxStatus.FAILED, IbcTxStatus.REFUNDED] },
    });

    // denom_all
    const denom_all = await this.ibcDenomModel.findCount();

    // base_denom_all
    const base_denom_all = await this.ibcBaseDenomModel.findCount();

    const parseCount = {
      tx_24hr_all,
      chains_24hr,
      channels_24hr,
      chain_all,
      channel_all,
      tx_all,
      tx_success,
      tx_failed,
      base_denom_all,
      denom_all,
    };

    StatisticsNames.forEach(async statistics_name => {
      const statisticsRecord = await this.ibcStatisticsModel.findStatisticsRecord(
        statistics_name,
      );
      if (!statisticsRecord) {
        // 如果不存在则新建
        await this.ibcStatisticsModel.insertManyStatisticsRecord({
          statistics_name,
          count: parseCount[statistics_name],
          create_at: dateNow,
          update_at: dateNow,
        });
      } else {
        // 否则更新
        statisticsRecord.count = parseCount[statistics_name];
        statisticsRecord.update_at = dateNow;
        await this.ibcStatisticsModel.updateStatisticsRecord(statisticsRecord);
      }
    });
  }
}
