import { PagingReqDto } from './base.dto';
import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';

export class IbcChainListReqDto extends PagingReqDto {
    @ApiPropertyOptional()
    chain_name?: string;
}

export class IbcChainListResDto {
    chain_id: string;
    icon: string;
    chain_name: string;
    lcd: string;
    ibc_info: object;
}