export enum ErrorCodes {
    failed = -1,
    success = 0,
    unauthorization = 401,
    InvalidRequest = 10001,//custom error code
    InvalidParameter = 40000,
}
// 40001,//未认证
// 40002,//参数转化异常
// 40003,//记录已存在
// 40004,//记录未找到
// 40005,//操作被拒绝
// 50001,


/**
 *
 *  key:       error code
 *  value:     description
 *
 * */

export const ResultCodesMaps = new Map<number, string>([
    [ErrorCodes.failed, 'failed'],
    [ErrorCodes.success, 'success'],
    [ErrorCodes.unauthorization, 'you have no permission to access'],
    [ErrorCodes.InvalidRequest, 'InvalidRequest'],
]);