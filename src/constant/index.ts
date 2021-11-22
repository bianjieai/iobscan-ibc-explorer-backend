export const RecordLimit = 100;

export const Delimiter = '|';

export const unAuth = 'Others';

export enum SubState {
  // recv_packet success tx not found
  SuccessRecvPacketNotFound = 1,
  // recv_packet ack failed
  RecvPacketAckFailed = 2,
  // timeout_packet success tx not found
  SuccessTimeoutPacketNotFound = 3,
}

export const StatisticsNames = [
  'tx_24hr_all',
  'chains_24hr',
  'channels_24hr',
  'chain_all',
  'channel_all',
  'channel_opened',
  'channel_closed',
  'tx_all',
  'tx_success',
  'tx_failed',
  'base_denom_all',
  'denom_all',
];

export enum TaskEnum {
  tx = 'ex_sync_tx',
  chain = 'ex_sync_chain',
  statistics = 'ex_sync_statistics',
  denom = 'ex_sync_denom',
  nft = 'ex_sync_nft',
  txServiceName = 'sync_tx_service_name',
  validators = 'sync_validators',
  identity = 'sync_identity',
  stakingSyncValidatorsInfo = 'staking_sync_validators_info',
  stakingSyncValidatorsMoreInfo = 'staking_sync_validators_more_info',
  stakingSyncParameters = 'staking_sync_parameters',
  tokens = 'tokens',
  proposal = 'ex_sync_proposal',
  account = 'ex_sync_account',
  accountInfo = 'ex_sync_account_info',
}

export const DefaultPaging = {
  pageNum: 1,
  pageSize: 10,
};

export enum ENV {
  development = 'development',
  production = 'production',
}

export enum TxType {
  // service
  define_service = 'define_service',
  bind_service = 'bind_service',
  call_service = 'call_service',
  respond_service = 'respond_service',
  update_service_binding = 'update_service_binding',
  disable_service_binding = 'disable_service_binding',
  enable_service_binding = 'enable_service_binding',
  refund_service_deposit = 'refund_service_deposit',
  pause_request_context = 'pause_request_context',
  start_request_context = 'start_request_context',
  kill_request_context = 'kill_request_context',
  update_request_context = 'update_request_context',
  service_set_withdraw_address = 'service/set_withdraw_address',
  withdraw_earned_fees = 'withdraw_earned_fees',
  // nft
  burn_nft = 'burn_nft',
  transfer_nft = 'transfer_nft',
  edit_nft = 'edit_nft',
  issue_denom = 'issue_denom',
  mint_nft = 'mint_nft',
  // Asset
  issue_token = 'issue_token',
  edit_token = 'edit_token',
  mint_token = 'mint_token',
  transfer_token_owner = 'transfer_token_owner',
  burn_token = 'burn_token',
  //Transfer
  send = 'send',
  multisend = 'multisend',
  //Crisis
  verify_invariant = 'verify_invariant',
  //Evidence
  submit_evidence = 'submit_evidence',
  //Staking
  begin_unbonding = 'begin_unbonding',
  edit_validator = 'edit_validator',
  create_validator = 'create_validator',
  delegate = 'delegate',
  begin_redelegate = 'begin_redelegate',
  // Slashing
  unjail = 'unjail',
  // Distribution
  set_withdraw_address = 'set_withdraw_address',
  withdraw_delegator_reward = 'withdraw_delegator_reward',
  withdraw_validator_commission = 'withdraw_validator_commission',
  fund_community_pool = 'fund_community_pool',
  // Gov
  deposit = 'deposit',
  vote = 'vote',
  submit_proposal = 'submit_proposal',
  // Coinswap
  add_liquidity = 'add_liquidity',
  remove_liquidity = 'remove_liquidity',
  swap_order = 'swap_order',
  // Htlc
  create_htlc = 'create_htlc',
  claim_htlc = 'claim_htlc',
  refund_htlc = 'refund_htlc',
  // Guardian
  add_profiler = 'add_profiler',
  delete_profiler = 'delete_profiler',
  add_trustee = 'add_trustee',
  delete_trustee = 'delete_trustee',
  add_super = 'add_super',
  // Oracle
  create_feed = 'create_feed',
  start_feed = 'start_feed',
  pause_feed = 'pause_feed',
  edit_feed = 'edit_feed',
  // IBC
  transfer = 'transfer',
  recv_packet = 'recv_packet',
  create_client = 'create_client',
  update_client = 'update_client',
  timeout_packet = 'timeout_packet',
  // Identity
  create_identity = 'create_identity',
  update_identity = 'update_identity',
  // Record
  create_record = 'create_record',
  // Random
  request_rand = 'request_rand',
}

export enum TxStatus {
  SUCCESS = 1,
  FAILED = 0,
}

export enum IbcTxStatus {
  SUCCESS = 1,
  FAILED = 2,
  PROCESSING = 3,
  REFUNDED = 4,
  SETTING = 5,
}

export enum IbcTaskRecordStatus {
  OPEN = 'open',
  CLOSE = 'close'
}

export enum LoggerLevel {
  ALL = 'ALL',
  TRACE = 'TRACE',
  DEBUG = 'DEBUG',
  INFO = 'INFO',
  WARN = 'WARN',
  ERROR = 'ERROR',
  FATAL = 'FATAL',
  MARK = 'MARK',
  OFF = 'OFF',
}
