import {TasksService} from "./task.service";
import {Test} from "@nestjs/testing";
import {AppModule} from "../app.module";
import {TaskEnum} from '../constant';
import {IbcMonitorService} from "../monitor/ibc_monitor.service";

describe('TasksService', () => {
    let taskService: TasksService;
    let ibcMonitorService: IbcMonitorService;
    beforeEach(async () => {
        const module = await Test.createTestingModule({
            imports: [
                AppModule
            ]
        }).compile();
        taskService = module.get<TasksService>(TasksService);
        ibcMonitorService = module.get<IbcMonitorService>(IbcMonitorService);
    })


    describe('ibcMonitorService', () => {
        it('ibcMonitorService Test', async () => {
            await taskService.handleDoTask(TaskEnum.monitor,ibcMonitorService.doTask)
            console.log('----')
        });
    });
})