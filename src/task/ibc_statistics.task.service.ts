/* eslint-disable @typescript-eslint/camelcase */
import { Injectable } from '@nestjs/common';
import { Connection } from 'mongoose';
import { InjectConnection } from '@nestjs/mongoose';
import { IbcStatisticsSchema } from '../schema/ibc_statistics.schema';
import { IbcChainConfigSchema } from '../schema/ibc_chain_config.schema';
import { IbcChainSchema } from '../schema/ibc_chain.schema';
import { IbcTxSchema } from '../schema/ibc_tx.schema';
import { IbcDenomSchema } from '../schema/ibc_denom.schema';
import { IbcChannelSchema } from 'src/schema/ibc_channel.schema';
import { IbcBaseDenomSchema } from '../schema/ibc_base_denom.schema';
import { TaskEnum, IbcTxStatus, StatisticsNames } from '../constant';

@Injectable()
export class IbcStatisticsTaskService {
  private ibcStatisticsModel;
  private chainConfigModel;
  private ibcChainModel;
  private ibcTxModel;
  private ibcDenomModel;
  private ibcBaseDenomModel;
  private ibcChannelModel;

  constructor(@InjectConnection() private connection: Connection) {
    this.getModels();
    this.doTask = this.doTask.bind(this);
  }

  async doTask(taskName?: TaskEnum): Promise<void> {
    const dateNow = String(Math.floor(new Date().getTime() / 1000));
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

    // chainModel
    this.ibcChainModel = await this.connection.model(
      'ibcChainModel',
      IbcChainSchema,
      'ibc_chain',
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
      'ex_ibc_tx',
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

  // sync count
  async parseIbcStatistics(dateNow): Promise<void> {
    // tx_24hr_all
    const tx_24hr_all = await this.ibcTxModel.countActive();

    // chains_24hr_all
    const chains_24hr = await this.ibcChainModel.countActive();

    // channels_24hr
    const channels_24hr = await this.ibcChannelModel.countActive();

    // chain_all
    const chain_all = await this.chainConfigModel.findCount();

    const chain_all_record = await this.chainConfigModel.findAll();
    const channels_all_record = [];
    chain_all_record.forEach(chain => {
      chain.ibc_info && chain.ibc_info.forEach(ibc_info_item => {
        ibc_info_item.paths.forEach(channel => {
          channels_all_record.push({
            channel_id: channel.channel_id,
            state: channel.state,
          });
        });
      });
    });

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
    const base_denom_all = await this.ibcDenomModel.findBaseDenomCount();

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

    StatisticsNames.forEach(async statistics_name => {
      const statisticsRecord = await this.ibcStatisticsModel.findStatisticsRecord(
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
        await this.ibcStatisticsModel.updateStatisticsRecord(statisticsRecord);
      }
    });
  }
}
