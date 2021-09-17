import { IsString, IsNotEmpty } from 'class-validator';
import { BaseReqDto, BaseResDto, PagingReqDto } from './base.dto';
import { Coin } from './common.res.dto';
import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
import { IDenomStruct } from '../types/schemaTypes/denom.interface';

/************************   request dto   ***************************/

// distribution/delegators/{delegatorAddr}/withdraw_address request dto
export class WithdrawAddressReqDto extends BaseReqDto {
    @ApiProperty()
    delegatorAddr: string;
}

// distribution/delegators/{delegatorAddr}/reward
export class DelegatorRewardsReqDto extends BaseReqDto {
    @ApiProperty()
    delegatorAddr: string;
}

// distribution/validators/:address reqdto
export class ValCommissionRewReqDto {
    @ApiProperty()
    @ApiPropertyOptional()
    address?: string
}
/************************   respond dto   ***************************/
// distribution/delegators/{delegatorAddr}/withdraw_address request dto
export class WithdrawAddressResDto extends BaseResDto {
    address:string;
    constructor(address:string) {
        super();
        this.address = address;
    }
}

// distribution/delegators/{delegatorAddr}/rewards
export class DelegatorRewardsResDto extends BaseResDto{
    rewards: Reward[];
    total: Coin[];
    constructor(value) {
        super();
        let { rewards, total } = value;
        this.rewards = Reward.bundleData(rewards);
        this.total = Coin.bundleData(total);
    }
}

export class Reward {
    validator_address:string;
    moniker:string;
    reward:Coin[];
    constructor(value) {
        let { validator_address, reward, moniker } = value;
        this.validator_address = validator_address || '';
        this.moniker = moniker || '';
        this.reward = reward && Coin.bundleData(reward);
    }

    static bundleData(value: any = []): Reward[] {
        let data: Reward[] = [];
        data = value.map((v: any) => {
            return new Reward(v);
        });
        return data;
    }
}

// distribution/validators/:address resdto
export class ValCommissionRewResDto extends BaseResDto {
    // operator_address: string;
    // self_bond_rewards: [];
    val_commission: object;

    constructor(commissionRewards) {
        super();
        // this.operator_address = commissionRewards.operator_address || ''
        // this.self_bond_rewards = commissionRewards.self_bond_rewards || {}
        this.val_commission = commissionRewards.val_commission || {}
    }
}