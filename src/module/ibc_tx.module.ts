import { Module } from '@nestjs/common';
import { MongooseModule } from '@nestjs/mongoose';
import { IbcTxController } from 'src/controller/ibc_tx.controller';
import { IbcTxService } from '../service/ibc_tx.service';
import { IbcTxSchema } from '../schema/ibc_tx.schema';
@Module({
  imports: [
    MongooseModule.forFeature([
      {
        name: 'IbcTx',
        schema: IbcTxSchema,
        collection: 'ex_ibc_tx',
      },
    ]),
  ],
  providers: [IbcTxService],
  controllers: [IbcTxController],
  exports: [IbcTxService],
})
export class IbcTxModule {}
