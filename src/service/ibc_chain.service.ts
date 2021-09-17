import { Injectable } from '@nestjs/common';
import { Model } from 'mongoose';
import { InjectModel } from '@nestjs/mongoose';
import { ListStruct } from '../api/ApiResult';
import { IbcChainType } from '../types/schemaTypes/ibc_chain.interface';
import { IbcChainListReqDto, IbcChainListResDto } from '../dto/ibc_chain.dto';
import { TxSchema } from '../schema/tx.schema';

import { params } from '../app.module';

@Injectable()
export class IbcChainService {
  constructor(
    @InjectModel('IbcChain') private ibcChainModel: Model<IbcChainType>,
  ) {}

  // 分页查询，用于前端请求
  async queryList(
    query: IbcChainListReqDto,
  ): Promise<ListStruct<IbcChainListResDto>> {
    const { pageNum, pageSize, chain_name } = query;
    const ibcChainDatas = await this.ibcChainModel.findList(
      pageNum,
      pageSize,
      chain_name,
    );
    const res: IbcChainListResDto = ibcChainDatas;
    return new ListStruct(res, pageNum, pageSize);
  }

  // 查询数据库
  async queryListFromDb(): Promise<any> {
    const ibcChainDatas = await this.ibcChainModel.findList();
    return ibcChainDatas;
  }
}
