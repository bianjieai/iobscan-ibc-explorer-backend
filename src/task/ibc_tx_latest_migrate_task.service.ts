import {Injectable} from '@nestjs/common';
import {TaskEnum, IbcTxTable, MaxMigrateBatchLimit, IbcTxStatus} from "../constant";
import {cfg} from "../config/config";
import {Connection} from 'mongoose';
import {InjectConnection} from "@nestjs/mongoose";
import {IbcTxSchema} from "../schema/ibc_tx.schema";
import {Logger} from "../logger";

@Injectable()
export class IbcTxLatestMigrateTaskService {
    private ibcTxModel;
    private ibcTxLatestModel;

    constructor(@InjectConnection() private readonly connection: Connection) {
        this.getModels();
        this.doTask = this.doTask.bind(this);
    }

    // getModels
    async getModels(): Promise<void> {
        // ibcTxModel
        this.ibcTxModel = await this.connection.model(
            'ibcTxModel',
            IbcTxSchema,
            IbcTxTable.IbcTxTableName,
        );

        // ibcTxLatestModel
        this.ibcTxLatestModel = await this.connection.model(
            'ibcTxLatestModel',
            IbcTxSchema,
            IbcTxTable.IbcTxLatestTableName,
        );

    }

    async doTask(taskName?: TaskEnum): Promise<void> {
        const txCount = await this.ibcTxLatestModel.countAll();
        if (txCount <= cfg.serverCfg.displayIbcRecordMax) {
            return
        }
        const migrateCnt = txCount - cfg.serverCfg.displayIbcRecordMax
        await this.startMigrate(migrateCnt)
        Logger.log("ibc migrate have finished,migrate count:", migrateCnt)
    }

    async startMigrate(migrateCount): Promise<void> {
        // migrate start condition value with timestamp random
        const value = (Math.floor(new Date().getTime() / 1000) % 10) * 10000
        if (migrateCount > value) {
            // migrate max bitch count limit
            if (migrateCount > MaxMigrateBatchLimit) {
                let batchNum = Math.floor(migrateCount / MaxMigrateBatchLimit)
                const batchSize = batchNum * MaxMigrateBatchLimit
                for (; batchNum > 0; batchNum--) {
                    await this.migrateData(MaxMigrateBatchLimit)
                }
                return await this.migrateData(migrateCount - batchSize)
            }
            return await this.migrateData(migrateCount)
        }
    }

    // limit max value is MaxMigrateBatchLimit
    async migrateData(limit): Promise<void> {
        if (limit <= 0) {
            return
        }
        const settingTxs = await this.ibcTxLatestModel.queryTxsByStatusLimit({
            status: IbcTxStatus.SETTING,
            limit: limit
        })
        // this migrate status[1,2,3,4] txs data is (limit - settingTxs.length)
        let limitNormalData = 0
        if (limit > settingTxs.length) {
            limitNormalData = limit - settingTxs.length
        }
        let batchTxs = [ ...settingTxs]
        if (limitNormalData > 0 ) {
            const Txs = await this.ibcTxLatestModel.queryTxsLimit(limitNormalData, 1)
            batchTxs = [...Txs]
        }

        const session = await this.connection.startSession()
        session.startTransaction()
        try {
            let recordIds = []
            for (const one of batchTxs) {
                recordIds.push(one.record_id)
            }
            await this.ibcTxModel.insertManyIbcTx(batchTxs, session)
            await this.ibcTxLatestModel.deleteManyIbcTx(recordIds, session)
            await session.commitTransaction();
            session.endSession();
        } catch (e) {
            Logger.log(e, 'transaction is error')
            await session.abortTransaction()
            session.endSession();
        }
    }


}