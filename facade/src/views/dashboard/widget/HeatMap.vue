<template>
  <a-card size="small" title="告警分布图" style="margin-right: 1rem;margin-bottom: 1rem;">
    <template #extra>
      <a href="#" @click="showSetting = true">设置</a>
    </template>
    <div id="heatMap" style="width: 800px;height: 400px;"></div>
  </a-card>
  <a-modal v-model:open="showSetting" title="告警分布图：设置" @ok="getData">
    <a-form>
      <a-form-item label="统计范围" name="range">
        <a-range-picker :show-time="{ format: 'HH:mm' }" format="YYYY-MM-DD HH:mm" @change="onRangeChange"
          :placeholder="['起始时间', '终止时间']" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup lang="ts">

import * as echarts from 'echarts';
import { onMounted, ref } from 'vue';
import request from "@/utils/request"
import dayjs, { Dayjs } from 'dayjs';

const showSetting = ref(false)
const now = new Date()
const timeRange = ref([dayjs(now).add(-24, "hour").valueOf(), dayjs(now).valueOf()])

const getData = () => {
  loadData().then((resp: any) => {
    const data = resp["data"]
    const xAxis = Array.from(new Set(data.map((item: any) => item["timetag"]).sort()))
    const yAxis = Array.from(new Set(data.map((item: any) => {
      const bizName = item["bizName"]
      if (bizName === "")
        return "未分组"
      else {
        return bizName
      }
    })))
    let min = 0, max = 0
    const heatData = data.map((item: any) => {
      const xData = item["timetag"]
      const yData = item["bizName"] == "" ? "未分组" : item["bizName"]
      const value = item["total"]
      if (value < min) min = value
      else if (value > max) max = value
      return [
        xAxis.findIndex((idx: any) => {
          return idx === xData
        }),
        yAxis.findIndex((idx: any) => {
          return idx === yData
        }),
        value
      ]
    })
    renderGraph(xAxis, yAxis, heatData, min, max)
  })
  showSetting.value = false
}

const renderGraph = (xAxis, yAxis, heatData, min, max) => {

  let container: any = document.querySelector(".graph")
  let chart_dom: any = document.getElementById('heatMap')

  // let chart = echarts.init(chart_dom, undefined, {
  //   width: container.offsetWidth,
  //   height: container.offsetHeight
  // });

  let chart = echarts.init(chart_dom)

  let option;

  option = {

    tooltip: {
      position: 'top'
    },
    grid: {
      height: '50%',
      top: '10%'
    },
    xAxis: {
      type: 'category',
      data: xAxis,
      splitArea: {
        show: true
      }
    },
    yAxis: {
      type: 'category',
      data: yAxis,
      splitArea: {
        show: true
      }
    },
    visualMap: {
      min: min,
      max: max,
      calculable: true,
      orient: 'horizontal',
      left: 'center',
      bottom: '15%'
    },
    series: [
      {
        name: '告警计数',
        type: 'heatmap',
        data: heatData,
        label: {
          show: true
        },
        emphasis: {
          itemStyle: {
            shadowBlur: 10,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          }
        }
      }
    ]
  };


  option && chart.setOption(option);
}

const onRangeChange = (value: [Dayjs, Dayjs]) => {
  const start = value[0].valueOf()
  const end = value[1].valueOf()

  timeRange.value = [start, end]
};


const loadData = () => {
  return request({
    url: "/alarm/stat_heat",
    method: "GET",
    params: {
      "from": timeRange.value[0],
      "to": timeRange.value[1]
    }
  })
}

onMounted(() => {
  getData()
})


</script>

<style scoped></style>