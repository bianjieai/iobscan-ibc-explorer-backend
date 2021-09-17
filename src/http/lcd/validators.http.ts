import { HttpServer, Injectable, HttpService } from '@nestjs/common';
import { cfg } from "../../config/config"
import { Logger } from '../../logger';
import { VValidatorDto } from '../../dto/http.dto';

@Injectable()

export class ValidatorsHttp {
    async queryValidatorsFromLcd(pageNum: number, pageSize: number,validatorStatus?: boolean): Promise<VValidatorDto[]>{
        let validatorsLcdUrl: string = validatorStatus === undefined ? `${cfg.serverCfg.lcdAddr}/node/validators?page=${pageNum}&limit=${pageSize}` : `${cfg.serverCfg.lcdAddr}/node/validators?jailed=${validatorStatus}&page=${pageNum}&limit=${pageSize}`;
        try {
            const validatorsData:any = await new HttpService().get(validatorsLcdUrl).toPromise().then(result => result.data)
            if(validatorsData && validatorsData.result){
                return VValidatorDto.bundleData(validatorsData.result);
            }else{
            	Logger.warn('api-error:', 'there is no result of validators from lcd');
            }
        }catch (e) {
            Logger.warn(`api-error from ${validatorsLcdUrl}`,e.message)
        }
    }
}
