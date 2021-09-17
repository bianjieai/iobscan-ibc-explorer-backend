export interface IGovProposal {
    id: string,
    content: object,
    status: string,
    final_tally_result: object,
    current_tally_result: object,
    submit_time	: number,
    deposit_end_time: number,
    total_deposit: object,
    initial_deposit: object,
    voting_start_time: number,
    voting_end_time: number,
    hash:string,
    proposer: string,
    is_deleted: boolean,
    min_deposit: string,
    quorum: string,
    threshold: string,
    veto_threshold: string,
    create_time: number,
    update_time: number
}

export interface IGovProposalQuery {
    pageNum?: string;
    pageSize?: string;
    useCount?: boolean | string;
    status?: string;
}
