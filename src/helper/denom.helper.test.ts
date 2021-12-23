import {getDcDenom} from "./denom.helper";


describe('getDcDenom', () => {

    describe('getDcDenom', () => {
        it('getDcDenom Test', async () => {
            const msg = {
                msg:{
                    'packet':{
                        "source_port" : "transfer",
                        "source_channel" : "channel-10",
                        "destination_port" : "transfer",
                        "destination_channel" : "channel-36",
                        "data" : {
                            "denom" : "transfer/channel-9/transfer/channel-54/uiris",
                            "amount" : 2,
                            "sender" : "iaa1t7ktapnp9ym4latfqxpfhklhsspu7t8xphtwzl",
                            "receiver" : "iaa17cjdg63thy2vfqvvgj5lfv5dp339t0lra2gq8u"
                        }
                    }}
            }
            const data = await getDcDenom(msg)
            console.log(data, '--result--')
        });
    });
})