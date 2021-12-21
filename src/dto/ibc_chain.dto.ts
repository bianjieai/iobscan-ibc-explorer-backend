import {IbcChainConfigType} from '../types/schemaTypes/ibc_chain_config.interface';
import {BaseResDto} from './base.dto';

export class IbcChainResDto extends BaseResDto {
    chain_id: string;
    chain_name: string;
    icon: string;

    constructor(value) {
        super();
        const {chain_id, chain_name, icon} = value;
        this.chain_id = chain_id;
        this.chain_name = chain_name;
        this.icon = icon;
    }

    static bundleData(value: any): IbcChainResDto[] {
        const datasSortChars = value.filter(
            (item: IbcChainConfigType) => {
                const sASC = item.chain_name.charCodeAt(0);
                return (sASC >= 65 && sASC <= 90) || (sASC >= 97 && sASC <= 122);
            },
        );
        const datasSortOthers = value.filter(
            (item: IbcChainConfigType) => {
                const sASC = item.chain_name.charCodeAt(0);
                return sASC < 65 || (90 < sASC && sASC < 97) || sASC > 122;
            },
        );
        const datas: IbcChainResDto[] = [...datasSortChars, ...datasSortOthers].map(
            (item: IbcChainConfigType) => {
                item.chain_id = item.chain_id.replace(new RegExp("\_", "g"),"-")
                return new IbcChainResDto(item);
            },
        );
        return datas;
    }
}

export class IbcChainResultResDto {
    all: IbcChainResDto[];
    active: IbcChainResDto[];
    inactive: IbcChainResDto[];

    constructor(value) {
        const {all, active, inactive} = value;
        this.all = all || [];
        this.active = active || [];
        this.inactive = inactive || [];
    }
}
