import { Module } from '@nestjs/common';
import { IbcStatisticsTaskService } from '../task/ibc_statistics.task.service';

@Module({
  imports: [ IbcStatisticsTaskService],
  providers: [IbcStatisticsTaskService],
  exports: [IbcStatisticsTaskService],
})
export class IbcStatisticsTaskModule {}
