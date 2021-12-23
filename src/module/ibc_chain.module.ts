import {Module} from '@nestjs/common';
import {IbcChainController} from '../controller/ibc_chain.controller';
import {IbcChainService} from '../service/ibc_chain.service';

@Module({
    providers: [IbcChainService],
    controllers: [IbcChainController],
})
export class IbcChainModule {
}
