import {InjectMetric, makeGaugeProvider} from "@willsoto/nestjs-prometheus";
import {Gauge} from "prom-client";
import {Injectable} from "@nestjs/common";

@Injectable()
export class LcdConnectionMetric {
    constructor(@InjectMetric("ibc_explorer_backend_lcd_connection_status") public gauge: Gauge<string>) {}

    async collect(chainId,value) {
        this.gauge.set({
            "chain_id":chainId,
        },value)
    }
}

export function LcdConnectionProvider() {
    return makeGaugeProvider({
        name: "ibc_explorer_backend_lcd_connection_status",
        help: "ibc explorer backend lcd connection status (1:Reachable  0:NotReachable)",
        labelNames: ['chain_id'],
    })
}
