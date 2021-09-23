import { Controller, Get, Query } from '@nestjs/common';
import { ApiTags } from '@nestjs/swagger';
import { Result } from '../api/ApiResult';
import { IbcTxService } from '../service/ibc_tx.service';
import { IbcTxListReqDto, IbcTxListResDto } from '../dto/ibc_tx.dto';
import { ListStruct } from '../api/ApiResult';

@ApiTags('IbcTxs')
@Controller('ibc')
export class IbcTxController {
  constructor(private readonly ibcTxService: IbcTxService) {}

  @Get('txs')
  async getRecordList(
    @Query() query: IbcTxListReqDto,
  ): Promise<Result<ListStruct<IbcTxListResDto[]>>> {
    const result = await this.ibcTxService.queryIbcTxList(query);
    return new Result(result);
  }
}
