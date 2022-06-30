/* eslint-disable @typescript-eslint/camelcase */
import { IbcTxStatus,AllChain } from '../constant';
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
    const chains:string[] = chain_id.split(",")
    switch (chains.length) {
      case 1:// transfer_chain or recv_chain
          if (!chain_id?.includes(AllChain)) {
            const $or = [];
            $or.push({sc_chain_id: chain_id});
            $or.push({dc_chain_id: chain_id});
            queryParams.$and.push({$or});
          }
        break;
      case 2://transfer_chain and recv_chain
          if (chain_id?.includes(AllChain)) {
            if ((chains[0] === chains[1]) && (chains[0] === AllChain)) {
              // nothing to do
              if (!token) {
                delete queryParams.$and
              }
            }else{
              const index = chain_id?.indexOf(AllChain)
              if (index > 0) { //chain-id,allchain
                queryParams.$and.push({sc_chain_id:chains[0]})
              }else{ //allchain,chain-id
                queryParams.$and.push({dc_chain_id:chains[1]})
              }
            }
          }else{
            queryParams.$and.push({sc_chain_id:chains[0],dc_chain_id:chains[1]})
          }
        break;
    }
  }
  if (token && token.length) {
    // const $or = [];
    // if (token[0] && typeof token[0] === 'string') {
    //   $or.push({ 'denoms.sc_denom': { $in: token } });
    //   $or.push({ 'denoms.dc_denom': { $in: token } });
    // } else {
    //   const sc_or = {
    //     $or: token.map(item => {
    //       return {
    //         'denoms.sc_denom': item.denom,
    //         sc_chain_id: item.chain_id,
    //       };
    //     }),
    //   };
    //   const dc_or = {
    //     $or: token.map(item => {
    //       return {
    //         'denoms.dc_denom': item.denom,
    //         dc_chain_id: item.chain_id,
    //       };
    //     }),
    //   };
    //   $or.push(sc_or, dc_or);
    // }
    queryParams.$and.push({base_denom:{ $in: token}});
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
