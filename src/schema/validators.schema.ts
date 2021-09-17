import * as mongoose from 'mongoose';
import { IValidatorsQueryParams } from "../types/validators.interface"
import { Logger } from '@nestjs/common';
import { IValidatorsStruct } from '../types/schemaTypes/validators.interface';
export const ValidatorSchema = new mongoose.Schema({
    name: String,
    pubkey: String,
    power: String,
    jailed: Boolean,
    operator: String,
    details:String,
    hash:String
})
ValidatorSchema.index({name: 1},{unique: true})
ValidatorSchema.index({jailed: 1})

ValidatorSchema.statics.findValidators = async function (query:IValidatorsQueryParams):Promise<{count?: number, data?:IValidatorsStruct[]}> {
    let result: { count?: number,data?: Array<IValidatorsStruct> } = { }
    let queryParams = {
        jailed: undefined,
    };
    queryParams.jailed = query.jailed;
    if(query && query.useCount){
        result.count = await  this.countDocuments(queryParams)
    }
    result.data = await this.find(queryParams).skip((Number(query.pageNum) - 1) * Number(query.pageSize))
      .limit(Number(query.pageSize)).select({'_id':0,'__v':0,'hash':0})
    return  result
}
ValidatorSchema.statics.findCount = async function (isJailed: boolean):Promise<number> {
    return await this.countDocuments({jailed:isJailed});
}

ValidatorSchema.statics.findAllValidators = async function ():Promise<IValidatorsStruct[]>{
    return await this.find({}).select({'_id':0,'__v':0})
}

ValidatorSchema.statics.saveValidator = async  function (insertValidatorList:IValidatorsStruct[]):Promise<IValidatorsStruct[]> {
   return await this.insertMany(insertValidatorList,{ordered: false})
}

ValidatorSchema.statics.updateValidator = async function (name:string,needUpdateValidator:IValidatorsStruct):Promise<IValidatorsStruct> {
    return await this.updateOne({name:name},needUpdateValidator)
}
ValidatorSchema.statics.deleteValidator = async function (validatorName:string) {
    return await this.deleteOne({name:validatorName});
}
