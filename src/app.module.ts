import { Module } from '@nestjs/common';
import { MongooseModule } from '@nestjs/mongoose';
import { APP_FILTER, APP_PIPE } from '@nestjs/core';
import { HttpExceptionFilter } from './exception/HttpExceptionFilter';
import ValidationPipe from './pipe/validation.pipe';
import { ScheduleModule } from '@nestjs/schedule';
import { TasksService } from './task/task.service';
import { cfg } from './config/config';
import { TaskDispatchModule } from './module/task.dispatch.module';
import { IbcTxTaskModule } from './module/ibc_tx.task.module';
import { IbcTxModule } from './module/ibc_tx.module';
import { IbcChainTaskModule } from './module/ibc_chain.task.module';
import { IbcChainModule } from './module/ibc_chain.module';
import { IbcStatisticsTaskModule } from './module/ibc_statistics.task.module';
import { IbcStatisticsModule } from './module/ibc_statistics.module';
import { IbcBaseDenomModule } from './module/ibc_base_denom.module';
// const url: string = `mongodb://${cfg.dbCfg.user}:${cfg.dbCfg.psd}@${cfg.dbCfg.dbAddr}/${cfg.dbCfg.dbName}`;
// const url: string = `mongodb://iris:irispassword@10.1.4.66:27018/rainbow-server`;
const url: string = `mongodb://localhost:27017/ibc-db`;
export const params = {
  imports: [
    MongooseModule.forRoot(url),
    ScheduleModule.forRoot(),
    TaskDispatchModule,
    IbcTxTaskModule,
    IbcTxModule,
    IbcChainTaskModule,
    IbcChainModule,
    IbcStatisticsTaskModule,
    IbcStatisticsModule,
    IbcBaseDenomModule,
  ],
  providers: <any>[
    {
      provide: APP_FILTER,
      useClass: HttpExceptionFilter,
    },
    {
      provide: APP_PIPE,
      useClass: ValidationPipe,
    },
  ],
};

params.providers.push(TasksService);

// if (cfg.env !== 'development') {
//     params.providers.push(TasksService);
// }

@Module(params)
export class AppModule {}
