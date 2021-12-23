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

export interface LcdChannelClientState{
  identified_client_state: {
    client_id: string;
    client_state: {
      chain_id: string;
    }
  }
}

export interface NodeInfoType {
  node_info: {
    network:string;
    version:string;
    listen_addr: string;
  }
}
