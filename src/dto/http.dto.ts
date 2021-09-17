import { IsString, IsInt, Length, Min, Max, IsOptional, Equals, MinLength, ArrayNotEmpty } from 'class-validator';
import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
import { BaseReqDto, BaseResDto, PagingReqDto } from './base.dto';
import { Coin } from './common.res.dto';
import { ApiError } from '../api/ApiResult';
import { ErrorCodes } from '../api/ResultCodes';
import { IBindTx } from '../types/tx.interface';
import { IDenomStruct } from '../types/schemaTypes/denom.interface';
import { validatorStatusFromLcd } from '../constant'
// lcd相关请求返回值数据模型

/************************   response dto   ***************************/
// /ibc/applications/transfer/v1beta1/denom_traces/{hash} response dto
export class IbcTraceDto {
  denom_trace: { path: string, base_denom: string; };
  constructor(value) {
    const { denom_trace } = value;
    this.denom_trace = denom_trace || '';
  }
}

// /nft/nfts/denoms response dto
export class DenomDto {
    id: string;
    name: string;
    schema:string;
    creator:string;
    constructor(value) {
        this.id = value.id || '';
        this.name = value.name || '';
        this.schema = value.schema || '';
        this.creator = value.creator || '';
    }

    static bundleData(value: any = []): DenomDto[] {
        let data: DenomDto[] = [];
        data = value.map((v: any) => {
            return new DenomDto(v);
        });
        return data;
    }
}

// /nft/nfts/collections/{denom} response dto

export class Nft {
    id: string;
    name: string;
    uri: string;
    data: string;
    owner: string;

    constructor(value) {
        this.id = value.id || '';
        this.name = value.name || '';
        this.uri = value.uri || '';
        this.data = value.data || '';
        this.owner = value.owner || '';
    }
}

export class NftCollectionDto {
    denom: DenomDto;
    nfts: Nft[];

    constructor(value) {
        const {denom, nfts} = value;
        this.denom = new DenomDto(denom);
        this.nfts = (nfts || []).map(item=>{
            return new Nft(item);
        });
    }
}

// /validator/validators response dto
export class VValidatorDto {
    id: string;
    name: string;
    pubkey: string;
    certificate: string;
    power: string;
    operator: string;
    jailed:boolean;
    details:string;

    constructor(value) {
        this.id = value.id || '';
        this.name = value.name || '';
        this.pubkey = value.pubkey || '';
        this.certificate = value.certificate || '';
        this.power = value.power || '';
        this.operator = value.operator || '';
        this.jailed = value.jailed || undefined;
        this.details = value.details || '';
    }

    static bundleData(value: any = []): VValidatorDto[] {
        let data: VValidatorDto[] = [];
        data = value.map((v: any) => {
            return new VValidatorDto(v);
        });
        return data;
    }
}


// /blocks/{height} or /blocks/latest response dto
export class BlockId {
    hash: string;
    parts: { total:number, hash:string};
    constructor(value) {
        const { hash, parts } = value;
        this.hash = hash || '';
        this.parts = parts;
    }
}

export class Signatures {
    block_id_flag: string;
    validator_address: string;
    timestamp: string;
    signature: string;

    constructor(value) {
        const { block_id_flag, validator_address, timestamp, signature } = value;
        this.block_id_flag = block_id_flag || '';
        this.validator_address = validator_address || '';
        this.timestamp = timestamp || '';
        this.signature = signature || '';
    }
    static bundleData(value: any = []): Signatures[] {
        let data: Signatures[] = [];
        data = value.map((v: any) => {
            return new Signatures(v);
        });
        return data;
    }
}

export class BlockHeader {
    version: { block: string, app:string };
    chain_id: string;
    height: number;
    time: string;
    last_block_id: BlockId;
    last_commit_hash: string;
    data_hash: string;
    validators_hash: string;
    next_validators_hash: string;
    consensus_hash: string;
    app_hash: string;
    last_results_hash: string;
    evidence_hash: string;
    proposer_address: string;

    constructor(value) {
        this.version = value.version;
        this.chain_id = value.chain_id || '';
        this.height = Number(value.height);
        this.time = value.time || '';
        this.last_block_id = new BlockId(value.last_block_id);
        this.last_commit_hash = value.last_commit_hash || '';
        this.data_hash = value.data_hash || '';
        this.validators_hash = value.validators_hash || '';
        this.next_validators_hash = value.next_validators_hash || '';
        this.consensus_hash = value.consensus_hash || '';
        this.app_hash = value.app_hash || '';
        this.last_results_hash = value.last_results_hash || '';
        this.evidence_hash = value.evidence_hash || '';
        this.proposer_address = value.proposer_address || '';
    }
}

export class Commit {
    height: number;
    round: number;
    block_id: BlockId;
    signatures: Signatures[];
    constructor(value) {
        const { height, round, block_id, signatures } = value;
        this.height = Number(height);
        this.round = round;
        this.block_id = new BlockId(block_id);
        this.signatures = Signatures.bundleData(signatures);
    }
}

export class Block {
    header: BlockHeader;
    data: { txs:  any[]};
    evidence: { evidence:  any[]};
    last_commit: Commit;

    constructor(value) {
        const { header, data, evidence, last_commit } = value;
        this.header = new BlockHeader(header);
        this.data =  data;
        this.evidence =  evidence;
        this.last_commit = new Commit(last_commit);
    }
}

export class BlockDto {
    block_id: BlockId;
    block: Block;
    constructor(value) {
        const { block_id,  block } = value;
        this.block_id = new BlockId(block_id);
        this.block = new Block(block);
    }
}

// distribution/delegators/{delegatorAddr}/withdraw_address response dto
export class WithdrawAddressDto {
    address: string;
    constructor(address:string) {
        this.address = address;
    }
}

// /distribution/delegators/{delegatorAddr}/rewards
export class DelegatorRewardsDto {
    rewards: Reward[];
    total: Coin[];
    constructor(value) {
        const { rewards, total } = value;
        this.rewards = Reward.bundleData(rewards);
        this.total = Coin.bundleData(total);
    }
}

export class Reward {
    validator_address:string;
    reward:Coin[];
    constructor(value) {
        const { validator_address, reward } = value;
        this.validator_address = validator_address || '';
        this.reward = reward ? Coin.bundleData(reward) : [];
    }

    static bundleData(value: any = []): Reward[] {
        let data: Reward[] = [];
        data = value.map((v: any) => {
            return new Reward(v);
        });
        return data;
    }
}
export class TokensLcdDto {
    '@type':string;
    symbol: string;
    name:string;
    scale:number;
    denom:string;
    initial_supply: string;
    max_supply: string;
    mintable:boolean;
    owner: string;
    is_main_token: boolean;
    total_supply: string;
    update_block_height: number;
    src_protocol: string;
    chain: string;

    constructor(value) {
        this['@type'] = value['@type'] || '';
        this.symbol = value.symbol || '';
        this.name = value.name || '';
        this.scale = value.scale || 0;
        this.denom = value.min_unit || '';
        this.initial_supply = value.initial_supply || '';
        this.max_supply = value.max_supply || '';
        this.mintable = value.mintable || false;
        this.owner = value.owner || '';
        this.is_main_token = value.is_main_token || false;
        this.total_supply = value.initial_supply || '';
        this.update_block_height = value.update_block_height || 0;
        this.src_protocol = value.src_protocol || '';
        this.chain = value.chain || '';
    }

    static bundleData(value: any = []): TokensLcdDto[] {
        let data: TokensLcdDto[] = [];
        data = value.map((v: any) => {
            return new TokensLcdDto(v);
        });
        return data;
    }
}
export class TokensStakingLcdToken {
        unbonding_time:string;
        max_validators:number;
        max_entries:number;
        historical_entries:number;
        bond_denom:string
    constructor(value) {
        this.unbonding_time = value.unbonding_time || '';
        this.max_validators = value.max_validators || 0;
        this.max_entries = value.max_entries || 0;
        this.historical_entries = value.historical_entries || 0;
        this.bond_denom = value.bond_denom || '';

    }
}
// staking/validators
export class StakingValidatorLcdDto {
    operator_address: string;
    consensus_pubkey: string | object;
    jailed: boolean;
    status: number;
    tokens: string;
    delegator_shares: string;
    description: object;
    unbonding_height: string;
    unbonding_time: string;
    commission: object;
    min_self_delegation: string;


    constructor(value) {
        this.operator_address = value.operator_address || '';
        this.consensus_pubkey = value.consensus_pubkey || '';
        this.jailed = value.jailed || false;
        this.status = value.status ?  validatorStatusFromLcd[value.status] : '';
        this.tokens = value.tokens || '';
        this.delegator_shares = value.delegator_shares || '';
        this.description = value.description || '';
        this.unbonding_height = value.unbonding_height || '0';
        this.unbonding_time = value.unbonding_time || '';
        this.commission = value.commission || null;
        this.min_self_delegation = value.min_self_delegation || '';
    }

    static bundleData(value: any = []): StakingValidatorLcdDto[] {
        let data: StakingValidatorLcdDto[] = [];
        data = value.map((v: any) => {
            return new StakingValidatorLcdDto(v);
        });
        return data;
    }
}

// /slashing/validators/${validatorPubkey}/signing_info
export class StakingValidatorSlashLcdDto {
    address: string;
    start_height?: string;
    index_offset: string;
    jailed_until: string;
    tombstoned?: boolean;
    missed_blocks_counter?: string;

    constructor(value) {
        this.address = value.address || '';
        this.start_height = value.start_height || '';
        this.index_offset = value.index_offset || '';
        this.jailed_until = value.jailed_until || '';
        this.tombstoned = value.tombstoned || false;
        this.missed_blocks_counter = value.missed_blocks_counter || '';
    }
}

// /staking/validators/${valOperatorAddr}/delegations
export class StakingValidatorDelegationLcdDto {
    total: number;
    result: Array<IDelegationLcd>;

    constructor(value) {
        this.total = value.total || 0;
        this.result = StakingValidatorDelegationLcdDto.bundleData(value.result) || [];
    }
    static bundleData(value: any = []): IDelegationLcd[] {
        let data: IDelegationLcd[] = [];
        data = value.map((v: any) => {
            return new IDelegationLcd(v);
        });
        return data;
    }
}



export class IDelegationLcd {
    delegation: {
        delegator_address: string;
        validator_address: string;
        shares: string
    };
    balance: {
        amount: string;
        denom: string
    };

    constructor(value) {
        this.delegation = value.delegation || {};
        this.balance = value.balance || {};
    }
}


// /validatorsets/{height}
export class Validatorset {
    address:string;
    pub_key:string | object;
    proposer_priority:string;
    voting_power:string;
    constructor(value) {
        const { address, pub_key, proposer_priority, voting_power } = value;
        this.address = address || '';
        this.pub_key = pub_key || '';
        this.proposer_priority = proposer_priority || '';
        this.voting_power = voting_power || '';
    }

    static bundleData(value: any = []): Validatorset[] {
        let data: Validatorset[] = [];
        data = value.map((v: any) => {
            return new Validatorset(v);
        });
        return data;
    }
}

export class IThemStruct {
    id?: string;
    pictures?: {primary?: {url?: string}};
}

// https://keybase.io/_/api/1.0/user/lookup.json?fields=pictures&key_suffix={identity}
export class IconUriLcdDto {
    status: {
        code: number,
        name: string
    };
    them?: IThemStruct[]

    constructor(value) {
        this.status = value.status || {};
        this.them = value.them || {};
    }

}


// /slashing/parameters
export class StakingValidatorParametersLcdDto {
    signed_blocks_window: string;
    min_signed_per_window: string;
    downtime_jail_duration: string;
    slash_fraction_double_sign: string;
    slash_fraction_downtime: string;

    constructor(value) {
        this.signed_blocks_window = value.signed_blocks_window || '';
        this.min_signed_per_window = value.min_signed_per_window || '';
        this.downtime_jail_duration = value.downtime_jail_duration.replace('s', '000000000') || '';
        this.slash_fraction_double_sign = value.slash_fraction_double_sign || '';
        this.slash_fraction_downtime = value.slash_fraction_downtime || '';
    }
}

export class ISelfBondRewards {
    amount: string;
    denom: string;

    constructor(value) {
        this.amount = value.amount || '';
        this.denom = value.denom || '';
    }
}

// /distribution/validators/${valAddress}
export class commissionRewardsLcdDto {
    // operator_address: string;
    // self_bond_rewards: [];
    val_commission: object;

    constructor(value) {
        // this.operator_address = value.operator_address || '';
        // this.self_bond_rewards = value.self_bond_rewards || [];
        this.val_commission = value.commission || {};
    }
}

// /cosmos/distribution/v1beta1/community_pool
export class communityPoolLcdDto {
    pool: Coin[];

    constructor(value) {
        this.pool = value.pool || [];
    }
}



// /staking/validators/${address}/unbonding_delegations


export class StakingValUnBondingDelLcdDto {
    total: number;
    result: Array<StakingValUnBondingDelLcdResultDto>;

    constructor(value) {
        this.total = value.total || 0;
        this.result = StakingValUnBondingDelLcdResultDto.bundleData(value.result) || [];
    }
}

export class StakingValUnBondingDelLcdResultDto {
    delegator_address: string;
    validator_address: string;
    entries: Array<UnBondingDel>

    constructor(value) {
        this.delegator_address = value.delegator_address || '';
        this.validator_address = value.validator_address || [];
        this.entries = value.entries || [];
    }

    static bundleData(value: any = []): StakingValUnBondingDelLcdResultDto[] {
        let data: StakingValUnBondingDelLcdResultDto[] = [];
        data = value.map((v: any) => {
            return new StakingValUnBondingDelLcdResultDto(v);
        });
        return data;
    }
}

export class UnBondingDel {
    creation_height: string;
    completion_time: string;
    initial_balance: string;
    balance: string
}

export class AddressBalancesLcdDto {
    amount: string
    denom: string

    constructor(value) {
        this.amount = value.amount || '';
        this.denom = value.denom || '';
    }

    static bundleData(value: any = []): AddressBalancesLcdDto[] {
        let data: AddressBalancesLcdDto[] = [];
        data = value.map((v: any) => {
            return new AddressBalancesLcdDto(v);
        });
        return data;
    }
}

export class DelegatorsDelegationLcdDto {
    height: string;
    total: number;
    next_key: string | null;
    result: DelegatorsResult[];
    constructor(value) {
        this.height = value.height || '';
        this.result = DelegatorsResult.bundleData(value.delegation_responses);
        this.total = (value.pagination && Number(value.pagination.total)) || 0;
        this.next_key = value.pagination && value.pagination.next_key;
    }
}

export class DelegatorsResult {
    delegation: {
        delegator_address: string;
        validator_address: string;
        shares:string
    };
    balance: {
        denom: string;
        amount: string
    };

    constructor(value) {
        this.delegation = value.delegation || {},
        this.balance = value.balance || {}
    }

    static bundleData(value: any = []): DelegatorsResult[] {
        let data: DelegatorsResult[] = [];
        data = value.map((v: any) => {
            return new DelegatorsResult(v);
        });
        return data;
    }
}

export class DelegatorsUndelegationLcdDto {
    height: string;
    next_key: string | null;
    result: UndelegatorsResult[];
    total: number;
    constructor(value) {
        this.height = value.height || '';
        this.result = UndelegatorsResult.bundleData(value.unbonding_responses);
        this.total = (value.pagination && Number(value.pagination.total)) || 0;
        this.next_key = value.pagination && value.pagination.next_key;
    }
}

export class UndelegatorsResult {
    delegator_address: string;
    validator_address: string;
    entries: {
        creation_height: string;
        completion_time: string;
        initial_balance: string;
        balance: string
    };
    constructor(value) {
        this.delegator_address = value.delegator_address || '',
        this.validator_address = value.validator_address || ''
        this.entries = value.entries || {}
    }

    static bundleData(value: any = []): UndelegatorsResult[] {
        let data: UndelegatorsResult[] = [];
        data = value.map((v: any) => {
            return new UndelegatorsResult(v);
        });
        return data;
    }
}

export class BondedTokensLcdDto {
    not_bonded_tokens:string;
    bonded_tokens:string;

    constructor(value) {
        this.not_bonded_tokens = value.not_bonded_tokens || 0;
        this.bonded_tokens = value.bonded_tokens || 0;
    }
}


export class TotalSupplyLcdDto {
    supply: Coin[];
    constructor(value) {
        const { supply } = value;
        this.supply = Coin.bundleData(supply);
    }
}

export class GovParamsLcdDto {
    voting_params: {
        voting_period: string;
    };
    deposit_params: {
        min_deposit: Coin[],
        max_deposit_period: string
    };
    tally_params: {
        quorum:string,
        threshold:string,
        veto_threshold:string,
    };
    constructor(value) {
        this.voting_params = value.voting_params || {};
        this.deposit_params = {
            min_deposit: Coin.bundleData(value.deposit_params.min_deposit) || [],
            max_deposit_period: value.deposit_params.max_deposit_period || ''
        }
        this.tally_params = value.tally_params || {};
    }
}

export class GovTallyParamsLcdDto {
    quorum: string;
    threshold: string;
    veto_threshold: string;
    constructor(value) {
        this.quorum = value.quorum || '';
        this.threshold = value.threshold || '';
        this.veto_threshold = value.veto_threshold || '';
    }
}

export class GovDepositParamsLcdDto {
    min_deposit: Coin[];
    max_deposit_period: string
    constructor(value) {
        this.min_deposit = Coin.bundleData(value.min_deposit),
        this.max_deposit_period = value.max_deposit_period || '';
    }
}

export class GovProposalLcdDto {
    id: number;
    content: object;
    status: string;
    final_tally_result: {
        yes: number,
        abstain: number,
        no: number,
        no_with_veto: number
    };
    submit_time: string;
    deposit_end_time: string;
    total_deposit: Coin[];
    voting_start_time: string;
    voting_end_time: string;
    is_deleted: boolean;

    constructor(value) {
        this.id = Number(value.proposal_id) || 0;
        this.content = value.content || {};
        this.status = value.status || '';
        this.final_tally_result = {
            yes: Number(value.final_tally_result && value.final_tally_result.yes) || 0,
            abstain: Number(value.final_tally_result && value.final_tally_result.abstain) || 0,
            no: Number(value.final_tally_result && value.final_tally_result.no) || 0,
            no_with_veto: Number(value.final_tally_result && value.final_tally_result.no_with_veto) || 0,
        };
        this.submit_time = value.submit_time || '';
        this.deposit_end_time = value.deposit_end_time || '';
        this.total_deposit = Coin.bundleData(value.total_deposit);
        this.voting_start_time = value.voting_start_time || '';
        this.voting_end_time = value.voting_end_time || '';
        this.is_deleted = false;
    }

    static bundleData(value: any): GovProposalLcdDto[] {
        let data: GovProposalLcdDto[] = [];
        data = value.map((v: any) => {
            return new GovProposalLcdDto(v);
        });
        return data;
    }
}
