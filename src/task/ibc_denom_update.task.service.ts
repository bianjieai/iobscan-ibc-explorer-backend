import {Connection} from 'mongoose';
import {Injectable} from "@nestjs/common";
import {InjectConnection} from "@nestjs/mongoose";
import {RecordLimit, TaskEnum} from "../constant";
import { cfg } from '../config/config';
import {IbcDenomSchema} from "../schema/ibc_denom.schema";
import {IbcBaseDenomSchema} from "../schema/ibc_base_denom.schema";
import {IbcDenomCaculateSchema} from "../schema/ibc_denom_caculate.schema";

@Injectable()
export class IbcDenomUpdateTaskService {
    private ibcDenomCaculateModel;
    private ibcDenomModel;
    private ibcBaseDenomModel;

    constructor(@InjectConnection() private connection: Connection) {
        this.getModels();
        this.doTask = this.doTask.bind(this);
    }

    async doTask(taskName?: TaskEnum): Promise<void> {
        this.handleChainDenoms();
    }

    async getModels(): Promise<void> {
        this.ibcDenomModel = await this.connection.model(
            'ibcDenomModel',
            IbcDenomSchema,
            'ibc_denom',
        );
        this.ibcBaseDenomModel = await this.connection.model(
            'ibcBaseDenomModel',
            IbcBaseDenomSchema,
            'ibc_base_denom',
        );
        this.ibcDenomCaculateModel = await this.connection.model(
            'ibcDenomCaculateModel',
            IbcDenomCaculateSchema,
            'ibc_denom_caculate',
        );
    }

    async collectChainDenomsMap(ibcDenomData) {
        let chainDenomsMap = new Map,baseDenoms = []
        ibcDenomData.forEach(item => {
            if (item?.chain_id) {
                if (item?.denom?.indexOf("ibc/") < 0 ) {
                    baseDenoms.push(item.denom)
                }else{
                    if (!chainDenomsMap.has(item.chain_id)) {
                        let denoms = []
                        denoms.push(item.denom)
                        chainDenomsMap.set(item.chain_id, denoms)
                    } else {
                        let arrayDenoms = chainDenomsMap.get(item.chain_id)
                        arrayDenoms.push(item.denom)
                        chainDenomsMap.set(item.chain_id, arrayDenoms)
                    }
                }
            }
        })
        return {chainDenomsMap,baseDenoms}
    }

    async getNeedhandleIbcDenoms(page_num,page_size) {
        return await this.ibcDenomModel.findUnAuthDenom(page_num,page_size)
    }

    async getIbcDenoms(chainId, denoms) {
        return await this.ibcDenomCaculateModel.findIbcDenoms(chainId, denoms)
    }

    async getBaseDenomMap() : Promise<any>{
        let  baseDenomMap = new Map
        const baseDenom = await this.ibcBaseDenomModel.findAllRecord()
        if (baseDenom?.length) {
            for (const item of baseDenom) {
                baseDenomMap.set(`${item?.denom}`,item?.symbol)
            }
        }
        return baseDenomMap
    }

    async handleChainDenoms() {
        let pageNum = 1,supportDenoms = []

        const baseDenomMap = await this.getBaseDenomMap()

        //judge  supportDenoms size for handle batch limit
        while(supportDenoms?.length < cfg.serverCfg.updateDenomBatchLimit) {
            const ibcDenoms = await this.getNeedhandleIbcDenoms(pageNum, RecordLimit)
            ibcDenoms.forEach(item => {
                const paths = item?.denom_path?.split("/")
                //only support one skip path
                if (paths?.length <= 2){

                    //only support denom which base_denom in ibc_base_denom
                    if (baseDenomMap && baseDenomMap?.has(`${item?.base_denom}`)) {
                        supportDenoms.push(item)
                    }
                }
            });

            // break when finish scan all the ibc denom which symbol is empty.
            if (ibcDenoms?.length < RecordLimit) {
                break;
            }
            pageNum++
        }


        const {chainDenomsMap,baseDenoms} = await this.collectChainDenomsMap(supportDenoms)
        let denomData = []

        if (baseDenoms) {
            const baseDenomDatas = await this.ibcBaseDenomModel.findByDenoms(baseDenoms)
            for (const item of baseDenomDatas) {
                await this.ibcDenomModel.updateDenomRecord({
                    chain_id: item.chain_id,
                    denom: item.denom,
                    symbol: item.symbol,
                })
            }
        }
        if (chainDenomsMap) {
            chainDenomsMap.forEach((value, key) => {
                denomData.push({chain_id: key, denoms: value})
            })

            for (const one of denomData) {
                const items = await this.getIbcDenoms(one.chain_id, one.denoms)
                for (const item of items) {
                    await this.ibcDenomModel.updateDenomRecord({
                        chain_id: one.chain_id,
                        denom: item.denom,
                        symbol: item.symbol,
                    })
                }
            }
        }
    }
}