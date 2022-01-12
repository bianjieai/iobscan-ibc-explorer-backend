import sha256 from 'sha256';

export function getDcDenom(msg) {
    let result = '';
    let dc_denom_path = '';
    const {
      source_port,
      source_channel,
      destination_port,
      destination_channel,
      data,
    } = msg.msg.packet;
    const dc_denom = data.denom;

    const prefix_sc = `${source_port}/${source_channel}/`;
    const prefix_dc = `${destination_port}/${destination_channel}/`;

    if (dc_denom.startsWith(prefix_sc)) {
      let dc_denom_clear_prefix = dc_denom.replace(prefix_sc, '');
      if (dc_denom_clear_prefix.indexOf('/') === -1) {
        result = dc_denom_clear_prefix;
        dc_denom_path = dc_denom_clear_prefix;
      } else {
        result = `ibc/${sha256(dc_denom_clear_prefix).toUpperCase()}`;
        dc_denom_path = dc_denom_clear_prefix;
      }
    } else {
      result = `ibc/${sha256(`${prefix_dc}${dc_denom}`).toUpperCase()}`;
      dc_denom_path = `${prefix_dc}${dc_denom}`;
    }

    return { dc_denom: result, dc_denom_path: dc_denom_path };
}
export function IbcDenom(denomPath :string, baseDenom :string) {
    if (!denomPath) {
        return baseDenom
    }
    return `ibc/${sha256(denomPath+'/'+baseDenom).toUpperCase()}`
}