import { IbcChainConfigType } from '../types/schemaTypes/ibc_chain_config.interface';
import { IbcChainType } from '../types/schemaTypes/ibc_chain.interface';
import { BaseResDto } from './base.dto'

export class IbcChainResDto extends BaseResDto {
  chain_id: string;
  chain_name: string;
  icon: string;

  constructor(value) {
    super()
    const { chain_id, chain_name, icon } = value;
    this.chain_id = chain_id;
    this.chain_name = chain_name;
    this.icon = icon;
  }

  static bundleData(value: any): IbcChainResDto[] {
    const datas: IbcChainResDto[] = value.map((item: IbcChainConfigType | IbcChainType) => {
      return new IbcChainResDto(item);
    });
    return datas;
  }
}

export class IbcChainResultResDto {
  all: IbcChainResDto[];
  active: IbcChainResDto[];
  inactive: IbcChainResDto[];
  constructor(value) {
    const { all, active, inactive } = value;
    this.all = all || [];
    this.active = active || [];
    this.inactive = inactive || [];
  }
}
