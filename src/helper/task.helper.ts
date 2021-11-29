import {Logger} from '../logger';
export async function getTaskStatus(chainId:any,taskModel: any,taskName): Promise<boolean>{
    let count: number = await taskModel.queryTaskCount()
    console.log(chainId,count,`${chainId} task count is ${count}`)
    if (count <= 0) {
        Logger.log(`${taskName}: Catch-up status task suspended`)
    }
    return count == 1
}
