import { ChainHttp } from './chain.http';

describe('ChainHttp', () => {
    describe('getIbcChannels', () => {
        it('getIbcChannels Test', async () => {
            const lcdAddr = "https://cosmoshub.stakesystems.io"
            const lcdApiPath = {channels_path:"/ibc/core/channel/v1beta1/channels?pagination.offset=OFFSET&pagination.limit=LIMIT&pagination.count_total=true"}
            const result = await ChainHttp.getIbcChannels(lcdAddr,lcdApiPath.channels_path)
            console.log(result,'----')
        });
    });

    describe('getDenomByLcdAndHash', () => {
        it('getDenomByLcdAndHash Test', async () => {
            const lcdAddr = "https://cosmoshub.stakesystems.io",ibcHash = "EC4B5D87917DD5668D9998146F82D70FDF86652DB333D04CE29D1EB18E296AF5"
            const result = await ChainHttp.getDenomByLcdAndHash(lcdAddr,ibcHash)
            console.log(result,'----')
        });
    });

    describe('getDcChainIdByScChannel', () => {
        it('getDcChainIdByScChannel Test', async () => {
            const lcdAddr = "https://osmosis.stakesystems.io",
                lcdApiPath = {client_state_path:"/ibc/core/channel/v1/channels/CHANNEL/ports/PORT/client_state"}
            const result = await ChainHttp.getDcChainIdByScChannel(lcdAddr,lcdApiPath.client_state_path,"transfer","channel-6")
            console.log(result,'----')
        });
    });
})