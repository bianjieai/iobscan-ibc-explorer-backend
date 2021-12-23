import { PrometheusController } from "@willsoto/nestjs-prometheus";
import { Controller, Get, Res } from "@nestjs/common";
import { Response } from "express";

@Controller()
export class MonitorController extends PrometheusController {
    @Get()
    async index(@Res() response: Response) {
        await super.index(response);
    }
}