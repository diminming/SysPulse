import request from "@/utils/request"

export class User {

    id: number
    username?: string
    password?: string
    isActive?: boolean
    roleLst: Array<Role>
    createTimestamp?: number
    updateTimestamp?: number

    constructor(info: { id: number, username?: string, password?: string, isActive?: boolean, createTimestamp?: number, updateTimestamp?: number, roleLst?: Array<Role> }) {
        this.id = info.id
        this.username = info.username
        this.password = info.password
        this.isActive = info.isActive
        this.roleLst = info.roleLst || new Array()
        this.createTimestamp = info.createTimestamp
        this.updateTimestamp = info.updateTimestamp
    }

    static loadPage(args: {
        pageNum: number, pageSize: number
    }) {
        return request({
            url: "/user/page",
            method: "GET",
            params: {
                "pageNum": args.pageNum,
                "pageSize": args.pageSize
            }
        })
    }

    save() {
        return request({
            url: this.id > 0 ? `/user/${this.id}` : '/user',
            method: this.id > 0 ? "PUT" : "POST",
            data: this
        })
    }

    load() {
        return request({
            url: `/user/${this.id}`,
            method: "GET",
        })
    }

    delete() {
        return request({
            url: `/user/${this.id}`,
            method: "DELETE"
        })
    }

}

export class Permission {
    id: number
    identity?: string
    name?: string
    url?: string
    method?: string
    createTimestamp?: number
    updateTimestamp?: number

    constructor(info: { id: number, identity?: string, name?: string, url?: string, method?: string, createTimestamp?: number, updateTimestamp?: number }) {
        this.id = info.id
        this.identity = info.identity
        this.name = info.name
        this.url = info.url
        this.method = info.method
        this.createTimestamp = info.createTimestamp
        this.updateTimestamp = info.updateTimestamp
    }

    static loadPage(args: {
        pageNum: number, pageSize: number
    }) {
        return request({
            url: "/permission/page",
            method: "GET",
            params: {
                "pageNum": args.pageNum,
                "pageSize": args.pageSize
            }
        })
    }

    save() {
        return request({
            url: "/permission",
            method: "post",
            data: this
        })
    }

    delete() {
        return request({
            url: `/permission/${this.id}`,
            method: "delete",
        })
    }
}

export const showPermissionLst = (permissionLst: any): string => {

    if (permissionLst === undefined || permissionLst === null || permissionLst.length === 0) {
        return ""
    }
    let txtLst = new Array<string>()
    permissionLst.forEach((permission: any) => {
        txtLst.push(permission.name)
    })
    if (txtLst.length > 0) {
        return txtLst.join(", ")
    } else {
        return ""
    }
}

export const showRoleLst = (roleLst: any): string => {

    if (roleLst === undefined || roleLst === null || roleLst.length === 0) {
        return ""
    }
    let txtLst = new Array<string>()
    roleLst.forEach((role: any) => {
        txtLst.push(role.name)
    })
    if (txtLst.length > 0) {
        return txtLst.join(", ")
    } else {
        return ""
    }
}

export class Role {
    id: number
    identity?: string
    name?: string
    permissionLst: Array<Permission>
    createTimestamp?: number
    updateTimestamp?: number

    constructor(info: { id: number, identity?: string, name?: string, permission?: Array<Permission>, createTimestamp?: number, updateTimestamp?: number }) {
        this.id = info.id
        this.identity = info.identity
        this.name = info.name
        this.permissionLst = info.permission || [] // 设置默认值为空数组
        this.createTimestamp = info.createTimestamp
        this.updateTimestamp = info.updateTimestamp
    }
    static loadPage(params: {
        pageNum: number, pageSize: number
    }) {
        return request({
            url: "/role/page",
            method: "GET",
            params
        })
    }

    delete() {
        return request({
            url: `/role/${this.id}`,
            method: "DELETE"
        })
    }

    save() {
        return request({
            url: "/role",
            method: "POST",
            data: this
        })
    }
}

export class Menu {
    id: number
    identity?: string
    title?: string
    url?: string
    type?: string
    index?: number
    parentId?: number
    parentTitle?: string
    permissionLst: Array<Permission>
    createTimestamp?: number
    updateTimestamp?: number

    constructor(info: {
        id: number,
        identity?: string,
        title?: string,
        url?: string,
        index?: number,
        type?: string,
        parentId?: number
        parentTitle?: string
        permissionLst?: Array<Permission>,
        createTimestamp?: number,
        updateTimestamp?: number
    }) {
        this.id = info.id,
            this.identity = info.identity
        this.title = info.title
        this.url = info.url
        this.type = info.type
        this.parentId = info.parentId
        this.parentTitle = info.parentTitle
        this.index = info.index
        this.permissionLst = info.permissionLst || new Array<Permission>()
        this.createTimestamp = info.createTimestamp
        this.updateTimestamp = info.updateTimestamp
    }

    static loadPage(params: {
        pageNum: number, pageSize: number
    }) {
        return request({
            url: "/menu/page",
            method: "GET",
            params
        })
    }

    save() {
        return request({
            url: "/menu",
            method: "post",
            data: this
        })
    }

    delete() {
        return request({
            url: `/menu/${this.id}`,
            method: "delete",
        })
    }

    showTypeString() {
        switch (this.type) {
            case "item":
                return "菜单项"
            case "group":
                return "菜单分组"
            default:
                return "未知类型"
        }
    }

}