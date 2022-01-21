import {Test} from '@nestjs/testing'
import {IbcStatisticsTaskService} from "./ibc_statistics.task.service";
import {AppModule} from "../app.module";


describe('IbcStatisticsTaskService', () => {
    let ibcStatisticsTaskService: IbcStatisticsTaskService;
    beforeEach(async () => {
        const module = await Test.createTestingModule({
            imports: [
                AppModule
            ]
        }).compile();
        ibcStatisticsTaskService = module.get<IbcStatisticsTaskService>(IbcStatisticsTaskService);
    })
    describe('parseIbcStatistics', () => {
        it('parseIbcStatistics', async () => {
            const dateNow = Math.floor(1637737022);
            await ibcStatisticsTaskService.parseIbcStatistics(dateNow)
            // console.log(ibcTxTaskService,'----')
        });
    });

    describe("AggregateFindSrcChannels", () => {
        it('AggregateFindSrcChannels ', async () => {
            const dateNow = Math.floor(1623955689);
            const chains = ["osmosis_1", "cosmoshub_4"]
            const data = await ibcStatisticsTaskService.aggregateFindSrcChannels(dateNow, chains)
            console.log(data)
        });
    });

    describe("AggregateFindDesChannels", () => {
        it('AggregateFindDesChannels Test ', async () => {
            const dateNow = Math.floor(1623955689);
            const chains = ["osmosis_1", "cosmoshub_4", "regen_1"]
            const data = await ibcStatisticsTaskService.aggregateFindDesChannels(dateNow, chains)
            console.log(data)
        });
    });

    describe("updateStaticsData", () => {
        it('updateDb Test', async () => {
            const record = await ibcStatisticsTaskService.findStatisticsRecord("chains_24hr")
            record.count = 10
            console.log(record)
            await ibcStatisticsTaskService.updateStatisticsRecord(record)
        });
    });

    describe("aggregateBaseDenomCnt", () => {
        it('aggregateBaseDenomCnt Test', async () => {
            const record = await ibcStatisticsTaskService.aggregateBaseDenomCnt()
            console.log(record)
        });

        it('getCountinfo Test', async () => {
            const record = await ibcStatisticsTaskService.getCountinfo()
            console.log(record)
        });
    });

})