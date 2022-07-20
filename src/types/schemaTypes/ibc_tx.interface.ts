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
  sc_tx_info: { time?:number; status?:number; hash?:string;height?:number;fee?:object;msg_amount?:object;msg?:object };
  dc_tx_info: { time?:number; hash?:string };
  refunded_tx_info?: { time?:number };
  substate: number;
  log: object;
  denoms: object;
  base_denom: string;
  create_at: number;
  update_at: number;
  tx_time: number;
  end_time?: number;
  retry_times?: number;
  next_try_time?: number;
}

export interface IbcTxQueryType {
  useCount?: boolean;
  date_range?: number[];
  chain_id?: string;
  status?: number[];
  token?: string[];
  page_num?: number;
  page_size?: number;
}

export  interface AggregateResult24hr {
  _id:{sc_chain_id?:string,dc_chain_id?:string,sc_channel?:string,dc_channel?:string};
}