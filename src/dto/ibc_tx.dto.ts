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
  sc_addr: string;
  dc_addr: string;
  status: number;
  sc_chain_id: string;
  dc_chain_id:string;
  sc_tx_info: object;
  dc_tx_info?: object;
  base_denom: string;
  denoms: string[];
  create_at: string;
  tx_time: string;

  constructor(value: any) {
    super();
    const {
      sc_addr,
      dc_addr,
      dc_chain_id,
      sc_chain_id,
      status,
      sc_tx_info,
      dc_tx_info,
      base_denom,
      denoms,
      create_at,
      tx_time,
    } = value;
    this.sc_addr = sc_addr || '';
    this.dc_addr = dc_addr || '';
    this.status = status;
    this.sc_chain_id = sc_chain_id || '';
    this.dc_chain_id = dc_chain_id || '';
    this.sc_tx_info = sc_tx_info || {};
    this.dc_tx_info = dc_tx_info || {};
    this.base_denom = base_denom || '';
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
