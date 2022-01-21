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

    // get already caculate ibc denom
    async getCaculateDenomMap(srcChainId) : Promise<any>{
        let  caculateDenomMap = new Map
        const caculateDenom = await this.ibcDenomCaculateModel.findCaculateDenom(srcChainId)
        for (const one of caculateDenom) {
            if (one?.denom_path && one?.chain_id && one?.base_denom) {
                //目标chain_id + denom_path + base_denom
                caculateDenomMap.set(`${one?.chain_id}${one?.denom_path}${one?.base_denom}`,'')
            }
        }
        return caculateDenomMap
    }

    async caculateBaseDenom(oneBaseDenom, channelMap) {
        const chainCfg: IbcChainConfigType = channelMap.get(oneBaseDenom.chain_id)
        // caculate ibc_info hash to compare with ibc_info_hash_caculate
        //获取最新的ibc_info的hashCode
        const hashCode = Md5.hashStr(JSON.stringify(chainCfg.ibc_info))
        let ibcDenomInfos = []

        //根据baseDenom的chain_id获取已经计算好的denom
        const  caculateDenomMap = await this.getCaculateDenomMap(oneBaseDenom.chain_id)

        if (hashCode !== oneBaseDenom?.ibc_info_hash_caculate && chainCfg?.ibc_info?.length > 0 && oneBaseDenom?.denom) {
            for (const ibcInfo of chainCfg.ibc_info) {
                if (ibcInfo?.chain_id && ibcInfo?.paths?.length > 0) {
                    ibcInfo.paths.forEach(item => {
                        if (item?.counterparty?.port_id && item?.counterparty?.channel_id) {
                            const denomPath = `${item?.counterparty?.port_id}/${item?.counterparty?.channel_id}`
                            const ibcDenom = IbcDenom(denomPath, oneBaseDenom.denom)

                            //目标chain_id + denom_path + base_denom
                            const existKey = `${ibcInfo.chain_id}${denomPath}${oneBaseDenom.denom}`

                            // check if not caculate before push caculate denom to ibcDenomInfos
                            if (!caculateDenomMap?.has(existKey)) {
                                ibcDenomInfos.push({
                                    symbol: oneBaseDenom.symbol,
                                    base_denom: oneBaseDenom.denom,
                                    denom: ibcDenom,
                                    denom_path: denomPath,
                                    chain_id: ibcInfo.chain_id,
                                    sc_chain_id: oneBaseDenom.chain_id,
                                })
                            }
                        }
                    })
                }
            }
            if (ibcDenomInfos?.length > 0) {
                const session = await this.connection.startSession()
                session.startTransaction()
                try {
                    oneBaseDenom.ibc_info_hash_caculate = hashCode
                    await this.ibcDenomCaculateModel.insertDenomCaculate(ibcDenomInfos,session)
                    await this.ibcBaseDenomModel.updateBaseDenomWithSession(oneBaseDenom,session)

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

    async handleChain() {
        const chainConfig = await this.findAllChainConfig()
        const baseDenom = await this.findAllBaseDenom()
        let channelMap = new Map

        for (const one of chainConfig) {
            channelMap.set(`${one.chain_id}`, one)
        }
        for (const one of baseDenom) {
            if (channelMap.has(one.chain_id)) {
                await this.caculateBaseDenom(one, channelMap)
            }
        }
    }
}