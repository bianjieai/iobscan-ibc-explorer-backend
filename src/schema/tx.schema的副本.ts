import { HttpException } from '@nestjs/common';
import * as mongoose from 'mongoose';
import {
    ITxsQuery,
    ITxsWithHeightQuery,
    ITxsWithAddressQuery,
    ITxsWithContextIdQuery,
    ITxsWithNftQuery,
    ITxsWithServiceNameQuery,
    IExFieldQuery, IIdentityTx,
    ITxStruct,
    ITxStructMsgs,
    ITxStructHash,
    ITxsWithAssetQuery,
    ITxVoteProposal,
    ITxSubmitProposal,
    ITxVoteProposalAll,
    ITxVoteALL
} from '../types/schemaTypes/tx.interface';
import { IBindTx, IServiceName, ITxsQueryParams } from '../types/tx.interface';
import { IListStruct,IQueryBase } from '../types';
import { INCREASE_HEIGHT, TxStatus, TxType,MAX_OPERATE_TX_COUNT } from '../constant';
import Cache from '../helper/cache';
import { dbRes } from '../helper/tx.helper';
import { cfg } from '../config/config';
import { PagingReqDto } from '../dto/base.dto';
import {
    stakingTypes,
    serviceTypes,
    declarationTypes,
    govTypes,
    coinswapTypes
} from '../helper/txTypes.helper';

export const TxSchema = new mongoose.Schema({
    time: Number,
    height: Number,
    tx_hash: String,
    memo: String,
    status: Number,
    log: String,
    complex_msg: Boolean,
    type: String,
    from: String,
    to: String,
    coins: Array,
    signer: String,
    events: Array,
    events_new: Array,
    msgs: Array,
    signers: Array,
    addrs: Array,
    fee: Object,
    tx_index: Number,
}, { versionKey: false });
TxSchema.index({ time: -1, "msgs.type": -1,status:-1 }, { background: true });
TxSchema.index({ addrs: -1, time: -1, status:-1 }, { background: true });
TxSchema.index({"msgs.type": -1,height:-1,"msgs.msg.ex.service_name":-1 }, { background: true });


//	csrb 浏览器交易记录过滤正则表达式
// function filterExTxTypeRegExp(): object {
//     let RegExpStr: string = Cache.supportTypes.join('|');
//     return new RegExp(RegExpStr || '//');
// }

// function filterTxTypeRegExp(types: TxType[]): object {
//     if (!Array.isArray(types) || types.length === 0) return;
//     let typeStr: string = ``;
//     types.forEach((type: TxType) => {
//         if (!typeStr) {
//             typeStr = `${type}`;
//         } else {
//             typeStr += `|${type}`;
//         }
//     });
//     return new RegExp(typeStr);
// }


// 	txs
TxSchema.statics.queryTxList = async function(query: ITxsQuery): Promise<IListStruct> {
    let result: IListStruct = {};
    let queryParameters: ITxsQueryParams = {};
    const { type } = query;
    if (query.type && query.type.length) {
        let typeArr = type.split(",");
        queryParameters['msgs.type'] = {
            $in: typeArr
        }
    } else {
        // queryParameters.$or = [{ 'msgs.type': filterExTxTypeRegExp() }];
        queryParameters['msgs.type'] = {
            $in: Cache.supportTypes || []
        }
    }
    if (query.status && query.status.length) {
        switch (query.status) {
            case '1':
                queryParameters.status = TxStatus.SUCCESS;
                break;
            case '2':
                queryParameters.status = TxStatus.FAILED;
                break;
        }
    }
    if ((query.beginTime && query.beginTime.length) || (query.endTime && query.endTime.length)) {
        queryParameters.time = {};
    }
    if (query.beginTime && query.beginTime.length) {
        queryParameters.time.$gte = Number(query.beginTime);
    }
    if (query.endTime && query.endTime.length) {
        queryParameters.time.$lte = Number(query.endTime);
    }
    if (query.address && query.address.length) {
        queryParameters['addrs'] = { $elemMatch: { $eq: query.address } };
    }
    result.data = await this.find(queryParameters, dbRes.txList)
        .sort({ time: -1 })
        .skip((Number(query.pageNum) - 1) * Number(query.pageSize))
        .limit(Number(query.pageSize));
    if (query.useCount && query.useCount == true) {
        result.count = await this.find(queryParameters).countDocuments();
    }
    return result;
};

//  txs/staking
TxSchema.statics.queryStakingTxList = async function(query: ITxsQuery): Promise<IListStruct> {
    const result: IListStruct = {};
    let queryParameters: any = {};
    if (query.type && query.type.length) {
        queryParameters['msgs.type'] = query.type;
    } else {
        queryParameters['msgs.type'] = {'$in':stakingTypes()};
    }
    if (query.status && query.status.length) {
        switch (query.status) {
            case '1':
                queryParameters.status = TxStatus.SUCCESS;
                break;
            case '2':
                queryParameters.status = TxStatus.FAILED;
                break;
        }
    }
    if (query.address && query.address.length) {
        queryParameters['addrs'] = { $elemMatch: { $eq: query.address } };
    }
    if ((query.beginTime && query.beginTime.length) || (query.endTime && query.endTime.length)) {
        queryParameters.time = {};
    }
    if (query.beginTime && query.beginTime.length) {
        queryParameters.time.$gte = Number(query.beginTime);
    }
    if (query.endTime && query.endTime.length) {
        queryParameters.time.$lte = Number(query.endTime);
    }
    result.data = await this.find(queryParameters, dbRes.delegations)
        .sort({ time: -1 })
        .skip((Number(query.pageNum) - 1) * Number(query.pageSize))
        .limit(Number(query.pageSize));
    if (query.useCount && query.useCount == true) {
        result.count = await this.find(queryParameters).countDocuments();
    }
    return result;
};

//  txs/coinswap
TxSchema.statics.queryCoinswapTxList = async function(query: ITxsQuery): Promise<IListStruct> {
    const result: IListStruct = {};
    let queryParameters: any = {};
    const { type } = query;
    if (query.type && query.type.length) {
        let typeArr = type.split(",");
        queryParameters['msgs.type'] = {
            $in: typeArr
        }
    } else {
        queryParameters['msgs.type'] = {'$in':coinswapTypes()};
    }
    if (query.status && query.status.length) {
        switch (query.status) {
            case '1':
                queryParameters.status = TxStatus.SUCCESS;
                break;
            case '2':
                queryParameters.status = TxStatus.FAILED;
                break;
        }
    }
    if (query.address && query.address.length) {
        queryParameters['addrs'] = { $elemMatch: { $eq: query.address } };
    }
    if ((query.beginTime && query.beginTime.length) || (query.endTime && query.endTime.length)) {
        queryParameters.time = {};
    }
    if (query.beginTime && query.beginTime.length) {
        queryParameters.time.$gte = Number(query.beginTime);
    }
    if (query.endTime && query.endTime.length) {
        queryParameters.time.$lte = Number(query.endTime);
    }
    result.data = await this.find(queryParameters)
        .sort({ time: -1 })
        .skip((Number(query.pageNum) - 1) * Number(query.pageSize))
        .limit(Number(query.pageSize));
    if (query.useCount && query.useCount == true) {
        result.count = await this.find(queryParameters).countDocuments();
    }
    return result;
};

//  txs/declaration
TxSchema.statics.queryDeclarationTxList = async function(query: ITxsQuery): Promise<IListStruct> {
    const result: IListStruct = {};
    let queryParameters: any = {};
    if (query.type && query.type.length) {
        queryParameters['msgs.type'] = query.type;
    } else {
        queryParameters['msgs.type'] = { '$in': declarationTypes() };
    }
    if (query.status && query.status.length) {
        switch (query.status) {
            case '1':
                queryParameters.status = TxStatus.SUCCESS;
                break;
            case '2':
                queryParameters.status = TxStatus.FAILED;
                break;
        }
    }
    if (query.address && query.address.length) {
        queryParameters['addrs'] = { $elemMatch: { $eq: query.address } };
    }
    if ((query.beginTime && query.beginTime.length) || (query.endTime && query.endTime.length)) {
        queryParameters.time = {};
    }
    if (query.beginTime && query.beginTime.length) {
        queryParameters.time.$gte = Number(query.beginTime);
    }
    if (query.endTime && query.endTime.length) {
        queryParameters.time.$lte = Number(query.endTime);
    }
    result.data = await this.find(queryParameters, dbRes.validations)
        .sort({time: -1})
        .skip((Number(query.pageNum) - 1) * Number(query.pageSize))
        .limit(Number(query.pageSize));
    if (query.useCount && query.useCount == true) {
        result.count = await this.find(queryParameters).countDocuments();
    }
    return result;
};

//  txs/gov
TxSchema.statics.queryGovTxList = async function(query: ITxsQuery): Promise<IListStruct> {
    const result: IListStruct = {};
    let queryParameters: any = {};
    if (query.type && query.type.length) {
        queryParameters['msgs.type'] = query.type;
    } else {
        queryParameters['msgs.type'] = { '$in': govTypes() };
    }
    if (query.status && query.status.length) {
        switch (query.status) {
            case '1':
                queryParameters.status = TxStatus.SUCCESS;
                break;
            case '2':
                queryParameters.status = TxStatus.FAILED;
                break;
        }
    }
    if (query.address && query.address.length) {
        queryParameters['addrs'] = { $elemMatch: { $eq: query.address } };
    }
    if ((query.beginTime && query.beginTime.length) || (query.endTime && query.endTime.length)) {
        queryParameters.time = {};
    }
    if (query.beginTime && query.beginTime.length) {
        queryParameters.time.$gte = Number(query.beginTime);
    }
    if (query.endTime && query.endTime.length) {
        queryParameters.time.$lte = Number(query.endTime);
    }
    result.data = await this.find(queryParameters, dbRes.govs)
        .sort({time: -1})
        .skip((Number(query.pageNum) - 1) * Number(query.pageSize))
        .limit(Number(query.pageSize));
    if (query.useCount && query.useCount == true) {
        result.count = await this.find(queryParameters).countDocuments();
    }
    return result;
};






//  txs/e 供edgeServer调用  返回数据不做过滤
TxSchema.statics.queryTxListEdge = async function(types:string, gt_height:number, pageNum:number, pageSize:number, useCount:boolean,status?:number,address?:string,include_event_addr?:boolean): Promise<IListStruct> {
    let result: IListStruct = {};
    let queryParameters: any = {};
    if (types && types.length) {
        queryParameters['msgs.type'] = {'$in':types.split(',')};
    }
    if (gt_height) {
        queryParameters['height'] = {'$gt':gt_height};
    }
    if (status || status === 0) {
        queryParameters['status'] = status;
    }
    if (include_event_addr && include_event_addr == true && address && address.length) {
        queryParameters['$or'] = [
            { 'events.attributes.value': address },
            { 'addrs': { $elemMatch: {'$in': address.split(',')} }}
        ]
    } else if (address && address.length) {
        queryParameters['addrs'] = { $elemMatch: { '$in': address.split(',') } };
    }
    result.data = await this.find(queryParameters)
    .sort({ height: 1 })
    .skip((Number(pageNum) - 1) * Number(pageSize))
    .limit(Number(pageSize));
    if (useCount && useCount == true) {
        result.count = await this.find(queryParameters).countDocuments();
    }
    return result;
};

//  供edgeServer调用  返回数据不做过滤
TxSchema.statics.queryTxListByHeightEdge = async function(height:number, pageNum:number, pageSize:number, useCount:boolean,status?:number): Promise<IListStruct> {
    let result: IListStruct = {};
    let queryParameters: any = { height: height };
    if (status || status === 0) {
        queryParameters['status'] = status;
    }
    result.data = await this.find(queryParameters)
        .skip((Number(pageNum) - 1) * Number(pageSize))
        .limit(Number(pageSize));

    if (useCount && useCount == true) {
        result.count = await this.find(queryParameters).countDocuments();
    }
    return result;
};


// 	txs/blocks
TxSchema.statics.queryTxWithHeight = async function(query: ITxsWithHeightQuery): Promise<IListStruct> {
    let result: IListStruct = {};
    // let queryParameters: { height?: number, $or: object[] } = { $or: [{ 'msgs.type': filterExTxTypeRegExp() }] };
    let queryParameters: { height?: number, 'msgs.type': object } = { 'msgs.type': { $in: Cache.supportTypes || [] } };
    if (query.height) {
        queryParameters.height = Number(query.height);
    }
    result.data = await this.find(queryParameters, dbRes.txList)
        .sort({ height: -1 })
        .skip((Number(query.pageNum) - 1) * Number(query.pageSize))
        .limit(Number(query.pageSize));
    if (query.useCount && query.useCount == true) {
        result.count = await this.find(queryParameters).countDocuments();
    }
    return result;
};

//  txs/addresses
TxSchema.statics.queryTxWithAddress = async function(query: ITxsWithAddressQuery): Promise<IListStruct> {
    let result: IListStruct = {};
    let queryParameters: any = {};
    if (query.address && query.address.length) {
        queryParameters = {
            // $or:[
            // 	{"from":query.address},
            // 	{"to":query.address},
            // 	{"signer":query.address},
            // ],
            addrs: { $elemMatch: { $eq: query.address } },
        };
    }
    if (query.type && query.type.length) {
        queryParameters['msgs.type'] = query.type;
    } else {
        // queryParameters.$or = [{ 'msgs.type': filterExTxTypeRegExp() }];
        queryParameters['msgs.type'] = {
            $in: Cache.supportTypes || []
        }
    }
    if (query.status && query.status.length) {
        switch (query.status) {
            case '1':
                queryParameters.status = TxStatus.SUCCESS;
                break;
            case '2':
                queryParameters.status = TxStatus.FAILED;
                break;
        }
    }
    result.data = await this.find(queryParameters, dbRes.txList)
        .sort({ height: -1 })
        .skip((Number(query.pageNum) - 1) * Number(query.pageSize))
        .limit(Number(query.pageSize));
    if (query.useCount && query.useCount == true) {
        result.count = await this.find(queryParameters).countDocuments();
    }
    return result;
};

//  txs/relevance
TxSchema.statics.queryTxWithContextId = async function(query: ITxsWithContextIdQuery): Promise<IListStruct> {
    let result: IListStruct = {};
    let queryParameters: any = {};
    if (query.contextId && query.contextId.length) {
        queryParameters = {
            $or: [
                { 'events.attributes.value': query.contextId },
                { 'msgs.msg.ex.request_context_id': query.contextId },
                { 'msgs.msg.request_context_id': query.contextId },
            ],
        };
    }
    if (query.type && query.type.length) {
        queryParameters['msgs.type'] = query.type;
    } else {
        // queryParameters.$or = [{ 'msgs.type': filterExTxTypeRegExp() }];
        queryParameters['msgs.type'] = {
            $in: Cache.supportTypes || []
        }
    }

    if (query.status && query.status.length) {
        switch (query.status) {
            case '1':
                queryParameters.status = TxStatus.SUCCESS;
                break;
            case '2':
                queryParameters.status = TxStatus.FAILED;
                break;
        }
    }
    result.data = await this.find(queryParameters, dbRes.service)
        .sort({ height: -1 })
        .skip((Number(query.pageNum) - 1) * Number(query.pageSize))
        .limit(Number(query.pageSize));
    if (query.useCount && query.useCount == true) {
        result.count = await this.find(queryParameters).countDocuments();
    }
    return result;
};

//  txs/nfts
TxSchema.statics.queryTxWithNft = async function(query: ITxsWithNftQuery): Promise<IListStruct> {
    let result: IListStruct = {};
    const nftTypesList: TxType[] = [
        TxType.mint_nft,
        TxType.edit_nft,
        TxType.transfer_nft,
        TxType.burn_nft,
    ];
    // let queryParameters: { 'msgs.msg.denom'?: string, 'msgs.msg.id'?: string, $or: object[] } = { $or: [{ 'msgs.type': filterTxTypeRegExp(nftTypesList) }] };
    let queryParameters: { 'msgs.msg.denom'?: string, 'msgs.msg.id'?: string, 'msgs.type': object } = { 'msgs.type': { $in: nftTypesList || [] }};
    if (query.denomId && query.denomId.length) {
        queryParameters['msgs.msg.denom'] = query.denomId;
    }
    if (query.tokenId && query.tokenId.length) {
        queryParameters['msgs.msg.id'] = query.tokenId;
    }
    result.data = await this.find(queryParameters, dbRes.txList)
        .sort({ height: -1 })
        .skip((Number(query.pageNum) - 1) * Number(query.pageSize))
        .limit(Number(query.pageSize));
    if (query.useCount && query.useCount == true) {
        result.count = await this.find(queryParameters).countDocuments();
    }
    return result;
};

// 废弃
TxSchema.statics.queryTxWithServiceName = async function(query: ITxsWithServiceNameQuery): Promise<IListStruct> {
    let result: IListStruct = {};
    let queryParameters: any = {};
    if (query.serviceName && query.serviceName.length) {
        queryParameters = {
            $or: [
                { 'msgs.msg.ex.service_name': query.serviceName },
                { 'msgs.type': { $in: Cache.supportTypes || [] } }
            ],
        };
    }
    result.data = await this.find(queryParameters)
        .sort({ height: -1 })
        .skip((Number(query.pageNum) - 1) * Number(query.pageSize))
        .limit(Number(query.pageSize));
    if (query.useCount && query.useCount == true) {
        result.count = await this.find(queryParameters).countDocuments();
    }
    return result;
};

//  txs/services/detail/{serviceName}
TxSchema.statics.queryTxDetailWithServiceName = async function(serviceName: string): Promise<ITxStruct> {
    return await this.findOne({ 'msgs.msg.name': serviceName, 'msgs.type': TxType.define_service, status:TxStatus.SUCCESS});
};

// ==> txs/services/call-service
TxSchema.statics.queryCallServiceWithConsumerAddr = async function(consumerAddr: string, pageNum: string, pageSize: string, useCount: boolean): Promise<IListStruct> {
    let result: IListStruct = {};
    let queryParameters: any = {
        'msgs.msg.consumer': consumerAddr,
        'msgs.type': TxType.call_service,
        status: TxStatus.SUCCESS,
    };
    result.data = await this.find(queryParameters, dbRes.service)
        .sort({ height: -1 })
        .skip((Number(pageNum) - 1) * Number(pageSize))
        .limit(Number(pageSize));
    if (useCount && useCount == true) {
        result.count = await this.find(queryParameters).countDocuments();
    }
    return result;
};

// ==> txs/services/call-service
TxSchema.statics.queryRespondServiceWithContextId = async function(ContextId: string): Promise<ITxStruct[]> {
    return await this.find({ 
        'msgs.msg.ex.request_context_id': ContextId, 
        'msgs.type': TxType.respond_service }, dbRes.service);
};

// ==> txs/services/respond-service
TxSchema.statics.queryBindServiceWithProviderAddr = async function(ProviderAddr: string, pageNum: string, pageSize: string, useCount: boolean): Promise<IListStruct> {
    let result: IListStruct = {};
    let queryParameters: any = {
        'msgs.msg.provider': ProviderAddr,
        'msgs.type': TxType.bind_service,
        status: TxStatus.SUCCESS,
    };
    result.data = await this.find(queryParameters, dbRes.service)
        .sort({ height: -1 })
        .skip((Number(pageNum) - 1) * Number(pageSize))
        .limit(Number(pageSize));
    if (useCount && useCount == true) {
        result.count = await this.find(queryParameters).countDocuments();
    }
    return result;
};

// ==> txs/services/respond-service
TxSchema.statics.queryRespondCountWithServceName = async function(servceName: string, providerAddr: string): Promise<ITxStruct[]> {
    return await this.find({
        'msgs.msg.ex.service_name': servceName,
        'msgs.msg.provider': providerAddr,
        'msgs.type': TxType.respond_service,
    }).countDocuments();
};

// ==> txs/services/respond-service
TxSchema.statics.querydisableServiceBindingWithServceName = async function(servceName: string, providerAddr: string): Promise<ITxStruct[]> {
    return await this.find({
        'msgs.msg.ex.service_name': servceName,
        'msgs.msg.provider': providerAddr,
        'msgs.type': TxType.disable_service_binding,
    },{time:1})
        .sort({ height: -1 })
        .limit(1);
};


//  txs/{hash}
TxSchema.statics.queryTxWithHash = async function(hash: string): Promise<ITxStruct> {
    return await this.findOne({ tx_hash: hash });
};

//  /statistics
TxSchema.statics.queryTxCountStatistics = async function(): Promise<number> {
    return await this.find({ 'msgs.type': { $in: Cache.supportTypes || [] } }).countDocuments();
};

TxSchema.statics.queryServiceCountStatistics = async function(): Promise<number> {
    return await this.find({ 'msgs.type': TxType.define_service, status: TxStatus.SUCCESS }).countDocuments();
};

//	获取指定条数的serviceName==null&&type == respond_service 的 tx
// TxSchema.statics.findRespondServiceTx = async function(pageSize?:number):Promise<ITxStructHash[]>{
// 	pageSize = pageSize || cfg.taskCfg.syncTxServiceNameSize;
// 	return await this.find({
// 							type:TxType.respond_service,
// 							'msgs.msg.ex.service_name':null
// 						},{tx_hash:1,'msgs.msg.request_id':1})
// 					 .sort({time:-1})
// 					 .limit(Number(pageSize));
// }

//	根据Request_Context_Id list && type == call_service 获取指定tx list
TxSchema.statics.findCallServiceTxWithReqContextIds = async function(reqContextIds: string[]): Promise<ITxStructMsgs[]> {
    if (!reqContextIds || !reqContextIds.length) {
        return [];
    }
    ;
    let query = {
        'msgs.type': TxType.call_service,
        'events.attributes.key': 'request_context_id',
        'events.attributes.value': { $in: reqContextIds },
    };
    return await this.find(query, {
        'events.attributes': 1,
        'msgs.msg.service_name': 1,
        'msgs.msg.consumer': 1,
        'tx_hash': 1,
    });
};

//	根据Request_Context_Id list && type == call_service 获取指定tx list
// TxSchema.statics.updateServiceNameToResServiceTxWithTxHash = async function(txHash:string, serviceName:string, requestContextId:string, callHash: string, consumer: string):Promise<ITxStruct>{
// 	if (!txHash || !txHash.length) {return null};
// 	let query = {
// 		tx_hash:txHash,
// 	};
// 	let updateParams: any = {
// 	    $set: {
// 	        'msgs.0.msg.ex.service_name': serviceName,
//             'msgs.0.msg.ex.request_context_id': requestContextId,
//             'msgs.0.msg.ex.call_hash': callHash,
//             'msgs.0.msg.ex.consumer': consumer,
// 	    }
// 	};
// 	return await this.findOneAndUpdate(query,updateParams);
// }

//定时任务, 查询所有关于service的tx
TxSchema.statics.findAllServiceTx = async function(pageSize?: number): Promise<ITxStruct[]> {
    pageSize = pageSize || cfg.taskCfg.syncTxServiceNameSize;
    const serviceTypesList: TxType[] = [
        TxType.define_service,
        TxType.bind_service,
        TxType.call_service,
        TxType.respond_service,
        TxType.update_service_binding,
        TxType.disable_service_binding,
        TxType.enable_service_binding,
        TxType.refund_service_deposit,
        TxType.pause_request_context,
        TxType.start_request_context,
        TxType.kill_request_context,
        TxType.update_request_context
    ];

    let queryParameters: any = {
        'msgs.type': { $in: serviceTypesList },
        'msgs.msg.ex.service_name': null,
    };
    return await this.find(queryParameters, dbRes.syncServiceTask).sort({ height: -1 }).limit(Number(pageSize));
};

//用request_context_id查询call_service的service_name
TxSchema.statics.queryServiceName = async function(requestContextId: string): Promise<string> {
    let queryParameters: any = {
        'msgs.type': TxType.call_service,
        'events.attributes.key': 'request_context_id',
        'events.attributes.value': requestContextId.toUpperCase(),
        'status': TxStatus.SUCCESS,
    };
    return await this.findOne(queryParameters);
};

//在msg结构中增加ex字段
TxSchema.statics.addExFieldForServiceTx = async function(ex: IExFieldQuery): Promise<string> {
    const { requestContextId, consumer, serviceName, callHash, hash, bind } = ex;
    let updateParams: any = {
        $set: {},
    };
    if (requestContextId && requestContextId.length) {
        updateParams['$set']['msgs.0.msg.ex.request_context_id'] = requestContextId;
    }
    if (consumer && consumer.length) {
        updateParams['$set']['msgs.0.msg.ex.consumer'] = consumer;
    }
    if (serviceName && serviceName.length) {
        updateParams['$set']['msgs.0.msg.ex.service_name'] = serviceName;
    }
    if (callHash && callHash.length) {
        updateParams['$set']['msgs.0.msg.ex.call_hash'] = callHash;
    }
    if (bind) {
        updateParams['$set']['msgs.0.msg.ex.bind'] = bind;
    }
    return await this.findOneAndUpdate({ tx_hash: hash }, updateParams);
};

//根据serviceName 查询define_service tx
TxSchema.statics.queryDefineServiceTxHashByServiceName = async function(serviceName: string, status?:TxStatus): Promise<ITxStruct> {
    let queryParameters: any = {
        'msgs.type': TxType.define_service,
        'msgs.msg.name': serviceName,
    };
    if (typeof status != 'undefined') {
        queryParameters.status = status;
    }
    return await this.findOne(queryParameters, { 'tx_hash': 1 });
};

// txs/services
TxSchema.statics.findServiceAllList = async function(
    pageNum: number,
    pageSize: number,
    useCount: boolean | undefined,
    nameOrDescription?: string,
): Promise<ITxStruct> {

    const queryParameters: any = {
        'msgs.type': TxType.define_service,
        status: TxStatus.SUCCESS,
    };
    if (nameOrDescription) {
        const reg = new RegExp(nameOrDescription, 'i');
        queryParameters['$or'] = [
            { 'msgs.msg.name': { $regex: reg } },
            { 'msgs.msg.description': { $regex: reg } },
        ];
    }
    return await this.find(queryParameters, {
        'msgs.msg.ex':1,
        'msgs.msg.description':1,
        'msgs.msg.service_name':1,
        'msgs.msg.name':1})
        .sort({
            'msgs.msg.ex.bind': -1,
            height: -1,
        })
        .skip((Number(pageNum) - 1) * Number(pageSize))
        .limit(Number(pageSize));
};

// /txs/services/providers 
TxSchema.statics.findBindServiceTxList = async function(
    serviceName: IServiceName,
    pageNum?: number,
    pageSize?: number,
): Promise<ITxStruct> {
    const queryParameters: any = {
        'msgs.type': TxType.bind_service,
        status: TxStatus.SUCCESS,
        'msgs.msg.ex.service_name': serviceName,
    };
    if (pageNum && pageSize) {
        return await this.find(queryParameters,{'msgs.msg.provider':1,time:1})
            .sort({ 'height': -1 })
            .skip((Number(pageNum) - 1) * Number(pageSize))
            .limit(Number(pageSize));
    } else {
        return await this.find(queryParameters).sort({ 'height': -1 });
    }

};

//查询某个provider在某个service下所有的响应次数
TxSchema.statics.findProviderRespondTimesForService = async function(serviceName: string, provider: string): Promise<number> {
    const queryParameters: any = {
        'msgs.type': TxType.respond_service,
        'msgs.msg.ex.service_name': serviceName,
        'msgs.msg.provider': provider,
    };
    return await this.countDocuments(queryParameters);
};

//查询某个consumer在某个service下所有的调用次数
// TxSchema.statics.findConsumerCallServiceForService = async function(serviceName: string, consumer: string): Promise<number> {
//     const queryParameters: any = {
//         'msgs.type': TxType.call_service,
//         'msgs.msg.service_name': serviceName,
//         'msgs.msg.consumer': consumer,
//     };
//     return await this.countDocuments(queryParameters);
// };

TxSchema.statics.findAllServiceCount = async function(nameOrDescription?: string): Promise<number> {
    const queryParameters: any = {
        'msgs.type': TxType.define_service,
        status: TxStatus.SUCCESS,
    };
    if (nameOrDescription) {
        const reg = new RegExp(nameOrDescription, 'i');
        queryParameters['$or'] = [
            { 'msgs.msg.name': { $regex: reg } },
            { 'msgs.msg.description': { $regex: reg } },
        ];
    }
    return await this.countDocuments(queryParameters);
};

TxSchema.statics.findServiceProviderCount = async function(serviceName): Promise<number> {
    const queryParameters: any = {
        'msgs.type': TxType.bind_service,
        status: TxStatus.SUCCESS,
        'msgs.msg.ex.service_name': serviceName,
    };
    return await this.countDocuments(queryParameters);
};

// /txs/services/tx
TxSchema.statics.findServiceTx = async function(
    serviceName: string,
    type: string,
    status: number,
    pageNum: number,
    pageSize: number,
): Promise<ITxStruct> {
    const queryParameters: any = {
        'msgs.msg.ex.service_name': serviceName,
    };
    if (type) {
        queryParameters['msgs.type'] = type;
    }else{
        queryParameters['msgs.type'] = {
            '$in':[
                TxType.define_service,
                TxType.bind_service,
                TxType.call_service,
                TxType.respond_service,
                TxType.update_service_binding,
                TxType.disable_service_binding,
                TxType.enable_service_binding,
                TxType.refund_service_deposit,
                TxType.pause_request_context,
                TxType.start_request_context,
                TxType.kill_request_context,
                TxType.update_request_context,
                // TxType.service_set_withdraw_address,
                // TxType.withdraw_earned_fees
            ]
        }
    }
    switch (status) {
        case 0:
            queryParameters.status = TxStatus.FAILED;
            break;
        case 1:
            queryParameters.status = TxStatus.SUCCESS;
            break;
        default:
            break;
    }
    return await this.find(queryParameters, dbRes.service)
        .sort({ 'height': -1 })
        .skip((Number(pageNum) - 1) * Number(pageSize))
        .limit(Number(pageSize));
};

TxSchema.statics.findServiceTxCount = async function(serviceName: string, type: string, status: number): Promise<number> {
    const queryParameters: any = {
        'msgs.msg.ex.service_name': serviceName,
    };
    if (type) {
        queryParameters['msgs.type'] = type;
    }else{
        queryParameters['msgs.type'] = {
            '$in':[
                TxType.define_service,
                TxType.bind_service,
                TxType.call_service,
                TxType.respond_service,
                TxType.update_service_binding,
                TxType.disable_service_binding,
                TxType.enable_service_binding,
                TxType.refund_service_deposit,
                TxType.pause_request_context,
                TxType.start_request_context,
                TxType.kill_request_context,
                TxType.update_request_context,
                // TxType.service_set_withdraw_address,
                // TxType.withdraw_earned_fees
            ]
        }
    }

    switch (status) {
        case 0:
            queryParameters.status = TxStatus.FAILED;
            break;
        case 1:
            queryParameters.status = TxStatus.SUCCESS;
            break;
        default:
            break;
    }
    return await this.countDocuments(queryParameters);
};

TxSchema.statics.findBindTx = async function(serviceName: string, provider: string): Promise<ITxStruct | null> {
    const queryParameters: any = {
        'msgs.msg.ex.service_name': serviceName,
        'msgs.msg.provider': provider,
        'msgs.type': TxType.bind_service,
        status: TxStatus.SUCCESS,
    };
    return await this.findOne(queryParameters);
};

TxSchema.statics.findServiceOwner = async function(serviceName: string): Promise<ITxStruct | null> {
    const queryParameters: any = {
        'msgs.msg.name': serviceName,
        'msgs.type': TxType.define_service,
        status: TxStatus.SUCCESS,
    };
    return await this.findOne(queryParameters);
};

// /txs/services/respond
TxSchema.statics.queryServiceRespondTx = async function(serviceName: string, provider: string, pageNum: number, pageSize: number): Promise<ITxStruct[]> {
    const queryParameters: any = {
        'msgs.type': TxType.respond_service,
    };
    if (serviceName && serviceName.length) {
        queryParameters['msgs.msg.ex.service_name'] = serviceName;
    }
    if (provider && provider.length) {
        queryParameters['msgs.msg.provider'] = provider;
    }
    return await this.find(queryParameters, {...dbRes.common, 'msgs.msg.ex':1})
        .sort({ 'height': -1 })
        .skip((Number(pageNum) - 1) * Number(pageSize))
        .limit(Number(pageSize));
};

TxSchema.statics.findRespondServiceCount = async function(serviceName: string, provider: string): Promise<ITxStruct[]> {
    const queryParameters: any = {
        'msgs.type': TxType.respond_service,
        // status: TxStatus.SUCCESS,
    };
    if (serviceName && serviceName.length) {
        queryParameters['msgs.msg.ex.service_name'] = serviceName;
    }
    if (provider) {
        queryParameters['msgs.msg.provider'] = provider;
    }
    return await this.countDocuments(queryParameters);
};

TxSchema.statics.queryDenomTx = async function(
    pageNum: number,
    pageSize: number,
    denomNameOrId?: string,
): Promise<ITxStruct[]> {
    const params = {
        'msgs.type': TxType.issue_denom,
        status: TxStatus.SUCCESS
    };
    if (denomNameOrId) {
        params['$or'] = [
            { 'msgs.msg.id': denomNameOrId },
            { 'msgs.msg.name': denomNameOrId },
        ];
    }
    return await this.find(params)
        .skip((Number(pageNum) - 1) * Number(pageSize))
        .limit(Number(pageSize))
        .sort({height: -1});
};

TxSchema.statics.queryDenomTxCount = async function (denomNameOrId?: string,):Promise<ITxStruct[]>{
    const params = {
        'msgs.type':TxType.issue_denom,
        status: TxStatus.SUCCESS
    };
    if (denomNameOrId) {
        params['$or'] = [
            { 'msgs.msg.id': denomNameOrId },
            { 'msgs.msg.name': denomNameOrId },
        ];
    }
    return await this.countDocuments(params);
};


TxSchema.statics.queryTxByDenom = async function(
    denom: string,
): Promise<ITxStruct | null> {
    const params = {
        'msgs.type':TxType.issue_denom,
        status: TxStatus.SUCCESS,
        'msgs.msg.id': denom,
    };
    return await this.findOne(params);
};

TxSchema.statics.queryTxByDenomIdAndNftId = async function (
    nftId: string,
    denomId: string
): Promise<ITxStruct | null> {
    const typesList: TxType[] = [
        TxType.transfer_nft,
        TxType.edit_nft,
        TxType.mint_nft
    ];

    const params = {
        status: TxStatus.SUCCESS,
        'msgs.msg.id': nftId,
        'msgs.msg.denom': denomId,
        // $or:[
        //     {
        //         'msgs.type': TxType.transfer_nft,
        //     },
        //     {
        //         'msgs.type': TxType.edit_nft,
        //     },
        //     {
        //         'msgs.type': TxType.mint_nft,
        //     }
        // ],
        'msgs.type': {
            $in: typesList
        }
    };
    return await this.find(params, { time: 1 })
        .sort({ height: -1 })
        .limit(1);
};

// sync Identity task
TxSchema.statics.queryListByCreateAndUpDateIdentity = async function(
  height: number,
  limitSize:number,
): Promise<ITxStruct | null>{
    const typesList: TxType[] = [
        TxType.create_identity,
        TxType.update_identity
    ];
    const params =  {
        height:{
            $gte:height
        },
        // $or:[
        //     {
        //         'msgs.type': TxType.create_identity,
        //         status: TxStatus.SUCCESS,
        //     },
        //     {
        //         'msgs.type': TxType.update_identity,
        //         status: TxStatus.SUCCESS,
        //     }
        // ],
        status: TxStatus.SUCCESS,
        'msgs.type': {
            $in: typesList
        }
    }
    return await this.find(params, dbRes.syncIdentityTask ).limit(limitSize).sort({'height':1})
}

// /txs/identity
TxSchema.statics.queryTxListByIdentity = async function (query:IIdentityTx){
    let result: IListStruct = {};
    const typesList: TxType[] = [
        TxType.create_identity,
        TxType.update_identity
    ];
    const params =  {
        'msgs.msg.id':query.id,
        // $or:[
        //     {
        //         'msgs.type': TxType.create_identity,
        //     },
        //     {
        //         'msgs.type': TxType.update_identity,
        //     }
        // ],
        'msgs.type': {
            $in: typesList
        }
    }
    result.data = await this.find(params, dbRes.txList)
      .sort({ height: -1 })
      .skip((Number(query.pageNum) - 1) * Number(query.pageSize))
      .limit(Number(query.pageSize));
    if (query.useCount && query.useCount == true) {
        result.count = await this.find(params).countDocuments();
    }
    return result;
}

TxSchema.statics.queryDepositsAndSubmitByAddress = async function (address: string) {
    let parameters: any = {
        'msgs.type': {'$in':[TxType.deposit, TxType.submit_proposal]},
        $or: [{'msgs.msg.depositor': address},
            {'msgs.msg.proposer': address}],
        status: 1
    }
    let result: any = {}
    result.data = await this.find(parameters)
    return result
}

// sync tokens task
TxSchema.statics.queryTxBySymbol = async function(
    symbol: string,
    height: number
): Promise<ITxStruct | null>{
      const params =  {
          'msgs.type': {
            $in:[TxType.mint_token,TxType.burn_token]
          },
        status: TxStatus.SUCCESS,
        'msgs.msg.symbol': symbol,
        height: { $gt: height }
      }
    return await this.find(params, {height:1,msgs:1}).sort({'height': 1}).limit(1000)
}

// 	txs/asset
TxSchema.statics.queryTxWithAsset = async function(query: ITxsWithAssetQuery): Promise<IListStruct> {
    let result: IListStruct = {};
    let queryParameters: {'msgs.type':string,'msgs.msg.symbol'?:string} = {
        'msgs.type': query.type,
    };
    if (query.symbol) {
        queryParameters['msgs.msg.symbol'] = query.symbol;
    }
    result.data = await this.find(queryParameters, dbRes.assetList)
        .sort({ height: -1 })
        .skip((Number(query.pageNum) - 1) * Number(query.pageSize))
        .limit(Number(query.pageSize));
    if (query.useCount && query.useCount == true) {
        result.count = await this.find(queryParameters).countDocuments();
    }
    return result;
};
TxSchema.statics.queryNftTxList = async function (lastBlockHeight: number): Promise<ITxStruct[]>  {
    let cond = [
        {
            $sort: {
                height:1,
                tx_index:1,
            },
        },
        {
            $match:{
                status: TxStatus.SUCCESS,
                'msgs.type':{
                    $in:[TxType.mint_nft, TxType.edit_nft, TxType.transfer_nft, TxType.burn_nft]
                },
                height: {$gt: lastBlockHeight, $lte: lastBlockHeight + INCREASE_HEIGHT}
            }
        },
        {
            $unwind:'$msgs',
        },

        {
            $match:{
                'msgs.type':{
                    $in:[TxType.mint_nft, TxType.edit_nft, TxType.transfer_nft, TxType.burn_nft]
                },
            }
        },
        {
            $project:{
                msgs:1,
                height: 1,
                time: 1,
            }
        },
    ];
    return await this.aggregate(cond);
};

TxSchema.statics.queryDenomTxList = async function (lastBlockHeight: number): Promise<ITxStruct[]>  {
    return await this.find(
        {
            'msgs.type': TxType.issue_denom,
            status: TxStatus.SUCCESS,
            height: {
                $gt: lastBlockHeight
            }
        }, { msgs: 1, height: 1, time: 1, tx_hash: 1 }).sort({ height: 1 }).limit(MAX_OPERATE_TX_COUNT);
};
TxSchema.statics.queryMaxNftTxList = async function (): Promise<ITxStruct[]>  {
    const typesList: TxType[] = [
        TxType.mint_nft,
        TxType.edit_nft,
        TxType.transfer_nft,
        TxType.burn_nft
    ];
    return await this.find({
        // $or:[
        //     {'msgs.type':TxType.mint_nft},
        //     {'msgs.type':TxType.edit_nft},
        //     {'msgs.type':TxType.transfer_nft},
        //     {'msgs.type':TxType.burn_nft},
        // ],
        'msgs.type': {
            $in: typesList
        },
        status: TxStatus.SUCCESS
    },{height: 1}).sort({height: -1}).limit(1);
};

TxSchema.statics.queryMaxDenomTxList = async function (): Promise<ITxStruct[]>  {
    return await this.find({ 'msgs.type': TxType.issue_denom,status: TxStatus.SUCCESS },{height: 1}).sort({height: -1}).limit(1);
};

TxSchema.statics.querySubmitProposalById = async function (id:string): Promise<ITxSubmitProposal>  {
    const params =  {
        'msgs.type': TxType.submit_proposal,
        'events.attributes.key':'proposal_id',
        'events.attributes.value': id,
        status:TxStatus.SUCCESS 
      }
    return await this.findOne(params).select({'tx_hash':1,'msgs.msg':1,'_id':0});
};

TxSchema.statics.queryVoteByProposalId = async function (id:number): Promise<ITxVoteProposal[]>  {
    const cond = [
        {
            $match: { 'msgs.type': 'vote', 'msgs.msg.proposal_id': id }
        },
        {
            $group: { _id: "$msgs.msg.voter", msg: { $first: "$msgs.msg" }, count: { $sum: 1 } }
        }
    ]
    return await this.aggregate(cond);
};

TxSchema.statics.queryVoteByProposalIdAll = async function (id:number): Promise<ITxVoteProposalAll[]>  {
    const params =  {
        'msgs.type': TxType.vote,
        'msgs.msg.proposal_id': id,
        status:TxStatus.SUCCESS 
      }
    return await this.find(params,dbRes.voteList).sort({height:1});
};

TxSchema.statics.queryVoteByTxhashs = async function (hash: string[], query?: PagingReqDto): Promise<IListStruct>  {
    const queryParameters =  {
        'tx_hash': {
            $in: hash
        },
        status:TxStatus.SUCCESS 
    }
    let data:IListStruct={};
    if (query) {
        data.data = await this.find(queryParameters, dbRes.voteList)
        .sort({ height: -1 })
        .skip((Number(query.pageNum) - 1) * Number(query.pageSize))
        .limit(Number(query.pageSize));
        if (query.useCount && query.useCount == true) {
            data.count = await this.find(queryParameters).countDocuments();
        }
    } else {
        data = await this.find(queryParameters,dbRes.voteList).sort({ height: -1 })
    }
    return data;
};

TxSchema.statics.queryVoteByTxhashsAndOptoin = async function (hash: string[], option:number): Promise<number>  {
    const queryParameters =  {
        'tx_hash': {
            $in: hash
        },
        status: TxStatus.SUCCESS,
        'msgs.msg.option':option
    }
    return await this.find(queryParameters).countDocuments();
};

TxSchema.statics.queryVoteByTxhashsAndAddress = async function (hash: string[], address:string[],query?: PagingReqDto): Promise<IListStruct>  {
    const queryParameters =  {
        'tx_hash': {
            $in: hash
        },
        status: TxStatus.SUCCESS,
        'msgs.msg.voter': {
            $in: address
        }
    }
    let data:IListStruct={};
    if (query) {
        data.data = await this.find(queryParameters, dbRes.voteList)
        .sort({ height: -1 })
        .skip((Number(query.pageNum) - 1) * Number(query.pageSize))
        .limit(Number(query.pageSize));
        if (query.useCount && query.useCount == true) {
            data.count = await this.find(queryParameters).countDocuments();
        }
    } else {
        data = await this.find(queryParameters,dbRes.voteList).sort({ height: -1 })
    }
    return data;
};


TxSchema.statics.queryDepositorById = async function (id:number,query: PagingReqDto): Promise<IListStruct>  {
    const queryParameters =  {
        'msgs.type': { $in: [TxType.deposit, TxType.submit_proposal] },
        $or: [
            { 'msgs.msg.proposal_id': Number(id) },
            { 'events.attributes.key': 'proposal_id', 'events.attributes.value': String(id) }
        ],
        status:TxStatus.SUCCESS 
    }
    let result: IListStruct = {};
    result.data = await this.find(queryParameters, dbRes.depositorList)
        .sort({ height: -1 })
        .skip((Number(query.pageNum) - 1) * Number(query.pageSize))
        .limit(Number(query.pageSize));
    if (query.useCount && query.useCount == true) {
        result.count = await this.find(queryParameters).countDocuments();
    }
    return result;
};

TxSchema.statics.queryVoteByAddr = async function (address:string): Promise<ITxVoteALL[]>  {
    const queryParameters =  {
        'msgs.type': TxType.vote,
        'msgs.msg.voter': address,
        status:TxStatus.SUCCESS 
    }
    return await this.find(queryParameters,dbRes.voteList).sort({ height: 1 });
};

TxSchema.statics.queryDepositsByAddress = async function(address:string,query:PagingReqDto): Promise<IListStruct> {
    let result: IListStruct = {};
    let queryParameters = {
        'msgs.msg.depositor': address,
        status:TxStatus.SUCCESS
    };
    result.data = await this.find(queryParameters, dbRes.depositList)
        .sort({ height: -1 })
        .skip((Number(query.pageNum) - 1) * Number(query.pageSize))
        .limit(Number(query.pageSize));
    if (query.useCount && query.useCount == true) {
        result.count = await this.find(queryParameters).countDocuments();
    }
    return result;
};

TxSchema.statics.queryAccountTxList = async function (lastSyncBlockHeight:number): Promise<ITxStruct[]>  {
    const typesList: TxType[] = [
        TxType.send,
        TxType.multisend,
        TxType.set_withdraw_address,
        TxType.add_super
    ];
    return await this.find({
        'msgs.type': {
            $in: typesList
        },
        height: {$gt: lastSyncBlockHeight, $lte: lastSyncBlockHeight + INCREASE_HEIGHT}
    },{addrs: 1,height: 1,_id:0}).sort({height:1});
};

TxSchema.statics.queryTxMaxHeight = async function (): Promise<ITxStruct[]>  {
    return await this.find({},{height: 1}).sort({height: -1}).limit(1);
};
