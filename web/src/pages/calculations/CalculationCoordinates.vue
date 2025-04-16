<template>
  <div class="grid grid-cols-1 md:grid-cols-2 gap-6 p-4 md:p-6">
    <ElCard class="box-card">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="font-semibold text-lg">Расчет координат спутника</span>
        </div>
      </template>
      <ElForm
        ref="formRef"
        :model="formModel"
        label-position="top"
        @submit.prevent="calculateCoordinates"
        v-loading="loading || satelliteStore.loading"
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
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </ElSelect>
          <small v-if="!satelliteStore.loading && !satelliteOptions.length" class="text-red-500 mt-1">
            Нет доступных спутников. Добавьте спутник в разделе "Управление спутниками".
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
            :loading="loading"
            :disabled="!formModel.satelliteId || loading || satelliteStore.loading"
          >
            Рассчитать
          </ElButton>
        </ElFormItem>
      </ElForm>
    </ElCard>

    <ElCard class="box-card" v-if="loading || error || results">
      <template #header>
        <div class="flex justify-between items-center">
          <span class="font-semibold text-lg">Результаты расчета</span>
        </div>
      </template>
      <div v-if="loading" class="p-4">
        <ElSkeleton :rows="4" animated />
      </div>

      <ElAlert v-if="error" :title="error" type="error" show-icon :closable="false" class="m-4" />

      <div v-if="results" class="space-y-3 p-4">
        <div class="mb-4" style="height: 300px;">
          <l-map
            ref="mapRef"
            :zoom="mapZoom"
            :center="mapCenter as any"
            :use-global-leaflet="false"
            class="rounded border border-gray-300"
            style="z-index: 0;"
          >
            <l-tile-layer
              url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
              layer-type="base"
              name="OpenStreetMap"
              :attribution="osmAttribution"
            ></l-tile-layer>
            <l-marker
              v-if="markerCoords"
              :lat-lng="markerCoords"
              :icon="satelliteIcon as unknown as import('leaflet').Icon<import('leaflet').IconOptions>"
            ></l-marker>
          </l-map>
        </div>

        <p><strong>Широта:</strong> {{ results.lat.toFixed(6) }}°</p>
        <p><strong>Долгота:</strong> {{ results.lon.toFixed(6) }}°</p>
        <p><strong>Высота:</strong> {{ results.alt.toFixed(3) }} км</p>
        <p v-if="results?.gmapsLink" class="flex items-center gap-2">
          <strong>Google Maps:</strong>
          <ElLink :href="results.gmapsLink" type="primary" target="_blank" rel="noopener noreferrer">
            {{ results.gmapsLink }}
          </ElLink>
        </p>
        <p v-else-if="results" class="flex items-center gap-2">
          <strong>Google Maps:</strong>
          <span class="text-gray-500">Ссылка недоступна</span>
        </p>
      </div>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, reactive, nextTick } from 'vue';
import {
  ElCard,
  ElForm,
  ElFormItem,
  ElSelect,
  ElOption,
  ElDatePicker,
  ElButton,
  ElAlert,
  ElLink,
  ElSkeleton,
} from 'element-plus';
import { useSatelliteStore } from '@/stores/satelliteStore';
import api from '@/api';
import type { SatellitePositionResponse } from '@/types/calculation';

// Импорты для Leaflet
import L from 'leaflet';
import "leaflet/dist/leaflet.css";
import { LMap, LTileLayer, LMarker, LIcon } from "@vue-leaflet/vue-leaflet";
import type { LatLngExpression } from "leaflet";

// Иконка спутника (SVG)
const satelliteIcon = L.divIcon({
  html: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="currentColor" class="w-8 h-8 text-blue-600"><path d="M12.001 3.00158C7.19277 3.00158 3.30029 6.5909 3.02573 11.2517L1.37448 9.60044L0.126953 10.8479L3.18745 13.9084L6.24795 10.8479L5.00043 9.60044L3.33762 11.2633C3.60543 7.30835 7.40042 4.25158 12.001 4.25158C16.8357 4.25158 20.751 8.16687 20.751 13.0016C20.751 17.8363 16.8357 21.7516 12.001 21.7516C7.40042 21.7516 3.60543 18.6948 3.33762 14.7398L5.00043 16.4026L6.24795 15.1551L3.18745 12.0946L0.126953 15.1551L1.37448 16.3969L3.02573 14.7457C3.30029 19.4065 7.19277 23.0016 12.001 23.0016C17.5238 23.0016 22.001 18.5248 22.001 13.0016C22.001 7.47835 17.5238 3.00158 12.001 3.00158ZM12 9.00158C10.3431 9.00158 9 10.3447 9 12.0016C9 13.6584 10.3431 15.0016 12 15.0016C13.6569 15.0016 15 13.6584 15 12.0016C15 10.3447 13.6569 9.00158 12 9.00158Z"/></svg>`,
  className: 'satellite-icon bg-transparent border-none', // Класс для возможной доп. стилизации + убираем фон/рамку по умолчанию
  iconSize: [32, 32],       // Размер иконки
  iconAnchor: [16, 16],     // Центрирование иконки
});

const satelliteStore = useSatelliteStore();

const formModel = reactive({
  satelliteId: undefined as number | undefined,
  timestamp: undefined as string | undefined,
});

const results = ref<SatellitePositionResponse | null>(null);
const loading = ref(false);
const error = ref<string | null>(null);

// --- Leaflet Map State ---
const mapRef = ref<any>(null); // Ссылка на компонент карты
const mapZoom = ref(3); // Начальный зум
const osmAttribution = '&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors';

// Центр карты, вычисляемый на основе результатов
const mapCenter = computed<LatLngExpression>(() => {
  if (results.value) {
    return [results.value.lat, results.value.lon];
  }
  return [0, 0]; // Центр по умолчанию, если нет результатов
});

// Координаты маркера спутника
const markerCoords = computed<LatLngExpression | null>(() => {
  if (results.value) {
    return [results.value.lat, results.value.lon];
  }
  return null; // Нет маркера, если нет результатов
});
// --- End Leaflet Map State ---

onMounted(() => {
  if (Object.keys(satelliteStore.satellites).length === 0) {
    satelliteStore.fetchSatellites();
  }
});

const satelliteOptions = computed(() => {
  return satelliteStore.satellites.map(sat => ({
    id: sat.id,
    name: `${sat.name} (ID: ${sat.id}, NORAD: ${sat.noradId || 'N/A'})`,
  }));
});

const calculateCoordinates = async () => {
  error.value = null;
  results.value = null;

  if (!formModel.satelliteId) {
    error.value = 'Пожалуйста, выберите спутник.';
    return;
  }

  console.log('Отправка запроса /calculate/ с satelliteId:', formModel.satelliteId);

  loading.value = true;

  let timestampUnix: number | undefined = undefined;
  if (formModel.timestamp) {
    try {
      timestampUnix = Math.floor(new Date(formModel.timestamp).getTime() / 1000);
    } catch (dateError) {
      console.error('Ошибка парсинга даты:', dateError);
      error.value = 'Некорректный формат времени.';
      loading.value = false;
      return;
    }
  }

  try {
    const response = await api.post<SatellitePositionResponse>('/calculate/', {
      satelliteId: formModel.satelliteId,
      timestamp: timestampUnix,
    });
    console.log('Ответ API /calculate/:', response.data);
    results.value = response.data;

    // После получения результатов, установим вид карты
    // Используем nextTick, чтобы карта успела отрендериться, если она была скрыта
    nextTick(() => {
      if (mapRef.value?.leafletObject && results.value) {
        const viewZoom = 5; // Увеличим зум при отображении спутника
        mapRef.value.leafletObject.setView([results.value.lat, results.value.lon], viewZoom);
      }
    });

  } catch (err: any) {
    console.error('Ошибка при расчете координат:', err);
    error.value =
      err.response?.data?.detail ||
      err.response?.data?.error ||
      err.message ||
      'Произошла неизвестная ошибка.';
  } finally {
    loading.value = false;
  }
};
</script>

<style>
/* Стили для иконки спутника */
.satellite-icon svg {
  width: 32px;  /* Синхронизируем с iconSize */
  height: 32px; /* Синхронизируем с iconSize */
  /* Можно добавить drop-shadow или другие стили */
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.3)); /* Добавим тень для лучшей видимости */
}

/* Убираем фон и рамку, которые Leaflet может добавлять к divIcon по умолчанию */
.leaflet-div-icon {
  background: none !important;
  border: none !important;
}
</style>
