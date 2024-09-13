import request from "@/utils/request";
import type { Linux } from "../linux/api";

export class Alarm {
    id: number;
    trigger?: string;
    timestamp?: number;
    createTimestamp?: number;
    ack?: boolean;
    linux?: Linux;

    constructor(id: number) {
        this.id = id
    }

    load() {

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