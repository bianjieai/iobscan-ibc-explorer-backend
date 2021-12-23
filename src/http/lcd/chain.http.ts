/* eslint-disable @typescript-eslint/camelcase */
import {HttpService, Injectable} from '@nestjs/common';
import {Logger} from '../../logger';
import {LcdChannelDto, LcdDenomDto} from '../../dto/http.dto'
import {LcdChannelType, DenomType,LcdChannelClientState} from '../../types/lcd.interface'
import {cfg} from "../../config/config";

@Injectable()
export class ChainHttp {
    static async getIbcChannels(lcdAddr, channelsPath: string) {
        let rgexOffset = "\\OFFSET", regexLimit = "\\LIMIT"
        channelsPath = channelsPath.replace(new RegExp(rgexOffset, "g"), <string>cfg.channels.offset);
        channelsPath = channelsPath.replace(new RegExp(regexLimit, "g"), <string>cfg.channels.limit);

        const ibcChannelsUrl = `${lcdAddr}${channelsPath}`;
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

    static async getDcChainIdByScChannel(lcdAddr, clientStatePath: string, port: string, channel: string) {
        let rgexPort = "\\PORT", regexChannel = "\\CHANNEL"
        clientStatePath = clientStatePath.replace(new RegExp(rgexPort, "g"), port);
        clientStatePath = clientStatePath.replace(new RegExp(regexChannel, "g"), channel);
        const ibcClientStateUrl = `${lcdAddr}${clientStatePath}`;
        try {
            const ibcClientState: LcdChannelClientState = await new HttpService()
                .get(ibcClientStateUrl)
                .toPromise()
                .then(result => result.data);
            if (ibcClientState) {
                return ibcClientState.identified_client_state.client_state.chain_id;
            } else {
                Logger.warn(
                    'api-error:',
                    'there is no result of total_supply from lcd',
                );
            }
        } catch (e) {
            // Logger.warn(`api-error from ${ibcChannelsUrl}`, e);
            Logger.warn(`api-error from ${ibcClientStateUrl} error`);
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
