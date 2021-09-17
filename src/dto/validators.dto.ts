import { BaseResDto,PagingReqDto} from './base.dto';
import { ApiError } from '../api/ApiResult';
import { ErrorCodes} from '../api/ResultCodes';
import { ApiPropertyOptional } from '@nestjs/swagger';
export class ValidatorsReqDto extends PagingReqDto{
  @ApiPropertyOptional()
  jailed: boolean | string
    static validate (value:any){
      super.validate(value)
      if(value.jailed.toString() !== 'true' && value.jailed.toString() !== 'false'){
        throw new ApiError(ErrorCodes.InvalidRequest,"jailed must be true or false")
      }
    }
}
export class ValidatorsResDto extends BaseResDto {
  name: string
  pubkey: string
  power: string
  jailed: boolean |string
  operator: boolean
  constructor(validatorsData) {
    super();
    //validatorsData 数据库查询出来的结果
    //定义返回结果的字段名称，处理返回结果
    this.name = validatorsData.name;
    this.pubkey = validatorsData.pubkey;
    this.operator = validatorsData.operator;
    this.power = validatorsData.power;
    this.jailed = validatorsData.jailed;
  }
  static bundleData(validatorData:any){
    let data: Array<ValidatorsResDto> = [];
    data = validatorData.map((item: any) => {
      return new ValidatorsResDto(item)
    })
    return data
  }
}
