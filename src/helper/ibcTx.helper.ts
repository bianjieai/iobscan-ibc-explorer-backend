/* eslint-disable @typescript-eslint/camelcase */
import { IbcTxStatus } from '../constant';
import { IbcTxQueryType } from '../types/schemaTypes/ibc_tx.interface';

interface IbcTxQueryParamsType {
  tx_time?: { $gte?; $lte? };
  status?: number | { $in: number[] };
  $and?: any[];
}

const parseQuery = (query: IbcTxQueryType): IbcTxQueryParamsType => {
  const { chain_id, status, token, date_range } = query;
  const queryParams: IbcTxQueryParamsType = {};
  if ((date_range && date_range[0] > 0) || (date_range && date_range[1] > 0)) {
    queryParams.tx_time = {};
  }
  if (date_range && date_range[0] > 0) {
    queryParams.tx_time.$gte = date_range[0];
  }
  if (date_range && date_range[1] > 0) {
    queryParams.tx_time.$lte = date_range[1];
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
  if (token && token.length) {
    const $or = [];
    $or.push({ 'denoms.sc_denom': { $in: token }});
    $or.push({ 'denoms.dc_denom': { $in: token } });
    queryParams.$and.push({ $or });
  }
  if (status) {
    queryParams.status = {
      $in: status,
    };
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
