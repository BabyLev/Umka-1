<template>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-6 p-4 md:p-6">
    <ElCard class="box-card">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="font-semibold text-lg">Расчет углов (азимут, элевация, расстояние)</span>
        </div>
      </template>
      <ElForm
        ref="formRef"
        :model="formModel"
        label-position="top"
        @submit.prevent="handleSubmit"
        v-loading="isLoading"
        class="space-y-4"
      >
        <ElFormItem label="Спутник" prop="satelliteId" required>
          <ElSelect
            v-model="formModel.satelliteId"
            placeholder="Выберите спутник"
            class="w-full"
            :disabled="satelliteStore.loading || !satelliteOptions.length"
          >
            <ElOption
              v-for="item in satelliteOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </ElSelect>
          <small v-if="!satelliteStore.loading && !satelliteOptions.length" class="text-red-500 mt-1">
            Нет доступных спутников. Добавьте спутник в разделе "Управление спутниками".
          </small>
        </ElFormItem>

        <ElFormItem label="Локация наблюдателя" prop="observerPositionId" required>
          <ElSelect
            v-model="formModel.observerPositionId"
            placeholder="Выберите локацию"
            class="w-full"
            :disabled="locationStore.loading || !locationOptions.length"
          >
            <ElOption
              v-for="item in locationOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </ElSelect>
          <small v-if="!locationStore.loading && !locationOptions.length" class="text-red-500 mt-1">
            Нет доступных локаций. Добавьте локацию в разделе "Локации".
          </small>
        </ElFormItem>

        <ElFormItem label="Время расчета" prop="timestamp">
          <ElDatePicker
            v-model="formModel.timestamp"
            type="datetime"
            placeholder="Текущее время"
            format="DD.MM.YYYY HH:mm:ss"
            value-format="YYYY-MM-DD HH:mm:ss"
            class="w-full"
          />
          <small class="text-gray-500 mt-1">Если оставить пустым, будет использовано текущее время.</small>
        </ElFormItem>

        <ElFormItem class="pt-2">
          <ElButton
            type="primary"
            native-type="submit"
            :loading="isLoading"
            :disabled="!formModel.satelliteId || !formModel.observerPositionId || isLoading"
          >
            Рассчитать
          </ElButton>
        </ElFormItem>
      </ElForm>
    </ElCard>

    <ElCard class="box-card" v-if="isLoading || error || result">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="font-semibold text-lg">Результаты расчета</span>
        </div>
      </template>
      <div v-if="isLoading" class="p-4">
        <ElSkeleton :rows="3" animated />
      </div>
      <ElAlert v-if="error" :title="error" type="error" show-icon :closable="false" class="m-4" />
      <div v-if="result && typeof (result as any).az === 'number' && typeof (result as any).el === 'number' && typeof (result as any).range === 'number'" class="p-4">
        <PlotlySkyPlot :azimuth="(result as any).az" :elevation="(result as any).el" />
        <p class="text-center mt-2"><strong>Расстояние:</strong> {{ (result as any).range.toFixed(2) }} км</p>
      </div>
      <div v-else-if="result" class="text-red-500 p-4">
        Ошибка: некорректный ответ от сервера.
      </div>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, reactive } from 'vue';
import {
  ElCard,
  ElForm,
  ElFormItem,
  ElSelect,
  ElOption,
  ElDatePicker,
  ElButton,
  ElAlert,
  ElSkeleton,
} from 'element-plus';
import { useSatelliteStore } from '@/stores/satelliteStore';
import { useLocationStore } from '@/stores/locationStore';
import { useCalculationStore } from '@/stores/calculationStore';
import type { LookAnglesResult } from '@/types';
import PlotlySkyPlot from '@/components/visualizations/PlotlySkyPlot.vue';

const satelliteStore = useSatelliteStore();
const locationStore = useLocationStore();
const calculationStore = useCalculationStore();

const formRef = ref();
const formModel = reactive({
  satelliteId: undefined as number | undefined,
  observerPositionId: undefined as number | undefined,
  timestamp: undefined as string | undefined,
});

const isLoading = computed(() => calculationStore.isLoading || satelliteStore.loading || locationStore.loading);
const error = computed(() => calculationStore.error);
const result = computed<LookAnglesResult | null>(() => calculationStore.lookAnglesResult);

const satelliteOptions = computed(() => satelliteStore.satellitesForSelect);
const locationOptions = computed(() => locationStore.locationsForSelect);

onMounted(() => {
  if (!satelliteStore.satellites.length) satelliteStore.fetchSatellites();
  if (!locationStore.locations.length) locationStore.fetchLocations();
});

function validateForm() {
  // Простая валидация (можно расширить)
  return !!formModel.satelliteId && !!formModel.observerPositionId;
}

async function handleSubmit() {
  calculationStore.clearResult('lookAngles');
  if (!validateForm()) return;
  let timestampUnix: number | undefined = undefined;
  if (formModel.timestamp) {
    try {
      timestampUnix = Math.floor(new Date(formModel.timestamp).getTime() / 1000);
    } catch (e) {
      // Ошибка парсинга времени
    }
  }
  await calculationStore.calculateLookAngles({
    satelliteId: formModel.satelliteId!,
    observerPositionId: formModel.observerPositionId!,
    timestamp: timestampUnix,
  });
}
</script>
