import { Module } from '@nestjs/common';
import { IbcBaseDenomController } from '../controller/ibc_base_denom.controller';
import { IbcBaseDenomService } from '../service/ibc_base_denom.service';
@Module({
  providers: [IbcBaseDenomService],
  controllers: [IbcBaseDenomController],
})
export class IbcBaseDenomModule {}
