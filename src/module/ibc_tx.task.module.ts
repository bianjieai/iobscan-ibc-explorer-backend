import { Module, forwardRef } from '@nestjs/common';
import { IbcTxTaskService } from '../task/ibc_tx.task.service';
import { IbcChainModule } from '../module/ibc_chain.module';
import { TxModule } from '../module/tx.module';
import { IbcDenomService } from '../service/ibc_denom.service'
@Module({
  imports: [IbcChainModule, TxModule, ],
  providers: [IbcTxTaskService, IbcDenomService],
  exports: [IbcTxTaskService],
})
export class IbcTxTaskModule {}
