import { Module } from '@nestjs/common';
import { IbcTxTaskService } from '../task/ibc_tx.task.service';
import { IbcChainModule } from '../module/ibc_chain.module';
@Module({
  imports: [IbcChainModule],
  providers: [IbcTxTaskService],
  exports: [IbcTxTaskService],
})
export class IbcTxTaskModule {}
