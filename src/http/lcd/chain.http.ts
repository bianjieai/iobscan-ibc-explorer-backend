/* eslint-disable @typescript-eslint/camelcase */
import { HttpService, Injectable } from '@nestjs/common';
import { Logger } from '../../logger';
import { LcdChannelDto, LcdDenomDto } from '../../dto/http.dto'
import { LcdChannelType, DenomType } from '../../types/lcd.interface'
import {cfg} from "../../config/config";
@Injectable()
export class ChainHttp {
  static async getIbcChannels(lcdAddr,channelsPath) {
    const ibcChannelsUrl = `${lcdAddr}${channelsPath}?pagination.offset=${cfg.channels.offset}&pagination.limit=${cfg.channels.limit}&pagination.count_total=true`;
    try {
      const ibcChannels: LcdChannelType[] = await new HttpService()
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
      // Logger.warn(`api-error from ${ibcChannelsUrl}`, e);
      Logger.warn(`api-error from ${ibcChannelsUrl} error`);
    }
  }

  static async getDenomByLcdAndHash(lcdAddr, hash) {
    const ibcDenomTraceUrl = `${lcdAddr}/ibc/applications/transfer/v1beta1/denom_traces/${hash}`;
    try {
      const denom_trace: DenomType = await new HttpService()
        .get(ibcDenomTraceUrl)
        .toPromise()
        .then(result => result.data.denom_trace);
      if (denom_trace) {
        return new LcdDenomDto(denom_trace);
      } else {
        Logger.warn(
          'api-error:',
          'there is no result of total_supply from lcd',
        );
      }
    } catch (e) {
      // Logger.warn(`api-error from ${ibcDenomTraceUrl}`, e);
      Logger.warn(`api-error from ${ibcDenomTraceUrl} error`);
    }
  }
}
