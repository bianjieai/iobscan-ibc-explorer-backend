import {IbcMonitorService} from "../monitor/ibc_monitor.service";
import {Test} from "@nestjs/testing";
import {AppModule} from "../app.module";

describe('IbcMonitorService', () => {
    let ibcMonitorService: IbcMonitorService;
    beforeEach(async () => {
        const module = await Test.createTestingModule({
            imports: [
                AppModule
            ]
        }).compile();
        ibcMonitorService = module.get<IbcMonitorService>(IbcMonitorService);
    })




    describe('getProcessingCnt', () => {
        it('getProcessingCnt Test', async () => {
            const result = await  ibcMonitorService.getProcessingCnt()
            console.log(result, '----')
        });
    });

    describe('getNodeInfo', () => {
        it('getNodeInfo Test', async () => {
            const result = await  ibcMonitorService.getNodeInfo("https://cosmoshub.stakesystems.io/","cosmoshub_4")
            console.log(result, '----')
        });
    });

})
