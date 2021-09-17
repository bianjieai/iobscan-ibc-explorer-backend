import { IsString, IsNotEmpty } from 'class-validator';
import { PagingReqDto } from './base.dto';
import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
import { IDenomStruct } from '../types/schemaTypes/denom.interface';

export class NftListReqDto extends PagingReqDto {
    @ApiPropertyOptional()
    denomId?: string;

    @ApiPropertyOptional()
    nftId?: string;

    @ApiPropertyOptional()
    owner?: string;

}

export class NftDetailReqDto {
    @IsString()
    @ApiProperty()
    @IsNotEmpty({ message: 'denom id is necessary' })
    denomId: string;

    @IsString()
    @ApiProperty()
    @IsNotEmpty({ message: 'nft id is necessary' })
    nftId: string;
}


export class NftResDto {
    denom_id: string;
    nft_id: string;
    owner: string;
    tokenUri: string;
    tokenData: string;
    denomDetail: IDenomStruct;
    denom_name: string;
    nft_name: string;

    constructor(
        denom_id: string,
        nft_id: string,
        owner: string,
        tokenUri: string,
        tokenData: string,
        denomDetail: IDenomStruct,
        denom_name: string,
        nft_name: string,
    ) {
        this.denom_id = denom_id;
        this.nft_id = nft_id;
        this.owner = owner;
        this.tokenUri = tokenUri;
        this.tokenData = tokenData;
        this.denomDetail = denomDetail;
        this.denom_name = denom_name;
        this.nft_name = nft_name;
    }
}

export class NftListResDto extends NftResDto {
    last_block_time: number;

    constructor(
        denom_id: string,
        nft_id: string,
        owner: string,
        tokenUri: string,
        tokenData: string,
        denomDetail: IDenomStruct,
        denom_name: string,
        nft_name: string,
        last_block_time: number
    ) {
        super(denom_id, nft_id, owner, tokenUri, tokenData, denomDetail, denom_name, nft_name);
        this.last_block_time = last_block_time;
    }

}

export class NftDetailResDto extends NftResDto {

    constructor(
        denom_id: string,
        nft_id: string,
        owner: string,
        tokenUri: string,
        tokenData: string,
        denomDetail: IDenomStruct,
        denom_name: string,
        nft_name: string,
    ) {
        super(denom_id, nft_id, owner, tokenUri, tokenData, denomDetail, denom_name, nft_name);
    }
}

