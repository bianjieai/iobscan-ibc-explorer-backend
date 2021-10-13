import { Injectable } from '@nestjs/common';
import { Connection } from 'mongoose';
import { InjectConnection } from '@nestjs/mongoose';
import { ConfigSchema } from '../schema/config.schema';
import { ConfigResDto } from '../dto/config.dto';
@Injectable()
export class ConfigService {
  private configModel;

  constructor(@InjectConnection() private connection: Connection) {
    this.getModels();
  }

  async getModels(): Promise<void> {
    this.configModel = await this.connection.model(
      'configModel',
      ConfigSchema,
      'ibc_config',
    );
  }

  // findOne
  async findOne(): Promise<ConfigResDto> {
    const result: ConfigResDto = new ConfigResDto(await this.configModel.findRecord())
    return result;
  }
}
