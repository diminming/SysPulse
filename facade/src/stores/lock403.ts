import { defineStore } from 'pinia'

export const lock403Store = defineStore("lock_403", {
  state: () => {
    return {
      locked: false
    }
  },
  getters: {
    
  },
  actions: {
    lock() {
      this.locked = true
    },
    unlock() {
      this.locked = false
    },
    isLocked() {
      return this.locked
    }
  }
})