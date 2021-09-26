import * as Path from 'path';
import * as Log4js from 'log4js';
import * as Util from 'util';
import moment from 'moment';
import * as StackTrace from 'stacktrace-js';
import Chalk from 'chalk';
import log4Config from '../config/log4js.config';
import { cfg } from '../config/config';
import { ENV, LoggerLevel } from '../constant';

export class ContextTrace {
  constructor(
        public readonly context: string,
        public readonly path?: string,
        public readonly lineNumber?: number,
        public readonly columnNumber?: number,
    ) {}
}

Log4js.addLayout('csrb-nest', (logConfig: any) => {
    return (logEvent: Log4js.LoggingEvent): string => {
        let moduleName: string = '';
        let position: string = '';
        
        const messageList: string[] = [];
        logEvent.data.forEach((value: any) => {
            if (value instanceof ContextTrace) {
                moduleName = value.context;
                
                if (value.lineNumber && value.columnNumber) {
                  position = `${value.lineNumber}, ${value.columnNumber}`;
                }
                return;
            }
            if (typeof value !== 'string') {
                value = Util.inspect(value, false, 3, true);
            }
            messageList.push(value);
        });
        
        const messageOutput: string = messageList.join(' ');
        const positionOutput: string = position ? ` [${position}]` : '';
        const typeOutput: string = `[${logConfig.type}] ${logEvent.pid.toString()}   - `;
        const dateOutput: string = `${moment(logEvent.startTime).format('YYYY-MM-DD HH:mm:ss')}`;
        const moduleOutput: string = moduleName ? `[${moduleName}] ` : '[LoggerService] ';
        let levelOutput: string = `[${logEvent.level}] ${messageOutput}`;

        switch (logEvent.level.toString()) {
            case LoggerLevel.DEBUG:
                levelOutput = Chalk.green(levelOutput);
                break;
            case LoggerLevel.INFO:
                levelOutput = Chalk.cyan(levelOutput);
                break;
            case LoggerLevel.WARN:
                levelOutput = Chalk.yellow(levelOutput);
                break;
            case LoggerLevel.ERROR:
                levelOutput = Chalk.red(levelOutput);
                break;
            case LoggerLevel.FATAL:
                levelOutput = Chalk.hex('#DD4C35')(levelOutput);
                break;
            default:
                levelOutput = Chalk.grey(levelOutput);
                break;
        }
        return `${Chalk.green(typeOutput)}${dateOutput}  ${Chalk.yellow(moduleOutput)}${levelOutput}${positionOutput}`;
    };
});

Log4js.configure(log4Config);

let logger_console = Log4js.getLogger('console');
logger_console.level = LoggerLevel.TRACE;

let logger_common = Log4js.getLogger('common');
logger_common.level = LoggerLevel.INFO;

let logger_http = Log4js.getLogger('http');
logger_http.level = LoggerLevel.INFO;

if (cfg.env == ENV.development || cfg.disableLog) {
    logger_common = logger_console;
    logger_http = logger_console;
}


export class Logger {
    static trace(...args) {
       logger_common.trace(Logger.getStackTrace(), ...args);
    }
    static debug(...args) {
       logger_common.debug(Logger.getStackTrace(), ...args);
    }
    static log(...args) {
       logger_common.log(Logger.getStackTrace(), ...args);
    }
    static info(...args) {
       logger_common.info(Logger.getStackTrace(), ...args);
    }
    static warn(...args) {
       logger_common.warn(Logger.getStackTrace(), ...args);
    }
    static error(...args) {
       logger_common.error(Logger.getStackTrace(), ...args);
    }
    static fatal(...args) {
       logger_common.fatal(Logger.getStackTrace(), ...args);
    }
    static access(...args) {
       logger_http.info(Logger.getStackTrace(), ...args);
    }

    static getStackTrace(deep: number = 2): string {
        const stackList: StackTrace.StackFrame[] = StackTrace.getSync();
        const stackInfo: StackTrace.StackFrame = stackList[deep];
        const lineNumber: number = stackInfo.lineNumber;
        const columnNumber: number = stackInfo.columnNumber;
        const fileName: string = stackInfo.fileName;
        const basename: string = Path.basename(fileName);
        return `${basename}(line: ${lineNumber}, column: ${columnNumber}): \n`;
    }
}

(<any>global).Logger = Logger;