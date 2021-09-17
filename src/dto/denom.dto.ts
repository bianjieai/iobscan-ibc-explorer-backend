import { IsString, IsNotEmpty } from 'class-validator';
import { PagingReqDto } from './base.dto';
import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';


export class DenomListReqDto extends PagingReqDto{
    @ApiPropertyOptional()
    denomNameOrId?: string;

    @ApiPropertyOptional()
    needAll?: boolean;


}

export class DenomResDto {
    denomName: string;
    denomId: string;
    hash: string;
    nftCount: number;
    sender: string;
    time: number;

    constructor(
        denomName: string,
        denomId: string,
        hash: string,
        nftCount: number,
        sender: string,
        time: number,
    ){
        this.denomName = denomName;
        this.denomId = denomId;
        this.hash = hash;
        this.nftCount = nftCount;
        this.sender = sender;
        this.time = time;
    }
}

export class DenomListResDto extends DenomResDto{
    constructor(
        denomName: string,
        denomId: string,
        hash: string,
        nftCount: number,
        sender: string,
        time: number,
    ){
        super(denomName, denomId, hash, nftCount, sender,time);
    }
}