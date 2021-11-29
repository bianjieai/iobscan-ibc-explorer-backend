import {TaskDispatchService} from "../service/task.dispatch.service";
import {Test} from "@nestjs/testing";
import {AppModule} from "../app.module";
import {TaskEnum} from '../constant';


describe('TaskDispatchService', () => {
    let taskService: TaskDispatchService;
    beforeEach(async () => {
        const module = await Test.createTestingModule({
            imports: [
                AppModule
            ]
        }).compile();
        taskService = module.get<TaskDispatchService>(TaskDispatchService);
    })

    describe('needDoTask', () => {
        it('needDoTask Test', async () => {
            const taskName = TaskEnum.tx
            const isNeed = await taskService.needDoTask(taskName)
            console.log(isNeed, '----')
        });
    });
})