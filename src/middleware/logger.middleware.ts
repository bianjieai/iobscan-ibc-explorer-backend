import { Injectable, NestMiddleware } from '@nestjs/common';
import { Request, Response } from 'express';
import { Logger } from '../logger';
import { getHttpRequestLog } from '../helper/log.helper';

export function LoggerMiddleware(req: Request, res: Response, next: () => any) {
  const code = res.statusCode;
  next();

  let urlTest = /^\/txs\?|^\/blocks\?|^\/statistics/;
  if (urlTest.test(req.originalUrl)) {
    return;
  }
  const logFormat = getHttpRequestLog(req, res);
  Logger.access(logFormat);
  if (code >= 500) {
    Logger.error(logFormat);
  } else if (code >= 400) {
    Logger.warn(logFormat);
  }
}
