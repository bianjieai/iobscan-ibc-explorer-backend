import { Module } from '@nestjs/common';
import { MongooseModule } from '@nestjs/mongoose';
import { IbcChainController } from 'src/controller/ibc_chain.controller';
import { IbcChainService } from '../service/ibc_chain.service';
import { IbcChainSchema } from '../schema/ibc_chain.schema';
export const ibcChainMongooseFeature = [
  {
    name: 'IbcChain',
    schema: IbcChainSchema,
    collection: 'chain_config',
  },
];
@Module({
  imports: [MongooseModule.forFeature(ibcChainMongooseFeature)],
  providers: [IbcChainService],
  controllers: [IbcChainController],
  exports: [IbcChainService],
})
export class IbcChainModule {}
