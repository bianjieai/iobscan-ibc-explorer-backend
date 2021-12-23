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

  async getModels(): Promise<void> {
    this.ibcBaseDenomModel = await this.connection.model(
      'ibcBaseDenomModel',
      IbcBaseDenomSchema,
      'ibc_base_denom',
    );
  }

  // findAllRecord
  async findAllRecord(): Promise<IbcBaseDenomResDto[]> {
    const result: IbcBaseDenomResDto[] = IbcBaseDenomResDto.bundleData(await this.ibcBaseDenomModel.findAllRecord())
    return result;
  }

  async insertBaseDenom(dto):Promise<void>{
      return await this.ibcBaseDenomModel.insertBaseDenom(dto)
  }
}
