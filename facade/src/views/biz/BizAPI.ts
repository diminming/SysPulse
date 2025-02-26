import request from "@/utils/request"

export class BizUtil {

    id: number
    name?: string
    idenetity?: string
    desc?: string

    constructor(id: number, name?: string, identity?: string, desc?: string) {
        this.id = id
        this.name = name
        this.idenetity = identity
        this.desc = desc
    }

    loadTopo(graphSetting: {}) {
        return request({
            url: `/biz/${this.id}/topo`,
            method: "POST",
            data: graphSetting
        })
    }

    load() {
        return request({
            url: `/biz/${this.id}`,
            method: "GET"
        })
    }

    countInst() {
        return request({
            url: `biz/${this.id}/count_inst`,
            method: "GET"
        })
    }

}