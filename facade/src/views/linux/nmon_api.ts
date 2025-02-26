import request from "@/utils/request"
import axios from 'axios'

const nmonService = axios.create({})

export class NMON {
    id: number
    hostname?: string
    from?: number
    to?: number
    source?: string
    path?: string
    constructor(id: number, hostname?: string, from?: number, to?: number, source?: string, path?: string) {
        this.id = id
        this.hostname = hostname
        this.from = from
        this.to = to
        this.source = source
        this.path = path
    }

    static getByPage({linuxId, page, size}: {linuxId: number, page: number, size: number}) {
        return request({
            url: "/nmon/page",
            method: "get",
            params: {
                page,
                size,
                linuxId
            }
        })
    }

    getCategories() {
        
        return nmonService({
            url: "/nmon/categories",
            method: "post",
            data: {
                "filePath": this.path
            }
        })
    }

    getNMONData(category: string, field: string) {
        return nmonService({
            url: "/nmon/data",
            method: "post",
            data: {
                filePath: this.path,
                category,
                field
            }
        })
    }
    
}