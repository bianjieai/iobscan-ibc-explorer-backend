import { Controller, Get } from '@nestjs/common';
import { ApiTags } from '@nestjs/swagger';
import { Result } from '../api/ApiResult';
import { ConfigService } from '../service/config.service';
import { ConfigResDto } from '../dto/config.dto';
@ApiTags('Config')
@Controller('ibc')
export class ConfigController {
  constructor(private readonly configService: ConfigService) {}

  @Get('config')
  async getAllRecord(): Promise<Result<ConfigResDto>> {
    const result: ConfigResDto = await this.configService.findOne();
    return new Result<ConfigResDto>(result);
  }
}
