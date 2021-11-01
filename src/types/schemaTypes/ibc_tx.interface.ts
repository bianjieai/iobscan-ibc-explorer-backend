export interface IbcTxType {
  record_id: string;
  sc_addr: string;
  dc_addr: string;
  sc_port: string;
  sc_channel: string;
  sc_chain_id: string;
  dc_port: string;
  dc_channel: string;
  dc_chain_id: string;
  sequence: string;
  status: number;
  sc_tx_info: object;
  dc_tx_info: object;
  refunded_tx_info?: object;
  log: object;
  denoms: object;
  base_denom: string;
  create_at: string;
  update_at: string;
  tx_time: string;
}

export interface IbcTxQueryType {
  useCount?: boolean;
  date_range?: number[];
  chain_id?: string;
  status?: number[];
  token?: { denom: string; chain_id: string }[];
  page_num?: number;
  page_size?: number;
}
