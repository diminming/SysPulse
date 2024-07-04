<template>
  <div class="linuxList">
    <div class="tile">
      <div id="myEChartsBar"></div>
      <div id="myEChartsLine"></div>
    </div>
  </div>
</template>

<!--需要vue2改写vue3-->
<script lang="ts">
import Tile from "@/views/dashboard/widget/Tile.vue"
import echarts from "@/utils/echarts";
export default {
  data() {
    return {};
  },
  methods: {
    // 基本柱形图
    changeBar() {
      const myEChart= echarts.init(document.getElementById("myEChartsBar"));
      const option = {
        xAxis: {
          data: ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"],
        },
        yAxis: {},
        series: [
          {
            type: "bar",
            data: [23, 24, 18, 25, 27, 28, 25],
          },
        ],
      };
      myEChart.setOption(option);
      // 根据页面大小自动响应图表大小
      window.addEventListener("resize", function () {
        myEChart.resize();
      });
    },
    // 折线图
    changeLine() {
      // 获取组件实例
      const myEChart= echarts.init(document.getElementById("myEChartsLine"));
      // 设置配置项
      const option = {
        xAxis: {
          data: ["A", "B", "C", "D", "E"],
        },
        yAxis: {},
        series: [
          {
            data: [10, 22, 28, 43, 49],
            type: "line",
            stack: "x",
          },
          {
            data: [5, 4, 3, 5, 10],
            type: "line",
            stack: "x",
          },
        ],
      };
      // 复制
      myEChart.setOption(option);
      // 根据页面大小自动响应图表大小
      window.addEventListener("resize", function () {
        myEChart.resize();
      });
    },
  },
  mounted() {
    this.changeBar();
    this.changeLine();
  },
};

</script>


<style scoped>
#myEChartsBar {
  min-width: 31.25rem;
  min-height: 31.25rem;
}
#myEChartsLine {
  max-height: 500px;
  height: 500px;
}

</style>