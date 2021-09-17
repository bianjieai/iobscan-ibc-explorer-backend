import { IsString, IsInt, Length, Min, Max, IsOptional, Equals, MinLength, ArrayNotEmpty, validate } from 'class-validator';
import {ApiProperty, ApiPropertyOptional} from '@nestjs/swagger';
import {BaseReqDto, BaseResDto, PagingReqDto} from './base.dto';
import {ApiError} from '../api/ApiResult';
import {ErrorCodes} from '../api/ResultCodes';
import {IBindTx,ExternalIBindTx} from '../types/tx.interface';

/************************   request dto   ***************************/
//txs request dto
export class TxListReqDto extends PagingReqDto {
    @ApiPropertyOptional()
    type?: string;

    @ApiPropertyOptional({description: '1:success  2:fail'})
    status?: string;

    @ApiPropertyOptional()
    address?: string;

    @ApiPropertyOptional()
    beginTime?: string;

    @ApiPropertyOptional()
    endTime?: string;

    static validate(value: any) {
        super.validate(value);
        if (value.status && value.status !== '1' && value.status !== '2') {
            throw new ApiError(ErrorCodes.InvalidParameter, 'status must be 1 or 2');
        }
    }
}

// txs/e
export class eTxListReqDto extends PagingReqDto {
    @ApiPropertyOptional()
    types?: string;

    @ApiPropertyOptional({description: 'Greater than block height'})
    height?: number;

    @ApiPropertyOptional()
    status?: number;

    @ApiPropertyOptional()
    address?: string;

    @ApiPropertyOptional({description:'true/false'})
    include_event_addr?: boolean;

    static convert(value: any): any {
        super.convert(value);
        if(!value.include_event_addr){
            value.include_event_addr = false;
        }else {
            if(value.include_event_addr === 'true'){
                value.include_event_addr = true;
            }else {
                value.include_event_addr = false;
            }
        }
        return value;
    }
}

//txs/blocks request dto
export class TxListWithHeightReqDto extends PagingReqDto {
    @ApiPropertyOptional()
    height?: string;
}

//txs/addresses request dto
export class TxListWithAddressReqDto extends PagingReqDto {
    @ApiPropertyOptional()
    address?: string;

    @ApiPropertyOptional()
    type?: string;

    @ApiPropertyOptional({description: '1:success  2:fail'})
    status?: string;

    static validate(value: any) {
        super.validate(value);
        if (value.status && value.status !== '1' && value.status !== '2') {
            throw new ApiError(ErrorCodes.InvalidParameter, 'status must be 1 or 2');
        }
    }
}

// txs/relevance
export class TxListWithContextIdReqDto extends PagingReqDto {
    @ApiPropertyOptional()
    contextId?: string;

    @ApiPropertyOptional()
    type?: string;

    @ApiPropertyOptional({description: '1:success  2:fail'})
    status?: string;

    static validate(value: any) {
        super.validate(value);
        if (value.status && value.status !== '1' && value.status !== '2') {
            throw new ApiError(ErrorCodes.InvalidParameter, 'status must be 1 or 2');
        }
    }
}

//txs/nfts request dto
export class TxListWithNftReqDto extends PagingReqDto {
    @ApiPropertyOptional()
    denom?: string;

    @ApiPropertyOptional()
    tokenId?: string;
}

//txs/services request dto
export class TxListWithServicesNameReqDto extends PagingReqDto {
    @ApiPropertyOptional()
    serviceName?: string;
}

//txs/services/detail/{serviceName} request dto
export class ServicesDetailReqDto extends BaseReqDto {
    @ApiProperty()
    serviceName: string;
}

//txs/service/call-service
export class TxListWithCallServiceReqDto extends PagingReqDto {
    @ApiProperty()
    @MinLength(1, {message: "consumerAddr is too short"})
    consumerAddr: string;
}

//txs/e/services/respond-service
export class ExternalQueryRespondServiceReqDto {
    @ApiProperty()
    serviceName: string;

    @ApiProperty()
    @MinLength(1, {message: "providerAddr is too short"})
    providerAddr: string;
}

//txs/service/respond-service
export class TxListWithRespondServiceReqDto extends PagingReqDto {
    @ApiProperty()
    @MinLength(1, {message: "providerAddr is too short"})
    providerAddr: string;
}

//Post txs/types request dto
export class PostTxTypesReqDto extends BaseReqDto {
    @ApiProperty()
    @ArrayNotEmpty()
    typeNames: Array<string>;
}

//put txs/types request dto
export class PutTxTypesReqDto extends BaseReqDto {
    @ApiProperty()
    @MinLength(1, {message: 'typeName is too short'})
    typeName: string;

    @ApiProperty()
    @MinLength(1, {message: 'newTypeName is too short'})
    newTypeName: string;
}

//Delete txs/types request dto
export class DeleteTxTypesReqDto extends BaseReqDto {
    @ApiProperty()
    @MinLength(1, {message: 'typeName is too short'})
    typeName: string;
}

//txs/{hash} request dto
export class TxWithHashReqDto extends BaseReqDto {
    @ApiProperty()
    hash: string;
}

export class ServiceListReqDto extends PagingReqDto {
    @ApiPropertyOptional()
    nameOrDescription?: string;
}


export class ServiceProvidersReqDto extends PagingReqDto {
    @ApiProperty()
    serviceName: string;
}


export class ServiceTxReqDto extends PagingReqDto {
    @ApiProperty()
    serviceName: string;

    @ApiPropertyOptional()
    type?: string;

    @ApiPropertyOptional()
    status?: string;

    static convert(value: any): any {
        value.status = Number(value.status);
        return value;
    }
}

export class ServiceBindInfoReqDto {
    @ApiProperty()
    serviceName: string;

    @ApiProperty()
    provider: string;
}


export class ServiceTxResDto {
    hash: string;
    type: string;
    height: number;
    time: number;
    status: number;
    msgs: any;
    events: any;
    fee: any;

    constructor(
        hash: string,
        type: string,
        height: number,
        time: number,
        status: number,
        msgs: any,
        events: any,
        fee: any
    ) {
        this.hash = hash;
        this.type = type;
        this.height = height;
        this.time = time;
        this.status = status;
        this.msgs = msgs;
        this.events = events;
        this.fee = fee;
    }
}

export class ServiceRespondReqDto extends PagingReqDto {
    @ApiPropertyOptional()
    serviceName: string;

    @ApiProperty()
    provider?: string;
}

export class IdentityTxReqDto extends PagingReqDto {
    @ApiProperty()
    id: string
}

export class IDepositsAddress {
    address: string
}

export class ServiceRespondResDto {
    respondHash: string;
    type: string;
    height: number;
    time: number;
    consumer: string;
    requestHash: string;
    requestContextId: string;
    serviceName: string;
    respondStatus: number;

    constructor(
        respondHash: string,
        type: string,
        height: number,
        time: number,
        consumer: string,
        requestHash: string,
        requestContextId: string,
        serviceName: string,
        respondStatus: number,
    ) {
        this.respondHash = respondHash;
        this.type = type;
        this.height = height;
        this.time = time;
        this.consumer = consumer;
        this.requestHash = requestHash;
        this.requestContextId = requestContextId;
        this.serviceName = serviceName;
        this.respondStatus = respondStatus;
    }
}

/************************   response dto   ***************************/
//txs response dto
export class TxResDto extends BaseResDto {
    time: string;
    height: string;
    tx_hash: string;
    memo: string;
    status: number;
    log: string;
    complex_msg: boolean;
    type: string;
    from: string;
    to: string;
    coins: Array<any>;
    signer: string;
    events: Array<any>;
    msgs: Array<any>;
    signers: Array<any>;
    fee: object;
    monikers: any[];
    addrs?: any[];
    ex?: object;
    proposal_link?: boolean;
    events_new?:any[];

    constructor(txData) {
        super();
        this.time = txData.time;
        this.height = txData.height;
        this.tx_hash = txData.tx_hash;
        this.memo = txData.memo;
        this.status = txData.status;
        this.log = txData.log;
        this.complex_msg = txData.complex_msg;
        this.type = txData.type;
        this.from = txData.from;
        this.to = txData.to;
        this.coins = txData.coins;
        this.signer = txData.signer;
        this.events = txData.events;
        this.events_new = txData.events_new;
        this.msgs = txData.msgs;
        this.signers = txData.signers;
        this.fee = txData.fee;
        this.monikers = txData.monikers || [];
        if (txData.ex) this.ex = txData.ex;
        if (txData.proposal_link) this.proposal_link = true;
        if (txData.addrs) this.addrs = txData.addrs;
    }

    static bundleData(value: any): TxResDto[] {
        let data: TxResDto[] = [];
        data = value.map((v: any) => {
            return new TxResDto(v);
        });
        return data;
    }
}

//txs/service/call-service
export class callServiceResDto extends TxResDto {
    respond: TxResDto[];

    constructor(txData) {
        super(txData);
        if (txData.respond && txData.respond.length) {
            this.respond = (txData.respond || []).map((item) => {
                return new TxResDto(item);
            });
        }
    }

    static bundleData(value: any): callServiceResDto[] {
        let data: callServiceResDto[] = [];
        data = value.map((v: any) => {
            return new callServiceResDto(v);
        });
        return data;
    }
}

//e/services/respond-service
export class ExternalQueryRespondServiceResDto {
    respondTimes: number

    constructor(value) {
        this.respondTimes = value || 0;
    }
}

//txs/service/respond-service
export class RespondServiceResDto extends TxResDto {
    respond_times: number;
    unbinding_time: number;

    constructor(txData) {
        super(txData);
        this.respond_times = txData.respond_times;
        this.unbinding_time = txData.unbinding_time;
    }

    static bundleData(value: any): RespondServiceResDto[] {
        let data: RespondServiceResDto[] = [];
        data = value.map((v: any) => {
            return new RespondServiceResDto(v);
        });
        return data;
    }
}

//txs/types response dto
export class TxTypeResDto extends BaseResDto {
    typeName: string;

    constructor(typeData) {
        super();
        this.typeName = typeData.type_name;
    }

    static bundleData(value: any): TxTypeResDto[] {
        let data: TxTypeResDto[] = [];
        data = value.map((v: any) => {
            return new TxTypeResDto(v);
        });
        return data;
    }
}

export class ServiceResDto {
    serviceName: string;
    description: string;
    bindList: IBindTx[];

    constructor(serviceName: string, description: string, bindList: IBindTx[]) {
        this.serviceName = serviceName;
        this.description = description;
        this.bindList = bindList;
    }
}

export class ExternalServiceResDto {
    serviceName: string;
    description: string;
    bindList: ExternalIBindTx[];

    constructor(serviceName: string,description: string,  bindList: ExternalIBindTx[]) {
        this.serviceName = serviceName;
        this.description = description;
        this.bindList = bindList;
    }
}

export class ServiceProvidersResDto implements IBindTx {
    provider: string;
    respondTimes?: number;
    bindTime: string;

    constructor(provider: string, respondTimes: number, bindTime: string) {
        this.provider = provider;
        this.respondTimes = respondTimes;
        this.bindTime = bindTime;
    }
}

export class ServiceBindInfoResDto {
    hash: string;
    owner: string;
    time: number;

    constructor(
        hash: string,
        owner: string,
        time: number,
    ) {
        this.hash = hash;
        this.time = time;
        this.owner = owner;
    }
}

//txs/blocks request dto
export class TxListWithAssetReqDto extends PagingReqDto {
    @ApiProperty()
    type: string;

    @ApiPropertyOptional()
    symbol?: string;
}
