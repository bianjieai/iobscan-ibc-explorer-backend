import {TasksService} from "./task.service";
import {Test} from "@nestjs/testing";
import {AppModule} from "../app.module";
import {TaskEnum} from '../constant';
import {IbcTxTaskService} from "./ibc_tx.task.service";

describe('TasksService', () => {
    let taskService: TasksService;
    let ibcTxTaskService: IbcTxTaskService;
    beforeEach(async () => {
        const module = await Test.createTestingModule({
            imports: [
                AppModule
            ]
        }).compile();
        taskService = module.get<TasksService>(TasksService);
        ibcTxTaskService = module.get<IbcTxTaskService>(IbcTxTaskService);
    })

    describe('handleDoTask', () => {
        it('handleDoTask Test', async () => {
            await taskService.handleDoTask(TaskEnum.tx,ibcTxTaskService.doTask)
            console.log('----')
        });
    });
})