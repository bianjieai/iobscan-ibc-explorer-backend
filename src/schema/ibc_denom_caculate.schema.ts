import * as mongoose from 'mongoose';
import {IbcDenomCaculateType} from "../types/schemaTypes/ibc_denom_hash.interface";

export const IbcDenomCaculateSchema = new mongoose.Schema(
    {
        chain_id: String,
        sc_chain_id: String,
        denom: String,
        base_denom: String,
        denom_path: String,
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

    },
    {versionKey: false},
);

IbcDenomCaculateSchema.index({chain_id: 1, denom: 1}, {unique: true});
IbcDenomCaculateSchema.statics = {
    async findCount(): Promise<number> {
        return this.count();
    },

    async findCaculateDenom(sc_chain_id): Promise<IbcDenomCaculateType[]> {
        return this.find({sc_chain_id:sc_chain_id});
    },

    async findIbcDenoms(chain_id,denoms): Promise<IbcDenomCaculateType[]> {
        return this.find({
            chain_id:chain_id,
            denom: {
                $in: denoms,
            }}, {_id: 0});
    },

    async insertDenomCaculate(ibcDenom, session): Promise<void>{
        return this.insertMany(ibcDenom,{ ordered: false }, (error) => {
            if(JSON.stringify(error).includes('E11000')){
                // Primary key conflict handling
            }else {
                if (error) {
                    console.log(error,'insertMany IbcTx error')
                }
            }
        },session);
    },
};