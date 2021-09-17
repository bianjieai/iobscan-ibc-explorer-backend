import {BaseResDto, PagingReqDto} from './base.dto';
import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
import {ArrayNotEmpty} from 'class-validator';
import { Coin } from './common.res.dto';

/***************Req***********************/


// export class ValCommissionRewReqDto {
//     @ApiProperty()
//     @ApiPropertyOptional()
//     address?: string
// }

export class CommissionInfoReqDto extends PagingReqDto {

}

export class ValidatorDelegationsReqDto {
    @ApiProperty()
    address: string
}

export class ValidatorDelegationsQueryReqDto extends PagingReqDto {
}

export class ValidatorUnBondingDelegationsReqDto {
    @ApiProperty()
    address: string
}

export class ValidatorUnBondingDelegationsQueryReqDto extends PagingReqDto {
}

export class allValidatorReqDto extends PagingReqDto {
    @ApiProperty({description:'jailed/candidate/active'})
    status: string
}

export class ValidatorDetailAddrReqDto {
    @ApiProperty()
    address: string
}

export class AccountAddrReqDto {
    @ApiProperty()
    address: string
}

export class DelegatorsDelegationsReqDto extends PagingReqDto {
}

export class DelegatorsDelegationsParamReqDto {
    @ApiProperty()
    delegatorAddr: string
}

export class DelegatorsUndelegationsReqDto extends PagingReqDto {
}

export class DelegatorsUndelegationsParamReqDto {
    @ApiProperty()
    delegatorAddr: string
}

//Post staking/blacks request dto
export class PostBlacksReqDto {
    @ApiProperty({description:`{"blacks": [{"iva_addr": "iva176dd0tgn38grpc8hpxfmwl6sl8jfmkneak3emy","moniker_m":"test1","is_block":true}]}`})
    @ArrayNotEmpty()
    blacks: String;
}

/***************Res*************************/

export class stakingValidatorResDto extends BaseResDto {
    operator_address: string;
    consensus_pubkey: string;
    jailed: boolean;
    status: number;
    tokens: string;
    delegator_shares: string;
    description: object;
    bond_height: string;
    unbonding_height: string;
    unbonding_time: string;
    commission: object;
    uptime: number;
    self_bond: object;
    delegator_num: number;
    proposer_addr: string;
    voting_power: number;
    icons: string;
    voting_rate: number;

    constructor(validator) {
        super();
        this.operator_address = validator.operator_address || ''
        this.consensus_pubkey = validator.consensus_pubkey || ''
        this.jailed = validator.jailed || false
        this.status = validator.status || 0
        this.tokens = validator.tokens || ''
        this.delegator_shares = validator.delegator_shares || ''
        this.description = validator.description || {}
        this.bond_height = validator.start_height || ''
        this.unbonding_height = validator.unbonding_height || ''
        this.commission = validator.commission || ''
        this.uptime = validator.uptime || 0
        this.self_bond = validator.self_bond || {}
        this.delegator_num = validator.delegator_num || 0
        this.proposer_addr = validator.proposer_addr || ''
        this.voting_power = validator.voting_power || 0
        this.icons = validator.icon || ''
        this.voting_rate = validator.voting_rate || 0
    }

    static bundleData(value: any): stakingValidatorResDto[] {
        let data: stakingValidatorResDto[] = [];
        data = value.map((v: any) => {
            return new stakingValidatorResDto(v);
        });
        return data;
    }
}


// export class ValCommissionRewResDto extends BaseResDto {
//     operator_address: string;
//     self_bond_rewards: [];
//     val_commission: object;

//     constructor(commissionRewards) {
//         super();
//         this.operator_address = commissionRewards.operator_address || ''
//         this.self_bond_rewards = commissionRewards.self_bond_rewards || {}
//         this.val_commission = commissionRewards.val_commission || {}
//     }
// }

export class CommissionInfoResDto extends BaseResDto {
    commission_rate: string;
    bonded_tokens: string;
    moniker: string;
    operator_address: string;

    constructor(CommissionInfo) {
        super();
        this.operator_address = CommissionInfo.operator_address || ''
        this.commission_rate = CommissionInfo.commission.commission_rates.rate || ''
        this.moniker = CommissionInfo.description.moniker || ''
        this.bonded_tokens = CommissionInfo.tokens || ''
    }

    static bundleData(value: any): CommissionInfoResDto[] {
        let data: CommissionInfoResDto[] = [];
        data = value.map((v: any) => {
            return new CommissionInfoResDto(v);
        });
        return data;
    }
}

export class ValidatorDelegationsResDto extends BaseResDto {
    address: string;
    moniker: string;
    amount: object;
    self_shares: string;
    total_shares: number;

    constructor(delegations) {
        super();
        this.address = delegations.address || ''
        this.moniker = delegations.moniker || ''
        this.amount = delegations.amount || ''
        this.total_shares = delegations.total_shares || ''
        this.self_shares = delegations.self_shares || ''
    }

    static bundleData(value: any): ValidatorDelegationsResDto[] {
        let data: ValidatorDelegationsResDto[] = [];
        data = value.map((v: any) => {
            return new ValidatorDelegationsResDto(v);
        });
        return data;
    }
}

export class ValidatorUnBondingDelegationsResDto extends BaseResDto {
    address: string;
    moniker: string;
    amount: object;
    block: number;
    until: number;

    constructor(unBondingDelegations) {
        super();
        this.address = unBondingDelegations.address || ''
        this.moniker = unBondingDelegations.moniker || ''
        this.amount = unBondingDelegations.amount || ''
        this.block = unBondingDelegations.block || ''
        this.until = unBondingDelegations.until || ''
    }

    static bundleData(value: any): ValidatorUnBondingDelegationsResDto[] {
        let data: ValidatorUnBondingDelegationsResDto[] = [];
        data = value.map((v: any) => {
            return new ValidatorUnBondingDelegationsResDto(v);
        });
        return data;
    }

}

export class ValidatorDetailResDto {
    total_power: number;
    self_power: number;
    status: string;
    bonded_tokens: string;
    self_bond: string;
    bonded_stake: string;
    delegator_shares: string;
    delegator_num: number;
    commission_rate: string;
    commission_update: number;
    commission_max_rate: string;
    commission_max_change_rate: string;
    bond_height: string;
    unbonding_height: string;
    jailed_until: string;
    missed_blocks_count: string;
    operator_addr: string;
    owner_addr: string;
    consensus_pubkey: string;
    description: object;
    icons: string;
    uptime: number;
    stats_blocks_window: string;

    constructor(validatorDetail) {
        this.total_power = validatorDetail.total_power || 0
        this.self_power = validatorDetail.tokens || 0
        this.status = validatorDetail.valStatus || ''
        this.bonded_tokens = validatorDetail.tokens || ''
        this.self_bond = validatorDetail.self_bond || ''
        this.bonded_stake = validatorDetail.bonded_stake || ''
        this.delegator_shares = validatorDetail.delegator_shares || ''
        this.delegator_num = validatorDetail.delegator_num || 0
        this.commission_rate = validatorDetail.commission.commission_rates.rate || ''
        this.commission_update = validatorDetail.commission.update_time || 0
        this.commission_max_rate = validatorDetail.commission.commission_rates.max_rate || ''
        this.commission_max_change_rate = validatorDetail.commission.commission_rates.max_change_rate || ''
        this.bond_height = validatorDetail.start_height || '0'
        this.unbonding_height = validatorDetail.unbonding_height || ''
        this.jailed_until = validatorDetail.jailed_until || ''
        this.missed_blocks_count = validatorDetail.missed_blocks_counter || '0'
        this.operator_addr = validatorDetail.operator_address || ''
        this.owner_addr = validatorDetail.owner_addr || ''
        this.consensus_pubkey = validatorDetail.consensus_pubkey || ''
        this.description = validatorDetail.description || {}
        this.icons = validatorDetail.icon || ''
        this.uptime = validatorDetail.uptime || ''
        this.stats_blocks_window = validatorDetail.stats_blocks_window || ''
    }
}

export class AccountAddrResDto {
    amount: {
        denom: string;
        amount: string
    };
    withdrawAddress: string;
    address: string;
    deposits: object;
    isProfiler: boolean;
    moniker: string;
    status: string;
    operator_address: string;

    constructor(account) {
        this.amount = account.amount || {}
        this.withdrawAddress = account.withdrawAddress || ''
        this.address = account.address || ''
        this.deposits = account.deposits || ''
        this.isProfiler = account.isProfiler || false
        this.moniker = account.moniker || ''
        this.status = account.status || ''
        this.operator_address = account.operator_address || ''
    }
}

export class DelegatorsDelegationsResDto extends BaseResDto {
    address: string;
    moniker: string;
    amount: {
        denom: string;
        amount: string|number
    };
    shares: string;
    height: string;

    constructor(delegations) {
        super();
        this.address = delegations.address || ''
        this.moniker = delegations.moniker || ''
        this.amount = delegations.amount || {}
        this.shares = delegations.shares || ''
        this.height = delegations.height || ''
    }

    static bundleData(value: any): DelegatorsDelegationsResDto[] {
        let data: DelegatorsDelegationsResDto[] = [];
        data = value.map((v: any) => {
            return new DelegatorsDelegationsResDto(v);
        });
        return data;
    }
}

export class DelegatorsUndelegationsResDto extends BaseResDto {
    address: string;
    moniker: string;
    amount: {
        denom: string;
        amount: string|number
    };
    height: string;
    end_time: string;

    constructor(delegations) {
        super();
        this.address = delegations.address || ''
        this.moniker = delegations.moniker || ''
        this.amount = delegations.amount || {}
        this.height = delegations.height || ''
        this.end_time = delegations.end_time || ''
    }

    static bundleData(value: any): DelegatorsUndelegationsResDto[] {
        let data: DelegatorsUndelegationsResDto[] = [];
        data = value.map((v: any) => {
            return new DelegatorsUndelegationsResDto(v);
        });
        return data;
    }
}

export class ValidatorVotesResDto extends BaseResDto {
    title: string;
    proposal_id: number;
    status: string;
    voted: string;
    tx_hash: string;
    proposal_link: boolean;

    constructor(vote) {
        super();
        this.title = vote.title || ''
        this.proposal_id = vote.proposal_id || ''
        this.status = vote.status || ''
        this.voted = vote.voted || ''
        this.tx_hash = vote.tx_hash || ''
        this.proposal_link = vote.proposal_link
    }

    static bundleData(value: any): ValidatorVotesResDto[] {
        let data: ValidatorVotesResDto[] = [];
        data = value.map((v: any) => {
            return new ValidatorVotesResDto(v);
        });
        return data;
    }
}

export class ValidatorDepositsResDto extends BaseResDto {
    proposal_id: number;
    proposer: string;
    amount: Coin[];
    submited: boolean;
    tx_hash: string;
    moniker: string;
    proposal_link: boolean;

    constructor(deposit) {
        super();
        this.proposal_id = deposit.proposal_id || 0;
        this.proposer = deposit.proposer || '';
        this.amount = Coin.bundleData(deposit.amount) || [];
        this.submited = deposit.submited || false;
        this.tx_hash = deposit.tx_hash || '';
        this.moniker = deposit.moniker || '';
        this.proposal_link = deposit.proposal_link;
    }

    static bundleData(value: any): ValidatorDepositsResDto[] {
        let data: ValidatorDepositsResDto[] = [];
        data = value.map((v: any) => {
            return new ValidatorDepositsResDto(v);
        });
        return data;
    }
}