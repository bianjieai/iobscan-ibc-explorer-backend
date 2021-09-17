import { Module } from '@nestjs/common';
import { ChainHttp } from '../http/lcd/chain.http';
import { IbcChainTaskService } from '../task/ibc_chain.task.service';

@Module({
  imports: [ChainHttp],
  providers: [IbcChainTaskService],
  exports: [IbcChainTaskService],
})
export class IbcChainTaskModule {}
