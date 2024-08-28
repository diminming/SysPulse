<template>
    <div style="text-align: right;margin-bottom: 8px;">
        <a-space>
            <a-range-picker style="width: 400px" show-time :format="dateTimeFormat" :presets="rangePresets"
                v-model:value="dateTimeRange" @change="onRangeChange" size="small" />
            自动刷新
            <a-switch v-model:checked="autoRefreshSetting.enable" size="small" @change="changeAutoRefreshStat" />
            <a-select size="small" style="width: 8rem" @select="autoRefresh" v-model:value="autoRefreshSetting.option"
                :disabled="!autoRefreshSetting.enable" :options="autoRefreshOptions.map(item => {
                    return {
                        value: item.key,
                        label: item.label
                    }
                })">
            </a-select>
        </a-space>
    </div>
    <a-row style="margin-bottom: 8px;">
        <a-space size='middle'>
            <a-card size="small" class="graph-tile">
                <div class="title">CPU使用率</div>
                <div class="value">{{ lastCpuUsage }} %</div>
                <div class="graph cpu-usage-mini"></div>
            </a-card>
            <a-card size="small" class="graph-tile">
                <div class="title">系统负载</div>
                <div class="value">{{ lastLoad }}</div>
                <div class="graph load-mini"></div>
            </a-card>
            <a-card size="small" class="graph-tile">
                <div class="title">可用内存</div>
                <div class="value">{{ (lastAvailableMemory / 1024 / 1024 / 1024).toFixed(2) + " GB" }}</div>
                <div class="graph mem-available-mini"></div>
            </a-card>
            <a-card size="small" class="graph-tile">
                <div class="title">Swap使用量</div>
                <div class="value">{{ (lastSwapUsed / 1024 / 1024 / 1024).toFixed(2) + " GB" }}</div>
                <div class="graph swap-mini"></div>
            </a-card>
            <a-card size="small" class="graph-tile">
                <div class="title">磁盘写入</div>
                <div class="value">{{ lastDiskWrite }}</div>
                <div class="graph disk-write-counter-mini"></div>
            </a-card>
            <a-card size="small" class="graph-tile">
                <div class="title">磁盘读取</div>
                <div class="value">{{ lastDiskRead }}</div>
                <div class="graph disk-read-counter-mini"></div>
            </a-card>
            <a-card size="small" class="graph-tile">
                <div class="title">网络上传</div>
                <div class="value">{{ (lastSent / 1024 / 1024 / 1024).toFixed(2) + " GB" }}</div>
                <div class="graph if-sent-mini"></div>
            </a-card>
            <a-card size="small" class="graph-tile">
                <div class="title">网络下载</div>
                <div class="value">{{ (lastRecv / 1024 / 1024 / 1024).toFixed(2) + " GB" }}</div>
                <div class="graph if-recv-mini"></div>
            </a-card>
        </a-space>
    </a-row>
    <a-row style="margin-bottom: 10px;" class="gallery">
        <a-space size="small" wrap>
            <a-card size="small" title="CPU使用情况" class="frame">
                <div class="graph cpu"></div>
            </a-card>

            <a-card size="small" title="系统负载情况" class="frame">
                <div class="graph load"></div>
            </a-card>

            <a-card size="small" title="内存使用量" class="frame">
                <div class="graph mem"></div>
            </a-card>

            <a-card size="small" title="内存交换情况" class="frame">
                <div class="graph mem exchange"></div>
            </a-card>

            <a-card size="small" title="Swap使用量" class="frame">
                <div class="graph swap"></div>
            </a-card>

            <a-card size="small" title="Swap交换情况" class="frame">
                <div class="graph swap exchange"></div>
            </a-card>

            <a-card size="small" title="文件系统使用量" class="frame">
                <div class="graph fs"></div>
            </a-card>

            <a-card size="small" title="inode使用量" class="frame">
                <div class="graph inode"></div>
            </a-card>

            <a-card size="small" title="磁盘吞吐情况" class="frame">
                <div class="graph disk throughput"></div>
            </a-card>
            <a-card size="small" title="磁盘IO性能" class="frame">
                <div class="graph disk io"></div>
            </a-card>
            <a-card size="small" title="网络吞吐情况" class="frame">
                <div class="graph net throughput"></div>
            </a-card>
            <a-card size="small" title="网络IO性能" class="frame">
                <div class="graph net io"></div>
            </a-card>
        </a-space>
    </a-row>
</template>
<script setup lang="ts">

import { ref, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import dayjs, { Dayjs } from 'dayjs';

import LinuxAPI, { CaclUsage, Linux } from "@/views/linux/api"
import { JsonResponse } from '@/utils/common';

type RangeValue = [Dayjs, Dayjs];
type AutoRefreshSetting = { option: string | undefined, handler: number | undefined, enable: boolean }
type AutoRefreshOption = { key: string, label: string, handler: Function }

const linuxId = ref<number>(0),
    lastLoad = ref<number>(0),
    lastAvailableMemory = ref<number>(0),
    lastSwapUsed = ref<number>(0),
    lastCpuUsage = ref<string>('0'),
    lastDiskWrite = ref<number>(0),
    lastDiskRead = ref<number>(0),
    lastSent = ref<number>(0),
    lastRecv = ref<number>(0),
    dateTimeRange = ref<RangeValue>([dayjs().add(-1, 'h'),
    dayjs()]),
    dateTimeFormat = "YYYY/MM/DD HH:mm:ss",
    autoRefreshSetting = ref<AutoRefreshSetting>({
        "option": undefined,
        "handler": undefined,
        "enable": false
    })


const onRangeChange = (dates: RangeValue, dateStrings: string[]) => {
    if (dates) {
        render(dates[0].valueOf() * 1000, dates[1].valueOf() * 1000)
    }
};

const rangePresets = ref([
    { label: '最后 1 小时', value: [dayjs().add(-1, 'h'), dayjs()] },
    { label: '最后 12 小时', value: [dayjs().add(-12, 'h'), dayjs()] },
    { label: '最后 1 周', value: [dayjs().add(-7, 'd'), dayjs()] },
    { label: '最后 2 周', value: [dayjs().add(-14, 'd'), dayjs()] },
    { label: '最后 30 天', value: [dayjs().add(-30, 'd'), dayjs()] },
]);

const changeAutoRefreshStat = (checked: boolean, event: Event) => {
    if (checked === true) {
        autoRefreshSetting.value.option = autoRefreshOptions[0]['key']
        autoRefresh(autoRefreshOptions[0]['key'], undefined)
    } else {
        if (autoRefreshSetting.value.handler) {
            window.clearInterval(autoRefreshSetting.value.handler)
            autoRefreshSetting.value.handler = undefined
        }
    }
}

const autoRefreshOptions: Array<AutoRefreshOption> = [{
    label: "最后5分钟",
    key: "last5Min",
    handler: () => {
        render(dayjs().add(-5, 'm').valueOf(), dayjs().valueOf())
    }
}, {
    key: "last15Min",
    label: "最后15分钟",
    handler: () => {
        render(dayjs().add(-15, 'm').valueOf(), dayjs().valueOf())
    }
}, {
    key: "last30Min",
    label: "最后30分钟",
    handler: () => {
        render(dayjs().add(-30, 'm').valueOf(), dayjs().valueOf())
    }
}, {
    key: "last1Hour",
    label: "最后1小时",
    handler: () => {
        render(dayjs().add(-1, 'h').valueOf(), dayjs().valueOf())
    }
}]

const autoRefresh = (value: string, _: any) => {
    autoRefreshSetting.value.handler = window.setInterval(function () {
        try {
            autoRefreshOptions.find(item => item.key === value)?.handler()
        } catch (exception) {
            console.error(exception)
        }
    }, 3000)

}

onMounted(() => {
    let route = useRoute()
    linuxId.value = parseInt(route.params.linuxId as string)
    render(dateTimeRange.value[0].valueOf(), dateTimeRange.value[1].valueOf())
})

const render = (start: number, end: number) => {
    const linux = new Linux(linuxId.value)
    linux.RenderCpuPerfChart(start, end)
    linux.RenderMemoryChart(start, end)
    linux.RenderLoadPerfChart(start, end)
    linux.RenderExchangePerfChart(start, end)
    linux.RenderFileSystemPerfChart(start, end)
    linux.RenderDiskIOPerfChart(start, end)
    linux.RenderNetIOPerfChart(start, end)

    LinuxAPI.GetCpuUsageLine(linuxId.value, start, end).then((resp: any) => {
        if (resp.data && resp.data != null) {
            let jsonResp = new JsonResponse(resp.data, resp.msg, resp.status);
            let item = resp.data[resp.data.length - 1]
            let user = item["user"]
            lastCpuUsage.value = CaclUsage(user, item)
            LinuxAPI.RenderCpuUsageLine(jsonResp)
        }
    })
    LinuxAPI.GetLoad1Line(linuxId.value, start, end).then((resp: any) => {
        if (resp.data && resp.data != null) {
            let jsonResp = new JsonResponse(resp.data, resp.msg, resp.status);
            lastLoad.value = resp.data[resp.data.length - 1].load1
            LinuxAPI.RenderLoad1Line(jsonResp)
        }
    })
    LinuxAPI.GetAvailableMemoryLine(linuxId.value, start, end).then((resp: any) => {
        if (resp.data && resp.data != null) {
            let jsonResp = new JsonResponse(resp.data, resp.msg, resp.status);
            lastAvailableMemory.value = resp.data[resp.data.length - 1].available
            LinuxAPI.RenderAvailableMemoryLine(jsonResp)
        }
    })
    LinuxAPI.GetSwapUsedLine(linuxId.value, start, end).then((resp: any) => {
        if (resp.data && resp.data != null) {
            let jsonResp = new JsonResponse(resp.data, resp.msg, resp.status);
            lastSwapUsed.value = resp.data[resp.data.length - 1].used
            LinuxAPI.RenderSwapUsedLine(jsonResp)
        }
    })
    LinuxAPI.GetDiskIOLine(linuxId.value, start, end).then((resp: any) => {
        if (resp.data && resp.data != null) {
            let mapping: Map<string, any[]> = new Map()
            resp.data.forEach((item: any) => {
                if (mapping.has(item.name)) {
                    (mapping.get(item.name) as any[]).push(item)
                } else {
                    mapping.set(item.name, [item])
                }
            })
            let writecount = 0, readcount = 0
            let lst = []
            for (const [key, valLst] of mapping) {
                let last = valLst[valLst.length - 1]
                writecount += last["writecount"]
                readcount += last["readcount"]
                if (lst.length < valLst.length) {
                    lst = new Array(valLst.length)
                    for (let i = 0; i < valLst.length; i++) {
                        lst[i] = valLst[i]
                    }
                } else {
                    for (let i = 0; i < valLst.length; i++) {
                        lst[i].writecount += valLst[i].writecount
                        lst[i].readcount += valLst[i].readcount
                    }
                }
            }
            lastDiskWrite.value = writecount
            lastDiskRead.value = readcount
            let jsonResp = new JsonResponse(lst, resp.msg, resp.status);
            LinuxAPI.RenderDiskIOCounterLine(jsonResp)
        }
    })
    LinuxAPI.GetIfIOLine(linuxId.value, start, end).then((resp: any) => {
        if (resp.data && resp.data != null) {
            let mapping: Map<string, any[]> = new Map()
            resp.data.forEach((item: any) => {
                if (mapping.has(item.name)) {
                    (mapping.get(item.name) as any[]).push(item)
                } else {
                    mapping.set(item.name, [item])
                }
            })
            let bytessent = 0, bytesrecv = 0
            let lst = []
            for (const [key, valLst] of mapping) {
                let last = valLst[valLst.length - 1]
                bytessent += last["bytessent"]
                bytesrecv += last["bytesrecv"]
                if (lst.length < valLst.length) {
                    lst = new Array(valLst.length)
                    for (let i = 0; i < valLst.length; i++) {
                        lst[i] = valLst[i]
                    }
                } else {
                    for (let i = 0; i < valLst.length; i++) {
                        lst[i].bytessent += valLst[i].bytessent
                        lst[i].bytesrecv += valLst[i].bytesrecv
                    }
                }
            }
            lastSent.value = bytessent
            lastRecv.value = bytesrecv

            let jsonResp = new JsonResponse(lst, resp.msg, resp.status);
            LinuxAPI.RenderIfIOCounterLine(jsonResp)
        }
    })

}

</script>

<style scoped>
.graph-tile {
    width: 11.5rem;
}

.graph-tile .graph {
    height: 40px;
}

.graph-tile .title {
    text-align: center;
    font-size: 0.8rem;
    font-weight: 600;
}

.graph-tile .value {
    text-align: center;
    font-size: 1rem;
}

.gallery .frame {
    width: 24.5rem;
}

.gallery .graph {
    width: 100%;
    height: 15rem
}
</style>