import * as echarts from "echarts/core";
import {BarChart, LineChart, PieChart} from "echarts/charts";
import {
    TitleComponent,
    TooltipComponent,
    GridComponent,
    DatasetComponent,
    TransformComponent,
    ToolboxComponent,
    LegendComponent,
} from "echarts/components";
import {LabelLayout, UniversalTransition} from "echarts/features";
import {CanvasRenderer} from "echarts/renderers";

echarts.use([
    TitleComponent,
    TooltipComponent,
    GridComponent,
    DatasetComponent,
    TransformComponent,
    ToolboxComponent,
    LegendComponent,
    LabelLayout,
    UniversalTransition,
    CanvasRenderer,
    BarChart,
    LineChart,
    PieChart,
]);

export default echarts;

