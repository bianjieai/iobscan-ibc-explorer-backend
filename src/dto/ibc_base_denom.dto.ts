export class IbcBaseDenomResDto {
  chain_id: string;
  denom: string;
  symbol: string;
  scale: number;
  icon: string;
  is_main_token: boolean;
  create_at: string;
  update_at: string;

  constructor(value) {
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
}
