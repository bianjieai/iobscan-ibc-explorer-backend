/* eslint-disable @typescript-eslint/camelcase */
import {BaseReqDto, PagingReqDto} from './base.dto';
import {ApiPropertyOptional} from '@nestjs/swagger';
import {BaseResDto} from './base.dto';

export class IbcTxListReqDto extends PagingReqDto {
    @ApiPropertyOptional()
    date_range?: string;
    @ApiPropertyOptional()
    status?: string;
    @ApiPropertyOptional()
    chain_id?: string;
    @ApiPropertyOptional()
    symbol?: string;
    @ApiPropertyOptional()
    denom?: string;
    @ApiPropertyOptional()
    page_num?: number;
    @ApiPropertyOptional()
    page_size?: number;
    @ApiPropertyOptional()
    start_time?: boolean

    static convert(value: any): any {
        return super.convert(value);
    }
}

export class TxWithHashReqDto extends BaseReqDto {
    @ApiPropertyOptional()
    hash: string;
}

export class IbcTxDetailsResDto extends BaseResDto{
    base_denom: string;
    sc_signers:string[];
    dc_signers:string[];
    dc_addr: string;
    dc_chain_id: string;
    dc_channel: string;
    dc_port: string;
    denoms: object;
    sc_addr: string;
    sc_chain_id: string;
    sc_channel: string;
    sc_port: string;
    sc_tx_info: object;
    sequence: string;
    status: string;
    tx_time: string;
    dc_tx_info:object;
    dc_connect:string;
    sc_connect:string;
    constructor(value: any) {
        super();
        const {
            base_denom,
            dc_addr,
            dc_chain_id,
            dc_channel,
            dc_port,
            denoms,
            sc_addr,
            sc_chain_id,
            sc_channel,
            sc_port,
            sc_tx_info,
            sequence,
            status,
            tx_time,
            dc_tx_info,
            sc_signers,
            dc_signers,
            dc_connect,
            sc_connect,
        } = value;
        this.base_denom = base_denom || '';
        this.dc_addr = dc_addr || '';
        this.dc_chain_id = dc_chain_id || '';
        this.dc_channel = dc_channel;
        this.dc_port = dc_port || '';
        this.denoms = denoms || '';
        this.sc_addr = sc_addr || '';
        this.sc_chain_id = sc_chain_id || '';
        this.sc_channel = sc_channel || '';
        this.sc_port = sc_port || [];
        this.sc_tx_info = sc_tx_info || {};
        this.sequence = sequence || '';
        this.status = status || '';
        this.tx_time = tx_time || '';
        this.dc_tx_info = dc_tx_info || {};
        this.sc_signers = sc_signers || [];
        this.dc_signers = dc_signers || [];
        this.dc_connect = dc_connect || '';
        this.sc_connect = sc_connect || '';
    }
    static bundleData(value: any): IbcTxDetailsResDto[] {
        const datas: IbcTxDetailsResDto[] = value.map((item: any) => {
            return new IbcTxDetailsResDto(item);
        });
        return datas;
    }
}

export class IbcTxResDto extends BaseResDto {
    record_id: string;
    sc_addr: string;
    dc_addr: string;
    status: number;
    sc_chain_id: string;
    dc_chain_id: string;
    dc_channel: string;
    sc_channel: string;
    sequence: string;
    sc_tx_info: object;
    dc_tx_info?: object;
    base_denom: string;
    denoms: string[];
    create_at: string;
    tx_time: string;

    constructor(value: any) {
        super();
        const {
            record_id,
            sc_addr,
            dc_addr,
            dc_chain_id,
            sc_chain_id,
            status,
            sequence,
            sc_tx_info,
            dc_tx_info,
            base_denom,
            denoms,
            create_at,
            tx_time,
        } = value;
        this.record_id = record_id || '';
        this.sc_addr = sc_addr || '';
        this.dc_addr = dc_addr || '';
        this.status = status;
        this.sequence = sequence || 0;
        this.sc_chain_id = sc_chain_id || '';
        this.dc_chain_id = dc_chain_id || '';
        this.sc_tx_info = sc_tx_info || {};
        this.dc_tx_info = dc_tx_info || {};
        this.base_denom = base_denom || '';
        this.denoms = denoms || [];
        this.create_at = create_at || '';
        this.tx_time = tx_time || '';
    }

    static bundleData(value: any): IbcTxResDto[] {
        const datas: IbcTxResDto[] = value.map((item: any) => {
            return new IbcTxResDto(item);
        });
        return datas;
    }
}
