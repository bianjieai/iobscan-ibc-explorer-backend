import { Controller, Get } from '@nestjs/common';
import { ApiTags } from '@nestjs/swagger';
import { Result } from '../api/ApiResult';
import { IbcBaseDenomService } from '../service/ibc_base_denom.service';
import { IbcBaseDenomResDto } from '../dto/ibc_base_denom.dto';
@ApiTags('IbcBaseDenoms')
@Controller('ibc')
export class IbcBaseDenomController {
  constructor(private readonly ibcBaseDenomService: IbcBaseDenomService) {}

  @Get('baseDenoms')
  async getAllRecord(): Promise<Result<IbcBaseDenomResDto[]>> {
    const result: IbcBaseDenomResDto[] = await this.ibcBaseDenomService.findAllRecord();
    return new Result<IbcBaseDenomResDto[]>(result);
  }
}
