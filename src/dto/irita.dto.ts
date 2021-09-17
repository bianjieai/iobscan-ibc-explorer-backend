import {IsString, IsInt, Length, Min, Max, IsOptional, Equals, MinLength, ArrayNotEmpty} from 'class-validator';
import {ApiProperty, ApiPropertyOptional} from '@nestjs/swagger';
import {BaseReqDto, BaseResDto, PagingReqDto} from './base.dto';
import {ApiError} from '../api/ApiResult';
import {ErrorCodes} from '../api/ResultCodes';
import {IBindTx} from '../types/tx.interface';
import { currentChain } from '../constant';

/************************   request dto   ***************************/
export class TokensReqDto extends BaseReqDto {
  @ApiProperty({ required: false })
  denom:string;
  
  @ApiProperty({ required: false })
  chain?: string;

  @ApiProperty({ required: false })
  key: string;

  static validate(value: any) {
    super.validate(value);
    if (value.chain && value.chain !== currentChain.cosmos && value.chain !== currentChain.iris && value.chain !== currentChain.binance) {
      throw new ApiError(ErrorCodes.InvalidParameter, 'chain must be one of iris, cosmos and binance');
    }
  }
}

/************************   response dto   ***************************/
//txs response dto
export class NetworkResDto extends BaseResDto{
    network_id: string;
    network_name: string;
    uri: string;
    is_main: boolean;

    constructor(value) {
        super();
        const { network_id, network_name, uri, is_main } = value;
        this.network_id = network_id;
        this.network_name = network_name;
        this.uri = uri;
        this.is_main = is_main || false;
    }

    static bundleData(value: any): NetworkResDto[] {
        let data: NetworkResDto[] = [];
        data = value.map((v: any) => {
            return new NetworkResDto(v);
        });
        return data;
    }
}

export class TokensResDto extends BaseResDto {
    symbol:string;
    denom:string;
    scale:number;
    is_main_token:boolean;
    initial_supply: string;
    max_supply: string;
    mintable: boolean;
    owner: string;
    name: string;
    total_supply: string;
    update_block_height: number;
    src_protocol: string;
    chain:string;
    constructor(value) {
        super();
        this.symbol = value.symbol;
        this.denom = value.denom;
        this.scale = value.scale;
        this.is_main_token = value.is_main_token;
        this.initial_supply = value.initial_supply;
        this.max_supply = value.max_supply;
        this.mintable = value.mintable;
        this.owner = value.owner;
        this.name = value.name;
        this.total_supply = value.total_supply;
        this.update_block_height = value.update_block_height;
        this.src_protocol = value.src_protocol;
        this.chain = value.chain;
    }
    static bundleData(value: any): TokensResDto[] {
        let data: TokensResDto[] = [];
        data = value.map((v: any) => {
            return new TokensResDto(v);
        });
        return data;
    }

}

export class StatusResDto {
    is_follow: boolean;
    constructor(value) {
        this.is_follow = value;
    }
}
