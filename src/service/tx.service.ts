import { Injectable } from '@nestjs/common';
import { Connection } from 'mongoose';
import { InjectConnection } from '@nestjs/mongoose';
import { TxSchema } from '../schema/tx.schema';

@Injectable()
export class TxService {
  constructor(@InjectConnection() private connection: Connection) {}

  async getTxModel(chain_id) {
    const txModel = await this.connection.model(
      'txModel',
      TxSchema,
      `sync_${chain_id}_tx`,
    );
    return txModel;
  }

  // txs
  async queryTxList(chain_id, query) {
    const txModel = await this.getTxModel(chain_id);
    return await txModel.queryTxList(query);
  }

}
