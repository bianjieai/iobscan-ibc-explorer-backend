import {InjectMetric, makeGaugeProvider} from "@willsoto/nestjs-prometheus";
import {Gauge} from "prom-client";
import {Injectable} from "@nestjs/common";

@Injectable()
export class IbcTxProcessingMetric {
    constructor(@InjectMetric("ibc_explorer_backend_processing_cnt") public gauge: Gauge<string>) {}
    async collect(value) {
        this.gauge.set(value)
    }
}

export function IbcTxProcessingProvider() {
    return makeGaugeProvider({
        name: "ibc_explorer_backend_processing_cnt",
        help: "ibc explorer backend processing cnt every scrape times",
    })
}
