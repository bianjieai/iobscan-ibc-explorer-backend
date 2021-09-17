import { BaseReqDto, PagingReqDto } from './base.dto';
import { ApiError } from '../api/ApiResult';
import { ErrorCodes } from '../api/ResultCodes';
import { ApiPropertyOptional } from '@nestjs/swagger';

export class ParametersListReqDto {
    @ApiPropertyOptional()
    module?: string;

    @ApiPropertyOptional({description:'example: key = key1,key2'})
    key?: string;
}

export class ParametersListResDto {
    module: string;
    key: string;
    cur_value: string;

    constructor(value) {
        this.module = value.module || '';
        this.key = value.key || '';
        this.cur_value = value.cur_value || '';
    }
}