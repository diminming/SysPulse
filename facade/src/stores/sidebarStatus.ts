import { defineStore } from 'pinia'

export const sidebarCollapsed = defineStore("sidebarCollapsed", {
    state: () => {
        return {
            collapsed: false
        }
    },
    actions: {
        toggle() {
            console.log("before:", this.collapsed)
            this.collapsed = !this.collapsed
            console.log("after:", this.collapsed)
        }
    }
})