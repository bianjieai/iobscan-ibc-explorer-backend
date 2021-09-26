import { BaseResDto } from './base.dto';
import { ApiProperty, ApiPropertyOptional } from '@nestjs/swagger';
export class IbcChainConfigResDto {
    chain_id: string;
    icon: string;
    chain_name: string;
    lcd: string;
    ibc_info: object;
}

export class IbcChainResDto {
    chain_id: string;
    chain_name: string;
    icon: string;
    create_at: string;
    update_at: string;
}

export class IbcChainResultResDto {
    all: IbcChainConfigResDto[];
    active: IbcChainResDto[];
    inactive: IbcChainConfigResDto[];
    constructor(value) {
        const { all, active, inactive } = value
        this.all = all || []
        this.active = active || []
        this.inactive = inactive || []
    }
}