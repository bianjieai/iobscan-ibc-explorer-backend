export interface LcdChannelType {
  state: string;
  counterparty: {
    port_id: string;
    channel_id: string;
  };
  port_id: string;
  channel_id: string;
}
