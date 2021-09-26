import { PagingReqDto } from './base.dto';
import { ApiPropertyOptional } from '@nestjs/swagger';

export class IbcTxListReqDto extends PagingReqDto {
    @ApiPropertyOptional()
    page_num?: number;
    @ApiPropertyOptional()
    page_size?: number
}

export class IbcTxResDto {
    record_id: string;
    sc_addr: string;
    dc_addr: string;
    sc_chain_id: string;
    dc_chain_id: string;
    status: number;
    sc_tx_info: object;
    dc_tx_info?: object;
    base_denom: string;
    update_at: string;
}