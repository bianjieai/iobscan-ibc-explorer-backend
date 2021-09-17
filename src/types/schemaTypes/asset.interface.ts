export interface IAssetStruct {
    symbol?:string,
    owner?:string,
    total_supply?:string,
    initial_supply?:string,
    max_supply?:string,
    mintable?: boolean,
    name?: string,
    scale?:string,
    src_protocol?: string;
    chain?: string;
}
