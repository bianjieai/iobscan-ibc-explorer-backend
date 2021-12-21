/* eslint-disable @typescript-eslint/camelcase */
import { PagingReqDto } from './base.dto';
import { ApiPropertyOptional } from '@nestjs/swagger';
import { BaseResDto } from './base.dto';

export class IbcTxListReqDto extends PagingReqDto {
  @ApiPropertyOptional()
  date_range?: string;
  @ApiPropertyOptional()
  status?: string;
  @ApiPropertyOptional()
  chain_id?: string;
  @ApiPropertyOptional()
  symbol?: string;
  @ApiPropertyOptional()
  denom?: string;
  @ApiPropertyOptional()
  page_num?: number;
  @ApiPropertyOptional()
  page_size?: number;
  @ApiPropertyOptional()
  start_time?: boolean

  static convert(value: any): any {
    return super.convert(value);
  }
}

export class IbcTxResDto extends BaseResDto {
  record_id: string;
  sc_addr: string;
  dc_addr: string;
  sc_chain_id: string;
  dc_chain_id: string;
  sc_port: string;
  dc_port: string;
  sc_channel: string;
  dc_channel: string;
  status: number;
  sc_tx_info: object;
  dc_tx_info?: object;
  refunded_tx_info?: object;
  base_denom: string;
  update_at: string;
  sequence: string;
  log: object;
  denoms: string[];
  create_at: string;
  tx_time: string;

  constructor(value: any) {
    super();
    const {
      record_id,
      sc_addr,
      dc_addr,
      sc_chain_id,
      dc_chain_id,
      sc_port,
      dc_port,
      sc_channel,
      dc_channel,
      status,
      sc_tx_info,
      dc_tx_info,
      refunded_tx_info,
      base_denom,
      update_at,
      sequence,
      log,
      denoms,
      create_at,
      tx_time,
    } = value;
    this.record_id = record_id;
    this.sc_addr = sc_addr || '';
    this.dc_addr = dc_addr || '';
    this.sc_chain_id = sc_chain_id || '';
    this.dc_chain_id = dc_chain_id || '';
    this.sc_port = sc_port || '';
    this.dc_port = dc_port || '';
    this.sc_channel = sc_channel || '';
    this.dc_channel = dc_channel || '';
    this.status = status;
    this.sc_tx_info = sc_tx_info || {};
    this.dc_tx_info = dc_tx_info || {};
    this.refunded_tx_info = refunded_tx_info || {};
    this.base_denom = base_denom || '';
    this.update_at = update_at || '';
    this.sequence = sequence || '';
    this.log = log || {};
    this.denoms = denoms || [];
    this.create_at = create_at || '';
    this.tx_time = tx_time || '';
  }

  static bundleData(value: any): IbcTxResDto[] {
    const datas: IbcTxResDto[] = value.map((item: any) => {
      return new IbcTxResDto(item);
    });
    return datas;
  }
}
