import { Module } from '@nestjs/common';
import { IbcDenomCaculateTaskService } from '../task/ibc_denom_caculate.task.service';

@Module({
    providers: [IbcDenomCaculateTaskService],
    exports: [IbcDenomCaculateTaskService],
})
export class IbcDenomCaculateTaskModule {}