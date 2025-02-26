import {
    DashboardOutlined,
    GroupOutlined,
    DesktopOutlined,
    UnorderedListOutlined,
    BellOutlined,
    TeamOutlined,
    SkinOutlined,
    LockOutlined,
} from "@ant-design/icons-vue";


export default [
    {
        name: "Dashboard",
        text: "仪表盘",
        href: "/main/dashboard",
        icon: DashboardOutlined
    },
    {
        name: "Infra",
        text: "基础设施",
        type: "group",
        children: [
            {
                name: "Biz",
                text: "业务系统",
                href: "/main/biz",
                icon: GroupOutlined
            },
            {
                name: "主机",
                text: "主机",
                href: "/main/linux",
                icon: DesktopOutlined
            },]
    },
    {
        name: "Notification",
        text: "消息中心",
        type: "group",
        children: [
            {
                name: "Alarm",
                text: "事件告警",
                href: "/main/notification",
                icon: BellOutlined
            },
        ]
    },
    {
        name: "Setting",
        text: "系统设置",
        type: "group",
        children: [
            {
                name: "Menu",
                text: "菜单管理",
                href: "/main/menu",
                icon: UnorderedListOutlined 
            },
            {
                name: "User",
                text: "用户管理",
                href: "/main/user",
                icon: TeamOutlined
            },
            {
                name: "Role",
                text: "角色管理",
                href: "/main/role",
                icon: SkinOutlined
            },
            {
                name: "Permission",
                text: "权限管理",
                href: "/main/permission",
                icon: LockOutlined
            },
        ]
    }
]