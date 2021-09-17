import { Controller } from '@nestjs/common';
import { ApiTags } from '@nestjs/swagger';
import { TxService } from '../service/tx.service';

@ApiTags('Txs')
@Controller('txs')
export class TxController {
  constructor(private readonly txService: TxService) {}
}
