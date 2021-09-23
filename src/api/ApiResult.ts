import { ErrorCodes } from './ResultCodes';
import { IResultBase, IListStructBase } from '../types';
import {
    HttpException,
    HttpStatus,
} from '@nestjs/common';
export class ListStruct<T> implements IListStructBase<T> {
    data: T;
    pageNum: number;
    pageSize: number;
    count?: number;
    constructor(data: T,pageNum: number, pageSize: number, count?: number) {
        this.data = data;
        this.pageNum = pageNum;
        this.pageSize = pageSize;
        if (count || count === 0) this.count = count;
    }
}

export class Result<T> implements IResultBase {
    public code: number = ErrorCodes.success;
    public data: T;

    constructor(data: T, code: number = ErrorCodes.success) {
        this.data = data;
        this.code = code;
    }
}

export class ApiError extends HttpException{
    constructor(code: number, message?: string){
        super({
            code,
            message,
        }, HttpStatus.OK)
    }
}




