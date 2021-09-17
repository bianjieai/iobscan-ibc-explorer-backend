
import { IQueryBase } from "."
export interface IValidatorsQueryParams extends IQueryBase {
    jailed: Boolean|string
}

export interface IValidatorListStruct {
    data?: any[],
    count?: number
}
