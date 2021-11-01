/* eslint-disable @typescript-eslint/camelcase */
import { Injectable } from '@nestjs/common';
import { Connection } from 'mongoose';
import { InjectConnection } from '@nestjs/mongoose';
import { ListStruct } from '../api/ApiResult';
import { IbcTxListReqDto, IbcTxResDto } from '../dto/ibc_tx.dto';
import { IbcDenomSchema } from '../schema/ibc_denom.schema';
import { IbcTxSchema } from '../schema/ibc_tx.schema';
import { unAuth } from '../constant';
@Injectable()
export class IbcTxService {
  private ibcDenomModel;
  private ibcTxModel;
  constructor(@InjectConnection() private connection: Connection) {
    this.getModels();
  }

  async getModels(): Promise<void> {
    this.ibcTxModel = await this.connection.model(
      'ibcTxModel',
      IbcTxSchema,
      'ex_ibc_tx',
    );
    this.ibcDenomModel = await this.connection.model(
      'ibcDenomModel',
      IbcDenomSchema,
      'ibc_denom',
    );
  }

  async queryIbcTxList(
    query: IbcTxListReqDto,
  ): Promise<ListStruct<IbcTxResDto[]> | number> {
    const { use_count, page_num, page_size, symbol, denom, start_time } = query;
    let token = undefined;
    if (symbol === unAuth) {
      const resultUnAuth = await this.ibcDenomModel.findRecordBySymbol('');
      token = resultUnAuth.map(item => {
        return {
          denom: item.denom,
          chain_id: item.chain_id
        };
      });
    } else if (symbol) {
      const result = await this.ibcDenomModel.findRecordBySymbol(symbol);
      token = result.map(item => {
        return {
          denom: item.denom,
          chain_id: item.chain_id
        };
      });
    }
    if (denom) {
      token = [denom];
    }
    if (start_time) {
      const startTx = await this.ibcTxModel.findFirstTx()
      return startTx?.tx_time;
    }
    if (use_count) {
      return await this.ibcTxModel.countTxList({ ...query, token });
    } else {
      const ibcTxDatas: IbcTxResDto[] = IbcTxResDto.bundleData(
        await this.ibcTxModel.findTxList({ ...query, token }),
      );
      return new ListStruct(ibcTxDatas, page_num, page_size);
    }
  }
}
