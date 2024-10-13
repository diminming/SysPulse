<template>
    <a-descriptions title="基本信息" bordered size="small" :column="3" :contentStyle="{ 'width': '20rem' }"
        :labelStyle="{ 'width': '10rem' }">

        <a-descriptions-item label="主机ID">
            {{ base["hostid"] }}
        </a-descriptions-item>

        <a-descriptions-item label="主机名">
            {{ base["hostname"] }}
        </a-descriptions-item>

        <a-descriptions-item label="启动时间">
            {{ base["bootTime"] }}
        </a-descriptions-item>

        <a-descriptions-item label="启动时长">
            {{ base["uptime"] }}
        </a-descriptions-item>

        <a-descriptions-item label="架构">
            {{ base["kernelArch"] }}
        </a-descriptions-item>

        <a-descriptions-item label="操作系统类型">
            {{ base["os"] }}
        </a-descriptions-item>

        <a-descriptions-item label="内核版本">
            {{ base["kernelVersion"] }}
        </a-descriptions-item>

        <a-descriptions-item label="发行版">
            {{ base["platform"] }}
        </a-descriptions-item>

        <a-descriptions-item label="操作系统版本">
            {{ base["platformVersion"] }}
        </a-descriptions-item>

    </a-descriptions>

    <a-descriptions title="CPU信息" bordered size="small" :column="1" :labelStyle="{ 'width': '10rem' }"
        style="margin-top: 2rem;width: 50rem;">

        <a-descriptions-item label="CPU数量">
            {{ cpuinfo["count"] }}
        </a-descriptions-item>

        <a-descriptions-item label="CPU详情">
            <ul>
                <li v-for="item in cpuinfo['content']">{{ item }}</li>
            </ul>
        </a-descriptions-item>

    </a-descriptions>

    <a-descriptions title="网卡信息" bordered size="small" :column="3" :labelStyle="{ 'width': '10rem' }"
        style="margin-top: 2rem;">

        <template v-for="item in ifLst">

            <a-descriptions-item label="网卡名称">
                {{ item["name"] }}
            </a-descriptions-item>

            <a-descriptions-item label="IP地址">
                <ul class="ip_lst">
                    <li style="display: inline" v-for="addr in item['addrs']">{{ addr["addr"] }}</li>
                </ul>
            </a-descriptions-item>

            <a-descriptions-item label="MAC地址">
                {{ item["hardwareAddr"] }}
            </a-descriptions-item>

        </template>

    </a-descriptions>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue';
import { Linux } from '@/views/linux/api';
import { useRoute } from 'vue-router';

const base = ref({}), cpuinfo = ref({}), ifLst = ref([])

onMounted(() => {
    let route = useRoute()
    let linux = new Linux(parseInt(route.params.linuxId as string))
    linux.loadDescription().then((resp) => {
        const data = resp["data"]
        base.value = data['base']
        let params = new Map()
        data['cpu'].forEach((item) => {
            const key = item["modelName"]
            if (params.has(key)) {
                params.set(key, params.get(key) + 1)
            } else {
                params.set(key, 1)
            }
        })
        cpuinfo.value = {
            "count": data['cpu'].length,
            "content": Array.from(params.entries()).map(item => {
                return item[0] + " x " + item[1]
            })
        }
        ifLst.value = data["ifLst"]
    })
})
</script>

<style lang="css" scoped>
ul {
    list-style: none;
    margin: 0 0;
    padding: 0 0;
}

.ip_lst li:not(:first-child) {
    margin-left: 2rem
}
</style>