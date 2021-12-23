import { Module } from '@nestjs/common';
import { IbcTxTaskService } from '../task/ibc_tx.task.service';
@Module({
  providers: [IbcTxTaskService],
  exports: [IbcTxTaskService],
})
export class IbcTxTaskModule {}
