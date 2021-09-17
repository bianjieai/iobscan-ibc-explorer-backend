import { addressPrefix } from '../constant'
//  todo: duanjie sdk待替换
let sdk = require('@irisnet/irishub-sdk');

export function getConsensusPubkey(value) {
    if (sdk && sdk.utils && sdk.types) {
        let pk = sdk.utils.Crypto.aminoMarshalPubKey({
            type: sdk.types.PubkeyType.ed25519,
            value: value
        })
        let pk_bech32  = sdk.utils.Crypto.encodeAddress(pk, addressPrefix.icp);
        return pk_bech32
    }
    return ''
}