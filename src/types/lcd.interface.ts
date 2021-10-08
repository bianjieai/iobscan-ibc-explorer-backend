export interface LcdChannelType {
  state: string;
  counterparty: {
    port_id: string;
    channel_id: string;
  };
  ordering: string;
  connection_hops: string[];
  port_id: string;
  channel_id: string;
  version: string;
}

export interface DenomType {
  path: string;
  base_denom: string;
}
