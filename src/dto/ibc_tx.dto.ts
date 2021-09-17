import { PagingReqDto } from './base.dto';
import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';

export class IbcTxListReqDto extends PagingReqDto {
    // @ApiPropertyOptional()
    // txId?: string;
}

export class IbcTxListResDto {}