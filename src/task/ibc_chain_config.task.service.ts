import { Injectable, Logger } from '@nestjs/common';
import { Connection } from 'mongoose';
import { InjectConnection } from '@nestjs/mongoose';
import { IbcChainConfigSchema } from '../schema/ibc_chain_config.schema';
import { TaskEnum, Delimiter } from '../constant';
import { ChainHttp } from '../http/lcd/chain.http';
import { groupBy } from 'lodash';
import {IbcChainConfigType} from "../types/schemaTypes/ibc_chain_config.interface";
import {Md5} from 'ts-md5/dist/md5'

@Injectable()
export class IbcChainConfigTaskService {
  private chainConfigModel;

  constructor(@InjectConnection() private connection: Connection) {
    this.getModels();
    this.doTask = this.doTask.bind(this);
  }

  async doTask(taskName?: TaskEnum): Promise<void> {
    // this.parseChainConfig();
    this.handleChain();
  }

  async getModels(): Promise<void> {
    this.chainConfigModel = await this.connection.model(
      'chainConfigModel',
      IbcChainConfigSchema,
      'chain_config',
    );
  }

  async findAllConfig():Promise<IbcChainConfigType[]> {
    return await this.chainConfigModel.findAll();
  }

  dcChainKey(chain,channel) :string{
      return `${chain.chain_id}${Delimiter}${channel.channel_id}/${channel.port_id}`
  }

  async getDcChainMap(chain){
        let scChainDcChainMap = new Map;
        for (const channel of chain['ibc_info']) {
            if (!chain.lcd_api_path.client_state_path) {
                Logger.error("get dc_chain_id fail from lcd for lcd_api_path.client_state_path is empty")
                continue
            }
            let dc_chain_id = await ChainHttp.getDcChainIdByScChannel(chain.lcd, chain.lcd_api_path.client_state_path, channel.port_id, channel.channel_id);
            if (dc_chain_id) {
                dc_chain_id = dc_chain_id.replace(new RegExp("\-", "g"),"_")
                scChainDcChainMap.set(this.dcChainKey(chain,channel),dc_chain_id)
            }
        }
        return scChainDcChainMap
  }


  async handleChain() {
      Logger.debug("Start to handle Chain Channels info")
      const allChains = await this.findAllConfig();
      // request configed allchannels
      Promise.all(
          allChains.map(async chain => {
              let channels = await ChainHttp.getIbcChannels(chain.lcd,chain.lcd_api_path.channels_path);
              channels && channels.map(channel => {
                  channel.sc_chain_id = chain.chain_id;
              });
              return Promise.resolve({
                  chain_id: chain.chain_id,
                  chain_name: chain.chain_name,
                  lcd: chain.lcd,
                  lcd_api_path: chain.lcd_api_path,
                  icon: chain.icon,
                  ibc_info: channels ? channels : [],
                  ibc_info_hash:"",
                  isManual: chain?.is_manual
              });
          }),
      ).then(allChains => {
          const channelsObj = {};
          let supportChainsIdMap = new Map;
          for(const chain of allChains) {
              supportChainsIdMap.set(chain.chain_id,chain.chain_id)
          }
          // set channelsObj datas
          allChains && allChains.forEach(chain => {
              channelsObj[chain['chain_id']] = {};
              chain['ibc_info'] && chain['ibc_info'].forEach(channel => {
                  channelsObj[chain['chain_id']][
                      `${channel.channel_id}/${channel.port_id}/${channel.counterparty.channel_id}/${channel.counterparty.port_id}`
                      ] = `${channel.sc_chain_id}${Delimiter}${channel.state}`;
              });
          });

          // get datas from channelsObj
          allChains && allChains.forEach(async chain => {

              //获取最新的ibc_info的hashCode
              const hashCode = Md5.hashStr(JSON.stringify(chain.ibc_info))
              //判断是否需要更新ibc_info信息
              if (hashCode !== chain.ibc_info_hash && chain.ibc_info.length > 0) {
                  chain.ibc_info_hash = Md5.hashStr(JSON.stringify(chain.ibc_info))
              }else{
                  Logger.log("this chain ibc_info no need update for ibc_info hashcode is not change")
                  return
              }

              const dcChainIdMap  = await this.getDcChainMap(chain)

              // console.log(dcChainIdMap,"========dcChainIdMap============>>")
              // console.log(supportChainsIdMap,"========supportChainsIdMap============>>")

              chain['ibc_info'] && chain['ibc_info'].forEach(channel => {
                  const key = this.dcChainKey(chain,channel)
                       // get dc_Chainid by key
                  const dc_Chainid = dcChainIdMap.get(key)
                  // console.log(key,"=============key====>>>")
                  // console.log(dc_Chainid,"=============dc_Chainid====>>>")
                  if (dcChainIdMap.size > 0 && dc_Chainid) {
                          if (supportChainsIdMap.has(dc_Chainid)){
                              // dcChainId find and only include config chain
                              channel['chain_id'] = dc_Chainid
                              const result =
                                  channelsObj[dc_Chainid][
                                      `${channel.counterparty.channel_id}/${channel.counterparty.port_id}/${channel.channel_id}/${channel.port_id}`
                                      ];
                              if (result) {
                                  channel['counterparty']['state'] = result.split(Delimiter)[1];
                              }
                          }

                      }
              });

              // filter unconfig channels
              chain['ibc_info'] = chain['ibc_info'] && chain['ibc_info'].filter(channel => {
                  return channel.hasOwnProperty('chain_id');
              });

              const ibcInfoGroupBy = groupBy(chain['ibc_info'], 'chain_id');
              const ibcInfo = [];
              // console.log(ibcInfoGroupBy,"========ibcInfoGroupBy============>>")
              Object.keys(ibcInfoGroupBy) && Object.keys(ibcInfoGroupBy).forEach(chain_id => {
                  ibcInfo.push({ chain_id, paths: ibcInfoGroupBy[`${chain_id}`] });
                  // console.log(ibcInfo,"========ibcInfo= in for===========>>")
              });
              chain['ibc_info'] = ibcInfo;

              // update
              if(!chain.isManual){
                  this.chainConfigModel.updateChain(chain);
              }
          });
      });
      Logger.debug("Quit to handle Chain Channels info")
  }

  // get and sync chainConfig datas
  async parseChainConfig() {
    const allChains = await this.findAllConfig();
    // request configed allchannels
    Promise.all(
      allChains.map(async chain => {
        let channels = await ChainHttp.getIbcChannels(chain.lcd,chain.lcd_api_path.channels_path);
        channels && channels.map(channel => {
          channel.sc_chain_id = chain.chain_id;
        });
        return Promise.resolve({
          chain_id: chain.chain_id,
          chain_name: chain.chain_name,
          lcd: chain.lcd,
          lcd_api_path: chain.lcd_api_path,
          icon: chain.icon,
          ibc_info: channels ? channels : [],
          isManual: chain?.is_manual
        });
      }),
    ).then(allChains => {
      const channelsObj = {};
      const allChainsId = allChains && allChains.map(chain => {
        return chain['chain_id'];
      });

      // set channelsObj datas
      allChains && allChains.forEach(chain => {
        channelsObj[chain['chain_id']] = {};
        chain['ibc_info'] && chain['ibc_info'].forEach(channel => {
          channelsObj[chain['chain_id']][
            `${channel.channel_id}/${channel.port_id}/${channel.counterparty.channel_id}/${channel.counterparty.port_id}`
          ] = `${channel.sc_chain_id}${Delimiter}${channel.state}`;
        });
      });

      // get datas from channelsObj
      allChains && allChains.forEach(async chain => {
        chain['ibc_info'] && chain['ibc_info'].forEach(channel => {
          allChainsId && allChainsId.forEach(chainId => {
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

        // filter unconfig channels
        chain['ibc_info'] = chain['ibc_info'] && chain['ibc_info'].filter(channel => {
          return channel.hasOwnProperty('chain_id');
        });

        // groupby datas
        const ibcInfoGroupBy = groupBy(chain['ibc_info'], 'chain_id');
        const ibcInfo = [];
        Object.keys(ibcInfoGroupBy) && Object.keys(ibcInfoGroupBy).forEach(chain_id => {
          ibcInfo.push({ chain_id, paths: ibcInfoGroupBy[`${chain_id}`] });
        });
        chain['ibc_info'] = ibcInfo;

        // update
        if(!chain.isManual){
          this.chainConfigModel.updateChain(chain);
        }
      });
    });
  }
}
