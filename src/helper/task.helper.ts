import {Logger} from '../logger';
export async function getTaskStatus(taskModel: any,taskName): Promise<boolean>{
    let count: number = await taskModel.queryTaskCount()
    if (count <= 0) {
        Logger.log(`${taskName}: Catch-up status task suspended`)
    }
    return count == 1
}
