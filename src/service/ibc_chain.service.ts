import {Injectable} from '@nestjs/common';
import {Connection} from 'mongoose';
import {InjectConnection} from '@nestjs/mongoose';
import {IbcChainConfigSchema} from '../schema/ibc_chain_config.schema';
import {IbcChainConfigType} from '../types/schemaTypes/ibc_chain_config.interface';
import {IbcChainResDto, IbcChainResultResDto} from '../dto/ibc_chain.dto';
import {IbcTxSchema} from "../schema/ibc_tx.schema";
import {AggregateResult24hr} from "../types/schemaTypes/ibc_tx.interface";

@Injectable()
export class IbcChainService {
    private ibcChainConfigModel;
    private ibcTxModel;

    constructor(@InjectConnection() private connection: Connection) {
        this.getModels();
    }

    async getModels(): Promise<void> {
        this.ibcChainConfigModel = await this.connection.model(
            'ibcChainConfigModel',
            IbcChainConfigSchema,
            'chain_config',
        );

        this.ibcTxModel = await this.connection.model(
            'ibcTxModel',
            IbcTxSchema,
            'ex_ibc_tx',
        );
    }

    async findActiveChains24hr(dateNow):Promise<Array<AggregateResult24hr>> {
        return await this.ibcTxModel.findActiveChains24hr(dateNow);
    }

    async handleActiveChains(dateNow,allIbcChainInfos: IbcChainConfigType[]): Promise<IbcChainConfigType[]> {
        const result24hrs = await this.findActiveChains24hr(dateNow);
        const chainMap = new Map();
        for (const element of result24hrs) {
            if (!chainMap.has(element._id.sc_chain_id) && !chainMap.has(element._id.dc_chain_id)) {
                chainMap.set(element._id.sc_chain_id, '')
            }
            if (!chainMap.has(element._id.sc_chain_id) && !chainMap.has(element._id.dc_chain_id)) {
                chainMap.set(element._id.dc_chain_id, '')
            }
        }
        return allIbcChainInfos.filter(
            (item: IbcChainConfigType) => {
                return chainMap.has(item.chain_id)
            },
        );

    }

    async getAllChainConfigs() :Promise<IbcChainConfigType[]>{
        return await this.ibcChainConfigModel.findList();
    }

    async queryChainsByDatetime(dateNow): Promise<IbcChainResultResDto> {
        const ibcChainAllDatas: IbcChainConfigType[] = await this.getAllChainConfigs()
        const ibcChainActiveDatas: IbcChainConfigType[] = await this.handleActiveChains(dateNow,ibcChainAllDatas)
        const ibcChainInActiveDatas: IbcChainConfigType[] = ibcChainAllDatas.filter(
            (item: IbcChainConfigType) => {
                return !ibcChainActiveDatas.some((subItem: IbcChainConfigType) => {
                    return subItem.chain_id === item.chain_id;
                });
            },
        );

        return new IbcChainResultResDto({
            all: IbcChainResDto.bundleData(ibcChainAllDatas),
            active: IbcChainResDto.bundleData(ibcChainActiveDatas),
            inactive: IbcChainResDto.bundleData(ibcChainInActiveDatas),
        });
    }
}
