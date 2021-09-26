import { Controller, Get, Query } from '@nestjs/common';
import { ApiTags } from '@nestjs/swagger';
import { Result } from '../api/ApiResult';
import { IbcTxService } from '../service/ibc_tx.service';
import { IbcTxListReqDto, IbcTxResDto } from '../dto/ibc_tx.dto';
import { ListStruct } from '../api/ApiResult';

@ApiTags('IbcTxs')
@Controller('ibc')
export class IbcTxController {
  constructor(private readonly ibcTxService: IbcTxService) {}

  @Get('txs')
  async getRecordList(
    @Query() query: IbcTxListReqDto,
  ): Promise<Result<ListStruct<IbcTxResDto[]>>> {
    const result: ListStruct<IbcTxResDto[]> = await this.ibcTxService.queryIbcTxList(query);
    return new Result<ListStruct<IbcTxResDto[]>>(result);
  }
}
