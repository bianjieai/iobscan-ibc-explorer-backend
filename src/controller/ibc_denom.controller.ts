import {Body, Controller,Headers, Get, Post} from '@nestjs/common';
import { ApiTags } from '@nestjs/swagger';
import { Result } from '../api/ApiResult';
import { IbcDenomService } from '../service/ibc_denom.service';
import { IbcDenomResDto } from '../dto/ibc_denom.dto';
import { IbcDenomDto } from '../types/schemaTypes/ibc_denom.interface';
import {cfg} from '../config/config';
@ApiTags('IbcDenoms')
@Controller('ibc')
export class IbcDenomController {
  constructor(private readonly ibcDenomService: IbcDenomService) {}

  @Get('denoms')
  async getAllRecord(): Promise<Result<IbcDenomResDto[]>> {
    const result: IbcDenomResDto[] = await this.ibcDenomService.findAllRecord();
    return new Result<IbcDenomResDto[]>(result);
  }

  @Post("denoms")
  async insertIbcDenom(@Body() dto:IbcDenomDto, @Headers() Headers):Promise<any> {
    const {executekey} = Headers
    if (executekey !== cfg.serverCfg.executeKey) {
      return {"message":"deny this operation for executekey is not right."}
    }
    return await this.ibcDenomService.createIbcDenom(dto)
  }
}
