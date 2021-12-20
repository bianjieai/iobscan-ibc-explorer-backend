export interface IbcDenomType {
  chain_id: string;
  denom: string;
  base_denom: string;
  denom_path: string;
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
  readonly base_denom: string;
  readonly denom_path: string;
  readonly is_source_chain: boolean;
  readonly create_at: string;
  readonly update_at: string;
}