import { Module } from '@nestjs/common';
import { IbcDenomUpdateTaskService } from '../task/ibc_denom_update.task.service';

@Module({
    providers: [IbcDenomUpdateTaskService],
    exports: [IbcDenomUpdateTaskService],
})
export class IbcDenomUpdateTaskModule {}