import { Module } from '@nestjs/common';
import { IbcDenomController } from 'src/controller/ibc_denom.controller';
import { IbcDenomService } from '../service/ibc_denom.service';
@Module({
  providers: [IbcDenomService],
  controllers: [IbcDenomController],
})
export class IbcDenomModule {}
