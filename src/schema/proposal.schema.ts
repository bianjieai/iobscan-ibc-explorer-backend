import * as mongoose from 'mongoose';
import {
    IGovProposal,
    IGovProposalQuery
} from "../types/schemaTypes/proposal.interface";
import {IListStruct} from "../types";
export const ProposalSchema = new mongoose.Schema({
    id: Number,
    content: Object,
    status: String,
    final_tally_result: Object,
    current_tally_result: Object,
    submit_time	: Number,
    deposit_end_time: Number,
    total_deposit: Object,
    initial_deposit: Object,
    voting_start_time: Number,
    voting_end_time: Number,
    hash:String,
    proposer: String,
    is_deleted: Boolean,
    min_deposit: String,
    quorum: String,
    threshold: String,
    veto_threshold: String,
    create_time: Number,
    update_time: Number
})
ProposalSchema.index({id: 1}, {unique: true})

ProposalSchema.statics = {
    async queryAllProposals() {
        return await this.find({is_deleted: false})
    },
    async updateProposals(ids) {
        return await this.updateMany({id: { $in: ids } }, {$set: {is_deleted: true}})
    },
    async insertProposal(insertProposal:IGovProposal) {
        let { id } = insertProposal
        const options = {upsert: true, new: false, setDefaultsOnInsert: true}
        await this.findOneAndUpdate({id}, insertProposal, options)
    },
    async queryProposals(query: IGovProposalQuery): Promise<IListStruct> {
        let queryParameters: any = {
            is_deleted: false
        };
        if (query.status) {
            queryParameters = {
                status: {
                    $in: query.status.split(",")
                },
                is_deleted: false
            };
        }
        const result: IListStruct = {}
        if (query.useCount && query.useCount == 'true') {
            result.count = await this.find(queryParameters).countDocuments();
        }
        if (query.status) {
            result.data = await this.find(queryParameters).sort({ id: -1 });
        } else {
            result.data = await this.find(queryParameters)
            .sort({id: -1})
            .skip((Number(query.pageNum) - 1) * Number(query.pageSize))
            .limit(Number(query.pageSize));
        }
        return result
    },
    async findOneById(id: number, is_deleted: boolean): Promise<IGovProposal> {
        let queryParameters = typeof is_deleted === 'undefined' ? { id: id } : { id: id, is_deleted };
        return await this.findOne(queryParameters).select({ '_id': 0, '__v': 0 });
    },
    async queryAllProposalsSelect() {
        return await this.find({}).select({_id: 0,id: 1,content: 1,is_deleted:1})
    },
    async queryAllProposalsDeletedID() {
        return await this.find({is_deleted:1}).select({_id: 0,id: 1})
    },
}
