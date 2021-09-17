import * as mongoose from 'mongoose';
import {moduleStakingBondDenom, signedBlocksWindow,govParams,moduleGov} from "../constant";
import {IParameters} from "../types/schemaTypes/parameters.interface";

export const ParametersSchema = new mongoose.Schema({
    module: String,
    key: String,
    cur_value: String,
    create_time: Number,
    update_time: Number,
})

ParametersSchema.index({module:1,key: 1},{unique: true})

ParametersSchema.statics = {

    async insertParameters(Parameters: IParameters) {
        await this.insertMany(Parameters,{ ordered: false })
    },

    async updateParameters(updateParameters: IParameters) {
        const { module, key } = updateParameters
        const options = {upsert: true, new: false, setDefaultsOnInsert: true}
        await this.findOneAndUpdate({module,key},updateParameters,options)
    },

    async queryAllParameters(){
        return await this.find({}).select({'_id':0,'__v':0})
    },

    async querySignedBlocksWindow(moduleName:string){
        return await this.findOne({module:moduleName,key:signedBlocksWindow}).select({'_id':0,'__v':0})
    },
    async queryStakingToken(moduleName:string){
        return await this.findOne({module:moduleName,key:moduleStakingBondDenom}).select({'_id':0,'__v':0})
    },
    async queryGovParams() {
        const param = {
            module: moduleGov,
            key: {
                $in: [govParams.min_deposit,govParams.quorum,govParams.threshold,govParams.veto_threshold]
            }
        }
        return await this.find(param).select({'_id':0,'__v':0})
    },
    async queryParams(query) {
        let { module, key } = query
        let param: { module?: string, key?: object } = {};
        if (module) {
            param['module'] = module
        }
        if (key) {
            param['key'] = {
                $in: key.split(',')
            }
        }
        return await this.find(param).select({'module':1,'key':1,'cur_value':1,'_id':0})
    }
}
