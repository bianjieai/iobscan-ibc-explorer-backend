/* eslint-disable @typescript-eslint/camelcase */
import { IbcTxStatus } from '../constant';
import { IbcTxQueryType } from '../types/schemaTypes/ibc_tx.interface';

interface IbcTxQueryParamsType {
  tx_time?: { $gte?; $lte? };
  status?: number | { $in: number[] };
  $and?: any[];
}

const parseQuery = (query: IbcTxQueryType): IbcTxQueryParamsType => {
  const { beginTime, endTime, chain_id, status, token } = query;
  const queryParams: IbcTxQueryParamsType = {};
  if ((beginTime && beginTime > 0) || (endTime && endTime > 0)) {
    queryParams.tx_time = {};
  }
  if (beginTime && beginTime > 0) {
    queryParams.tx_time.$gte = beginTime;
  }
  if (endTime && endTime > 0) {
    queryParams.tx_time.$lte = endTime;
  }
  if (chain_id || token) {
    queryParams.$and = [];
  }
  if (chain_id) {
    const $or = [];
    $or.push({ sc_chain_id: chain_id });
    $or.push({ dc_chain_id: chain_id });
    queryParams.$and.push({ $or });
  }
  if (token) {
    const $or = [];
    $or.push({ 'denoms.sc_denom': token });
    $or.push({ 'denoms.dc_denom': token });
    queryParams.$and.push({ $or });
  }
  if (status) {
    queryParams.status = status;
  } else {
    queryParams.status = {
      $in: [
        IbcTxStatus.SUCCESS,
        IbcTxStatus.FAILED,
        IbcTxStatus.PROCESSING,
        IbcTxStatus.REFUNDED,
      ],
    };
  }
  return queryParams;
};

export { parseQuery };
