import { Module } from '@nestjs/common';
import { ChainHttp } from '../http/lcd/chain.http';
import { IbcChainConfigTaskService } from '../task/ibc_chain_config.task.service';

@Module({
  imports: [ChainHttp],
  providers: [IbcChainConfigTaskService],
  exports: [IbcChainConfigTaskService],
})
export class IbcChainConfigTaskModule {}
