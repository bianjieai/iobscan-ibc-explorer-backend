// @ts-ignore
import {IbcTxHandler} from "../util/IbcTxHandler";
import {Test} from "@nestjs/testing";
import {AppModule} from "../app.module";

describe('IbcTxHandler', () => {
    let ibcTxHandler: IbcTxHandler;
    beforeEach(async () => {
        const module = await Test.createTestingModule({
            imports: [
                AppModule
            ]
        }).compile();
        ibcTxHandler = module.get<IbcTxHandler>(IbcTxHandler);
    })

    describe('parseIbcTx', () => {
        it('parseIbcTxBychain Test', async () => {
            const dateNow =  new Date().getTime() / 1000;
            const txmodel = await ibcTxHandler.getIbcTxModel()
            await ibcTxHandler.parseIbcTx(txmodel,dateNow)
            console.log('----')
        });
    });

    describe('getRecordLimitTx', () => {
        it('getRecordLimitTx Test', async () => {
            let chain_id = "emoney_3"
            const result  = await ibcTxHandler.getRecordLimitTx(chain_id,13106,10)
            console.log(result,'----')
        });
    });
    describe('checkTaskFollowingStatus', () => {
        it('checkTaskFollowingStatus Test', async () => {
            let chain_id = "irishub-test"
            const result  = await ibcTxHandler.checkTaskFollowingStatus(chain_id)
            console.log(result,'-checkTaskFollowingStatus---')
        });
    });

    describe('handlerSourcesTx', () => {
        it('handlerSourcesTx Test', async () => {
            let chain_id = "osmosis_1"
            const txs  = await ibcTxHandler.getRecordLimitTx(chain_id,1858626,10)
            const {allChainsMap,allChainsDenomPathsMap} = await ibcTxHandler.getAllChainsMap()
            const dateNow = Math.floor(new Date().getTime() / 1000);
            let denomMap = await ibcTxHandler.getDenomRecordByChainId(chain_id)
            const {handledTx,denoms}  = await ibcTxHandler.handlerSourcesTx(txs,chain_id,dateNow,allChainsMap,allChainsDenomPathsMap,denomMap)
            console.log(handledTx,'--ibcTxs--')
            console.log(denoms,'--denoms--')
        });
    });

    describe('changeIbcTxState', () => {
        it('changeIbcTxState Test', async () => {
            const dateNow = Math.floor(new Date().getTime() / 1000);
            jest.setTimeout(100000000)
            const txmodel = await ibcTxHandler.getIbcTxLatestModel()
            await ibcTxHandler.changeIbcTxState(txmodel,dateNow,[3],false,[])
            return null
        });
    });
    describe('getDenomRecordByChainId', () => {
        it('getDenomRecordByChainId Test', async () => {
            let data = await ibcTxHandler.getDenomRecordByChainId("osmosis_1")
            console.log(data,":<==============")
            return null
        });
    });

    describe('getProcessingTxs', () => {
        it('getProcessingTxs Test', async () => {
            jest.setTimeout(100000000)
            const txmodel = ibcTxHandler.getIbcTxModel()
            const ibcTxs = await ibcTxHandler.getProcessingTxs(txmodel,0)
            console.log(ibcTxs,'--ibcTxs--')
            let packetIdArr = ibcTxs?.length ? await ibcTxHandler.getPacketIds(ibcTxs) : [];
            console.log(packetIdArr,'--packetIdArr--')
            return null
        });
    });
})