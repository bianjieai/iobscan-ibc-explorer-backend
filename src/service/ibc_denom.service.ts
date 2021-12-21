import { Injectable } from '@nestjs/common';
import { Connection } from 'mongoose';
import { InjectConnection } from '@nestjs/mongoose';
import { IbcDenomSchema } from '../schema/ibc_denom.schema';
import { IbcDenomResDto } from '../dto/ibc_denom.dto';
@Injectable()
export class IbcDenomService {
  private ibcDenomModel;

  constructor(@InjectConnection() private connection: Connection) {
    this.getModels();
  }

  async getModels(): Promise<void> {
    this.ibcDenomModel = await this.connection.model(
      'ibcDenomModel',
      IbcDenomSchema,
      'ibc_denom',
    );
  }

  // findAllRecord
  async findAllRecord(): Promise<IbcDenomResDto[]> {
    const result: IbcDenomResDto[] = IbcDenomResDto.bundleData(await this.ibcDenomModel.findAllRecord())
    return result;
  }

  async updateIbcDenom(ibcDenom) :Promise<void>{
    return await this.ibcDenomModel.updateDenomRecord(ibcDenom)
  }
}
