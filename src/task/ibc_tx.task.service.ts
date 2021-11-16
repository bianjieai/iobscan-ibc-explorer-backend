/* eslint-disable @typescript-eslint/camelcase */
import { Injectable } from '@nestjs/common';
import { Connection } from 'mongoose';
import { InjectConnection } from '@nestjs/mongoose';
import { IbcChainConfigSchema } from '../schema/ibc_chain_config.schema';
import { IbcChainSchema } from '../schema/ibc_chain.schema';
import { IbcDenomSchema } from '../schema/ibc_denom.schema';
import { IbcTxSchema } from '../schema/ibc_tx.schema';
import { TxSchema } from '../schema/tx.schema';
import { IbcBlockSchema } from '../schema/ibc_block.schema';
import { IbcTaskRecordSchema } from '../schema/ibc_task_record.schema';
import { IbcChannelSchema } from '../schema/ibc_channel.schema';
import { ChainHttp } from '../http/lcd/chain.http';
import { IbcTxType } from '../types/schemaTypes/ibc_tx.interface';
import { JSONparse } from '../util/util';
import { getDcDenom } from '../helper/denom.helper';

import {
  TaskEnum,
  TxType,
  TxStatus,
  IbcTxStatus,
  RecordLimit,
  IbcTaskRecordStatus,
} from '../constant';

@Injectable()
export class IbcTxTaskService {
  private ibcTaskRecordModel;
  private chainConfigModel;
  private ibcChainModel;
  private ibcTxModel;
  private ibcDenomModel;
  private ibcChannelModel;

  constructor(@InjectConnection() private readonly connection: Connection) {
    this.getModels();
    this.doTask = this.doTask.bind(this);
  }

  async doTask(taskName?: TaskEnum): Promise<void> {
    const dateNow = Math.floor(new Date().getTime() / 1000);
    this.parseIbcTx(dateNow);
    this.changeIbcTxState(dateNow);
  }

  // getModels
  async getModels(): Promise<void> {
    // ibcTaskRecordModel
    this.ibcTaskRecordModel = await this.connection.model(
      'ibcTaskRecordModel',
      IbcTaskRecordSchema,
      'ibc_task_record',
    );

    // chainConfigModel
    this.chainConfigModel = await this.connection.model(
      'chainConfigModel',
      IbcChainConfigSchema,
      'chain_config',
    );

    // ibcChainModel
    this.ibcChainModel = await this.connection.model(
      'ibcChainModel',
      IbcChainSchema,
      'ibc_chain',
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

    // ibcChannelModel
    this.ibcChannelModel = await this.connection.model(
      'ibcChannelModel',
      IbcChannelSchema,
      'ibc_channel',
    );
  }

  // ibcTx first（transfer）
  async parseIbcTx(dateNow): Promise<void> {
    const allChains = await this.chainConfigModel.findAll();
    for (const { chain_id } of allChains) {
      // get taskRecord by chain_id
      let taskRecord = await this.ibcTaskRecordModel.findTaskRecord(chain_id);
      if (!taskRecord) {
        await this.ibcTaskRecordModel.insertManyTaskRecord({
          task_name: `sync_${chain_id}_transfer`,
          status: IbcTaskRecordStatus.OPEN,
          height: 0,
          create_at: dateNow,
          update_at: dateNow,
        });
        taskRecord = await this.ibcTaskRecordModel.findTaskRecord(chain_id);
      } else {
        if (taskRecord.status === IbcTaskRecordStatus.CLOSE) return;
      }
      const txModel = await this.connection.model(
          'txModel',
          TxSchema,
          `sync_${chain_id}_tx`,
      );

      let txs = [];
      //根据块高排序 查询最后限制条数的交易
      const txsByLimit = await txModel.queryTxListSortHeight({
        type: TxType.transfer,
        height: taskRecord.height,
        limit: RecordLimit,
      });
      // 根据块高查询限制条数的最后一条交易
      const txsByHeight = txsByLimit.length
          ? await txModel.queryTxListByHeight(
              TxType.transfer,
              txsByLimit[txsByLimit.length - 1].height,
          )
          : [];
      //这里为啥做这样的一步操作
      const hash = {};
      txs = [...txsByLimit, ...txsByHeight].reduce((txsResult, next) => {
        hash[next.tx_hash]
            ? ''
            : (hash[next.tx_hash] = true) && txsResult.push(next);
        return txsResult;
      }, []);

      txs.forEach((tx, txIndex) => {
        const height = tx.height;
        const log = tx.log;
        const time = tx.time;
        const hash = tx.tx_hash;
        const status = tx.status;
        const fee = tx.fee;
        tx.msgs.forEach(async (msg, msgIndex) => {
          if (msg.type === TxType.transfer) {
            const ibcTx: IbcTxType = {
              record_id: '',
              sc_addr: '',
              dc_addr: '',
              sc_port: '',
              sc_channel: '',
              sc_chain_id: '',
              dc_port: '',
              dc_channel: '',
              dc_chain_id: '',
              sequence: '',
              status: 0,
              sc_tx_info: {},
              dc_tx_info: {},
              refunded_tx_info: {},
              log: {},
              denoms: {
                sc_denom: '',
                dc_denom: '',
              },
              base_denom: '',
              create_at: '',
              update_at: '',
              tx_time: '',
            };
            switch (tx.status) {
              case TxStatus.SUCCESS:
                ibcTx.status = IbcTxStatus.PROCESSING;
                break;
              case TxStatus.FAILED:
                ibcTx.status = IbcTxStatus.FAILED;
                break;
              default:
                break;
            }

            const sc_chain_id = chain_id;
            const sc_port = msg.msg.source_port;
            const sc_channel = msg.msg.source_channel;
            let dc_chain_id = '';
            let dc_channel = '';
            let dc_port = '';
            const sc_addr = msg.msg.sender;
            const dc_addr = msg.msg.receiver;
            const sc_denom = msg.msg.token.denom;
            const msg_amount = msg.msg.token;
            let base_denom = '';
            let denom_path = '';
            let sequence = '';
            let lcd = '';
            //根据 sc_chain_id  sc_port sc_channel 查询目标链的信息
            /*
            * 能否通过管道直接查询出目标链的数据
            * 不要在循环里查询数据库 单独拎出来
            * */
            const dcChainConfig = await this.chainConfigModel.findDcChain({
              sc_chain_id,
              sc_port,
              sc_channel,
            });

            if (
                dcChainConfig &&
                dcChainConfig.ibc_info &&
                dcChainConfig.ibc_info.length
            ) {
              lcd = dcChainConfig.lcd;
              dcChainConfig.ibc_info.forEach(info_item => {
                info_item.paths.forEach(path_item => {
                  if (
                      path_item.channel_id === sc_channel &&
                      path_item.port_id === sc_port
                  ) {
                    dc_chain_id = info_item.chain_id;
                    dc_channel = path_item.counterparty.channel_id;
                    dc_port = path_item.counterparty.port_id;
                  }
                });
              });
            }
            /*
            * 通过lcd 查询的话可以不可以在起个定时任务去查询
            * */

            if (ibcTx.status === IbcTxStatus.FAILED) {
              // get base_denom、denom_path from lcd API
              if (sc_denom.indexOf('ibc') !== -1 && lcd) {
                const result = await ChainHttp.getDenomByLcdAndHash(
                    lcd,
                    sc_denom.replace('ibc/', ''),
                );
                base_denom = result?.base_denom;
                denom_path = result?.denom_path;
              } else {
                base_denom = sc_denom;
              }
            } else {
              // get base_denom、denom_path、sequence from events
              const event_msg = this.getIbcInfoFromEventsMsg(tx, msgIndex);
              dc_port = event_msg.dc_port;
              dc_channel = event_msg.dc_channel;
              base_denom = event_msg.base_denom;
              denom_path = event_msg.denom_path;
              sequence = event_msg.sequence;
            }

            if (!dc_chain_id && ibcTx.status !== IbcTxStatus.FAILED) {
              ibcTx.status = IbcTxStatus.SETTING;
            }

            ibcTx.record_id = `${sc_port}${sc_channel}${dc_port}${dc_channel}${sequence}${sc_chain_id}${hash}${msgIndex}`;
            ibcTx.sc_addr = sc_addr;
            ibcTx.dc_addr = dc_addr;
            ibcTx.sc_port = sc_port;
            ibcTx.sc_channel = sc_channel;
            ibcTx.sc_chain_id = sc_chain_id;
            ibcTx.dc_port = dc_port;
            ibcTx.dc_channel = dc_channel;
            ibcTx.dc_chain_id = dc_chain_id;
            ibcTx.sequence = sequence;
            ibcTx.denoms['sc_denom'] = sc_denom;
            ibcTx.base_denom = base_denom;
            ibcTx.create_at = dateNow;
            ibcTx.update_at = dateNow;
            ibcTx.tx_time = tx.time;
            ibcTx.sc_tx_info = {
              hash,
              status,
              time,
              height,
              fee,
              msg_amount,
              msg,
            };
            ibcTx.log['sc_log'] = log;
/*
* TODO 采用批量插入的方式做 用事务做height的更新及数据的插入
* */

            await this.ibcTxModel.insertManyIbcTx(ibcTx, async err => {
              taskRecord.height = height;

              await this.ibcTaskRecordModel.updateTaskRecord(taskRecord);

              let is_base_denom = true;
              if (Boolean(denom_path) && denom_path.split('/').length > 1) {
                const dc_port = denom_path.split('/')[0];
                const dc_channel = denom_path.split('/')[1];
                const chainConfigRecord = await this.chainConfigModel.findScChain(
                    {
                      dc_chain_id: sc_chain_id,
                      dc_port,
                      dc_channel,
                    },
                );
                if (chainConfigRecord && chainConfigRecord.chain_id) {
                  is_base_denom = false;
                }
              }
              // parse denom
              const real_denom = ibcTx.status !== IbcTxStatus.FAILED;
              await this.parseDenom(
                  ibcTx.sc_chain_id,
                  sc_denom,
                  ibcTx.base_denom,
                  denom_path,
                  !Boolean(denom_path),
                  dateNow,
                  dateNow,
                  tx.time,
                  is_base_denom,
                  real_denom,
              );
              if (ibcTx.status !== IbcTxStatus.FAILED) {
                // parse channel
                ibcTx.dc_chain_id &&
                ibcTx.status !== IbcTxStatus.FAILED &&
                (await this.parseChannel(
                    sc_chain_id,
                    sc_channel,
                    dateNow,
                    dateNow,
                    tx.time,
                ));

                // parse chain
                ibcTx.dc_chain_id &&
                ibcTx.status !== IbcTxStatus.FAILED &&
                (await this.parseChain(
                    sc_chain_id,
                    dateNow,
                    dateNow,
                    tx.time,
                ));
              }
            });
          }
        });
      });
    }
/*    allChains.forEach(async ({ chain_id }) => {

    });*/

  }

  // ibcTx second（recv_packet || timoout_packet）
  async changeIbcTxState(dateNow): Promise<any> {
    const ibcTxs = await this.ibcTxModel.queryTxList({
      status: IbcTxStatus.PROCESSING,
      limit: RecordLimit,
    });
    //获取最新块高（chain_config中的所有链）
    return Promise.all(
      ibcTxs.map(async ibcTx => {
        if (!ibcTx.dc_chain_id) return;

        const txModel = await this.connection.model(
          'txModel',
          TxSchema,
          `sync_${ibcTx.dc_chain_id}_tx`,
        );

        //批量查询，组装map结构
        const counter_party_tx = await txModel.queryTxListByPacketId({
          type: TxType.recv_packet,
          status: TxStatus.SUCCESS,
          packet_id: ibcTx.sc_tx_info.msg.msg.packet_id,
        });

        // txs have status is success's tx?
        if (counter_party_tx) {
          counter_party_tx.msgs.forEach(async (msg, msgIndex) => {
            if (
              msg.type === TxType.recv_packet &&
              ibcTx.sc_tx_info.msg.msg.packet_id === msg.msg.packet_id
            ) {
              const { dc_denom, dc_denom_origin } = getDcDenom(msg);

              // add write_acknowledgement solution， value type is string;
              let result = '';
              const counter_party_tx_events = counter_party_tx.events_new.find(
                event_new => {
                  return event_new.msg_index === msgIndex;
                },
              );
              counter_party_tx_events &&
                counter_party_tx_events.events.forEach(event => {
                  if (event.type === 'write_acknowledgement') {
                    event.attributes.forEach(attribute => {
                      if (attribute.key === 'packet_ack') {
                        const resultObj = JSONparse(attribute.value);
                        result = resultObj.hasOwnProperty('error')
                          ? 'false'
                          : 'true';
                      }
                    });
                  }
                });
              ibcTx.status =
                result === 'false' ? IbcTxStatus.FAILED : IbcTxStatus.SUCCESS;
              ibcTx.dc_tx_info = {
                hash: counter_party_tx.tx_hash,
                status: counter_party_tx.status,
                time: counter_party_tx.time,
                height: counter_party_tx.height,
                fee: counter_party_tx.fee,
                msg_amount: msg.msg.token,
                msg,
              };
              ibcTx.update_at = dateNow;
              // ibcTx.tx_time = counter_party_tx.time;
              ibcTx.denoms['dc_denom'] = dc_denom;
              const denom_path =
                dc_denom_origin === ibcTx.base_denom
                  ? ''
                  : dc_denom_origin.replace(`/${ibcTx.base_denom}`, '');
              await this.ibcTxModel.updateIbcTx(ibcTx);

              // parse denom
              const real_denom = ibcTx.status === IbcTxStatus.SUCCESS;
              await this.parseDenom(
                ibcTx.dc_chain_id,
                dc_denom,
                ibcTx.base_denom,
                denom_path,
                !Boolean(denom_path),
                dateNow,
                dateNow,
                counter_party_tx.time,
                false,
                real_denom,
              );

              // parse Channel 【移除】
              await this.parseChannel(
                ibcTx.dc_chain_id,
                ibcTx.dc_channel,
                dateNow,
                dateNow,
                counter_party_tx.time,
              );

              // parse Chain 【移除】
              await this.parseChain(
                ibcTx.dc_chain_id,
                dateNow,
                dateNow,
                counter_party_tx.time,
              );
            }
          });
        } else {
          const blockModel = await this.connection.model(
            'blockModel',
            IbcBlockSchema,
            `sync_${ibcTx.dc_chain_id}_block`,
          );

          //放在开始位置查询，根据map的key找到对应chain的块高
          const lastBlock = await blockModel.findLatestBlock();
          if (!lastBlock) return;
          const { height, time } = lastBlock;
          const ibcHeight =
            ibcTx.sc_tx_info.msg.msg.timeout_height.revision_height;
          const ibcTime = ibcTx.sc_tx_info.msg.msg.timeout_timestamp;
          if (ibcHeight < height || ibcTime < time) {
            const txModel = await this.connection.model(
              'txModel',
              TxSchema,
              `sync_${ibcTx.sc_chain_id}_tx`,
            );
            // 取最新一条成功的timeout_packet，order by heiht -1,
            const refunded_tx_arr = await txModel.queryTxListByPacketId({
              type: TxType.timeout_packet,
              limit: RecordLimit,
              status: TxStatus.SUCCESS,
              packet_id: ibcTx.sc_tx_info.msg.msg.packet_id,
            });
            const refunded_tx = refunded_tx_arr[0];
            refunded_tx &&
              refunded_tx.msgs.forEach(msg => {
                if (
                  msg.type === TxType.timeout_packet &&
                  ibcTx.sc_tx_info.msg.msg.packet_id === msg.msg.packet_id
                ) {
                  ibcTx.status = IbcTxStatus.REFUNDED;
                  ibcTx.refunded_tx_info = {
                    hash: refunded_tx.tx_hash,
                    status: refunded_tx.status,
                    time: refunded_tx.time,
                    height: refunded_tx.height,
                    fee: refunded_tx.fee,
                    msg_amount: msg.msg.token,
                    msg,
                  };
                  ibcTx.update_at = dateNow;
                  // ibcTx.tx_time = refunded_tx.time;
                  this.ibcTxModel.updateIbcTx(ibcTx);
                }
              });
          }
        }
      }),
    );
  }

  // get dc_port、dc_channel、sequence
  getIbcInfoFromEventsMsg(
    tx,
    msgIndex,
  ): {
    dc_port: string;
    dc_channel: string;
    sequence: string;
    base_denom: string;
    denom_path: string;
  } {
    const msg = {
      dc_port: '',
      dc_channel: '',
      sequence: '',
      base_denom: '',
      denom_path: '',
    };

    tx.events_new[msgIndex] &&
      tx.events_new[msgIndex].events.forEach(evt => {
        if (evt.type === 'send_packet') {
          evt.attributes.forEach(attr => {
            switch (attr.key) {
              case 'packet_dst_port':
                msg.dc_port = attr.value;
                break;
              case 'packet_dst_channel':
                msg.dc_channel = attr.value;
                break;
              case 'packet_sequence':
                msg.sequence = attr.value;
                break;
              case 'packet_data':
                const packet_data = JSONparse(attr.value);
                const denomOrigin = packet_data.denom;
                const denomOriginSplit = denomOrigin.split('/');
                msg.base_denom = denomOriginSplit[denomOriginSplit.length - 1];
                msg.denom_path = denomOriginSplit
                  .slice(0, denomOriginSplit.length - 1)
                  .join('/');
              default:
                break;
            }
          });
        }
      });
    return msg;
  }

  // parse Denom
  async parseDenom(
    chain_id,
    denom,
    base_denom,
    denom_path,
    is_source_chain,
    create_at,
    update_at,
    tx_time,
    is_base_denom,
    real_denom,
  ): Promise<void> {
    const ibcDenomRecord = await this.ibcDenomModel.findDenomRecord(
      chain_id,
      denom,
    );

    if (!ibcDenomRecord) {
      const ibcDenom = {
        chain_id,
        denom,
        base_denom,
        denom_path,
        is_source_chain,
        create_at,
        update_at,
        tx_time,
        is_base_denom,
        real_denom,
      };
      await this.ibcDenomModel.insertManyDenom(ibcDenom);
    } else {
      ibcDenomRecord.update_at = update_at;
      ibcDenomRecord.real_denom = ibcDenomRecord.real_denom || real_denom;
      const currentTime = ibcDenomRecord.tx_time;
      if (tx_time > currentTime) {
        ibcDenomRecord.tx_time = tx_time;
      }
      await this.ibcDenomModel.updateDenomRecord(ibcDenomRecord);
    }
  }

  // parse Channel
  async parseChannel(
    chain_id,
    channel_id,
    create_at,
    update_at,
    tx_time,
  ): Promise<void> {
    const channels_all_record = await this.getChannelsConfig();
    const isFindRecord = channels_all_record.find(channel => {
      return channel.record_id === `${chain_id}${channel_id}`;
    });

    if (!isFindRecord) return;

    const ibcChannelRecord = await this.ibcChannelModel.findChannelRecord(
      `${chain_id}${channel_id}`,
    );

    if (!ibcChannelRecord) {
      const ibcChannel = {
        ...isFindRecord,
        update_at,
        create_at,
        tx_time,
      };
      await this.ibcChannelModel.insertManyChannel(ibcChannel);
    } else {
      ibcChannelRecord.update_at = update_at;
      const currentTime = ibcChannelRecord.tx_time;
      if (tx_time > currentTime) {
        ibcChannelRecord.tx_time = tx_time;
      }
      await this.ibcChannelModel.updateChannelRecord(ibcChannelRecord);
    }
  }
  // parse Chain
  async parseChain(
    chain_id,
    create_at,
    update_at,
    tx_time,
    num?,
  ): Promise<void> {
    const ibcChainRecord = await this.ibcChainModel.findById(chain_id);

    if (!ibcChainRecord) {
      const allChainsConfig = await this.chainConfigModel.findAll();
      const findChainConfig = allChainsConfig.find(chainConfig => {
        return chainConfig.chain_id === chain_id;
      });
      if (!findChainConfig) return;
      const ibcChain = {
        chain_id,
        chain_name: findChainConfig ? findChainConfig.chain_name : '',
        icon: findChainConfig ? findChainConfig.icon : '',
        create_at,
        update_at,
        tx_time,
      };

      await this.ibcChainModel.insertManyChain(ibcChain);
    } else {
      ibcChainRecord.update_at = update_at;
      const currentTime = ibcChainRecord.tx_time;
      if (tx_time > currentTime) {
        ibcChainRecord.tx_time = tx_time;
      }

      await this.ibcChainModel.updateChainRecord(ibcChainRecord);
    }
  }

  // get configed channels
  async getChannelsConfig() {
    const channels_all_record = [];

    const allChains = await this.chainConfigModel.findAll();
    allChains.forEach(chain => {
      chain.ibc_info.forEach(ibc_info_item => {
        ibc_info_item.paths.forEach(path_item => {
          channels_all_record.push({
            channel_id: path_item.channel_id,
            record_id: `${chain.chain_id}${path_item.channel_id}`,
            state: path_item.state,
          });
        });
      });
    });

    return channels_all_record;
  }
}
