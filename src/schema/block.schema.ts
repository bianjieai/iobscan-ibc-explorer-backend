import * as mongoose from 'mongoose';
import { IBlockStruct } from '../types/schemaTypes/block.interface';

export const BlockSchema = new mongoose.Schema({
    height: Number,
    hash: String,
    txn: Number,
    time: Number,
}, { versionKey: false });

BlockSchema.statics = {
    async findList(pageNum: number, pageSize: number): Promise<IBlockStruct[]> {
        return await this.find({})
            .select({
                _id: 0,
                'txn-revno': 0,
                'txn-queue': 0,
            })
            .sort({ height: -1 })
            .skip((pageNum - 1) * pageSize)
            .limit(pageSize).exec();
    },

    async findCount(): Promise<number> {
        return await this.find({}).countDocuments().exec();
    },

    async findOneByHeight(height: number): Promise<IBlockStruct | null> {
        return await this.findOne({ height }).select({
            _id: 0,
            'txn-revno': 0,
            'txn-queue': 0,
        });
    },

    async findOneByHeightDesc(): Promise<IBlockStruct | null> {
        const res: IBlockStruct[] = await this.find({})
            .sort({ height: -1 })
            .select({
                _id: 0,
                'txn-revno': 0,
                'txn-queue': 0,
            })
            .limit(1);
        if (res && res.length > 0) {
            return res[0];
        } else {
            return null;
        }
    },

    async findNum100Height(): Promise<IBlockStruct | null> {
        const res: IBlockStruct[] = await this.find({})
            .select({
                _id: 0,
                'txn-revno': 0,
                'txn-queue': 0,
            })
            .sort({ height: -1 })
            .limit(100);
        if (res && res.length > 0) {
            return res[res.length - 1];
        } else {
            return null;
        }
    },


};