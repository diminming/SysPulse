import request from "@/utils/request";
import type { Linux } from "../linux/api";

export class Alarm {
    id: number;
    trigger?: string;
    triggerId?: string;
    timestamp?: number;
    createTimestamp?: number;
    ack?: boolean;
    linux?: Linux;
    msg?: string;
    perfData?: {};

    constructor(id: number) {
        this.id = id
    }

    load() {
        return request({
            url: `/alarm/${this.id}`,
            method: "GET"
        })
    }

    static loadPage(pagination: any) {
        return request({
            url: "/alarm/page",
            method: "GET",
            params: {
                page: pagination.page,
                pageSize: pagination.pageSize,
            },
        })
    }
}
