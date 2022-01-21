export interface IbcDenomType {
  chain_id: string;
  denom: string;
  base_denom: string;
  denom_path: string;
  symbol: string;
  is_source_chain: boolean;
  create_at: string;
  update_at: string;
}

export  interface AggregateBaseDenomCnt {
  _id:{base_denom?:string,chain_id?:string};
}

export class IbcDenomDto {
  readonly chain_id: string;
  readonly denom: string;
  readonly symbol: string;
}