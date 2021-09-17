import { HttpService, Injectable } from '@nestjs/common';
import { cfg } from '../../config/config';
import { Logger } from '../../logger';
import { NftCollectionDto } from '../../dto/http.dto';

@Injectable()
export class NftHttp {
    constructor(){

    }
    async queryNftsFromLcdByDenom(denom: string): Promise<NftCollectionDto> {
        const url: string = `${cfg.serverCfg.lcdAddr}/nft/nfts/collections/${denom}`;
        try {
            const data: any = await new HttpService().get(url).toPromise().then(res => res.data);
            if(data && data.result){
                let { denom, nfts } = data.result;
                return new NftCollectionDto({denom, nfts}) ;
            }else{
                Logger.warn('api-error:', 'there is no result of nft from lcd');
            }

        } catch (e) {
            Logger.warn(`api-error from ${url}:`, e.message);
            // cron jobs error should not throw errors;
        }
    }
}