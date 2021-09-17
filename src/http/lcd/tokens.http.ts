import { HttpService, Injectable } from '@nestjs/common';
import {cfg} from "../../config/config";
import {TokensLcdDto, IbcTraceDto} from "../../dto/http.dto";
import {Logger} from "../../logger";
@Injectable()
export class TokensHttp {
    async getTokens () {
        const TokensUrl = `${cfg.serverCfg.lcdAddr}/irismod/token/tokens`
        try {
            const TokensData: any = await new HttpService().get(TokensUrl).toPromise().then(result => result.data)
            if (TokensData && TokensData.Tokens) {
                return TokensLcdDto.bundleData(TokensData.Tokens);
            } else {
                Logger.warn('api-error:', 'there is no result of tokens from lcd');
            }
        } catch (e) {
            Logger.warn(`api-error from ${TokensUrl}`, e)
        }
    }

    async getCirculationtTokens () {
        const CirculationtTokensUrl = `https://rpc.irisnet.org/token-stats/circulation`
        try {
            const CirculationtTokens: any = await new HttpService().get(CirculationtTokensUrl).toPromise().then(result => result.data)
            if (CirculationtTokens) {
                return CirculationtTokens;
            } else {
                Logger.warn('api-error:', `there is no result of circulationt_tokens from ${CirculationtTokensUrl}`);
            }
        } catch (e) {
            Logger.warn(`api-error from ${CirculationtTokensUrl}`, e)
        }
    }

    async getIbcTraces(hash: string): Promise<IbcTraceDto> {
      const ibcTracesUrl = `${cfg.serverCfg.lcdAddr}/ibc/applications/transfer/v1beta1/denom_traces/${hash}`
      try {
        const TracesData: IbcTraceDto = await new HttpService().get(ibcTracesUrl).toPromise().then(result => result.data)
        if (TracesData) {
          return new IbcTraceDto(TracesData);
        } else {
          Logger.warn('api-error:', `there is no result of ibcTraces from ${ibcTracesUrl}`);
        }
      } catch (error) {
        Logger.warn(`api-error from ${ibcTracesUrl}`, error)
      }
    }
}
