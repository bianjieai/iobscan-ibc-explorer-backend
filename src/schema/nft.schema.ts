import * as mongoose from 'mongoose';
import { Logger } from '../logger';
import {
    IDeleteQuery,
    INftCountQueryParams,
    INftStruct,
} from '../types/schemaTypes/nft.interface';
import { IListStruct } from '../types';
import { getTimestamp } from '../util/util';

export const NftSchema = new mongoose.Schema({
    denom_id: String,
    denom_name: String,
    nft_id: String,
    nft_name: String,
    owner: String,
    uri: String,
    data: String,
    last_block_height: Number,
    last_block_time: Number,
    create_time: Number,
    update_time: Number,
    hash: String,
    time: Number
}, { versionKey: false });
NftSchema.index({ denom_id: 1, nft_id: 1 }, { unique: true });
// 新增
NftSchema.index({ last_block_height: 1, denom_id: 1, nft_id: 1 }, { background:true});
NftSchema.index({ owner: 1, last_block_height: 1 }, { background:true });

NftSchema.statics = {
    async findList(
        pageNum: number,
        pageSize: number,
        denomId?: string,
        nftId?: string,
        owner?: string,
        useCount?:boolean,
    ): Promise<IListStruct> {
        let result: IListStruct = {};
        // const condition: any[] = [
            // {
            //     $lookup: {
            //         from: 'ex_sync_denom',
            //         localField: 'denom_id',
            //         foreignField: 'denom_id',
            //         as: 'denomDetail',
            //     },
            // }, {
            //     $project: {
            //         'denomDetail._id': 0,
            //         'denomDetail.update_time': 0,
            //         'denomDetail.create_time': 0,
            //     },
            // },
        // ];

        let queryParameters:any = {};
        if (denomId || nftId || owner) {
            if (denomId) queryParameters.denom_id = denomId;
            if (nftId) queryParameters['$or']= [
                {'nft_name': nftId},
                {'nft_id': nftId},
            ];
            if (owner) queryParameters.owner = owner;
            // condition.push({'$match': queryParameters});
        }
        result.data = await this.find(queryParameters)
            .sort({last_block_height:-1})
            .skip((Number(pageNum) - 1) * Number(pageSize))
            .limit(Number(pageSize));

        // result.data = await this.aggregate(condition)
        //     .sort({create_time:-1, nft_id:-1})
        //     .skip((Number(pageNum) - 1) * Number(pageSize))
        //     .limit(Number(pageSize));
        if (useCount) {
            result.count = await this.find(queryParameters).countDocuments();
        }
        return result; 
    },

    async findOneByDenomAndNftId(denomId: string, nftId: string): Promise<INftStruct | null> {
        return await this.findOne({denom_id:denomId,nft_id:nftId})
    },

    async findCount(denomId: string, 
        nftIdOrName: string, 
        nftName: string, 
        owner: string, 
    ): Promise<number> {

        let query: any = {};
        if (denomId){
            query.denom_id = denomId;
        }
        if (nftIdOrName){
            const reg = new RegExp(nftIdOrName, 'i');
            query['$or'] = [
                { 'nft_name': nftIdOrName },
                { 'nft_id': nftIdOrName },
            ];
        }
        if (owner){
            query.owner = owner;
        }
        
        return await this.find(query).countDocuments();
    },

    async findListByName(denomId: string): Promise<INftStruct> {
        return await this.find({ denom_id: denomId }).exec();
    },

    saveBulk(nfts: INftStruct[]): Promise<INftStruct[]> {
        return this.insertMany(nfts, { ordered: false });
    },

    async deleteOneByDenomAndId(nft: IDeleteQuery): Promise<INftStruct> {
        return await this.deleteOne(nft, (e) => {
            if (e) Logger.error('mongo-error:', e.message);
        });
    },

    updateOneById(nft: INftStruct): Promise<INftStruct> {
        const {denom_id ,nft_id, owner, data, uri, denom_name, nft_name, hash, time } = nft;
        return this.updateOne({
            nft_id,
            denom_id
        }, {
            owner,
            data,
            uri,
            denom_name,
            nft_name,
            hash,
            update_time: getTimestamp(),
            time
        }, 
        (e) => {
            if (e) Logger.error('mongo-error:', e.message);
        });
    },

    async queryNftCount(denomId: string): Promise<INftStruct>{
        return await this.find({denom_id:denomId}).countDocuments().exec();
    },

    async queryLastBlockHeight(): Promise<INftStruct>{
        return await this.find({},{last_block_height: 1}).sort({last_block_height: -1}).limit(1);
    },

    updateNft(nft: INftStruct): Promise<INftStruct>{
        let cond = {
            denom_id:nft.denom_id,
            nft_id:nft.nft_id,
        };
        return this.findOneAndUpdate(cond,nft,{ upsert:true,new: true});
    },

    deleteNft(nft: INftStruct): Promise<INftStruct>{
        let cond = {
            denom_id:nft.denom_id,
            nft_id:nft.nft_id,
        };
        return this.deleteOne(cond);
    },






};