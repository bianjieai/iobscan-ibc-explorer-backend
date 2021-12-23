import { CallHandler, ExecutionContext, Injectable, NestInterceptor } from '@nestjs/common';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { Logger } from '../logger';
import { getHttpRespondLog } from '../helper/log.helper'
@Injectable()
export class LoggerInterceptor implements NestInterceptor {
  	intercept(context: ExecutionContext, next: CallHandler): Observable<any> {
	    const req = context.getArgByIndex(1).req;
	    let urlTest = /^\/txs\?|^\/blocks\?|^\/statistics/;
	    return next.handle().pipe(
		    map(data => {
		        const 	logFormat = getHttpRespondLog(req, data);
		        if (!urlTest.test(req.originalUrl)) {Logger.access(logFormat)}
		        return data;
		    }),
	    );
  	}
}