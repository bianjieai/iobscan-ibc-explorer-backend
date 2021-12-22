/* eslint-disable @typescript-eslint/camelcase */
import { IbcDenomType } from '../types/schemaTypes/ibc_denom.interface';
import { BaseResDto } from './base.dto'

export class IbcDenomResDto extends BaseResDto {
    chain_id: string;
    denom: string;
    base_denom: string;
    denom_path: string;
    is_source_chain: boolean;
    create_at: string;
    update_at: string;
    symbol: string;

  constructor(value) {
    super()
    const {
      chain_id,
      denom,
      base_denom,
      denom_path,
      is_source_chain,
      create_at,
      update_at,
      symbol,
    } = value;
    this.chain_id = chain_id;
    this.denom = denom;
    this.base_denom = base_denom;
    this.denom_path = denom_path;
    this.is_source_chain = is_source_chain;
    this.symbol = symbol;
    this.create_at = create_at;
    this.update_at = update_at;
  }

  static bundleData(value: IbcDenomType[] = []) {
    const datas: IbcDenomResDto[] = value.map((item: any) => {
      item.chain_id = item.chain_id.replace(new RegExp("\_", "g"),"-")
      return {
        chain_id: item.chain_id,
        denom: item.denom,
        base_denom: item.base_denom,
        denom_path: item.denom_path,
        is_source_chain: item.is_source_chain,
        symbol: item.symbol,
        create_at: item.create_at,
        update_at: item.update_at,
      }
    })
    return datas
  }
}
