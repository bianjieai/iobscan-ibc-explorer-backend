import { LcdChannelType } from '../types/lcd.interface';
export class LcdChannelDto {
  channels: LcdChannelType[];
  constructor(value) {
    const { channels } = value;
    this.channels = channels || [];
  }
}
