import { Module } from '@nestjs/common';
import { IbcStatisticsService } from '../service/ibc_statistics.service';
import { IbcStatisticsController } from '../controller/ibc_statistics.controller';

@Module({
  providers: [IbcStatisticsService],
  controllers: [IbcStatisticsController],
})

export class IbcStatisticsModule {}
