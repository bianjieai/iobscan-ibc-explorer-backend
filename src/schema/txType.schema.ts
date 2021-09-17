import * as mongoose from 'mongoose';
import { ITxTypeStruct } from '../types/schemaTypes/txType.interface';
import { getTimestamp } from '../util/util';
import { 
    stakingTypes,
    serviceTypes,
	declarationTypes,
	govTypes
} from '../helper/txTypes.helper';
import { TxType } from '../constant';
export const TxTypeSchema = new mongoose.Schema({
    type_name:{type:String, required:true, unique: true},
    create_time:{
    	type:Number,
    	default:getTimestamp(),
    },
    update_time:{
    	type:Number,
    	default:getTimestamp(),
    }
},{versionKey: false});

// txs/types
TxTypeSchema.statics.queryTxTypeList = async function ():Promise<ITxTypeStruct[]>{
	return await this.find({},{type_name:1})
};

// txs/types/service
TxTypeSchema.statics.queryServiceTxTypeList = async function ():Promise<ITxTypeStruct[]>{
    let queryParameters: any = {
        type_name:{'$in':serviceTypes()}
    };
    return await this.find(queryParameters,{type_name:1});
};

// txs/types/staking
TxTypeSchema.statics.queryStakingTxTypeList = async function ():Promise<ITxTypeStruct[]>{
    let queryParameters: any = {
        type_name:{'$in':stakingTypes()}
    };
    return await this.find(queryParameters,{type_name:1});
};

// txs/types/declaration
TxTypeSchema.statics.queryDeclarationTxTypeList = async function ():Promise<ITxTypeStruct[]>{
    let queryParameters: any = {
        type_name:{'$in':declarationTypes()}
    };
    return await this.find(queryParameters,{type_name:1});
};


// post txs/types
TxTypeSchema.statics.insertTxTypes = async function (types:string[]):Promise<ITxTypeStruct[]>{
	if (types && types.length) {
		let data = types.map((t)=>{
			let item = {
				type_name:t,
			}
		    return new this(item);
		});
		return await this.insertMany(data);
	}else{
		return [];
	}
}

// put txs/types
TxTypeSchema.statics.updateTxType = async function (type:string, newType:string):Promise<ITxTypeStruct>{
	if (type && type.length && newType && newType.length) {
		return await this.findOneAndUpdate({
			type_name:type,
		},{
			type_name:newType,
			update_time:getTimestamp(),
		});
	}else{
		return null;
	}
}

// delete txs/types
TxTypeSchema.statics.deleteTxType = async function (type:string):Promise<ITxTypeStruct>{
	return await this.findOneAndRemove({type_name:type});
}

// txs/types/staking
TxTypeSchema.statics.queryGovTxTypeList = async function ():Promise<ITxTypeStruct[]>{
    let queryParameters: any = {
        type_name:{'$in':govTypes()}
    };
    return await this.find(queryParameters,{type_name:1});
};
