import {Module } from '@nestjs/common';
import { MongooseModule } from '@nestjs/mongoose';
import { TaskDispatchSchema } from '../schema/task.dispatch.schema';
import { TaskDispatchService } from '../service/task.dispatch.service';

@Module({
    imports:[
        MongooseModule.forFeature([{
            name: 'TaskDispatch',
            schema: TaskDispatchSchema,
            collection: 'ex_task_dispatch'
        }])
    ],
    providers:[TaskDispatchService],
    exports: [TaskDispatchService]
})
export class TaskDispatchModule{}