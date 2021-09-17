import * as path from 'path';
import { LoggerLevel } from '../constant';
const logDir = path.resolve(__dirname, '../../logs');
const log4jsConfig = {
    appenders: {
        console: {
            type: 'console',
        },
        access: {
            type: 'dateFile', 
            filename: `${logDir}/access/access.log`,
            alwaysIncludePattern: true,
            pattern: 'yyyyMMdd',
            daysToKeep: 15,
            maxLogSize: 20971520,
            numBackups: 15,
            category: 'http',
            keepFileExt: true,
        },
        app: {
            type: 'dateFile',
            filename: `${logDir}/app-out/app.log`,
            alwaysIncludePattern: true,
            layout: {
                type: 'pattern',
                pattern: '{"date":"%d","level":"%p","category":"%c","host":"%h","pid":"%z","data":\'%m\'}',
            },
            pattern: 'yyyyMMdd',
            daysToKeep: 15,
            maxLogSize: 10485760,
            numBackups: 15,
            keepFileExt: true,
        },
        errorFile: {
            type: 'dateFile',
            filename: `${logDir}/errors/error.log`,
            alwaysIncludePattern: true,
            layout: {
                type: 'pattern',
                pattern: '{"date":"%d","level":"%p","category":"%c","host":"%h","pid":"%z","data":\'%m\'}',
            },
            pattern: 'yyyyMMdd',
            daysToKeep: 15,
            maxLogSize: 10485760,
            numBackups: 15,
            keepFileExt: true,
        },
        errors: {
          type: 'logLevelFilter',
          level: LoggerLevel.ERROR,
          appender: 'errorFile',
        },
    },
    categories: {
        default: { appenders: ['console'], level: LoggerLevel.DEBUG},
        console: { appenders: ['console'], level: LoggerLevel.TRACE},
        common: { appenders: ['console', 'app', 'errors'], level: LoggerLevel.INFO },
        http: { appenders: ['access', 'console'], level: LoggerLevel.INFO },
    }
};
export default log4jsConfig;