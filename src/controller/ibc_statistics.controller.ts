import { Controller, Get, Query } from '@nestjs/common';
import { ApiTags } from '@nestjs/swagger';
import { Result } from '../api/ApiResult';
import { IbcStatisticsService } from '../service/ibc_statistics.service';

@ApiTags('IbcStatistics')
@Controller('ibc')
export class IbcStatisticsController {
  constructor(private readonly ibcStatisticsService: IbcStatisticsService) {}

  @Get('statistics')
  async getAllRecord() {
    const result = await this.ibcStatisticsService.findAllRecord();
    return new Result(result, 200);
  }
}
