import { Controller, Get, Query } from '@nestjs/common';
import { ApiTags } from '@nestjs/swagger';
import { Result } from '../api/ApiResult';
import { IbcStatisticsService } from '../service/ibc_statistics.service';
import { IbcStatisticsResDto } from '../dto/ibc_statistics.dto';
@ApiTags('IbcStatistics')
@Controller('ibc')
export class IbcStatisticsController {
  constructor(private readonly ibcStatisticsService: IbcStatisticsService) {}

  @Get('statistics')
  async getAllRecord(): Promise<Result<IbcStatisticsResDto[]>> {
    const result: IbcStatisticsResDto[] = await this.ibcStatisticsService.findAllRecord();
    return new Result(result);
  }
}
