import {HttpService, Injectable} from "@nestjs/common";
import {InjectConnection} from "@nestjs/mongoose";
import {Connection} from 'mongoose';
import {IbcChainConfigSchema} from "../schema/ibc_chain_config.schema";
import {IbcTxSchema} from "../schema/ibc_tx.schema";
import {TaskEnum} from "../constant";
import {NodeInfoType} from "../types/lcd.interface";
import {Logger} from "../logger";
import {LcdConnectionMetric} from "../monitor/metrics/ibc_chain_lcd_connection.metric";
import {IbcTxProcessingMetric} from "../monitor/metrics/ibc_tx_processing_cnt.metric";
import {IbcTxTable} from "../constant/index";

@Injectable()
export class IbcMonitorService {
    private chainConfigModel;
    private ibcTxLatestModel;

    constructor(@InjectConnection() private readonly connection: Connection,
                private readonly lcdConnectionMetric: LcdConnectionMetric,
                private readonly ibcTxProcessMetric: IbcTxProcessingMetric) {
        this.getModels();
        this.doTask = this.doTask.bind(this);
    }

    async doTask(taskName?: TaskEnum): Promise<void> {
        const allChains = await this.chainConfigModel.findAll();
        for (const one of allChains) {
           await this.getNodeInfo(one.lcd,one.chain_id)
        }
        this.getProcessingCnt()
    }
    // getModels
    async getModels(): Promise<void> {
        // chainConfigModel
        this.chainConfigModel = await this.connection.model(
            'chainConfigModel',
            IbcChainConfigSchema,
            'chain_config',
        );


        // ibcTxModel
        this.ibcTxLatestModel = await this.connection.model(
            'ibcTxLatestModel',
            IbcTxSchema,
            IbcTxTable.IbcTxLatestTableName,
        );
    }

    async getNodeInfo(lcdAddr,chain) {
        const nodeInfoUrl = `${lcdAddr}/node_info`;
        try {
            const nodeInfo: NodeInfoType= await new HttpService()
                .get(nodeInfoUrl)
                .toPromise()
                .then(result => result.data);
            if (nodeInfo) {
                //  monitor code
                await this.lcdConnectionMetric.collect(chain,1)
                return nodeInfo
            } else {
                // monitor code
                await this.lcdConnectionMetric.collect(chain,0)
                Logger.warn(
                    'api-error:',
                    'there is no result of data from lcd',
                );
            }
        } catch (e) {
            // monitor code
            await this.lcdConnectionMetric.collect(chain,0)
            Logger.warn(`api-error from ${nodeInfoUrl} error`);
        }
    }



    async getProcessingCnt() {
        const processingCnt = await this.ibcTxLatestModel.countProcessing();
        await this.ibcTxProcessMetric.collect(processingCnt)
    }


}