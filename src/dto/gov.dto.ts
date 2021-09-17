import { BaseReqDto, PagingReqDto, BaseResDto } from './base.dto';
import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
import { Coin } from './common.res.dto';

/***************Req***********************/

export class proposalsReqDto {
    @ApiPropertyOptional()
    pageNum?: number;

    @ApiPropertyOptional()
    pageSize?: number;

    @ApiPropertyOptional({description:'true/false'})
    useCount?: boolean;

    @ApiPropertyOptional({description: 'status: "DepositPeriod,VotingPeriod,Passed,Rejected" 以,分割的字符串'})
    status?: string;
}

export class ProposalDetailReqDto extends BaseReqDto {
    @ApiProperty()
    id: number;
}

export class proposalsVoterReqDto extends PagingReqDto {
    @ApiPropertyOptional({description: 'type: all/validator/delegator'})
    voterType?: string;
}

/***************Res*************************/

export class govProposalResDto extends BaseResDto {
    id: string;
    content: object;
    status: string;
    final_tally_result: object;
    current_tally_result: object;
    tally_details: any;
    submit_time: number;
    deposit_end_time: number;
    total_deposit: object;
    initial_deposit: object;
    voting_end_time: number;
    min_deposit: string;
    quorum: string;
    threshold: string;
    veto_threshold: string;

    constructor(proposal) {
        super();
        this.id = proposal.id || '';
        this.content = proposal.content || {};
        this.status = proposal.status || '';
        this.final_tally_result = proposal.final_tally_result || {};
        this.current_tally_result = proposal.current_tally_result || {};
        this.tally_details = proposal.tally_details || [];
        this.submit_time = proposal.submit_time || {};
        this.deposit_end_time = proposal.deposit_end_time || 0;
        this.total_deposit = proposal.total_deposit || {};
        this.initial_deposit = proposal.initial_deposit || {};
        this.voting_end_time = proposal.voting_end_time || 0;
        this.min_deposit = proposal.min_deposit || '';
        this.quorum = proposal.quorum || '';
        this.threshold = proposal.threshold || '';
        this.veto_threshold = proposal.veto_threshold || '';
    }

    static bundleData(value: any): govProposalResDto[] {
        let data: govProposalResDto[] = [];
        data = value.map((v: any) => {
            return new govProposalResDto(v);
        });
        return data;
    }
}

export class govProposalDetailResDto extends govProposalResDto {
    hash: string;
    burned_rate: string | null;
    proposer: string;
    voting_start_time: string;

    constructor(proposal) {
        super(proposal);
        this.hash = proposal.hash || '';
        this.burned_rate = proposal.burned_rate || null;
        this.proposer = proposal.proposer || '';
        this.voting_start_time = proposal.voting_start_time || '';
    }
}

export class govProposalVoterResDto extends BaseResDto {
    voter: string;
    address: string;
    moniker: string;
    option: string;
    hash: string;
    timestamp: number;
    height: number;
    isValidator: boolean;

    constructor(value) {
        super();
        this.voter = value.voter || '';
        this.address = value.address || '';
        this.moniker = value.moniker || '';
        this.option = value.option || '';
        this.hash = value.hash || '';
        this.timestamp = value.timestamp || 0;
        this.height = value.height || 0;
        this.isValidator = value.isValidator || false;
    }

    static bundleData(value: any): govProposalVoterResDto[] {
        let data: govProposalVoterResDto[] = [];
        data = value.map((v: any) => {
            return new govProposalVoterResDto(v);
        });
        return data;
    }
}

export class govProposalDepositorResDto extends BaseResDto {
    hash: string;
    moniker: string;
    address: string;
    amount: Coin[];
    type: string;
    timestamp: number;

    constructor(value) {
        super();
        this.hash = value.hash || '';
        this.moniker = value.moniker || '';
        this.address = value.address || '';
        this.amount = Coin.bundleData(value.amount) || [];
        this.type = value.type || '';
        this.timestamp = value.timestamp || 0;
    }

    static bundleData(value: any): govProposalDepositorResDto[] {
        let data: govProposalDepositorResDto[] = [];
        data = value.map((v: any) => {
            return new govProposalDepositorResDto(v);
        });
        return data;
    }
}