import * as mongoose from 'mongoose'
import { IIdentityInfoQuery,IIdentityInfoResponse, IIdentityByAddressQuery} from '../types/schemaTypes/identity.interface';
import { Logger } from '../logger';
import { ITXWithIdentity } from '../types/schemaTypes/tx.interface';
import { IListStruct } from '../types';
import {hubDefaultEmptyValue} from "../constant";

export const IdentitySchema = new mongoose.Schema({
  identities_id: String,
  owner: String,
  credentials: String,
  'create_block_height': Number,
  'create_block_time': Number,
  'create_tx_hash': String,
  'update_block_height': Number,
  'update_block_time': Number,
  'update_tx_hash': String,
  create_time: Number,
  update_time: Number,
})
IdentitySchema.index({identities_id: 1},{unique: true})
IdentitySchema.index({update_block_height: -1,owner:-1})
IdentitySchema.statics = {
  async queryIdentityList(query:ITXWithIdentity):Promise<IListStruct> {
    const result: IListStruct = {}
    const queryParameters: any = {};
    if(query.search && query.search !== ''){
      //单条件模糊查询使用$regex $options为'i' 不区分大小写
      queryParameters.$or = [
        {identities_id:{ $regex: query.search,$options:'i' }},
        {owner:{ $regex: query.search,$options:'i' }}
      ]
    }
    if (query.useCount && query.useCount == true) {
      result.count = await this.find(queryParameters).countDocuments();
    }
    result.data = await this.find(queryParameters)
        .skip((Number(query.pageNum) - 1) * Number(query.pageSize))
        .limit(Number(query.pageSize)).sort({'update_block_height':-1});
    return result;
  },

  async queryIdentityCount(query:any){
    return await this.find(query || {}).countDocuments();
  },

  async queryHeight() {
      const height = await this.findOne({}).sort({'update_block_height': -1})
      const blockHeight = height ? height.update_block_height : 0
      return blockHeight
  },

  async insertIdentityInfo(IdentityInfo) {
      await  this.insertMany(IdentityInfo,{ ordered: false })
  },
  // base information
  async updateIdentityInfo(updateIdentityData) {
    const {identities_id,update_block_time,update_block_height,update_tx_hash,update_time} = updateIdentityData
    if(updateIdentityData.credentials && updateIdentityData.credentials !== hubDefaultEmptyValue){
        const { credentials } = updateIdentityData;
        await this.updateOne({identities_id},{credentials,update_block_time,update_block_height,update_tx_hash,update_time});
      }else {
        await this.updateOne({identities_id},{update_block_time,update_block_height,update_tx_hash,update_time});
      }
  },
  async queryIdentityInfo(id:IIdentityInfoQuery):Promise<IIdentityInfoResponse> {
    const queryId = {identities_id:id.id}
    const infoData:IIdentityInfoResponse = await this.findOne(queryId)
    return  infoData
  },
  // owner
  async queryIdentityByAddress(query: IIdentityByAddressQuery):Promise<IListStruct>{
    const result: IListStruct = {}
    const queryParameters: any = {};
    queryParameters.owner = query.address

    result.data = await this.find(queryParameters)
        .skip((Number(query.pageNum) - 1) * Number(query.pageSize))
        .limit(Number(query.pageSize)).sort({'update_block_height':-1});
    if (query.useCount && query.useCount == true) {
      result.count = await this.find(queryParameters).countDocuments();
    }
    return result
  }
}

