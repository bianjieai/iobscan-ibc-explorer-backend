import { HttpStatus} from '@nestjs/common';

export function getHttpRequestLog(req, res):string{
    let logFormat: string = ` Request original url: ${req.originalUrl}
                              Method: ${req.method}
                              IP: ${req.ip}
                              Status code: ${res.statusCode}
                              Parmas: ${JSON.stringify(req.params)}
                              Query: ${JSON.stringify(req.query)}
                              Body: ${JSON.stringify(req.body)} \n`;
    return logFormat;
}

export function getHttpRespondLog(req, resData):string{
    let resBodyStr: string = JSON.stringify(resData.data) || '';
    let logFormat: string = ` Request original url: ${req.originalUrl}
                            Method: ${req.method}
                            IP: ${req.ip}
                            User: ${JSON.stringify(req.user || '{}')}
                            Response data:\n ${resBodyStr.length>1000 ? resBodyStr.substr(0,1000) : resBodyStr} \n `;
    return logFormat;
}

export function getExceptionLog(req, exception):string{
    const status: string = typeof exception.getStatus === 'function' ? exception.getStatus() : HttpStatus.INTERNAL_SERVER_ERROR;
    const logFormat: string = `
                        Request original url: ${req.originalUrl}
                        Method: ${req.method}
                        IP: ${req.ip}
                        Status code: ${status}
                        Parmas: ${JSON.stringify(req.params)}
                        Query: ${JSON.stringify(req.query)}
                        Body: ${JSON.stringify(req.body)}
                        Response: ${exception.stack} \n`;
    return logFormat;
}
