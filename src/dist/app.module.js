"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
exports.__esModule = true;
exports.AppModule = exports.params = void 0;
var common_1 = require("@nestjs/common");
var mongoose_1 = require("@nestjs/mongoose");
var core_1 = require("@nestjs/core");
var HttpExceptionFilter_1 = require("./exception/HttpExceptionFilter");
var validation_pipe_1 = require("./pipe/validation.pipe");
var schedule_1 = require("@nestjs/schedule");
var task_service_1 = require("./task/task.service");
var config_1 = require("./config/config");
var task_dispatch_module_1 = require("./module/task.dispatch.module");
var ibc_tx_task_module_1 = require("./module/ibc_tx.task.module");
var ibc_tx_module_1 = require("./module/ibc_tx.module");
var ibc_chain_config_task_module_1 = require("./module/ibc_chain_config.task.module");
var ibc_chain_module_1 = require("./module/ibc_chain.module");
var ibc_statistics_task_module_1 = require("./module/ibc_statistics.task.module");
var ibc_statistics_module_1 = require("./module/ibc_statistics.module");
var ibc_base_denom_module_1 = require("./module/ibc_base_denom.module");
var url = "mongodb://" + config_1.cfg.dbCfg.user + ":" + config_1.cfg.dbCfg.psd + "@" + config_1.cfg.dbCfg.dbAddr + "/" + config_1.cfg.dbCfg.dbName;
// const url: string = `mongodb://localhost:27017/ibc-db`;
exports.params = {
    imports: [
        mongoose_1.MongooseModule.forRoot(url),
        schedule_1.ScheduleModule.forRoot(),
        task_dispatch_module_1.TaskDispatchModule,
        ibc_tx_task_module_1.IbcTxTaskModule,
        ibc_tx_module_1.IbcTxModule,
        ibc_chain_config_task_module_1.IbcChainConfigTaskModule,
        ibc_chain_module_1.IbcChainModule,
        ibc_statistics_task_module_1.IbcStatisticsTaskModule,
        ibc_statistics_module_1.IbcStatisticsModule,
        ibc_base_denom_module_1.IbcBaseDenomModule,
    ],
    providers: [
        {
            provide: core_1.APP_FILTER,
            useClass: HttpExceptionFilter_1.HttpExceptionFilter
        },
        {
            provide: core_1.APP_PIPE,
            useClass: validation_pipe_1["default"]
        },
    ]
};
exports.params.providers.push(task_service_1.TasksService);
// if (cfg.env !== 'development') {
//     params.providers.push(TasksService);
// }
var AppModule = /** @class */ (function () {
    function AppModule() {
    }
    AppModule = __decorate([
        common_1.Module(exports.params)
    ], AppModule);
    return AppModule;
}());
exports.AppModule = AppModule;
