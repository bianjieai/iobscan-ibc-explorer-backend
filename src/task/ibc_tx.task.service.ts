import { Injectable } from '@nestjs/common';
import { Connection } from 'mongoose';
import { InjectConnection } from '@nestjs/mongoose';
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
} from '../constant';

@Injectable()
export class IbcTxTaskService {
  private ibcTaskRecordModel;
  private chainModel;
  private ibcTxModel;
  private ibcDenomModel;
  private ibcChannelModel;

  constructor(
    @InjectConnection() private readonly connection: Connection,
    private ibcDenomService: IbcDenomService,
  ) {
    this.getModels();
    this.doTask = this.doTask.bind(this);
  }

  async doTask(taskName?: TaskEnum): Promise<void> {
    const dateNow = String(new Date().getTime());
    this.parseIbcTx(dateNow);
    this.changeIbcTxState(dateNow);
  }

  // 获取所有静态模型
  async getModels() {
    // ibcTaskRecordModel
    this.ibcTaskRecordModel = await this.connection.model(
      'ibcTaskRecordModel',
      IbcTaskRecordSchema,
      'ibc_task_record',
    );

    // chainModel
    this.chainModel = await this.connection.model(
      'chainModel',
      IbcChainSchema,
      'chain_config',
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

    // ibcChannelModel
    this.ibcChannelModel = await this.connection.model(
      'ibcChannelModel',
      IbcChannelSchema,
      'ibc_channel',
    );
  }

  // ibcTx 第一阶段（transfer）
  async parseIbcTx(dateNow) {
    const allChains = await this.chainModel.findAll();

    allChains.forEach(async ({ chain_id }) => {
      // get taskRecord by chain_id
      const taskRecord = await this.ibcTaskRecordModel.findTaskRecord(chain_id);

      if (!taskRecord) {
        // 如果没有定时任务记录则新建
        await this.ibcTaskRecordModel.insertManyTaskRecord({
          task_name: `sync_${chain_id}_transfer`,
          status: 'open',
          height: 0,
          create_at: `${dateNow}`,
          update_at: `${dateNow}`,
        });
      } else {
        const txModel = await this.connection.model(
          'txModel',
          TxSchema,
          `sync_${chain_id}_tx`,
        );

        const txs = await txModel.queryTxListSortHeight({
          type: TxType.transfer,
          height: taskRecord.height,
          limit: RecordLimit,
        });

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
                log: '',
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
              const msg_amount = msg;

              const {
                dc_port,
                dc_channel,
                sequence,
                base_denom,
                denom_path,
              } = this.getDcMsg(tx, msgIndex);
              // search dc_chain_id by sc_chain_id、sc_port、sc_channel、dc_port、dc_channel
              const dc_chain_id = await this.chainModel.findDcChainId({
                sc_chain_id,
                sc_port,
                sc_channel,
                dc_port,
                dc_channel,
              });

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
              };
              ibcTx.log = log;
              // 如果没找到对应的 dc_chain_id
              if (!dc_chain_id) {
                ibcTx.status = IbcTxStatus.SETTING;
              }
              // 插入 ibcTx
              await this.ibcTxModel.insertTx(ibcTx, async err => {
                taskRecord.height = height;

                // 记录最后一条交易记录内最后一个msg的高度
                if (
                  txIndex === RecordLimit - 1 &&
                  msgIndex === tx.msgs.length - 1 &&
                  !err
                ) {
                  // 更新 taskRecord
                  const result = await this.ibcTaskRecordModel.updateTaskRecord(
                    taskRecord,
                  );

                  if (ibcTx.status !== IbcTxStatus.FAILED) {
                    // 记录denom (ibc交易类型为只要不是Failed都会统计第一段)
                    this.parseDenom(
                      ibcTx.sc_chain_id,
                      sc_denom,
                      ibcTx.base_denom,
                      denom_path,
                      !Boolean(denom_path),
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
                  }
                }
              });
            }
          });
        });
      }
    });
  }

  // ibcTx 第二阶段（recv_packet || timoout_packet）
  async changeIbcTxState(dateNow) {
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
        packet_id: ibcTx.sc_tx_info.msg_amount.msg.packet_id,
      });

      // txs have status is success's tx?
      if (txs.length) {
        const counter_party_tx =
          txModel &&
          (
            await txModel.queryTxListByPacketId({
              type: TxType.recv_packet,
              limit: RecordLimit,
              packet_id: ibcTx.sc_tx_info.msg_amount.msg.packet_id,
              status: TxStatus.SUCCESS,
            })
          )[0];
        counter_party_tx &&
          counter_party_tx.msgs.forEach(msg => {
            if (
              msg.type === TxType.recv_packet &&
              ibcTx.sc_tx_info.msg_amount.msg.packet_id === msg.msg.packet_id
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
                msg_amount: msg,
              };
              ibcTx.update_at = dateNow;
              // 插入目标链denom
              ibcTx.denoms.push(dc_denom);
              const denom_path = dc_denom_origin.replace(
                `/${ibcTx.base_denom}`,
                '',
              );
              // 更新ibcTx
              const result = this.ibcTxModel.updateIbcTx(ibcTx);

              // 统计denom（当ibc交易成功时统计第二段）
              this.parseDenom(
                ibcTx.dc_chain_id,
                dc_denom,
                ibcTx.base_denom,
                denom_path,
                !Boolean(denom_path),
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
          ibcTx.sc_tx_info.msg_amount.msg.timeout_height.revision_height;
        const ibcTime = ibcTx.sc_tx_info.msg_amount.msg.timeout_timestamp;
        if (ibcHeight > height || ibcTime > time) {
          const txModel = await this.connection.model(
            'txModel',
            TxSchema,
            `sync_${ibcTx.sc_chain_id}_tx`,
          );
          const refunded_tx = await txModel.queryTxListByPacketId({
            type: TxType.timeout_packet,
            limit: RecordLimit,
            status: TxStatus.SUCCESS,
            packet_id: ibcTx.sc_tx_info.msg_amount.msg.packet_id,
          })[0];
          refunded_tx &&
            refunded_tx.msgs.forEach(msg => {
              if (
                msg.type === TxType.timeout_packet &&
                ibcTx.sc_tx_info.msg_amount.msg.packet_id === msg.msg.packet_id
              ) {
                ibcTx.status = IbcTxStatus.REFUNDED;
                ibcTx.refunded_tx_info = {
                  hash: refunded_tx.tx_hash,
                  status: refunded_tx.status,
                  time: refunded_tx.time,
                  height: refunded_tx.height,
                  fee: refunded_tx.fee,
                  msg_amount: msg,
                };
                ibcTx.update_at = dateNow;
                // 更新ibcTx
                const result = this.ibcTxModel.updateIbcTx(ibcTx);
              }
            });
        }
      }
    });
  }

  // 获取dc_port、dc_channel、sequence
  getDcMsg(tx, msgIndex) {
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
                const packet_data = JSON.parse(attr.value);
                const denomOrigin = packet_data.denom;
                const denomOriginSplit = denomOrigin.split('/');
                msg.base_denom = denomOriginSplit[denomOriginSplit.length - 1];
                msg.denom_path = denomOriginSplit
                  .slice(0, denomOriginSplit.length - 1)
                  .join('');
              default:
                break;
            }
          });
        }
      });
    return msg;
  }

  // 统计Denom
  parseDenom(
    chain_id,
    denom,
    base_denom,
    denom_path,
    is_source_chain,
    create_at,
    update_at,
  ) {
    const ibcDenom = {
      chain_id,
      denom,
      base_denom,
      base_denom_chain_id: '',
      denom_path,
      is_source_chain,
      create_at,
      update_at,
    };
    this.ibcDenomModel.insertManyDenom(ibcDenom, err => {
      err && err.code === 11000 && console.log('denom重复');
    });
  }

  // 统计Channel
  parseChannel(sc_chain_id, dc_chain_id, channel_id, dateNow) {
    const ibcChannelRecord = this.ibcChannelModel.findChannelRecord(
      `${sc_chain_id}${dc_chain_id}${channel_id}`,
    );

    if (!ibcChannelRecord) {
      const ibcChannel = {
        channel_id: channel_id,
        record_id: `${sc_chain_id}${dc_chain_id}${channel_id}`,
        update_at: dateNow,
        create_at: dateNow,
      };
      this.ibcChannelModel.insertManyChannel(ibcChannel, err => {
        err && err.code === 11000 && console.log('channel重复');
      });
    } else {
      ibcChannelRecord.update_at = dateNow;
      this.ibcChannelModel.updateChannelRecord(ibcChannelRecord);
    }
  }
}
