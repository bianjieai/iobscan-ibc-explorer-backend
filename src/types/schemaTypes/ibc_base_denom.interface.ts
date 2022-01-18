export interface IbcBaseDenomType {
  chain_id: string;
  denom: string;
  symbol: string;
  scale: number;
  icon: string;
  is_main_token: boolean;
  ibc_info_hash_caculate?: string;
  create_at: string;
  update_at: string;
}

export class IbcBaseDenomDto {
  readonly chain_id: string;
  readonly denom: string;
  readonly symbol: string;
  readonly scale: number;
  readonly icon: string;
  readonly is_main_token: boolean;
  readonly create_at: string;
  readonly update_at: string;
}