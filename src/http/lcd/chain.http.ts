import { HttpService, Injectable } from '@nestjs/common';
import {Logger} from "../../logger";
@Injectable()
export class ChainHttp {
    static async getIbcChannels (lcdAddr) {
        const ibcChannelsUrl = `${lcdAddr}/ibc/core/channel/v1beta1/channels`
        try {
            let ibcChannels: any = await new HttpService().get(ibcChannelsUrl).toPromise().then(result => result.data.channels)
            if (ibcChannels) {
                return ibcChannels
            } else {
                Logger.warn('api-error:', 'there is no result of total_supply from lcd');
            }
        } catch (e) {
            Logger.warn(`api-error from ${ibcChannelsUrl}`, e)
        }
    }
}
