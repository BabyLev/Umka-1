<template>
  <div ref="plotlyDiv" style="width: 100%; height: 400px;"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, toRefs, nextTick, computed } from 'vue';
import Plotly from 'plotly.js-dist-min';

interface Props {
  azimuth: number;
  elevation: number;
}

const props = defineProps<Props>();
const { azimuth, elevation } = toRefs(props);
const plotlyDiv = ref<HTMLDivElement | null>(null);

// Преобразование элевации в радиальную координату для Plotly:
// 90° (зенит) -> r = 0 (центр)
// 0° (горизонт) -> r = 90 (край)
const plotData = computed(() => {
    const radialValue = 90 - elevation.value;
    // Убедимся, что значение не выходит за пределы [0, 90]
    const clampedRadialValue = Math.max(0, Math.min(90, radialValue));

    return [{
        type: 'scatterpolar',
        mode: 'markers',
        r: [clampedRadialValue],
        theta: [azimuth.value],
        marker: {
            color: 'rgb(255, 65, 54)', // Красный маркер
            size: 12,
            symbol: 'circle',
        },
        name: 'Спутник',
    }] as Plotly.Data[]; // Указываем тип для Plotly.Data[]
});

const layout = computed<Partial<Plotly.Layout>>(() => ({ // Используем Partial<Plotly.Layout>
    polar: {
        angularaxis: {
            direction: "clockwise",
            thetaunit: "degrees",
            rotation: 90, // Север (0°) наверху
            tickmode: 'array',
            tickvals: [0, 45, 90, 135, 180, 225, 270, 315],
            ticktext: ['0° N', '45° NE', '90° E', '135° SE', '180° S', '225° SW', '270° W', '315° NW'],
            gridcolor: 'rgba(0,0,0,0.2)',
            linecolor: 'darkgrey',
        },
        radialaxis: {
             // Диапазон радиальной оси от 0 (центр) до 90 (край)
             range: [0, 90],
             // Значения тиков на радиальной оси (соответствуют 90-элевация)
             tickvals: [0, 15, 30, 45, 60, 75, 90],
             // Отображаемые метки - реальные значения элевации
             ticktext: ['90°', '75°', '60°', '45°', '30°', '15°', '0°'],
             // Угол для оси (не текста) - располагает ее вертикально
             angle: 90,
             // Угол для текста меток - делает их горизонтальными
             tickangle: 90,
             showline: true,
             showticklabels: true,
             gridcolor: 'rgba(0,0,0,0.2)',
             linecolor: 'darkgrey',
        },
        gridshape: 'linear', // Можно использовать 'circular'
        bgcolor: 'rgba(255, 255, 255, 0.9)',
    },
    showlegend: false,
    margin: { l: 60, r: 60, t: 50, b: 50 }, // Скорректированы отступы
    // title: 'Положение спутника (Азимут/Элевация)' // Заголовок можно добавить при необходимости
}));

async function renderPlot() {
  if (plotlyDiv.value && Plotly) {
    await nextTick(); // Убедимся, что div отрисован
    try {
      await Plotly.react(plotlyDiv.value, plotData.value, layout.value);
    } catch (error) {
        console.error("Ошибка рендеринга Plotly:", error);
    }
  } else {
      console.warn("Plotly div или библиотека Plotly не найдены.");
  }
}

onMounted(() => {
  renderPlot();
});

// Наблюдаем за изменением пропсов для обновления графика
watch([azimuth, elevation], () => {
  renderPlot();
}, { immediate: false }); // Не вызываем сразу при монтировании, т.к. onMounted уже это делает

</script>

<style scoped>
/* Стили при необходимости */
</style> 