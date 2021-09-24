import { Injectable } from '@nestjs/common';
import { Connection } from 'mongoose';
import { InjectConnection } from '@nestjs/mongoose';
import { IbcChainConfigSchema } from '../schema/ibc_chain_config.schema';
import { IbcChainSchema } from '../schema/ibc_chain.schema';
import { ListStruct } from '../api/ApiResult';
import { IbcChainListReqDto, IbcChainListResDto } from '../dto/ibc_chain.dto';
@Injectable()
export class IbcChainService {
  private ibcChainConfigModel;
  private ibcChainModel;
  constructor(@InjectConnection() private connection: Connection) {
    this.getModels();
  }

  async getModels(): Promise<void> {
    this.ibcChainConfigModel = await this.connection.model(
      'ibcChainConfigModel',
      IbcChainConfigSchema,
      'chain_config',
    );

    this.ibcChainModel = await this.connection.model(
      'ibcChainModel',
      IbcChainSchema,
      'ibc_chain',
    );
  }

  // todo 接口处理结果未使用 res dot
  // 分页查询，用于前端请求
  async queryList(): Promise<any> {
    const ibcChainAllDatas = await this.ibcChainConfigModel.findList();
    const ibcChainActiveDatas = await this.ibcChainModel.findActive();
    const ibcChainInActiveDatas = ibcChainAllDatas.filter(item => {
      return ibcChainActiveDatas.find(subItem => {
        return subItem.chain_id !== item.chain_id;
      });
    });
    return {
      all: ibcChainAllDatas,
      active: ibcChainActiveDatas,
      inactive: ibcChainInActiveDatas,
    };
  }
}
