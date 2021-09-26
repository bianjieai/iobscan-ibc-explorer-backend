import { HttpService, Injectable } from '@nestjs/common';
import { Logger } from '../../logger';
import { LcdChannelDto } from '../../dto/http.dto'
import { LcdChannelType } from '../../types/lcd.interface'
// todo 需要对lcd 增加dto
@Injectable()
export class ChainHttp {
  static async getIbcChannels(lcdAddr) {
    const ibcChannelsUrl: string = `${lcdAddr}/ibc/core/channel/v1beta1/channels`;
    try {
      let ibcChannels: LcdChannelType[] = await new HttpService()
        .get(ibcChannelsUrl)
        .toPromise()
        .then(result => result.data.channels);
      if (ibcChannels) {
        return LcdChannelDto.bundleData(ibcChannels);
      } else {
        Logger.warn(
          'api-error:',
          'there is no result of total_supply from lcd',
        );
      }
    } catch (e) {
      Logger.warn(`api-error from ${ibcChannelsUrl}`, e);
    }
  }
}
