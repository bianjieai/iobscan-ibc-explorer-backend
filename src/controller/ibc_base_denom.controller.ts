import { Controller, Get } from '@nestjs/common';
import { ApiTags } from '@nestjs/swagger';
import { Result } from '../api/ApiResult';
import { IbcBaseDenomService } from '../service/ibc_base_denom.service';

@ApiTags('IbcBaseDenoms')
@Controller('ibc')
export class IbcBaseDenomController {
  constructor(private readonly ibcBaseDenomService: IbcBaseDenomService) {}

  @Get('baseDenoms')
  async getAllRecord() {
    const result = await this.ibcBaseDenomService.findAllRecord();
    return new Result(result);
  }
}
