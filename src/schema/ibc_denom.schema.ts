/* eslint-disable @typescript-eslint/camelcase */
import * as mongoose from 'mongoose';
import { IbcDenomType } from '../types/schemaTypes/ibc_denom.interface';
import {Logger} from "../logger";
import {AggregateBaseDenomCnt} from "../types/schemaTypes/ibc_denom.interface";

export const IbcDenomSchema = new mongoose.Schema(
    {
        chain_id: String,
        denom: String,
        base_denom: String,
        denom_path: String,
        is_source_chain: Boolean,
        is_base_denom: Boolean,
        symbol: {
            type: String,
            default: '',
        },
        create_at: {
            type: Number,
            default: Math.floor(new Date().getTime() / 1000),
        },
        update_at: {
            type: Number,
            default: Math.floor(new Date().getTime() / 1000),
        },
        tx_time: {
            type: Number,
            default: Math.floor(new Date().getTime() / 1000),
        },
        real_denom: {
            type: Boolean,
            default: false,
        },
    },
    {versionKey: false},
);

IbcDenomSchema.index({chain_id: 1, denom: 1}, {unique: true});
IbcDenomSchema.index({symbol: -1}, {background: true});

IbcDenomSchema.statics = {
    async findAllRecord(): Promise<IbcDenomType[]> {
        return this.find({});
    },

    async findRecordBySymbol(symbol: string): Promise<IbcDenomType[]> {
        return this.find({symbol});
    },

    async findCount(): Promise<number> {
        return this.count({});
    },

    async findBaseDenomCount(): Promise<Array<AggregateBaseDenomCnt>> {
        // return this.count({
        //     is_base_denom: true,
        // });
        return this.aggregate([
            {
                $match: {
                    is_base_denom: true,
                }
            },
            {
                $group: {
                    _id: {base_denom: "$base_denom", chain_id: "$chain_id"}
                }
            }]);
    },

    async findDenomRecord(chain_id, denom): Promise<IbcDenomType> {
        return this.findOne({chain_id, denom}, {_id: 0});
    },
     async findAllDenomRecord(chain_id): Promise<IbcDenomType> {
        return this.find({chain_id}, {_id: 0});
    },
    // async findAllDenomRecord(): Promise<IbcDenomType> {
    //     return this.findOne({});
    // },
     // async updateDenomRecord(denomRecord): Promise<void> {
        //const {chain_id, denom} = denomRecord;
        //const options = {upsert: true, new: false, setDefaultsOnInsert: true};
        // return this.findOneAndUpdate({ chain_id, denom }, denomRecord, options);
    // },

    async insertManyDenom(ibcDenom): Promise<void> {
        return this.insertMany(ibcDenom, { ordered: false },(error) => {
            if(JSON.stringify(error).includes('E11000 duplicate key error collection')){
            }else {
                Logger.error(error)
            }
        });
    },
    async createDenom(ibcDenom): Promise<void> {
        return this.create(ibcDenom);
    },
};
