import { LcdChannelType } from '../types/lcd.interface';
export class LcdChannelDto {
  state: string;
  counterparty: {
    port_id: string;
    channel_id: string;
  };
  port_id: string;
  channel_id: string;
  sc_chain_id: string;

  constructor(value: any) {
    const { state, counterparty, port_id, channel_id } = value;
    this.state = state || '';
    this.counterparty = counterparty || {
      port_id: '',
      channel_id: '',
    };
    this.port_id = port_id || '';
    this.channel_id = channel_id || '';
  }

  static bundleData(value: LcdChannelType[]): LcdChannelDto[] {
    const datas: LcdChannelDto[] = value.map(item => {
      return new LcdChannelDto(item);
    });
    return datas;
  }
}
