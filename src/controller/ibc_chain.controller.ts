import { Controller, Get, Query } from '@nestjs/common';
import { ApiTags } from '@nestjs/swagger';
import { Result } from '../api/ApiResult';
import { IbcChainService } from '../service/ibc_chain.service';
import { IbcChainResultResDto } from '../dto/ibc_chain.dto';

@ApiTags('IbcChains')
@Controller('ibc')
export class IbcChainController {
  constructor(private readonly ibcChainService: IbcChainService) {}

  @Get('chains')
  async queryChains(): Promise<Result<IbcChainResultResDto>> {
    const dateNow = Math.floor(new Date().getTime() / 1000);
    const result: IbcChainResultResDto | null = await this.ibcChainService.queryChainsByDatetime(dateNow);
    return new Result<IbcChainResultResDto | null>(result);
  }
}
