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
  async queryList(): Promise<Result<IbcChainResultResDto>> {
    const result: IbcChainResultResDto | null = await this.ibcChainService.queryList();
    return new Result<IbcChainResultResDto | null>(result);
  }
}
