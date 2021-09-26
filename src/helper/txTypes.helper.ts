import { TxType } from '../constant';
// todo 无效变量待删除

export function stakingTypes():string[]{
    return [
        TxType.delegate,
        TxType.begin_redelegate,
        TxType.begin_unbonding,
        TxType.set_withdraw_address,
        TxType.withdraw_delegator_reward,
        TxType.fund_community_pool
    ];
}

export function coinswapTypes():string[]{
    return [
        TxType.add_liquidity,
        TxType.remove_liquidity,
        TxType.swap_order
    ];
}


export function serviceTypes():string[]{
    return [
        TxType.define_service,
        TxType.bind_service,
        TxType.call_service,
        TxType.respond_service,
        TxType.update_service_binding,
        TxType.disable_service_binding,
        TxType.enable_service_binding,
        TxType.refund_service_deposit,
        TxType.pause_request_context,
        TxType.start_request_context,
        TxType.kill_request_context,
        TxType.update_request_context
    ];
}

export function declarationTypes():string[]{
    return [
        TxType.create_validator,
        TxType.edit_validator,
        TxType.unjail,
        TxType.withdraw_validator_commission
    ];
}

export function govTypes():string[]{
    return [
        TxType.deposit,
        TxType.vote,
        TxType.submit_proposal
    ];
}