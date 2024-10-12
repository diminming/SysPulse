import {
    DashboardOutlined,
    GroupOutlined,
    DesktopOutlined,
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
                name: "Linux",
                text: "Linux",
                href: "/main/linux",
                icon: DesktopOutlined
            },]
    },

    // {
    //     name: "Database",
    //     text: "数据库",
    //     href: "/main/database",
    //     icon: DesktopOutlined
    // },
    // {
    //     name: "Cache",
    //     text: "缓存",
    //     href: "/main/cache",
    //     icon: DesktopOutlined
    // },
    // {
    //     name: "Queue",
    //     text: "队列",
    //     href: "/main/queue",
    //     icon: DesktopOutlined
    // },
    {
        name: "Notification",
        text: "消息中心",
        type: "group",
        children: [
            {
                name: "Alarm",
                text: "事件告警",
                href: "/main/notification",
                icon: DesktopOutlined
            },
        ]
    },
    
    // {
    //     name: "Users",
    //     text: "用户管理",
    //     href: "/main/users",
    //     icon: DesktopOutlined
    // },
    // {
    //     name: "Setting",
    //     text: "系统设置",
    //     href: "/main/setting",
    //     icon: DesktopOutlined
    // }
]