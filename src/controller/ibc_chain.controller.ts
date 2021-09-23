import { Controller, Get, Query } from '@nestjs/common';
import { ApiTags } from '@nestjs/swagger';
import { Result, ListStruct } from '../api/ApiResult';
import { IbcChainService } from '../service/ibc_chain.service';
import { IbcChainListReqDto, IbcChainListResDto } from '../dto/ibc_chain.dto';

@ApiTags('IbcChains')
@Controller('ibc')
export class IbcChainController {
  constructor(private readonly ibcChainService: IbcChainService) {}

  @Get('chains')
  async queryList(): Promise<Result<IbcChainListResDto>> {
    const result = await this.ibcChainService.queryList();
    return new Result(result);
  }
}
