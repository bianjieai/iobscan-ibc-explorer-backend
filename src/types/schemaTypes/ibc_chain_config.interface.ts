export interface IbcChainConfigType {
    chain_id: string;
    icon: string;
    chain_name: string;
    lcd: string;
    lcd_api_path: { channels_path: string;client_state_path: string;};
    ibc_info?: any;
    ibc_info_hash_lcd?: string;
    ibc_info_hash_caculate?: string;
    is_manual?: boolean
}
