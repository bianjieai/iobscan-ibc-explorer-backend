import {TasksService} from "./task.service";
import {Test} from "@nestjs/testing";
import {AppModule} from "../app.module";
import {TaskEnum} from '../constant';
import {IbcTxTaskService} from "./ibc_tx.task.service";
import {IbcMonitorService} from "../monitor/ibc_monitor.service";

describe('TasksService', () => {
    let taskService: TasksService;
    let ibcTxTaskService: IbcTxTaskService;
    let ibcMonitorService: IbcMonitorService;
    beforeEach(async () => {
        const module = await Test.createTestingModule({
            imports: [
                AppModule
            ]
        }).compile();
        taskService = module.get<TasksService>(TasksService);
        ibcTxTaskService = module.get<IbcTxTaskService>(IbcTxTaskService);
        ibcMonitorService = module.get<IbcMonitorService>(IbcMonitorService);
    })

    describe('handleDoTask', () => {
        it('handleDoTask Test', async () => {
            await taskService.handleDoTask(TaskEnum.tx,ibcTxTaskService.doTask)
            console.log('----')
        });
    });

    describe('ibcMonitorService', () => {
        it('ibcMonitorService Test', async () => {
            await taskService.handleDoTask(TaskEnum.monitor,ibcMonitorService.doTask)
            console.log('----')
        });
    });
})