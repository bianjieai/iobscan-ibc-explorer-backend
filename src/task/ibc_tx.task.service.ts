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
import { IbcChannelSchema } from 'src/schema/ibc_channel.schema';
import { IbcTxType } from '../types/schemaTypes/ibc_tx.interface';
import { IbcDenomService } from 'src/service/ibc_denom.service';
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

  // 获取所有配置过的channels
  private channels_all_record;

  constructor(
    @InjectConnection() private readonly connection: Connection,
    private ibcDenomService: IbcDenomService,
  ) {
    this.getModels();
    this.doTask = this.doTask.bind(this);
  }

  async doTask(taskName?: TaskEnum): Promise<void> {
    const dateNow = String(Math.floor(new Date().getTime() / 1000));
    this.parseIbcTx(dateNow);
    this.changeIbcTxState(dateNow);
  }

  // 获取所有静态模型
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

  // ibcTx 第一阶段（transfer）
  async parseIbcTx(dateNow): Promise<void> {
    const allChains = await this.chainConfigModel.findAll();

    // todo 此处最好每条链设置一个定时任务
    allChains.forEach(async ({ chain_id }) => {
      // get taskRecord by chain_id
      let taskRecord = await this.ibcTaskRecordModel.findTaskRecord(chain_id);
      if (!taskRecord) {
        // 如果没有定时任务记录则新建
        await this.ibcTaskRecordModel.insertManyTaskRecord({
          task_name: `sync_${chain_id}_transfer`,
          status: IbcTaskRecordStatus.OPEN,
          height: 0,
          create_at: `${dateNow}`,
          update_at: `${dateNow}`,
        });
        // 此处需重新赋值taskRecord
        taskRecord = await this.ibcTaskRecordModel.findTaskRecord(chain_id);
      } else {
        // 如果定时任务记录状态是close则跳过
        if (taskRecord.status === IbcTaskRecordStatus.CLOSE) return;
      }
      // todo 此处需要改为初始化判断一次record
      const txModel = await this.connection.model(
        'txModel',
        TxSchema,
        `sync_${chain_id}_tx`,
      );

      let txs = [];
      const txsByLimit = await txModel.queryTxListSortHeight({
        type: TxType.transfer,
        height: taskRecord.height,
        limit: RecordLimit,
      });

      const txsByHeight = txsByLimit.length
        ? await txModel.queryTxListByHeight(
            TxType.transfer,
            txsByLimit[txsByLimit.length - 1].height,
          )
        : [];

      // 对象数组去重
      const hash = {};
      txs = [...txsByLimit, ...txsByHeight].reduce((txsResult, next) => {
        hash[next.tx_hash]
          ? ''
          : (hash[next.tx_hash] = true) && txsResult.push(next);
        return txsResult;
      }, []);

      // 遍历tx
      txs.forEach((tx, txIndex) => {
        const height = tx.height;
        const log = tx.log;
        const time = tx.time;
        const hash = tx.tx_hash;
        const status = tx.status;
        const fee = tx.fee;
        // 遍历msg
        tx.msgs.forEach(async (msg, msgIndex) => {
          // 判断msg Type为transfer
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
              denoms: [],
              base_denom: '',
              create_at: '',
              update_at: '',
            };
            // 根据tx 状态判断ibcTx状态
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

            const sc_chain_id = JSON.parse(JSON.stringify(tx)).chain_id;
            const sc_port = msg.msg.source_port;
            const sc_channel = msg.msg.source_channel;
            const sc_addr = msg.msg.sender;
            const dc_addr = msg.msg.receiver;
            const sc_denom = msg.msg.token.denom;
            const msg_amount = msg.msg.token;
            const {
              dc_port,
              dc_channel,
              sequence,
              base_denom,
              denom_path,
            } = this.getIbcInfoFromEventsMsg(tx, msgIndex);

            // search dc_chain_id by sc_chain_id、sc_port、sc_channel、dc_port、dc_channel
            let dc_chain_id = '';
            const result = await this.chainConfigModel.findDcChain({
              sc_chain_id,
              sc_port,
              sc_channel,
              dc_port,
              dc_channel,
            });

            if (result && result.ibc_info && result.ibc_info.length) {
              result.ibc_info.forEach(info_item => {
                info_item.paths.forEach(path_item => {
                  if (
                    path_item.channel_id === sc_channel &&
                    path_item.port_id === sc_port &&
                    path_item.counterparty.channel_id === dc_channel &&
                    path_item.counterparty.port_id === dc_port
                  ) {
                    dc_chain_id = info_item.chain_id;
                  }
                });
              });
            } else {
              dc_chain_id = '';
            }

            ibcTx.record_id = `${sc_port}${sc_channel}${dc_port}${dc_channel}${sequence}${sc_chain_id}`;
            ibcTx.sc_addr = sc_addr;
            ibcTx.dc_addr = dc_addr;
            ibcTx.sc_port = sc_port;
            ibcTx.sc_channel = sc_channel;
            ibcTx.sc_chain_id = sc_chain_id;
            ibcTx.dc_port = dc_port;
            ibcTx.dc_channel = dc_channel;
            ibcTx.dc_chain_id = dc_chain_id;
            ibcTx.sequence = sequence;
            ibcTx.denoms.push(sc_denom);
            ibcTx.base_denom = base_denom;
            ibcTx.create_at = dateNow;
            ibcTx.update_at = dateNow;
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
            // 如果没找到对应的 dc_chain_id
            if (!dc_chain_id) {
              ibcTx.status = IbcTxStatus.SETTING;
            }
            // 插入 ibcTx
            await this.ibcTxModel.insertManyIbcTx(ibcTx, async err => {
              taskRecord.height = height;

              // 更新 taskRecord
              await this.ibcTaskRecordModel.updateTaskRecord(taskRecord);

              if (ibcTx.status !== IbcTxStatus.FAILED) {
                // 统计denom (ibc交易类型为只要不是Failed都会统计第一段)
                this.parseDenom(
                  ibcTx.sc_chain_id,
                  sc_denom,
                  ibcTx.base_denom,
                  denom_path,
                  !Boolean(denom_path),
                  dateNow,
                  dateNow,
                  dateNow,
                );

                // 统计channel (ibc交易类型为只要不是Failed都会统计第一段)
                this.parseChannel(
                  sc_chain_id,
                  dc_chain_id,
                  sc_channel,
                  dateNow,
                );

                // 统计chain (ibc交易类型为只要不是Failed都会统计第一段)
                this.parseChain(sc_chain_id, dateNow);
              }
            });
          }
        });
      });
    });
  }

  // ibcTx 第二阶段（recv_packet || timoout_packet）
  async changeIbcTxState(dateNow): Promise<void> {
    const ibcTxs = await this.ibcTxModel.queryTxList({
      status: IbcTxStatus.PROCESSING,
      limit: RecordLimit,
    });

    ibcTxs.forEach(async ibcTx => {
      if (!ibcTx.dc_chain_id) return;

      const txModel = await this.connection.model(
        'txModel',
        TxSchema,
        `sync_${ibcTx.dc_chain_id}_tx`,
      );

      const txs = await txModel.queryTxListByPacketId({
        type: TxType.recv_packet,
        limit: RecordLimit,
        status: TxStatus.SUCCESS,
        packet_id: ibcTx.sc_tx_info.msg.msg.packet_id,
      });

      // txs have status is success's tx?
      if (txs.length) {
        const counter_party_tx = txs[0];
        counter_party_tx &&
          counter_party_tx.msgs.forEach(msg => {
            if (
              msg.type === TxType.recv_packet &&
              ibcTx.sc_tx_info.msg.msg.packet_id === msg.msg.packet_id
            ) {
              const {
                dc_denom,
                dc_denom_origin,
              } = this.ibcDenomService.getDcDenom(msg);
              ibcTx.status = IbcTxStatus.SUCCESS;
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
              // 插入目标链denom
              ibcTx.denoms.push(dc_denom);
              const denom_path = dc_denom_origin.replace(
                `/${ibcTx.base_denom}`,
                '',
              );
              // 更新ibcTx
              this.ibcTxModel.updateIbcTx(ibcTx);

              // 统计denom（当ibc交易成功时统计第二段）
              this.parseDenom(
                ibcTx.dc_chain_id,
                dc_denom,
                ibcTx.base_denom,
                denom_path,
                !Boolean(denom_path),
                dateNow,
                dateNow,
                dateNow,
              );

              // 统计Channel（当ibc交易成功时统计第二段）
              this.parseChannel(
                ibcTx.sc_chain_id,
                ibcTx.dc_chain_id,
                ibcTx.dc_channel,
                dateNow,
              );

              // 统计Chain
              this.parseChain(ibcTx.dc_channel, dateNow);
            }
          });
      } else {
        const blockModel = await this.connection.model(
          'blockModel',
          IbcBlockSchema,
          `sync_${ibcTx.dc_chain_id}_block`,
        );
        // 超时时间  目标链时间
        // 超时高度  目标链高度
        const { height, time } = await blockModel.findBlockByLastHeight();
        const ibcHeight =
          ibcTx.sc_tx_info.msg.msg.timeout_height.revision_height;
        const ibcTime = ibcTx.sc_tx_info.msg.msg.timeout_timestamp;
        if (ibcHeight < height || ibcTime < time) {
          const txModel = await this.connection.model(
            'txModel',
            TxSchema,
            `sync_${ibcTx.sc_chain_id}_tx`,
          );
          const refunded_tx = await txModel.queryTxListByPacketId({
            type: TxType.timeout_packet,
            limit: RecordLimit,
            status: TxStatus.SUCCESS,
            packet_id: ibcTx.sc_tx_info.msg.msg.packet_id,
          })[0];
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
                // 更新ibcTx
                this.ibcTxModel.updateIbcTx(ibcTx);
              }
            });
        }
      }
    });
  }

  // 获取dc_port、dc_channel、sequence
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
                // todo JSON.parse 容错
                const packet_data = JSON.parse(attr.value);
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

  // 统计Denom
  async parseDenom(
    chain_id,
    denom,
    base_denom,
    denom_path,
    is_source_chain,
    create_at,
    update_at,
    dateNow,
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
      };
      await this.ibcDenomModel.insertManyDenom(ibcDenom);
    } else {
      ibcDenomRecord.update_at = dateNow;
      await this.ibcDenomModel.updateDenomRecord(ibcDenomRecord);
    }
  }

  // 统计Channel
  async parseChannel(
    sc_chain_id,
    dc_chain_id,
    channel_id,
    dateNow,
  ): Promise<void> {
    const channels_all_record = await this.getChannelsConfig();
    const isFindRecord = channels_all_record.find(channel => {
      return channel.record_id === `${sc_chain_id}${dc_chain_id}${channel_id}`;
    });

    // 如果当前channel_id没有配置过则跳过
    if (!isFindRecord) return;

    const ibcChannelRecord = await this.ibcChannelModel.findChannelRecord(
      `${sc_chain_id}${dc_chain_id}${channel_id}`,
    );

    if (!ibcChannelRecord) {
      const ibcChannel = {
        ...isFindRecord,
        update_at: dateNow,
        create_at: dateNow,
      };
      await this.ibcChannelModel.insertManyChannel(ibcChannel);
    } else {
      ibcChannelRecord.update_at = dateNow;
      await this.ibcChannelModel.updateChannelRecord(ibcChannelRecord);
    }
  }

  // 统计Chain
  async parseChain(chain_id, dateNow) {
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
        create_at: dateNow,
        update_at: dateNow,
      };

      this.ibcChainModel.insertManyChain(ibcChain);
    } else {
      ibcChainRecord.update_at = dateNow;
      this.ibcChainModel.updateChainRecord(ibcChainRecord);
    }
  }

  // 获取配置过的channels
  async getChannelsConfig() {
    const channels_all_record = [];

    const allChains = await this.chainConfigModel.findAll();
    allChains.forEach(chain => {
      chain.ibc_info.forEach(ibc_info_item => {
        ibc_info_item.paths.forEach(path_item => {
          // 统计转出channel_id
          channels_all_record.push({
            channel_id: path_item.channel_id,
            record_id: `${chain.chain_id}${ibc_info_item.chain_id}${path_item.channel_id}`,
            state: path_item.state,
          });
        });
      });
    });

    return channels_all_record;
  }
}
