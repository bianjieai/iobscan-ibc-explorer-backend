import * as mongoose from 'mongoose';
import { getTimestamp } from '../util/util';
import { IDenomMapStruct, IDenomStruct } from '../types/schemaTypes/denom.interface';

export const DenomSchema = new mongoose.Schema({
    name: String,
    json_schema: String,
    denom_id: { type: String, unique: true },
    creator: String,
    tx_hash: String,
    height: Number,
    time:Number,
    create_time: Number,
    last_block_height: Number,
    last_block_time: Number,
}, { versionKey: false });
// 新增
DenomSchema.index({ height: 1}, { background:true});

DenomSchema.statics = {
    async findList(
        pageNum: number,
        pageSize: number,
        denomNameOrId?: string,
        needAll?: string,
    ): Promise<IDenomStruct[]> {
        if (needAll) {
            return await this.find({});
        } else {
            const params = {};
            if(denomNameOrId){
                const reg = new RegExp(denomNameOrId, 'i');
                params['$or'] = [
                    { 'name': { $regex: reg } },
                    { 'denom_id': { $regex: reg } },
                ];
            }
            return await this.find(params)
                .skip((Number(pageNum) - 1) * Number(pageSize))
                .limit(Number(pageSize))
                .sort({ height: -1 });
        }
    },
    async queryDenomCount(denomNameOrId?: string){
        const params = {};
        if(denomNameOrId){
            const reg = new RegExp(denomNameOrId, 'i');
            params['$or'] = [
                { 'name': { $regex: reg } },
                { 'denom_id': { $regex: reg } },
            ]
        }
        return this.countDocuments(params);
    },
    async queryAllCount(){
        return this.countDocuments({});
    },
    async findOneByDenomId(denomId:string): Promise<IDenomStruct> {
        return await this.findOne({denom_id:denomId});
    },
    async saveDenom(denoms: IDenomMapStruct): Promise<IDenomStruct[]> {
        return await this.create({
            name: denoms.name,
            denom_id: denoms.denomId,
            json_schema: denoms.jsonSchema,
            creator: denoms.creator,
            tx_hash: denoms.txHash,
            height: denoms.height,
            time: denoms.createTime,
            create_time: getTimestamp(),
            update_time: getTimestamp(),
        });
    },
    async findAllNames(): Promise<IDenomStruct[]> {
        return await this.find({}, { denom_id: 1, name: 1 }).exec();
    },

    async updateDenom(denom: IDenomMapStruct): Promise<IDenomStruct> {
        return await this.findOneAndUpdate({
            denom_id:denom.denomId,
        }, {
            tx_hash: denom.txHash,
            height: denom.height,
            time:denom.createTime,
            update_time: getTimestamp(),
        });
    },

    async findOneByDenomAndNftIdFromDenom(denomId: string): Promise<IDenomStruct> {
        return await this.findOne({ denom_id: denomId}, {'_id': 0,'update_time': 0,'create_time': 0})
    },

    async findLastBlockHeight(): Promise<IDenomStruct[]> {
        return await this.find({}, { last_block_height: 1 }).sort({last_block_height:-1}).limit(1)
    },

    async insertManyDenom(denomList): Promise<IDenomStruct[]>{
       return await this.insertMany(denomList,{ ordered: false })
    },
};