import { IbcBaseDenomType } from '../types/schemaTypes/ibc_base_denom.interface';
import { BaseResDto } from './base.dto'

export class IbcBaseDenomResDto extends BaseResDto {
  chain_id: string;
  denom: string;
  symbol: string;
  scale: number;
  icon: string;
  is_main_token: boolean;
  create_at: string;
  update_at: string;

  constructor(value) {
    super()
    const {
      chain_id,
      denom,
      symbol,
      scale,
      icon,
      is_main_token,
      create_at,
      update_at,
    } = value;
    this.chain_id = chain_id;
    this.denom = denom;
    this.symbol = symbol;
    this.scale = scale;
    this.icon = icon;
    this.is_main_token = is_main_token;
    this.create_at = create_at;
    this.update_at = update_at;
  }

  static bundleData(value: IbcBaseDenomType[] = []) {
    const datas: IbcBaseDenomResDto[] = value.map((item: any) => {
      return {
        chain_id: item.chain_id,
        denom: item.denom,
        symbol: item.symbol,
        scale: item.scale,
        icon: item.icon,
        is_main_token: item.is_main_token,
        create_at: item.create_at,
        update_at: item.update_at,
      }
    })
    return datas
  }
}
