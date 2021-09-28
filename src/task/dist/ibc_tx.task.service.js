"use strict";
var __assign = (this && this.__assign) || function () {
    __assign = Object.assign || function(t) {
        for (var s, i = 1, n = arguments.length; i < n; i++) {
            s = arguments[i];
            for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
                t[p] = s[p];
        }
        return t;
    };
    return __assign.apply(this, arguments);
};
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __param = (this && this.__param) || function (paramIndex, decorator) {
    return function (target, key) { decorator(target, key, paramIndex); }
};
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __generator = (this && this.__generator) || function (thisArg, body) {
    var _ = { label: 0, sent: function() { if (t[0] & 1) throw t[1]; return t[1]; }, trys: [], ops: [] }, f, y, t, g;
    return g = { next: verb(0), "throw": verb(1), "return": verb(2) }, typeof Symbol === "function" && (g[Symbol.iterator] = function() { return this; }), g;
    function verb(n) { return function (v) { return step([n, v]); }; }
    function step(op) {
        if (f) throw new TypeError("Generator is already executing.");
        while (_) try {
            if (f = 1, y && (t = op[0] & 2 ? y["return"] : op[0] ? y["throw"] || ((t = y["return"]) && t.call(y), 0) : y.next) && !(t = t.call(y, op[1])).done) return t;
            if (y = 0, t) op = [op[0] & 2, t.value];
            switch (op[0]) {
                case 0: case 1: t = op; break;
                case 4: _.label++; return { value: op[1], done: false };
                case 5: _.label++; y = op[1]; op = [0]; continue;
                case 7: op = _.ops.pop(); _.trys.pop(); continue;
                default:
                    if (!(t = _.trys, t = t.length > 0 && t[t.length - 1]) && (op[0] === 6 || op[0] === 2)) { _ = 0; continue; }
                    if (op[0] === 3 && (!t || (op[1] > t[0] && op[1] < t[3]))) { _.label = op[1]; break; }
                    if (op[0] === 6 && _.label < t[1]) { _.label = t[1]; t = op; break; }
                    if (t && _.label < t[2]) { _.label = t[2]; _.ops.push(op); break; }
                    if (t[2]) _.ops.pop();
                    _.trys.pop(); continue;
            }
            op = body.call(thisArg, _);
        } catch (e) { op = [6, e]; y = 0; } finally { f = t = 0; }
        if (op[0] & 5) throw op[1]; return { value: op[0] ? op[1] : void 0, done: true };
    }
};
var __spreadArrays = (this && this.__spreadArrays) || function () {
    for (var s = 0, i = 0, il = arguments.length; i < il; i++) s += arguments[i].length;
    for (var r = Array(s), k = 0, i = 0; i < il; i++)
        for (var a = arguments[i], j = 0, jl = a.length; j < jl; j++, k++)
            r[k] = a[j];
    return r;
};
exports.__esModule = true;
exports.IbcTxTaskService = void 0;
var common_1 = require("@nestjs/common");
var mongoose_1 = require("@nestjs/mongoose");
var ibc_chain_config_schema_1 = require("../schema/ibc_chain_config.schema");
var ibc_chain_schema_1 = require("../schema/ibc_chain.schema");
var ibc_denom_schema_1 = require("../schema/ibc_denom.schema");
var ibc_tx_schema_1 = require("../schema/ibc_tx.schema");
var tx_schema_1 = require("../schema/tx.schema");
var ibc_block_schema_1 = require("../schema/ibc_block.schema");
var ibc_task_record_schema_1 = require("../schema/ibc_task_record.schema");
var ibc_channel_schema_1 = require("src/schema/ibc_channel.schema");
var util_1 = require("../util/util");
var denom_helper_1 = require("../helper/denom.helper");
var constant_1 = require("../constant");
var IbcTxTaskService = /** @class */ (function () {
    function IbcTxTaskService(connection) {
        this.connection = connection;
        this.getModels();
        this.doTask = this.doTask.bind(this);
    }
    IbcTxTaskService.prototype.doTask = function (taskName) {
        return __awaiter(this, void 0, Promise, function () {
            var dateNow;
            return __generator(this, function (_a) {
                dateNow = String(Math.floor(new Date().getTime() / 1000));
                this.parseIbcTx(dateNow);
                this.changeIbcTxState(dateNow);
                return [2 /*return*/];
            });
        });
    };
    // getModels
    IbcTxTaskService.prototype.getModels = function () {
        return __awaiter(this, void 0, Promise, function () {
            var _a, _b, _c, _d, _e, _f;
            return __generator(this, function (_g) {
                switch (_g.label) {
                    case 0:
                        // ibcTaskRecordModel
                        _a = this;
                        return [4 /*yield*/, this.connection.model('ibcTaskRecordModel', ibc_task_record_schema_1.IbcTaskRecordSchema, 'ibc_task_record')];
                    case 1:
                        // ibcTaskRecordModel
                        _a.ibcTaskRecordModel = _g.sent();
                        // chainConfigModel
                        _b = this;
                        return [4 /*yield*/, this.connection.model('chainConfigModel', ibc_chain_config_schema_1.IbcChainConfigSchema, 'chain_config')];
                    case 2:
                        // chainConfigModel
                        _b.chainConfigModel = _g.sent();
                        // ibcChainModel
                        _c = this;
                        return [4 /*yield*/, this.connection.model('ibcChainModel', ibc_chain_schema_1.IbcChainSchema, 'ibc_chain')];
                    case 3:
                        // ibcChainModel
                        _c.ibcChainModel = _g.sent();
                        // ibcTxModel
                        _d = this;
                        return [4 /*yield*/, this.connection.model('ibcTxModel', ibc_tx_schema_1.IbcTxSchema, 'ex_ibc_tx')];
                    case 4:
                        // ibcTxModel
                        _d.ibcTxModel = _g.sent();
                        // ibcDenomModel
                        _e = this;
                        return [4 /*yield*/, this.connection.model('ibcDenomModel', ibc_denom_schema_1.IbcDenomSchema, 'ibc_denom')];
                    case 5:
                        // ibcDenomModel
                        _e.ibcDenomModel = _g.sent();
                        // ibcChannelModel
                        _f = this;
                        return [4 /*yield*/, this.connection.model('ibcChannelModel', ibc_channel_schema_1.IbcChannelSchema, 'ibc_channel')];
                    case 6:
                        // ibcChannelModel
                        _f.ibcChannelModel = _g.sent();
                        return [2 /*return*/];
                }
            });
        });
    };
    // ibcTx first（transfer）
    IbcTxTaskService.prototype.parseIbcTx = function (dateNow) {
        return __awaiter(this, void 0, Promise, function () {
            var allChains;
            var _this = this;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, this.chainConfigModel.findAll()];
                    case 1:
                        allChains = _a.sent();
                        allChains.forEach(function (_a) {
                            var chain_id = _a.chain_id;
                            return __awaiter(_this, void 0, void 0, function () {
                                var taskRecord, txModel, txs, txsByLimit, txsByHeight, _b, hash;
                                var _this = this;
                                return __generator(this, function (_c) {
                                    switch (_c.label) {
                                        case 0: return [4 /*yield*/, this.ibcTaskRecordModel.findTaskRecord(chain_id)];
                                        case 1:
                                            taskRecord = _c.sent();
                                            if (!!taskRecord) return [3 /*break*/, 4];
                                            return [4 /*yield*/, this.ibcTaskRecordModel.insertManyTaskRecord({
                                                    task_name: "sync_" + chain_id + "_transfer",
                                                    status: constant_1.IbcTaskRecordStatus.OPEN,
                                                    height: 0,
                                                    create_at: "" + dateNow,
                                                    update_at: "" + dateNow
                                                })];
                                        case 2:
                                            _c.sent();
                                            return [4 /*yield*/, this.ibcTaskRecordModel.findTaskRecord(chain_id)];
                                        case 3:
                                            taskRecord = _c.sent();
                                            return [3 /*break*/, 5];
                                        case 4:
                                            if (taskRecord.status === constant_1.IbcTaskRecordStatus.CLOSE)
                                                return [2 /*return*/];
                                            _c.label = 5;
                                        case 5: return [4 /*yield*/, this.connection.model('txModel', tx_schema_1.TxSchema, "sync_" + chain_id + "_tx")];
                                        case 6:
                                            txModel = _c.sent();
                                            txs = [];
                                            return [4 /*yield*/, txModel.queryTxListSortHeight({
                                                    type: constant_1.TxType.transfer,
                                                    height: taskRecord.height,
                                                    limit: constant_1.RecordLimit
                                                })];
                                        case 7:
                                            txsByLimit = _c.sent();
                                            if (!txsByLimit.length) return [3 /*break*/, 9];
                                            return [4 /*yield*/, txModel.queryTxListByHeight(constant_1.TxType.transfer, txsByLimit[txsByLimit.length - 1].height)];
                                        case 8:
                                            _b = _c.sent();
                                            return [3 /*break*/, 10];
                                        case 9:
                                            _b = [];
                                            _c.label = 10;
                                        case 10:
                                            txsByHeight = _b;
                                            hash = {};
                                            txs = __spreadArrays(txsByLimit, txsByHeight).reduce(function (txsResult, next) {
                                                hash[next.tx_hash]
                                                    ? ''
                                                    : (hash[next.tx_hash] = true) && txsResult.push(next);
                                                return txsResult;
                                            }, []);
                                            txs.forEach(function (tx, txIndex) {
                                                var height = tx.height;
                                                var log = tx.log;
                                                var time = tx.time;
                                                var hash = tx.tx_hash;
                                                var status = tx.status;
                                                var fee = tx.fee;
                                                var update_at = '';
                                                tx.msgs.forEach(function (msg, msgIndex) { return __awaiter(_this, void 0, void 0, function () {
                                                    var ibcTx_1, sc_chain_id_1, sc_port_1, sc_channel_1, sc_addr, dc_addr, sc_denom_1, msg_amount, _a, dc_port_1, dc_channel_1, sequence, base_denom, denom_path_1, dc_chain_id_1, result;
                                                    var _this = this;
                                                    return __generator(this, function (_b) {
                                                        switch (_b.label) {
                                                            case 0:
                                                                if (!(msg.type === constant_1.TxType.transfer)) return [3 /*break*/, 3];
                                                                ibcTx_1 = {
                                                                    record_id: '',
                                                                    sc_addr: '',
                                                                    dc_addr: '',
                                                                    sc_port: '',
                                                                    sc_channel: '',
                                                                    sc_chain_id: '',
                                                                    dc_port: '',
                                                                    dc_channel: '',
                                                                    dc_chain_id: '',
                                                                    sequence: '',
                                                                    status: 0,
                                                                    sc_tx_info: {},
                                                                    dc_tx_info: {},
                                                                    refunded_tx_info: {},
                                                                    log: {},
                                                                    denoms: [],
                                                                    base_denom: '',
                                                                    create_at: '',
                                                                    update_at: ''
                                                                };
                                                                switch (tx.status) {
                                                                    case constant_1.TxStatus.SUCCESS:
                                                                        ibcTx_1.status = constant_1.IbcTxStatus.PROCESSING;
                                                                        break;
                                                                    case constant_1.TxStatus.FAILED:
                                                                        ibcTx_1.status = constant_1.IbcTxStatus.FAILED;
                                                                        break;
                                                                    default:
                                                                        break;
                                                                }
                                                                sc_chain_id_1 = chain_id;
                                                                sc_port_1 = msg.msg.source_port;
                                                                sc_channel_1 = msg.msg.source_channel;
                                                                sc_addr = msg.msg.sender;
                                                                dc_addr = msg.msg.receiver;
                                                                sc_denom_1 = msg.msg.token.denom;
                                                                msg_amount = msg.msg.token;
                                                                _a = this.getIbcInfoFromEventsMsg(tx, msgIndex), dc_port_1 = _a.dc_port, dc_channel_1 = _a.dc_channel, sequence = _a.sequence, base_denom = _a.base_denom, denom_path_1 = _a.denom_path;
                                                                dc_chain_id_1 = '';
                                                                return [4 /*yield*/, this.chainConfigModel.findDcChain({
                                                                        sc_chain_id: sc_chain_id_1,
                                                                        sc_port: sc_port_1,
                                                                        sc_channel: sc_channel_1,
                                                                        dc_port: dc_port_1,
                                                                        dc_channel: dc_channel_1
                                                                    })];
                                                            case 1:
                                                                result = _b.sent();
                                                                if (result && result.ibc_info && result.ibc_info.length) {
                                                                    result.ibc_info.forEach(function (info_item) {
                                                                        info_item.paths.forEach(function (path_item) {
                                                                            if (path_item.channel_id === sc_channel_1 &&
                                                                                path_item.port_id === sc_port_1 &&
                                                                                path_item.counterparty.channel_id === dc_channel_1 &&
                                                                                path_item.counterparty.port_id === dc_port_1) {
                                                                                dc_chain_id_1 = info_item.chain_id;
                                                                            }
                                                                        });
                                                                    });
                                                                }
                                                                else {
                                                                    dc_chain_id_1 = '';
                                                                }
                                                                ibcTx_1.record_id = "" + sc_port_1 + sc_channel_1 + dc_port_1 + dc_channel_1 + sequence + sc_chain_id_1;
                                                                ibcTx_1.sc_addr = sc_addr;
                                                                ibcTx_1.dc_addr = dc_addr;
                                                                ibcTx_1.sc_port = sc_port_1;
                                                                ibcTx_1.sc_channel = sc_channel_1;
                                                                ibcTx_1.sc_chain_id = sc_chain_id_1;
                                                                ibcTx_1.dc_port = dc_port_1;
                                                                ibcTx_1.dc_channel = dc_channel_1;
                                                                ibcTx_1.dc_chain_id = dc_chain_id_1;
                                                                ibcTx_1.sequence = sequence;
                                                                ibcTx_1.denoms.push(sc_denom_1);
                                                                ibcTx_1.base_denom = base_denom;
                                                                ibcTx_1.create_at = dateNow;
                                                                ibcTx_1.update_at = tx.time;
                                                                ibcTx_1.sc_tx_info = {
                                                                    hash: hash,
                                                                    status: status,
                                                                    time: time,
                                                                    height: height,
                                                                    fee: fee,
                                                                    msg_amount: msg_amount,
                                                                    msg: msg
                                                                };
                                                                ibcTx_1.log['sc_log'] = log;
                                                                if (!dc_chain_id_1 && ibcTx_1.status !== constant_1.IbcTxStatus.FAILED) {
                                                                    ibcTx_1.status = constant_1.IbcTxStatus.SETTING;
                                                                }
                                                                return [4 /*yield*/, this.ibcTxModel.insertManyIbcTx(ibcTx_1, function (err) { return __awaiter(_this, void 0, void 0, function () {
                                                                        return __generator(this, function (_a) {
                                                                            switch (_a.label) {
                                                                                case 0:
                                                                                    taskRecord.height = height;
                                                                                    return [4 /*yield*/, this.ibcTaskRecordModel.updateTaskRecord(taskRecord)];
                                                                                case 1:
                                                                                    _a.sent();
                                                                                    if (ibcTx_1.status !== constant_1.IbcTxStatus.FAILED) {
                                                                                        // parse denom
                                                                                        this.parseDenom(ibcTx_1.sc_chain_id, sc_denom_1, ibcTx_1.base_denom, denom_path_1, !Boolean(denom_path_1), dateNow, dateNow, dateNow);
                                                                                        // parse channel
                                                                                        this.parseChannel(sc_chain_id_1, dc_chain_id_1, sc_channel_1, dateNow);
                                                                                        // parse chain
                                                                                        this.parseChain(sc_chain_id_1, dateNow);
                                                                                    }
                                                                                    return [2 /*return*/];
                                                                            }
                                                                        });
                                                                    }); })];
                                                            case 2:
                                                                _b.sent();
                                                                _b.label = 3;
                                                            case 3: return [2 /*return*/];
                                                        }
                                                    });
                                                }); });
                                            });
                                            return [2 /*return*/];
                                    }
                                });
                            });
                        });
                        return [2 /*return*/];
                }
            });
        });
    };
    // ibcTx second（recv_packet || timoout_packet）
    IbcTxTaskService.prototype.changeIbcTxState = function (dateNow) {
        return __awaiter(this, void 0, Promise, function () {
            var ibcTxs;
            var _this = this;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, this.ibcTxModel.queryTxList({
                            status: constant_1.IbcTxStatus.PROCESSING,
                            limit: constant_1.RecordLimit
                        })];
                    case 1:
                        ibcTxs = _a.sent();
                        ibcTxs.forEach(function (ibcTx) { return __awaiter(_this, void 0, void 0, function () {
                            var txModel, txs, counter_party_tx_1, blockModel, _a, height, time, ibcHeight, ibcTime, txModel_1, refunded_tx_1;
                            var _this = this;
                            return __generator(this, function (_b) {
                                switch (_b.label) {
                                    case 0:
                                        if (!ibcTx.dc_chain_id)
                                            return [2 /*return*/];
                                        return [4 /*yield*/, this.connection.model('txModel', tx_schema_1.TxSchema, "sync_" + ibcTx.dc_chain_id + "_tx")];
                                    case 1:
                                        txModel = _b.sent();
                                        return [4 /*yield*/, txModel.queryTxListByPacketId({
                                                type: constant_1.TxType.recv_packet,
                                                limit: constant_1.RecordLimit,
                                                status: constant_1.TxStatus.SUCCESS,
                                                packet_id: ibcTx.sc_tx_info.msg.msg.packet_id
                                            })];
                                    case 2:
                                        txs = _b.sent();
                                        if (!txs.length) return [3 /*break*/, 3];
                                        counter_party_tx_1 = txs[0];
                                        counter_party_tx_1 &&
                                            counter_party_tx_1.msgs.forEach(function (msg) {
                                                if (msg.type === constant_1.TxType.recv_packet &&
                                                    ibcTx.sc_tx_info.msg.msg.packet_id === msg.msg.packet_id) {
                                                    var _a = denom_helper_1.getDcDenom(msg), dc_denom = _a.dc_denom, dc_denom_origin = _a.dc_denom_origin;
                                                    ibcTx.status = constant_1.IbcTxStatus.SUCCESS;
                                                    ibcTx.dc_tx_info = {
                                                        hash: counter_party_tx_1.tx_hash,
                                                        status: counter_party_tx_1.status,
                                                        time: counter_party_tx_1.time,
                                                        height: counter_party_tx_1.height,
                                                        fee: counter_party_tx_1.fee,
                                                        msg_amount: msg.msg.token,
                                                        msg: msg
                                                    };
                                                    ibcTx.update_at = counter_party_tx_1.time;
                                                    ibcTx.denoms.push(dc_denom);
                                                    var denom_path = dc_denom_origin.replace("/" + ibcTx.base_denom, '');
                                                    _this.ibcTxModel.updateIbcTx(ibcTx);
                                                    // parse denom
                                                    _this.parseDenom(ibcTx.dc_chain_id, dc_denom, ibcTx.base_denom, denom_path, !Boolean(denom_path), dateNow, dateNow, dateNow);
                                                    // parse Channel
                                                    _this.parseChannel(ibcTx.sc_chain_id, ibcTx.dc_chain_id, ibcTx.dc_channel, dateNow);
                                                    // parse Chain
                                                    _this.parseChain(ibcTx.dc_channel, dateNow);
                                                }
                                            });
                                        return [3 /*break*/, 8];
                                    case 3: return [4 /*yield*/, this.connection.model('blockModel', ibc_block_schema_1.IbcBlockSchema, "sync_" + ibcTx.dc_chain_id + "_block")];
                                    case 4:
                                        blockModel = _b.sent();
                                        return [4 /*yield*/, blockModel.findLatestBlock()];
                                    case 5:
                                        _a = _b.sent(), height = _a.height, time = _a.time;
                                        ibcHeight = ibcTx.sc_tx_info.msg.msg.timeout_height.revision_height;
                                        ibcTime = ibcTx.sc_tx_info.msg.msg.timeout_timestamp;
                                        if (!(ibcHeight < height || ibcTime < time)) return [3 /*break*/, 8];
                                        return [4 /*yield*/, this.connection.model('txModel', tx_schema_1.TxSchema, "sync_" + ibcTx.sc_chain_id + "_tx")];
                                    case 6:
                                        txModel_1 = _b.sent();
                                        return [4 /*yield*/, txModel_1.queryTxListByPacketId({
                                                type: constant_1.TxType.timeout_packet,
                                                limit: constant_1.RecordLimit,
                                                status: constant_1.TxStatus.SUCCESS,
                                                packet_id: ibcTx.sc_tx_info.msg.msg.packet_id
                                            })[0]];
                                    case 7:
                                        refunded_tx_1 = _b.sent();
                                        refunded_tx_1 &&
                                            refunded_tx_1.msgs.forEach(function (msg) {
                                                if (msg.type === constant_1.TxType.timeout_packet &&
                                                    ibcTx.sc_tx_info.msg.msg.packet_id === msg.msg.packet_id) {
                                                    ibcTx.status = constant_1.IbcTxStatus.REFUNDED;
                                                    ibcTx.refunded_tx_info = {
                                                        hash: refunded_tx_1.tx_hash,
                                                        status: refunded_tx_1.status,
                                                        time: refunded_tx_1.time,
                                                        height: refunded_tx_1.height,
                                                        fee: refunded_tx_1.fee,
                                                        msg_amount: msg.msg.token,
                                                        msg: msg
                                                    };
                                                    ibcTx.update_at = refunded_tx_1.time;
                                                    _this.ibcTxModel.updateIbcTx(ibcTx);
                                                }
                                            });
                                        _b.label = 8;
                                    case 8: return [2 /*return*/];
                                }
                            });
                        }); });
                        return [2 /*return*/];
                }
            });
        });
    };
    // get dc_port、dc_channel、sequence
    IbcTxTaskService.prototype.getIbcInfoFromEventsMsg = function (tx, msgIndex) {
        var msg = {
            dc_port: '',
            dc_channel: '',
            sequence: '',
            base_denom: '',
            denom_path: ''
        };
        tx.events_new[msgIndex] &&
            tx.events_new[msgIndex].events.forEach(function (evt) {
                if (evt.type === 'send_packet') {
                    evt.attributes.forEach(function (attr) {
                        switch (attr.key) {
                            case 'packet_dst_port':
                                msg.dc_port = attr.value;
                                break;
                            case 'packet_dst_channel':
                                msg.dc_channel = attr.value;
                                break;
                            case 'packet_sequence':
                                msg.sequence = attr.value;
                                break;
                            case 'packet_data':
                                var packet_data = util_1.JSONparse(attr.value);
                                var denomOrigin = packet_data.denom;
                                var denomOriginSplit = denomOrigin.split('/');
                                msg.base_denom = denomOriginSplit[denomOriginSplit.length - 1];
                                msg.denom_path = denomOriginSplit
                                    .slice(0, denomOriginSplit.length - 1)
                                    .join('/');
                            default:
                                break;
                        }
                    });
                }
            });
        return msg;
    };
    // parse Denom
    IbcTxTaskService.prototype.parseDenom = function (chain_id, denom, base_denom, denom_path, is_source_chain, create_at, update_at, dateNow) {
        return __awaiter(this, void 0, Promise, function () {
            var ibcDenomRecord, ibcDenom;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, this.ibcDenomModel.findDenomRecord(chain_id, denom)];
                    case 1:
                        ibcDenomRecord = _a.sent();
                        if (!!ibcDenomRecord) return [3 /*break*/, 3];
                        ibcDenom = {
                            chain_id: chain_id,
                            denom: denom,
                            base_denom: base_denom,
                            denom_path: denom_path,
                            is_source_chain: is_source_chain,
                            create_at: create_at,
                            update_at: update_at
                        };
                        return [4 /*yield*/, this.ibcDenomModel.insertManyDenom(ibcDenom)];
                    case 2:
                        _a.sent();
                        return [3 /*break*/, 5];
                    case 3:
                        ibcDenomRecord.update_at = dateNow;
                        return [4 /*yield*/, this.ibcDenomModel.updateDenomRecord(ibcDenomRecord)];
                    case 4:
                        _a.sent();
                        _a.label = 5;
                    case 5: return [2 /*return*/];
                }
            });
        });
    };
    // parse Channel
    IbcTxTaskService.prototype.parseChannel = function (sc_chain_id, dc_chain_id, channel_id, dateNow) {
        return __awaiter(this, void 0, Promise, function () {
            var channels_all_record, isFindRecord, ibcChannelRecord, ibcChannel;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, this.getChannelsConfig()];
                    case 1:
                        channels_all_record = _a.sent();
                        isFindRecord = channels_all_record.find(function (channel) {
                            return channel.record_id === "" + sc_chain_id + dc_chain_id + channel_id;
                        });
                        if (!isFindRecord)
                            return [2 /*return*/];
                        return [4 /*yield*/, this.ibcChannelModel.findChannelRecord("" + sc_chain_id + dc_chain_id + channel_id)];
                    case 2:
                        ibcChannelRecord = _a.sent();
                        if (!!ibcChannelRecord) return [3 /*break*/, 4];
                        ibcChannel = __assign(__assign({}, isFindRecord), { update_at: dateNow, create_at: dateNow });
                        return [4 /*yield*/, this.ibcChannelModel.insertManyChannel(ibcChannel)];
                    case 3:
                        _a.sent();
                        return [3 /*break*/, 6];
                    case 4:
                        ibcChannelRecord.update_at = dateNow;
                        return [4 /*yield*/, this.ibcChannelModel.updateChannelRecord(ibcChannelRecord)];
                    case 5:
                        _a.sent();
                        _a.label = 6;
                    case 6: return [2 /*return*/];
                }
            });
        });
    };
    // parse Chain
    IbcTxTaskService.prototype.parseChain = function (chain_id, dateNow) {
        return __awaiter(this, void 0, void 0, function () {
            var ibcChainRecord, allChainsConfig, findChainConfig, ibcChain;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0: return [4 /*yield*/, this.ibcChainModel.findById(chain_id)];
                    case 1:
                        ibcChainRecord = _a.sent();
                        if (!!ibcChainRecord) return [3 /*break*/, 3];
                        return [4 /*yield*/, this.chainConfigModel.findAll()];
                    case 2:
                        allChainsConfig = _a.sent();
                        findChainConfig = allChainsConfig.find(function (chainConfig) {
                            return chainConfig.chain_id === chain_id;
                        });
                        if (!findChainConfig)
                            return [2 /*return*/];
                        ibcChain = {
                            chain_id: chain_id,
                            chain_name: findChainConfig ? findChainConfig.chain_name : '',
                            icon: findChainConfig ? findChainConfig.icon : '',
                            create_at: dateNow,
                            update_at: dateNow
                        };
                        this.ibcChainModel.insertManyChain(ibcChain);
                        return [3 /*break*/, 4];
                    case 3:
                        ibcChainRecord.update_at = dateNow;
                        this.ibcChainModel.updateChainRecord(ibcChainRecord);
                        _a.label = 4;
                    case 4: return [2 /*return*/];
                }
            });
        });
    };
    // get configed channels
    IbcTxTaskService.prototype.getChannelsConfig = function () {
        return __awaiter(this, void 0, void 0, function () {
            var channels_all_record, allChains;
            return __generator(this, function (_a) {
                switch (_a.label) {
                    case 0:
                        channels_all_record = [];
                        return [4 /*yield*/, this.chainConfigModel.findAll()];
                    case 1:
                        allChains = _a.sent();
                        allChains.forEach(function (chain) {
                            chain.ibc_info.forEach(function (ibc_info_item) {
                                ibc_info_item.paths.forEach(function (path_item) {
                                    channels_all_record.push({
                                        channel_id: path_item.channel_id,
                                        record_id: "" + chain.chain_id + ibc_info_item.chain_id + path_item.channel_id,
                                        state: path_item.state
                                    });
                                });
                            });
                        });
                        return [2 /*return*/, channels_all_record];
                }
            });
        });
    };
    IbcTxTaskService = __decorate([
        common_1.Injectable(),
        __param(0, mongoose_1.InjectConnection())
    ], IbcTxTaskService);
    return IbcTxTaskService;
}());
exports.IbcTxTaskService = IbcTxTaskService;
