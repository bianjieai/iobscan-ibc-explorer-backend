export interface IQueryBase {
    pageNum?: string;
    pageSize?: string;
    useCount?: boolean | string;
}

export interface IListStructBase<T> {
    data?: T;
    pageNum?: number;
    pageSize?: number;
    count?: number;
}

export interface IListStruct {
    data?: any[],
    count?: number
}

export interface IResultBase {
    code: number;
    data?: any;
    message?: string;
}

export interface IRandomKey {
    key: string;
    step: number;
}