import { Injectable } from '@nestjs/common';
import { Connection } from 'mongoose';
import { InjectConnection } from '@nestjs/mongoose';
import { IbcChainConfigSchema } from '../schema/ibc_chain_config.schema';
import { IbcChainSchema } from '../schema/ibc_chain.schema';
import { IbcChainConfigType } from '../types/schemaTypes/ibc_chain_config.interface';
import { IbcChainType } from '../types/schemaTypes/ibc_chain.interface';
import { IbcChainResDto, IbcChainResultResDto } from '../dto/ibc_chain.dto';
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

  async queryList(): Promise<IbcChainResultResDto> {
    const ibcChainAllDatas: IbcChainConfigType[] = await this.ibcChainConfigModel.findList();
    const ibcChainActiveDatas: IbcChainType[] = await this.ibcChainModel.findActive();
    const ibcChainInActiveDatas: IbcChainConfigType[] = ibcChainAllDatas.filter((item: IbcChainConfigType) => {
      return ibcChainActiveDatas.find((subItem: IbcChainType) => {
        return subItem.chain_id !== item.chain_id;
      });
    });

    return new IbcChainResultResDto({
      all: IbcChainResDto.bundleData(ibcChainAllDatas),
      active: IbcChainResDto.bundleData(ibcChainActiveDatas),
      inactive: IbcChainResDto.bundleData(ibcChainInActiveDatas),
    });
  }
}
