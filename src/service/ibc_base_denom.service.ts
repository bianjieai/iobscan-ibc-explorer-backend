import { Injectable } from '@nestjs/common';
import { Connection } from 'mongoose';
import { InjectConnection } from '@nestjs/mongoose';
import { IbcBaseDenomSchema } from '../schema/ibc_base_denom.schema';
import { IbcBaseDenomResDto } from '../dto/ibc_base_denom.dto';
@Injectable()
export class IbcBaseDenomService {
  private ibcBaseDenomModel;

  constructor(@InjectConnection() private connection: Connection) {
    this.getModels();
  }

  // 获取所有静态模型
  async getModels(): Promise<void> {
    // ibcStatisticsModel
    this.ibcBaseDenomModel = await this.connection.model(
      'ibcBaseDenomModel',
      IbcBaseDenomSchema,
      'ibc_base_denom',
    );
  }

  // 获取所有记录
  async findAllRecord(): Promise<IbcBaseDenomResDto[]> {
    const result: IbcBaseDenomResDto[] = this.ibcBaseDenomModel.findAllRecord()
    return result;
  }
}
