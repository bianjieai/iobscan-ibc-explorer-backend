import { Controller, Get, Query } from '@nestjs/common';
import { ApiTags } from '@nestjs/swagger';
import { Result } from '../api/ApiResult';
import { IbcChainService } from '../service/ibc_chain.service';
import { IbcChainListReqDto, IbcChainListResDto } from '../dto/ibc_chain.dto';

@ApiTags('IbcChains')
@Controller('ibcChain')
export class IbcChainController {
  constructor(private readonly ibcChainService: IbcChainService) {}
  @Get('list')
  async queryList(@Query() query: IbcChainListReqDto): Promise<Result<IbcChainListResDto>> {
    const result = await this.ibcChainService.queryList(query)
    return new Result(result);
  }
}
