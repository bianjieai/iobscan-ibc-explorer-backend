export class StatisticsResDto {
    private avgBlockTime: number;
    private serviceCount: number;
    private validatorCount: number;
    private assetCount: number;
    private identityCount: number;
    private denomCount: number;
    private validatorNumCount: number;
    private moniker: string;
    private validator_icon: string;
    private operator_addr: string;
    
    constructor(Detail) {
        this.avgBlockTime = Detail.avgBlockTime;
        this.serviceCount = Detail.serviceCount;
        this.validatorCount = Detail.validatorCount;
        this.assetCount = Detail.assetCount;
        this.identityCount = Detail.identityCount;
        this.denomCount = Detail.denomCount;
        this.validatorNumCount = Detail.validatorNumCount;
    }
}

export class NetworkStatisticsResDto {
    blockHeight: number;
    moniker: string;
    validator_icon: string;
    operator_addr: string;
    latestBlockTime:number;
    txCount: number;
    bonded_tokens: string;
    total_supply: string;
    constructor(Detail) {
        this.blockHeight = Detail.block && Detail.block.height;
        this.moniker = Detail.block && Detail.block.moniker;
        this.validator_icon =Detail.block && Detail.block.validator_icon;
        this.operator_addr = Detail.block && Detail.block.operator_addr;
        this.latestBlockTime = Detail.block && Detail.block.latestBlockTime;
        this.txCount = Detail.txCount;
        this.bonded_tokens = Detail.bondedTokensInformation && Detail.bondedTokensInformation.bonded_tokens;
        this.total_supply = Detail.bondedTokensInformation && Detail.bondedTokensInformation.total_supply;
    }
}
