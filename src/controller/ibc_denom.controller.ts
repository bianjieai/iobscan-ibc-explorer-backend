import { Controller, Get } from '@nestjs/common';
import { ApiTags } from '@nestjs/swagger';
import { Result } from '../api/ApiResult';
import { IbcDenomService } from '../service/ibc_denom.service';
import { IbcDenomResDto } from '../dto/ibc_denom.dto';
@ApiTags('IbcDenoms')
@Controller('ibc')
export class IbcDenomController {
  constructor(private readonly ibcDenomService: IbcDenomService) {}

  @Get('denoms')
  async getAllRecord(): Promise<Result<IbcDenomResDto[]>> {
    const result: IbcDenomResDto[] = await this.ibcDenomService.findAllRecord();
    return new Result<IbcDenomResDto[]>(result);
  }
}
