/* eslint-disable @typescript-eslint/consistent-type-assertions */
import {Module} from '@nestjs/common';
import {MongooseModule} from '@nestjs/mongoose';
import {APP_FILTER, APP_PIPE} from '@nestjs/core';
import {HttpExceptionFilter} from './exception/HttpExceptionFilter';
import ValidationPipe from './pipe/validation.pipe';
import {ScheduleModule} from '@nestjs/schedule';
import {TasksService} from './task/task.service';
import {cfg} from './config/config';
import {TaskDispatchModule} from './module/task.dispatch.module';
import {IbcTxModule} from './module/ibc_tx.module';
import {IbcChainConfigTaskModule} from './module/ibc_chain_config.task.module';
import {IbcChainModule} from './module/ibc_chain.module';
import {IbcStatisticsTaskModule} from './module/ibc_statistics.task.module';
import {IbcStatisticsModule} from './module/ibc_statistics.module';
import {IbcBaseDenomModule} from './module/ibc_base_denom.module';
import {IbcDenomModule} from './module/ibc_denom.module';
import {ConfigModule} from './module/config.module';
import {MonitorModule} from './module/monitor.task.module';
import {IbcSyncTransferTxTaskModule} from "./module/ibc_sync_transfer_tx_task.module";
import {IbcUpdateProcessingTxModule} from "./module/ibc_update_processing_tx_task.module";
import {IbcUpdateSubstateTxTaskModule} from "./module/ibc_update_substate_tx_task.module";
import {IbcTxDataUpdateModule} from "./module/ibc_tx_data_update_task.module";
import {IbcTxLatestMigrateTaskModule} from "./module/ibc_tx_latest_migrate_task.module";
import {IbcDenomCaculateTaskModule} from "./module/ibc_denom_caculate.task.module";
import {IbcDenomUpdateTaskModule} from "./module/ibc_denom_update.task.module";

const url = `mongodb://${cfg.dbCfg.user}:${cfg.dbCfg.psd}@${cfg.dbCfg.dbAddr}/${cfg.dbCfg.dbName}`;
// const url: string = `mongodb://localhost:27017/ibc-db`;
export const params = {
    imports: [
        MongooseModule.forRoot(url),
        ScheduleModule.forRoot(),
        TaskDispatchModule,
        IbcTxModule,
        IbcChainConfigTaskModule,
        IbcChainModule,
        IbcStatisticsTaskModule,
        IbcStatisticsModule,
        IbcBaseDenomModule,
        IbcDenomModule,
        IbcDenomCaculateTaskModule,
        IbcDenomUpdateTaskModule,
        IbcSyncTransferTxTaskModule,
        IbcUpdateProcessingTxModule,
        IbcUpdateSubstateTxTaskModule,
        IbcTxDataUpdateModule,
        IbcTxLatestMigrateTaskModule,
        ConfigModule,
        MonitorModule,
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
export class AppModule {
}
