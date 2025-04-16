<template>
  <div class="p-4 md:p-6 space-y-6">
    <h1 class="text-2xl font-semibold text-gray-800">Расчет Видимости Спутника</h1>

    <ElCard class="box-card max-w-2xl">
       <ElForm
         ref="formRef" 
         :model="formModel" 
         label-position="top"
         @submit.prevent="calculateVisibility"
         v-loading="isLoadingExternal"
         class="space-y-4"
       >
        <!-- Спутник -->
        <ElFormItem label="Спутник" prop="satelliteId" required>
          <ElSelect
            v-model="formModel.satelliteId"
            placeholder="Выберите спутник"
            class="w-full"
            filterable  
            :disabled="isLoadingExternal || !satelliteOptions.length"
          >
            <ElOption
              v-for="item in satelliteOptions" 
              :key="item.value" 
              :label="item.label" 
              :value="item.value"
            />
          </ElSelect>
            <small v-if="!satelliteStore.loading && !satelliteOptions.length && !satelliteStore.error" class="text-gray-500 mt-1">
              Нет доступных спутников для выбора.
          </small>
           <small v-if="satelliteStore.error" class="text-red-500 mt-1">
              Ошибка загрузки спутников: {{ satelliteStore.error }}
          </small>
        </ElFormItem>

        <!-- Время -->
         <ElFormItem label="Начальное время поиска" prop="timestamp">
           <ElDatePicker
            v-model="formModel.timestamp"
            type="datetime"
            placeholder="Текущее время"
            format="DD.MM.YYYY HH:mm:ss"
            value-format="x" 
            class="w-full"
           />
           <small class="text-gray-500 mt-1">Если оставить пустым, будет использовано текущее время.</small>
         </ElFormItem>

        <!-- Локация -->
         <ElFormItem label="Локация наблюдателя" prop="selectedLocationId">
           <ElSelect
             v-model="formModel.selectedLocationId"
             placeholder="Выберите локацию или введите вручную"
             class="w-full"
             :disabled="isLoadingExternal || !locationOptions.length"
           >
              <ElOption
                v-for="item in locationOptions" 
                :key="item.value" 
                :label="item.label" 
                :value="item.value"
              />
           </ElSelect>
            <small v-if="!locationStore.loading && locationOptions.length <= 1 && !locationStore.error" class="text-gray-500 mt-1">
              Нет сохраненных локаций.
            </small>
            <small v-if="locationStore.error" class="text-red-500 mt-1">
                Ошибка загрузки локаций: {{ locationStore.error }}
            </small>
         </ElFormItem>

        <!-- Координаты (ручной ввод) -->
         <div v-if="formModel.selectedLocationId === -1" class="grid grid-cols-1 md:grid-cols-3 gap-4 border border-gray-300 border-dashed p-4 rounded mt-2">
           <p class="md:col-span-3 text-sm text-gray-500 mb-2">Введите координаты вручную:</p>
           <ElFormItem label="Долгота (°)" prop="observerLon" required>
             <ElInputNumber
               v-model="formModel.observerLon"
               controls-position="right"
               :precision="6" 
               placeholder="37.6173"
               class="w-full"
             />
           </ElFormItem>
           <ElFormItem label="Широта (°)" prop="observerLat" required>
              <ElInputNumber
               v-model="formModel.observerLat"
               controls-position="right"
               :precision="6"
               placeholder="55.7558"
               class="w-full"
             />
           </ElFormItem>
           <ElFormItem label="Высота (км)" prop="observerAlt" required>
             <ElInputNumber
               v-model="formModel.observerAlt"
               controls-position="right"
               :precision="4" 
               placeholder="0.15"
               class="w-full"
             />
           </ElFormItem>
         </div>
          <!-- Отображение координат выбранной локации -->
         <div v-else-if="selectedLocationData" class="text-sm text-gray-600 mt-2 border border-gray-300 p-3 rounded bg-gray-50">
           Выбрана локация: <strong class="font-medium">{{ selectedLocationData.name }}</strong>
           (Lon: {{ selectedLocationData.lon }}, Lat: {{ selectedLocationData.lat }}, Alt: {{ selectedLocationData.alt }} km)
         </div>

        <!-- Количество интервалов -->
         <ElFormItem label="Количество интервалов" prop="countOfTimeRanges" required>
            <ElInputNumber
              v-model="formModel.countOfTimeRanges"
              :min="1"
              :max="100" 
              controls-position="right"
              placeholder="1"
              class="w-full"
            />
         </ElFormItem>

        <!-- Кнопки -->
         <ElFormItem class="pt-2">
           <ElButton 
             type="primary" 
             native-type="submit" 
             :loading="isLoading" 
             :disabled="isSubmitDisabled"
            >
             Рассчитать
           </ElButton>
           <ElButton 
             @click="resetForm" 
             :disabled="isLoading"
             class="ml-2"
            >
             Сбросить
           </ElButton>
         </ElFormItem>
       </ElForm>
     </ElCard>

     <!-- Результаты -->
     <ElCard v-if="calculationPerformed || isLoading" class="box-card max-w-2xl">
         <template #header>
             <div class="flex justify-between items-center">
             <span class="font-semibold text-lg">Результаты расчета</span>
             </div>
         </template>
        <div v-if="isLoading" class="p-4">
            <ElSkeleton :rows="5" animated />
        </div>
         <ElAlert v-else-if="error" :title="error" type="error" show-icon :closable="true" @close="error=null" class="m-4" />
         <div v-else-if="visibilityRanges && visibilityRanges.length > 0" class="p-4 space-y-4">
            <!-- Интерактивная временная шкала -->
            <div ref="timelineContainer" class="border border-gray-300 rounded shadow-sm" style="height: 400px;"></div>

            <!-- Детали выбранного интервала -->
            <div v-if="selectedTimelineItem && selectedTimelineItem.start && selectedTimelineItem.end" class="mt-4 p-3 border border-gray-200 rounded bg-gray-50 text-sm">
                <p class="font-medium text-gray-800 mb-1">Выбранный интервал:</p>
                <p><span class="font-semibold">С:</span> {{ formatDateTime(new Date(selectedTimelineItem.start)) }}</p>
                <p><span class="font-semibold">До:</span> {{ formatDateTime(new Date(selectedTimelineItem.end)) }}</p>
                <p><span class="font-semibold">Длительность:</span> {{ calculateDuration(new Date(selectedTimelineItem.start), new Date(selectedTimelineItem.end)) }}</p>
                <!-- TODO: Add Max Elevation if available -->
            </div>
             <div v-else class="mt-4 text-center text-gray-500 text-sm">
                Кликните на интервал на шкале для просмотра деталей.
            </div>
        </div>
        <div v-else class="mt-6 text-center text-gray-500 py-4">
            Интервалы видимости не найдены для заданных параметров.
        </div>
     </ElCard>

  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch, reactive, nextTick, onBeforeUnmount } from 'vue';
import {
  ElCard,
  ElForm,
  ElFormItem,
  ElSelect,
  ElOption,
  ElDatePicker,
  ElInputNumber,
  ElButton,
  ElAlert,
  ElSkeleton,
} from 'element-plus';
import { useSatelliteStore } from '@/stores/satelliteStore';
import { useLocationStore } from '@/stores/locationStore';
import { useCalculationStore } from '@/stores/calculationStore';
import type { TimeRange, Location, TimeRangesParams } from '@/types';
import { formatDateTime } from '@/utils/dateFormatter';
import { calculateDuration } from '@/utils/timelineUtils';
import { Timeline, DataSet } from 'vis-timeline/standalone';
import type { DataItem, TimelineOptions, IdType } from 'vis-timeline/standalone';
import 'vis-timeline/styles/vis-timeline-graph2d.min.css';

const satelliteStore = useSatelliteStore();
const locationStore = useLocationStore();
const calculationStore = useCalculationStore();
const formRef = ref<any>(null);

// Используем reactive для модели формы
const formModel = reactive({
  satelliteId: undefined as number | undefined,
  timestamp: undefined as number | undefined,
  selectedLocationId: undefined as number | undefined,
  observerLon: undefined as number | undefined,
  observerLat: undefined as number | undefined,
  observerAlt: undefined as number | undefined,
  countOfTimeRanges: 1,
});

// Access state properties directly from $state
const isLoading = computed(() => calculationStore.$state.isLoading);
const error = ref<string | null>(null);
const calculationPerformed = ref(false);
const selectedTimelineItem = ref<DataItem | null>(null);

// Вычисляемое состояние для общей загрузки данных (спутники, локации)
const isLoadingExternal = computed(() => satelliteStore.loading || locationStore.loading);

// Загрузка данных при монтировании
onMounted(() => {
  satelliteStore.fetchSatellites();
  locationStore.fetchLocations();
});

// Используем геттеры из сторов
const satelliteOptions = computed(() =>
  satelliteStore.satellites.map(sat => ({ value: sat.id, label: sat.name }))
);

const locationOptions = computed(() => [
  { value: -1, label: 'Ввести координаты вручную' },
  ...locationStore.locations.map(loc => ({ value: loc.id, label: loc.name }))
]);

// Данные выбранной локации для отображения
const selectedLocationData = computed(() => {
  if (formModel.selectedLocationId === null || formModel.selectedLocationId === -1) {
    return null;
  }
  return locationStore.locations.find(loc => loc.id === formModel.selectedLocationId);
});

// Следим за выбором локации и обновляем координаты, если выбрана сохраненная
watch(() => formModel.selectedLocationId, (newId) => {
  const loc = locationStore.locations.find(l => l.id === newId);
  if (loc) {
    formModel.observerLon = loc.lon;
    formModel.observerLat = loc.lat;
    formModel.observerAlt = loc.alt;
  } else if (newId !== -1 && newId !== undefined) {
      formModel.observerLon = undefined;
      formModel.observerLat = undefined;
      formModel.observerAlt = undefined;
  }
});

const isManualInput = computed(() => formModel.selectedLocationId === -1);

// Проверка, можно ли отправлять форму
const isSubmitDisabled = computed(() => {
  if (isLoading.value || isLoadingExternal.value) return true;
  if (!formModel.satelliteId) return true;
  if (isManualInput.value) {
    return formModel.observerLon === undefined || formModel.observerLat === undefined || formModel.observerAlt === undefined;
  } else {
    // Ensure a location (not manual or undefined) is selected
    return formModel.selectedLocationId === undefined || formModel.selectedLocationId === -1;
  }
});

// Access results property directly from $state
const visibilityRanges = computed(() => calculationStore.$state.timeRangesResult);

// Обработчик отправки
const calculateVisibility = async () => {
  error.value = null;
  calculationPerformed.value = true;
  selectedTimelineItem.value = null;
  if (timeline) {
    timeline.destroy();
    timeline = null;
  }

  await formRef.value?.validate(async (valid: boolean) => {
    if (!valid) {
      error.value = "Пожалуйста, заполните все обязательные поля корректно.";
      return;
    }

    let obsCoords: { lon: number; lat: number; alt: number; } | undefined;

    if (isManualInput.value) {
        if (typeof formModel.observerLon === 'number' && typeof formModel.observerLat === 'number' && typeof formModel.observerAlt === 'number') {
             obsCoords = { lon: formModel.observerLon, lat: formModel.observerLat, alt: formModel.observerAlt };
        }
    } else if (selectedLocationData.value) {
         obsCoords = { lon: selectedLocationData.value.lon, lat: selectedLocationData.value.lat, alt: selectedLocationData.value.alt };
    }

    if (typeof formModel.satelliteId !== 'number' || !obsCoords) {
      error.value = 'Не удалось определить координаты наблюдателя или спутник.';
      return;
    }

    try {
        // Передаем параметры в соответствии с определением TimeRangesParams
        const params: TimeRangesParams = {
            satelliteId: formModel.satelliteId,
            lon: obsCoords.lon,         // Поле lon напрямую
            lat: obsCoords.lat,         // Поле lat напрямую
            alt: obsCoords.alt,         // Поле alt напрямую
            timestamp: formModel.timestamp ? Math.floor(Number(formModel.timestamp) / 1000) : undefined, // timestamp в секундах
            countOfTimeRanges: formModel.countOfTimeRanges,
        };

        await calculationStore.calculateTimeRanges(params);

        if (calculationStore.$state.error) {
            error.value = `Ошибка расчета: ${calculationStore.$state.error}`;
        } else {
            await nextTick();
            createTimeline();
        }
    } catch (err: any) {
        console.error("Calculation failed:", err);
        error.value = calculationStore.$state.error || `Неизвестная ошибка при расчете: ${err.message || err}`;
    }
  });
};

// Сброс формы
const resetForm = () => {
    formRef.value?.resetFields();
    // Reset reactive form model properties to undefined
    formModel.satelliteId = undefined;
    formModel.timestamp = undefined;
    formModel.selectedLocationId = undefined;
    formModel.observerLon = undefined;
    formModel.observerLat = undefined;
    formModel.observerAlt = undefined;
    formModel.countOfTimeRanges = 1; // Reset count to default

    calculationPerformed.value = false;
    error.value = null;
    calculationStore.$reset(); // Use $reset for Pinia store reset
    selectedTimelineItem.value = null;
    if (timeline) {
        timeline.destroy();
        timeline = null;
    }
};

const timelineContainer = ref<HTMLElement | null>(null);
let timeline: Timeline | null = null;

const createTimeline = () => {
  // Check if visibilityRanges.value is null or empty
   if (!timelineContainer.value || !visibilityRanges.value || visibilityRanges.value.length === 0) {
     if (timeline) {
        timeline.destroy();
        timeline = null;
    }
    return;
  }
    if (timeline) {
        timeline.destroy();
        timeline = null;
    }

  // Map from store state
  const items = new DataSet<DataItem>(
    // Add null check for visibilityRanges.value before mapping
    (visibilityRanges.value || []).map((range: TimeRange, index: number) => ({
      id: index,
      content: `Интервал ${index + 1}`,
      start: new Date(range.from),
      end: new Date(range.to),
    }))
  );

  const options: TimelineOptions = {
    stack: true,
    zoomable: true,
    moveable: true,
    selectable: true,
    margin: { item: { vertical: 5, horizontal: 2 } },
    orientation: 'top',
    height: '100%',
    maxHeight: 400,
    tooltip: {
            followMouse: true,
            overflowMethod: 'flip'
        },
        type: 'range',
  };

  timeline = new Timeline(timelineContainer.value, items, options);

  timeline.on('select', (properties: { items: IdType[] }) => {
    const selectedId = properties.items[0];
    if (selectedId !== undefined) {
      const selected = items.get(selectedId);
      selectedTimelineItem.value = selected;
    } else {
      selectedTimelineItem.value = null;
    }
  });
};

watch(visibilityRanges, async (newRanges) => {
    if (newRanges && newRanges.length > 0) {
       await nextTick();
        createTimeline();
    } else {
         if (timeline) {
            timeline.destroy();
            timeline = null;
        }
         selectedTimelineItem.value = null;
    }
}, { deep: true });

// Clean up timeline instance when component is unmounted
onBeforeUnmount(() => {
    if (timeline) {
        timeline.destroy();
    }
});

</script>

<style scoped>
/* Add any specific styles for timeline or other elements if needed */
.box-card {
  /* Maybe add transition for appearance */
}

/* Improve vis-timeline default styling if needed */
:deep(.vis-item) {
  background-color: #3b82f6; /* Example: Blue-500 */
  border-color: #1d4ed8; /* Example: Blue-700 */
  color: white;
  border-radius: 4px;
}
:deep(.vis-item.vis-selected) {
  background-color: #1d4ed8; /* Darker blue when selected */
  border-color: #1e3a8a; /* Even darker border */
}

:deep(.vis-timeline) {
    border: 1px solid #e5e7eb; /* Match border style */
     border-radius: 0.25rem; /* Match rounded style */
}

</style>
