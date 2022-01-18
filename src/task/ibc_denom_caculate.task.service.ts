import { Connection } from 'mongoose';
import {Injectable} from "@nestjs/common";
import {InjectConnection} from "@nestjs/mongoose";
import {TaskEnum} from "../constant";
import {IbcChainConfigSchema} from "../schema/ibc_chain_config.schema";
import {IbcBaseDenomSchema} from "../schema/ibc_base_denom.schema";
import {IbcDenomCaculateSchema} from "../schema/ibc_denom_caculate.schema";
import {IbcChainConfigType} from "../types/schemaTypes/ibc_chain_config.interface";
import {IbcBaseDenomType} from "../types/schemaTypes/ibc_base_denom.interface";
import {IbcDenom} from "../helper/denom.helper";
import {Md5} from "ts-md5";
import {Logger} from "../logger";

@Injectable()
export class IbcDenomCaculateTaskService {
    private chainConfigModel;
    private ibcBaseDenomModel;
    private ibcDenomCaculateModel;

    constructor(@InjectConnection() private connection: Connection) {
        this.getModels();
        this.doTask = this.doTask.bind(this);
    }

    async doTask(taskName?: TaskEnum): Promise<void> {
        this.handleChain();
    }

    async getModels(): Promise<void> {
        this.ibcBaseDenomModel = await this.connection.model(
            'ibcBaseDenomModel',
            IbcBaseDenomSchema,
            'ibc_base_denom',
        );
        this.chainConfigModel = await this.connection.model(
            'chainConfigModel',
            IbcChainConfigSchema,
            'chain_config',
        );
        this.ibcDenomCaculateModel = await this.connection.model(
            'ibcDenomCaculateModel',
            IbcDenomCaculateSchema,
            'ibc_denom_caculate',
        );
    }

    async findAllChainConfig(): Promise<IbcChainConfigType[]> {
        return await this.chainConfigModel.findAll();
    }

    async findAllBaseDenom(): Promise<IbcBaseDenomType[]> {
        return await this.ibcBaseDenomModel.findAllRecord();
    }

    async handleChain() {
        const chainConfig = await this.findAllChainConfig()
        const baseDenom = await this.findAllBaseDenom()
        let denomMap = new Map, channelMap = new Map

        for (const one of chainConfig) {
            channelMap.set(`${one.chain_id}`, one)
        }
        for (const one of baseDenom) {
            if (denomMap.has(`${one.chain_id}${one.denom}`)) {
                continue
            }
            let ibcDenomInfos = []
            denomMap.set(`${one.chain_id}${one.denom}`, one)
            if (channelMap.has(one.chain_id)) {
                const chainCfg: IbcChainConfigType = channelMap.get(one.chain_id)
                // caculate ibc_info hash to compare with ibc_info_hash_caculate
                //获取最新的ibc_info的hashCode
                const hashCode = Md5.hashStr(JSON.stringify(chainCfg.ibc_info))

                if (hashCode !== one?.ibc_info_hash_caculate && chainCfg?.ibc_info?.length > 0) {
                    for (const ibcInfo of chainCfg.ibc_info) {
                        if (ibcInfo?.chain_id && ibcInfo?.paths?.length > 0) {
                            ibcInfo.paths.forEach(item => {
                                if (item?.counterparty?.port_id && item?.counterparty?.channel_id) {
                                    const denomPath = `${item?.counterparty?.port_id}/${item?.counterparty?.channel_id}`
                                    const ibcDenom = IbcDenom(denomPath, one.denom)
                                    ibcDenomInfos.push({
                                        symbol: one.symbol,
                                        base_denom: one.denom,
                                        denom: ibcDenom,
                                        denom_path: denomPath,
                                        chain_id: ibcInfo.chain_id,
                                        sc_chain_id: one.chain_id,
                                    })
                                }
                            })
                        }
                    }
                    if (ibcDenomInfos?.length > 0) {
                        const session = await this.connection.startSession()
                        session.startTransaction()
                        try {
                            one.ibc_info_hash_caculate = hashCode
                            await this.ibcDenomCaculateModel.insertDenomCaculate(ibcDenomInfos,session)
                            await this.ibcBaseDenomModel.updateBaseDenomWithSession(one,session)

                            await session.commitTransaction();
                            session.endSession();
                        } catch (e) {
                            Logger.log(e, 'transaction is error')
                            await session.abortTransaction()
                            session.endSession();
                        }
                    }
                }
            }

        }
    }
}