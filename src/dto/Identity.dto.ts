
import { BaseReqDto, BaseResDto, PagingReqDto } from './base.dto';
import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
/****************Req******************/
export class IdentityPubKeyAndCertificateReqDto extends PagingReqDto{
    @ApiPropertyOptional()
    @ApiProperty()
    id:string
}
//tx identity
export class IdentityReqDto extends PagingReqDto {
  @ApiPropertyOptional()
  search?: string;
}

export class IdentityInfoReqDto extends BaseReqDto {
  @ApiProperty()
  id: string;
}
export class IdentityByAddressReqDto extends PagingReqDto{
  @ApiProperty()
  address: string;
}
/****************Res******************/
export class IdentityResDto extends BaseResDto{
  identities_id: string;
  owner: string;
  credentials: string;
  create_block_height: string;
  create_block_time: string;
  create_tx_hash: string;
  update_block_height: string;
  update_tx_hash:string;
  update_block_time: string;
  constructor(txIdentitiesData){
    super();
    this.identities_id = txIdentitiesData.identities_id;
    this.owner = txIdentitiesData.owner;
    this.credentials = txIdentitiesData.credentials;
    this.create_block_height = txIdentitiesData.create_block_height;
    this.create_block_time = txIdentitiesData.create_block_time;
    this.create_tx_hash = txIdentitiesData.create_tx_hash;
    this.update_block_height = txIdentitiesData.update_block_height;
    this.update_tx_hash = txIdentitiesData.update_tx_hash;
    this.update_block_time = txIdentitiesData.update_block_time;
  }
  static bundleData(value: any): IdentityResDto[] {
    let data: IdentityResDto[] = [];
    data = value.map((v: any) => {
      return new IdentityResDto(v);
    });
    return data;
  }
}

export class IdentityPubKeyResDto extends BaseResDto{
  identities_id: string
  hash: string
  height: number
  time: number
  pubkey: object
  constructor(IdentityData){
    super();
    this.identities_id = IdentityData.identities_id;
    this.hash = IdentityData.hash;
    this.height = IdentityData.height;
    this.pubkey = IdentityData.pubkey;
    this.time = IdentityData.time;
  }
  static bundleData(value: any): IdentityPubKeyResDto[] {
    let data: IdentityPubKeyResDto[] = [];
    data = value.map((v: any) => {
      return new IdentityPubKeyResDto(v);
    });
    return data;
  }
}

export class IdentityCertificateResDto extends BaseResDto{
  identities_id:string
  hash: string
  height: number
  time: number
  certificate: string
  constructor(IdentityData){
    super();
    this.identities_id = IdentityData.identities_id;
    this.hash = IdentityData.hash;
    this.height = IdentityData.height;
    this.certificate = IdentityData.certificate;
    this.time = IdentityData.time;
  }
  static bundleData(value: any): IdentityCertificateResDto[] {
    let data: IdentityCertificateResDto[] = [];
    data = value.map((v: any) => {
      return new IdentityCertificateResDto(v);
    });
    return data;
  }
}
export class IdentityInfoResDto extends BaseResDto {
  identities_id: string
  owner: string
  credentials: string
  create_block_height: number
  create_block_time: number
  create_tx_hash: string
  constructor(IdentityData){
    super();
    this.identities_id = IdentityData.identities_id;
    this.owner = IdentityData.owner;
    this.credentials = IdentityData.credentials;
    this.create_block_height = IdentityData.create_block_height;
    this.create_block_time = IdentityData.create_block_time;
    this.create_tx_hash = IdentityData.create_tx_hash;
  }
}
