import * as mongoose from 'mongoose';
import { INetworkStruct } from '../types/schemaTypes/network.interface';
import { getTimestamp } from '../util/util';

export const NetworkSchema = new mongoose.Schema({
    network_id: {type:String, required:true, unique: true},
    network_name: String,
    uri:String,
    is_main: Boolean,
    create_time: {
    	type:Number,
    	default:getTimestamp(),
    },
    update_time: {
    	type:Number,
    	default:getTimestamp(),
    }
},{versionKey: false});

NetworkSchema.statics.queryNetworkList = async function ():Promise<INetworkStruct[]>{
	return await this.find();
};


// post networks (未实现)
NetworkSchema.statics.insertNetwork = async function (networks:INetworkStruct[]):Promise<INetworkStruct[]>{
	if (networks && networks.length) {
		let data = networks.map((item)=>{
			return {
				network_id:item.network_id,
                network_name:item.network_name,
                uri:item.uri,
                is_main:item.is_main
			}
		});
		return await this.insertMany(data);
	}else{
		return [];
	}
}

// put networks  (未实现)
NetworkSchema.statics.updateNetwork = async function (network:INetworkStruct):Promise<INetworkStruct>{
	if (network) {
		return await this.findOneAndUpdate({
			network_id:network.network_id,
		},{
			...network,
			update_time:getTimestamp(),
		});
	}else{
		return null;
	}
}

// delete networks  (未实现)
NetworkSchema.statics.deleteNetwork = async function (network_id:string):Promise<INetworkStruct>{
	return await this.findOneAndRemove({network_id:network_id});
}
