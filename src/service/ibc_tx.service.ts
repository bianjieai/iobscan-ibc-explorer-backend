import { Injectable } from '@nestjs/common';
import { Model } from 'mongoose';
import { InjectModel } from '@nestjs/mongoose';
import { ListStruct } from '../api/ApiResult';
import { IbcTxType } from '../types/schemaTypes/ibc_tx.interface';
import { ITxStruct } from '../types/schemaTypes/tx.interface';
import { IbcTxListReqDto, IbcTxListResDto } from '../dto/ibc_tx.dto';
@Injectable()
export class IbcTxService {
  constructor(
    // @InjectModel('tx') private txModel: Model<ITxStruct>,
    @InjectModel('IbcTx') private ibcTxModel: Model<IbcTxType>
  ) {}

  async queryIbcTxList(
    query: IbcTxListReqDto,
  ): Promise<ListStruct<IbcTxListResDto[]>> {
    const { pageNum, pageSize } = query;
    const ibcTxDatas = await this.ibcTxModel.findTxList(
      pageNum,
      pageSize,
    );
    return new ListStruct(ibcTxDatas, pageNum, pageSize);
  }
}
