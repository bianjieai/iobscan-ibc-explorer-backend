export interface IAccountStruct {
    address?: string,
    account_total?: number,
    total?: Object,
    balance?: Array<object>,
    delegation?: Object,
    unbonding_delegation?: Object,
    rewards?: Object,
    create_time?: number,
    update_time?: number,
    handled_block_height?: number,
}

export interface ITokenTotal {
    _id: any,
    account_totals: number
}

export interface ITokenTotal {
    _id: any,
    account_totals: number
}

export interface CommissionReward {
    val_commission?: {
        commission?: Object[]
    }
}
