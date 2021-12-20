import {Body, Controller,Headers, Get, Post} from '@nestjs/common';
import {ApiTags} from '@nestjs/swagger';
import { Result } from '../api/ApiResult';
import { IbcBaseDenomService } from '../service/ibc_base_denom.service';
import { IbcBaseDenomResDto } from '../dto/ibc_base_denom.dto';
import {IbcBaseDenomDto} from "../types/schemaTypes/ibc_base_denom.interface";
import {cfg} from '../config/config';
@ApiTags('IbcBaseDenoms')
@Controller('ibc')
export class IbcBaseDenomController {
  constructor(private readonly ibcBaseDenomService: IbcBaseDenomService) {}

  @Get('baseDenoms')
  async getAllRecord(): Promise<Result<IbcBaseDenomResDto[]>> {
    const result: IbcBaseDenomResDto[] = await this.ibcBaseDenomService.findAllRecord();
    return new Result<IbcBaseDenomResDto[]>(result);
  }
  @Post("baseDenoms")
  async insertIbcDenom(@Body() dto:IbcBaseDenomDto, @Headers() Headers):Promise<any> {
    const {executekey} = Headers
    if (executekey !== cfg.serverCfg.executeKey) {
      return {"message":"deny this operation for executekey is not right."}
    }
    return await this.ibcBaseDenomService.insertBaseDenom(dto)
  }
}
