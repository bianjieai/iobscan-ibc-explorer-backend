import { Injectable } from '@nestjs/common';
import { Connection } from 'mongoose';
import { InjectConnection } from '@nestjs/mongoose';
import { IbcChainSchema } from '../schema/ibc_chain.schema';
import { ListStruct } from '../api/ApiResult';
import { IbcChainListReqDto, IbcChainListResDto } from '../dto/ibc_chain.dto';
@Injectable()
export class IbcChainService {
  private ibcChainModel;
  constructor(@InjectConnection() private connection: Connection) {
    this.getModels();
  }

  async getModels(): Promise<void> {
    this.ibcChainModel = await this.connection.model(
      'ibcChainModel',
      IbcChainSchema,
      'chain_config',
    );
  }

  // 分页查询，用于前端请求
  async queryList(
    query: IbcChainListReqDto,
  ): Promise<ListStruct<IbcChainListResDto>> {
    const { page_num, page_size, chain_name } = query;
    const ibcChainDatas = await this.ibcChainModel.findList(
      page_num,
      page_size,
      chain_name,
    );
    const res: IbcChainListResDto = ibcChainDatas;
    return new ListStruct(res, page_num, page_size);
  }
}
