import {Connection} from 'mongoose';
import {Injectable} from "@nestjs/common";
import {InjectConnection} from "@nestjs/mongoose";
import {RecordLimit, TaskEnum} from "../constant";
import {IbcDenomSchema} from "../schema/ibc_denom.schema";
import {IbcDenomCaculateSchema} from "../schema/ibc_denom_caculate.schema";

@Injectable()
export class IbcDenomUpdateTaskService {
    private ibcDenomCaculateModel;
    private ibcDenomModel;

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
        this.ibcDenomCaculateModel = await this.connection.model(
            'ibcDenomCaculateModel',
            IbcDenomCaculateSchema,
            'ibc_denom_caculate',
        );
    }

    async collectChainDenomsMap(ibcDenomData) {
        let chainDenomsMap = new Map
        ibcDenomData.forEach(item => {
            if (item.chain_id) {
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
        })
        return chainDenomsMap
    }

    async getNeedhandleIbcDenoms() {
        return await this.ibcDenomModel.findUnAuthDenom(RecordLimit)
    }

    async getIbcDenoms(chainId, denoms) {
        return await this.ibcDenomCaculateModel.findIbcDenoms(chainId, denoms)
    }

    async handleChainDenoms() {
        const ibcDenoms = await this.getNeedhandleIbcDenoms()
        const chainDenomsMap = await this.collectChainDenomsMap(ibcDenoms)
        let denomData = []
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