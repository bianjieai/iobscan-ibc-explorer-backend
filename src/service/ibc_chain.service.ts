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

  // todo 接口处理结果未使用 res dot
  // 分页查询，用于前端请求
  async queryList(): Promise<IbcChainListResDto> {
    const ibcChainDatas = await this.ibcChainModel.findList();
    const res: IbcChainListResDto = ibcChainDatas;
    return res;
  }
}
