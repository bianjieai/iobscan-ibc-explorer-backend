import {InjectMetric, makeGaugeProvider} from "@willsoto/nestjs-prometheus";
import {Gauge} from "prom-client";
import {Injectable} from "@nestjs/common";

@Injectable()
export class TransferTaskStatusMetric {
    constructor(@InjectMetric("ibc_explorer_backend_transfer_task_status") public gauge: Gauge<string>) {}

    async collect(value) {
        this.gauge.set(value)
    }
}

export function TransferTaskStatusProvider() {
    return makeGaugeProvider({
        name: "ibc_explorer_backend_transfer_task_status",
        help: "ibc explorer backend transfer task status (1:Working  0:Notwork)",
    })
}
