import { Module } from '@nestjs/common';
import { PrometheusModule} from "@willsoto/nestjs-prometheus";
import { IbcMonitorService } from "../monitor/ibc_monitor.service";
import { MonitorController } from "../controller/monitor.controller";
import {LcdConnectionMetric,LcdConnectionProvider} from "../monitor/metrics/ibc_chain_lcd_connection.metric";
import {IbcTxProcessingMetric,IbcTxProcessingProvider} from "../monitor/metrics/ibc_tx_processing_cnt.metric";

@Module({
    imports: [PrometheusModule.register({
        controller: MonitorController,
    })],
    providers: [
        IbcMonitorService,
        LcdConnectionMetric,
        LcdConnectionProvider(),
        IbcTxProcessingMetric,
        IbcTxProcessingProvider(),
    ],
    exports: [
        IbcMonitorService,
    ]
})
export class MonitorModule {}

