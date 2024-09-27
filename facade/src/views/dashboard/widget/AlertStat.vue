<template>
    <a-card size="small" title="告警趋势图" style="margin-right: 1rem;margin-bottom: 1rem;">
        <template #extra>
            <a href="#" @click="showSetting = true">设置</a>
        </template>
        <div id="alertStat" style="width: 800px;height: 400px;"></div>
    </a-card>
    <a-modal v-model:open="showSetting" title="告警趋势图：设置" @ok="getData">
        <a-form>
            <a-form-item label="统计范围" name="range">
                <a-range-picker :show-time="{ format: 'HH:mm' }" format="YYYY-MM-DD HH:mm" @change="onRangeChange" :placeholder="['起始时间', '终止时间']" />
            </a-form-item>
        </a-form>
    </a-modal>
</template>

<script setup lang="ts">
import request from '@/utils/request';
import * as echarts from 'echarts';
import { onMounted, ref } from 'vue';
import dayjs, { Dayjs } from 'dayjs';

const now = new Date()
const showSetting = ref(false);
const timeRange = ref([dayjs(now).add(-24, "hour").valueOf(), dayjs(now).valueOf()])

const onRangeChange = (value: [Dayjs, Dayjs]) => {
    
    const start = value[0].valueOf()
    const end = value[1].valueOf()

    timeRange.value = [start, end]

};

const render = (xAxis, valueLst, min, max) => {
    let chart_dom: any = document.getElementById('alertStat')

    let chart = echarts.init(chart_dom)
    let option;

    option = {
        visualMap: [
            {
                show: false,
                type: 'continuous',
                seriesIndex: 0,
                min: min,
                max: max
            },
        ],
        tooltip: {
            trigger: 'axis'
        },
        xAxis: [
            {
                data: xAxis
            },
        ],
        yAxis: [
            {},
        ],
        series: [
            {
                type: 'line',
                showSymbol: false,
                data: valueLst
            }
        ]
    };

    option && chart.setOption(option);
}

const loadData = () => {
    return request({
        url: "/alarm/stat_trend",
        method: "GET",
        params: {
            "from": timeRange.value[0],
            "to": timeRange.value[1]
        }
    })
}

const getData = () => {
    loadData().then(resp => {
        const data = resp["data"]
        data.sort((a, b) => {
            return parseInt(a["timetag"]) - parseInt(b["timetag"])
        })
        let min = 0, max = 0;
        const xAxis = data.map(item => item["timetag"])
        const valueLst = data.map(item => {
            const value = item["total"]
            if (value < min) min = value
            else if (value > max) max = value
            return value
        })
        render(xAxis, valueLst, min, max)
        showSetting.value = false
    })
}

onMounted(() => {
    getData()
})

</script>

<style scoped></style>