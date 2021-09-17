import { Module } from '@nestjs/common';
import { TxService } from '../service/tx.service';
@Module({
  providers: [TxService],
  exports: [TxService],
})
export class TxModule {}
