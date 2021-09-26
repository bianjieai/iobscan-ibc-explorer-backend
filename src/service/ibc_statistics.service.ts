import { Injectable } from '@nestjs/common';
import { Connection } from 'mongoose';
import { InjectConnection } from '@nestjs/mongoose';
import { IbcStatisticsSchema } from '../schema/ibc_statistics.schema';
import { IbcStatisticsResDto } from '../dto/ibc_statistics.dto';
@Injectable()
export class IbcStatisticsService {
  private ibcStatisticsModel;

  constructor(@InjectConnection() private connection: Connection) {
    this.getModels();
  }

  async getModels(): Promise<void> {
    // ibcStatisticsModel
    this.ibcStatisticsModel = await this.connection.model(
      'ibcStatisticsModel',
      IbcStatisticsSchema,
      'ibc_statistics',
    );
  }

  // getAllRecord
  async findAllRecord(): Promise<IbcStatisticsResDto[]> {
    const result: IbcStatisticsResDto[] = IbcStatisticsResDto.bundleData(
      await this.ibcStatisticsModel.findAllRecord(),
    );
    return result;
  }
}
