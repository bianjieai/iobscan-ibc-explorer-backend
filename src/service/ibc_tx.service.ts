import { Injectable } from '@nestjs/common';
import { Model } from 'mongoose';
import { InjectModel } from '@nestjs/mongoose';
import { ListStruct } from '../api/ApiResult';
import { IbcTxType } from '../types/schemaTypes/ibc_tx.interface';
import { IbcTxListReqDto, IbcTxResDto } from '../dto/ibc_tx.dto';
@Injectable()
export class IbcTxService {
  constructor(
    @InjectModel('IbcTx') private ibcTxModel: Model<IbcTxType>
  ) {}

  async queryIbcTxList(
    query: IbcTxListReqDto,
  ): Promise<ListStruct<IbcTxResDto[]>> {
    const { page_num, page_size } = query;
    const ibcTxDatas : IbcTxResDto[] = await this.ibcTxModel.findTxList(
      page_num,
      page_size,
    );
    return new ListStruct(ibcTxDatas, page_num, page_size);
  }
}
