import { LoggerService } from '@nestjs/common';
import { Logger } from './index'

export class MyLogger implements LoggerService {
  log(message: string) {
   Logger.log(message);
  }

  error(message: string, trace: string) {
   Logger.error(message, trace);
  }

  warn(message: string) {
   Logger.warn(message);
  }

  debug(message: string) {
   Logger.debug(message);
  }
}