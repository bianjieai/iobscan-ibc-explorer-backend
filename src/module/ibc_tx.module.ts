import { Module } from '@nestjs/common';
import { IbcTxController } from 'src/controller/ibc_tx.controller';
import { IbcTxService } from '../service/ibc_tx.service';
@Module({
  providers: [IbcTxService],
  controllers: [IbcTxController],
  exports: [IbcTxService],
})
export class IbcTxModule {}
