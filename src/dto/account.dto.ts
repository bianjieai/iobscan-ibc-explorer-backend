import { BaseReqDto, PagingReqDto, BaseResDto } from './base.dto';
import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
import { Coin } from './common.res.dto';

/***************Req***********************/
export class genesisAccountsReqDto {
    @ApiProperty({description: 'type: IRIS/COSMOS'})
    chain?: string;
}

/***************Res*************************/

export class accountsListData {
    rank: string;
    address: string;
    balance: Coin[];
    percent: number;

    constructor(account) {
        this.rank = account.rank || '';
        this.address = account.address || '';
        this.balance = account.balance || [];
        this.percent = account.percent || '';
    }
    static bundleData(value: any): accountsListResDto[] {
        let data: accountsListResDto[] = [];
        data = value.map((v: any) => {
            return new accountsListResDto(v);
        });
        return data;
    }
}

export class accountsListResDto {
    updated_time: number;
    data: accountsListData[]

    constructor(value) {
        this.updated_time = value.updated_time || 0;
        this.data = value.data || [];
    }
}

export class tokenStatsResDto {
    total_supply_tokens: Coin;
    circulation_tokens: Coin;
    bonded_tokens: Coin;
    community_tax: Coin;

    constructor(value) {
        this.total_supply_tokens = value.total_supply_tokens || {};
        this.circulation_tokens = value.circulation_tokens || {};
        this.bonded_tokens = value.bonded_tokens || {};
        this.community_tax = value.community_tax || {};
    }
}

export class accountTotal {
    total_amount: object;
    percent: number;
    constructor(value) {
        this.total_amount = value.total_amount || { denom: '',amount: 0 };
        this.percent = value.percent || 0;
    }
}

export class accountTotalResDto {
    "1-5": accountTotal;
    "6-10": accountTotal;
    "11-50": accountTotal;
    "51-100": accountTotal;
    "101-500": accountTotal;
    "501-1000": accountTotal;
    "1001-": accountTotal;

    constructor(value) {
        this["1-5"] = new accountTotal(value["1-5"]);
        this["6-10"] = new accountTotal(value["6-10"]);
        this["11-50"] = new accountTotal(value["11-50"]);
        this["51-100"] = new accountTotal(value["51-100"]);
        this["101-500"] = new accountTotal(value["101-500"]);
        this["501-1000"] = new accountTotal(value["501-1000"]);
        this["1001-"] = new accountTotal(value["1001-"]);
    }
}

