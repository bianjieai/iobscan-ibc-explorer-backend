import {IsOptional} from 'class-validator';
import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
import {ApiError} from '../api/ApiResult';
import {ErrorCodes} from '../api/ResultCodes';
import { DefaultPaging } from '../constant';

export class Coin {
    denom: string;
    amount: string;
    constructor(value) {
        let { denom, amount } = value;
        this.denom = denom || '';
        this.amount = amount || '';
    }

    static bundleData(value: any = []): Coin[] {
        let data: Coin[] = [];
        data = value.map((v: any) => {
            return new Coin(v);
        });
        return data;
    }
}

