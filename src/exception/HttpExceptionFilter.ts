import { ExceptionFilter, Catch, ArgumentsHost, HttpException, HttpStatus} from '@nestjs/common';
import {ErrorCodes} from '../api/ResultCodes';
import {ResultCodesMaps} from '../api/ResultCodes';
import { getExceptionLog } from '../helper/log.helper'

@Catch()
export class HttpExceptionFilter implements ExceptionFilter {
    catch(exception: any, host: ArgumentsHost) {
        const ctx = host.switchToHttp();
        const response = ctx.getResponse();
        const request = ctx.getRequest();
        const status = exception instanceof HttpException ? exception.getStatus() : HttpStatus.INTERNAL_SERVER_ERROR;

        const logFormat = getExceptionLog(request, exception);
        (global as any).Logger.error(logFormat);

        let code: number = exception.code || ErrorCodes.failed, message: string = ResultCodesMaps.get(exception.code) || (exception.errmsg || exception.message);
        if(exception.response && exception.response.code){
            code = exception.response.code;
            message = exception.response.message || (ResultCodesMaps.get(code) || '');
        }
        response
            .status(status)
            .json({
                code,
                message,
            });
    }
}
