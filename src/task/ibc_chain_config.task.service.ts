import { Injectable, Logger } from '@nestjs/common';
import { Connection } from 'mongoose';
import { InjectConnection } from '@nestjs/mongoose';
import { IbcChainConfigSchema } from '../schema/ibc_chain_config.schema';
import { TaskEnum, Delimiter } from '../constant';
import { ChainHttp } from '../http/lcd/chain.http';
import { groupBy } from 'lodash';

@Injectable()
export class IbcChainConfigTaskService {
  private chainConfigModel;

  constructor(@InjectConnection() private connection: Connection) {
    this.doTask = this.doTask.bind(this);
  }

  async doTask(taskName?: TaskEnum): Promise<void> {
    await this.getModels();
    await this.parseChainConfig();
  }

  async getModels(): Promise<void> {
    this.chainConfigModel = await this.connection.model(
      'chainConfigModel',
      IbcChainConfigSchema,
      'chain_config',
    );
  }

  // 获取并同步chainConfig配置表数据
  async parseChainConfig() {
    const allChains = await this.chainConfigModel.findAll();
    // 请求所有链配置的channels
    Promise.all(
      allChains.map(async chain => {
        let channels = await ChainHttp.getIbcChannels(chain.lcd);
        channels.map(channel => {
          channel.sc_chain_id = chain.chain_id;
        });
        return Promise.resolve({
          chain_id: chain.chain_id,
          chain_name: chain.chain_name,
          lcd: chain.lcd,
          icon: chain.icon,
          ibc_info: channels,
        });
      }),
    ).then(allChains => {
      const channelsObj = {};
      const allChainsId = allChains.map(chain => {
        return chain['chain_id'];
      });

      // 为channelsObj设值
      allChains.forEach(chain => {
        channelsObj[chain['chain_id']] = {};
        chain['ibc_info'].forEach(channel => {
          channelsObj[chain['chain_id']][
            `${channel.channel_id}/${channel.port_id}/${channel.counterparty.channel_id}/${channel.counterparty.port_id}`
          ] = `${channel.sc_chain_id}${Delimiter}${channel.state}`;
        });
      });

      // 从channelsObj取值
      allChains.forEach(async chain => {
        chain['ibc_info'].forEach(channel => {
          allChainsId.forEach(chainId => {
            if (chainId !== chain['chain_id']) {
              const result =
                channelsObj[chainId][
                  `${channel.counterparty.channel_id}/${channel.counterparty.port_id}/${channel.channel_id}/${channel.port_id}`
                ];
              if (result) {
                channel['chain_id'] = result.split(Delimiter)[0];
                channel['counterparty']['state'] = result.split(Delimiter)[1];
              }
            }
          });
        });

        // 过滤未配置的channels
        chain['ibc_info'] = chain['ibc_info'].filter(channel => {
          return channel.hasOwnProperty('chain_id');
        });

        // 分组数据
        const ibcInfoGroupBy = groupBy(chain['ibc_info'], 'chain_id');
        const ibcInfo = [];
        Object.keys(ibcInfoGroupBy).forEach(chain_id => {
          ibcInfo.push({ chain_id, paths: ibcInfoGroupBy[`${chain_id}`] });
        });
        chain['ibc_info'] = ibcInfo;

        // 更新数据库
        this.chainConfigModel.updateChain(chain);
      });
    });
  }
}
