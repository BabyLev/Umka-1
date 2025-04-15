<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-2xl font-semibold mb-6">{{ pageTitle }}</h1>

    <form @submit.prevent="handleSubmit" class="max-w-lg mx-auto bg-white p-8 rounded shadow-md">

      <div v-if="formError" class="mb-4 bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
        <strong class="font-bold">Ошибка!</strong>
        <span class="block sm:inline"> {{ formError }}</span>
      </div>

      <!-- Имя локации -->
      <div class="mb-4">
        <label for="name" class="block text-gray-700 text-sm font-bold mb-2">Имя локации:</label>
        <input
          type="text"
          id="name"
          v-model="formData.name"
          required
          class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          :class="{ 'border-red-500': validationErrors.name }"
        />
        <p v-if="validationErrors.name" class="text-red-500 text-xs italic">{{ validationErrors.name }}</p>
      </div>

      <!-- Карта Leaflet -->
      <div class="mb-4" style="height: 300px;"> 
        <l-map
          ref="mapRef"
          :zoom="mapZoom"
          :center="mapCenter as any"
          :use-global-leaflet="false" 
          @click="handleMapClick"
          class="rounded"
        >
          <l-tile-layer
            url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            layer-type="base"
            name="OpenStreetMap"
            :attribution="osmAttribution"
          ></l-tile-layer>
          <l-marker 
            :lat-lng="markerCoords" 
            :draggable="true" 
            @update:lat-lng="handleMarkerDrag"
          ></l-marker>
        </l-map>
        <p class="text-xs text-gray-600 mt-1">Кликните на карту или перетащите маркер для выбора координат.</p>
      </div>

      <!-- Координаты -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
        <!-- Долгота -->
        <div>
          <label for="lon" class="block text-gray-700 text-sm font-bold mb-2">Долгота (°):</label>
          <input
            type="number"
            step="any"
            id="lon"
            v-model.number="formData.location.lon"
            required
            @input="updateMapFromInput"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
             :class="{ 'border-red-500': validationErrors.lon }"
          />
           <p v-if="validationErrors.lon" class="text-red-500 text-xs italic">{{ validationErrors.lon }}</p>
        </div>

        <!-- Широта -->
        <div>
          <label for="lat" class="block text-gray-700 text-sm font-bold mb-2">Широта (°):</label>
          <input
            type="number"
            step="any"
            id="lat"
            v-model.number="formData.location.lat"
            required
            @input="updateMapFromInput"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
             :class="{ 'border-red-500': validationErrors.lat }"
          />
          <p v-if="validationErrors.lat" class="text-red-500 text-xs italic">{{ validationErrors.lat }}</p>
        </div>

        <!-- Высота -->
        <div class="mb-6 md:mb-0">  
          <label for="alt" class="block text-gray-700 text-sm font-bold mb-2">Высота (км):</label>
          <input
            type="number"
            step="any"
            id="alt"
            v-model.number="formData.location.alt"
            required
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            :class="{ 'border-red-500': validationErrors.alt }"
          />
          <p v-if="validationErrors.alt" class="text-red-500 text-xs italic">{{ validationErrors.alt }}</p>
        </div>
      </div>

      <!-- Кнопки -->
      <div class="flex items-center justify-between">
        <button
          type="submit"
          :disabled="isLoading"
          class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline disabled:opacity-50"
        >
          {{ isLoading ? 'Сохранение...' : 'Сохранить' }}
        </button>
        <router-link
          to="/locations"
          class="inline-block align-baseline font-bold text-sm text-blue-500 hover:text-blue-800"
        >
          Отмена
        </router-link>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, reactive, watch, nextTick } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useLocationStore } from '@/stores/locationStore';
import type { LocationCreate, LocationUpdatePayload, Location } from '@/types';

// Импорты для Leaflet и поиска
import L from 'leaflet';
import "leaflet/dist/leaflet.css";
import { LMap, LTileLayer, LMarker } from "@vue-leaflet/vue-leaflet";
import type { LatLngExpression, LeafletMouseEvent, PointExpression } from "leaflet";
import "leaflet-search/dist/leaflet-search.min.css";
import "leaflet-search";

const route = useRoute();
const router = useRouter();
const locationStore = useLocationStore();

const mode = computed(() => route.meta.mode as 'new' | 'edit');
const locationId = computed(() => route.params.id ? Number(route.params.id) : null);

const pageTitle = computed(() => mode.value === 'new' ? 'Добавить новую локацию' : 'Редактировать локацию');

// Форма для v-model. Используем reactive для вложенной структуры location
const formData = reactive<LocationCreate>({ 
  name: '',
  location: {
    lon: 0,
    lat: 0,
    alt: 0,
  }
});

const isLoading = ref(false);
const formError = ref<string | null>(null); // Ошибка от API
const validationErrors = reactive<Record<string, string>>({}); // Ошибки валидации полей

// --- Leaflet Map State ---
const mapRef = ref<any>(null); // Ссылка на компонент карты
const mapZoom = ref(3); // Начальный зум
const initialCenter: LatLngExpression = [50, 10]; // Начальный центр карты
const osmAttribution = '&copy; <a href="http://osm.org/copyright">OpenStreetMap</a> contributors';

// --- Добавим состояние загрузки данных для формы ---
const isFetchingLocation = ref(false);
// ---

// Вычисляемый центр карты (реагирует на изменение formData)
const mapCenter = computed<LatLngExpression>(() => [
  formData.location.lat || initialCenter[0],
  formData.location.lon || initialCenter[1],
]);

// Координаты маркера (реагируют на изменение formData)
const markerCoords = computed<LatLngExpression>(() => [
  formData.location.lat,
  formData.location.lon
]);

// Обновляем formData при клике на карту
function handleMapClick(event: LeafletMouseEvent) {
  if (event.latlng) {
    formData.location.lat = parseFloat(event.latlng.lat.toFixed(6));
    formData.location.lon = parseFloat(event.latlng.lng.toFixed(6));
  }
}

// Обновляем formData при перетаскивании маркера
function handleMarkerDrag(event: { lat: number; lng: number }) { // Тип события от LMarker
   formData.location.lat = parseFloat(event.lat.toFixed(6));
   formData.location.lon = parseFloat(event.lng.toFixed(6));
}

// Обновляем центр карты при изменении координат в полях ввода 
function updateMapFromInput() {
  // Карта автоматически центрируется через computed `mapCenter`
  // Дополнительно можно плавно переместить вид, если нужно:
  // if (mapRef.value?.leafletObject) {
  //   mapRef.value.leafletObject.panTo([formData.location.lat, formData.location.lon]);
  // }
}
// --- End Leaflet Map State ---

// Загрузка данных для режима редактирования и ИНИЦИАЛИЗАЦИЯ КАРТЫ И ПОИСКА
onMounted(async () => { // Делаем onMounted асинхронным
  isFetchingLocation.value = true; // Начинаем загрузку данных формы
  formError.value = null; // Сбрасываем предыдущие ошибки

  // Функция для инициализации карты и поиска (вынесена для читаемости)
  const initializeMapAndSearch = () => {
    const map = mapRef.value.leafletObject; // Получаем объект карты Leaflet
    if (!map) return; // Убедимся, что карта есть

    // Инициализация контрола поиска
    const searchControl = new (L.Control as any).Search({
      url: 'https://nominatim.openstreetmap.org/search?format=json&q={s}', // URL для Nominatim
      jsonpParam: 'json_callback',
      propertyName: 'display_name', // Ключ в ответе JSON с названием места
      propertyLoc: ['lat', 'lon'],    // Ключи с координатами
      marker: false, // Не будем добавлять отдельный маркер поиска, используем наш основной
      autoCollapse: true,          // Сворачивать после поиска
      autoType: false,
      minLength: 2,                // Минимальная длина запроса
      zoom: 12,                   // Зум при нахождении места
      textPlaceholder: 'Поиск города/места...', // Текст в поле ввода
      textErr: 'Место не найдено',
    });

    // Обработчик события нахождения места
    searchControl.on('search:locationfound', (e: any) => {
      console.log('Search result:', e);
      if (e.latlng) {
         // Обновляем наши данные формы
         formData.location.lat = parseFloat(e.latlng.lat.toFixed(6));
         formData.location.lon = parseFloat(e.latlng.lng.toFixed(6));
         // Можно и имя обновить, если нужно
         // formData.name = e.layer.options.title || formData.name;
      }
    });

    map.addControl(searchControl); // Добавляем контрол на карту
  };

  // Ждем, пока карта будет готова, затем инициализируем
  const checkMapInterval = setInterval(() => {
    if (mapRef.value?.leafletObject) {
      clearInterval(checkMapInterval);
      initializeMapAndSearch(); // Инициализируем карту и поиск

      // --- Логика загрузки данных для редактирования ---
      if (mode.value === 'edit' && locationId.value !== null) {
        loadLocationForEdit(locationId.value); // Вызываем функцию загрузки
      } else {
         isFetchingLocation.value = false; // Если не режим редактирования, загрузка не нужна
      }
      // --- Конец логики загрузки ---
    }
  }, 100); // Проверяем каждые 100мс
});

// --- Новая функция для загрузки данных локации ---
async function loadLocationForEdit(id: number) {
  isFetchingLocation.value = true;
  formError.value = null;
  try {
    // --- Используем новое действие для загрузки одной локации ---
    // console.log('Редактирование: Загрузка локации ID:', id);
    const locationToEdit = await locationStore.fetchLocationById(id);
    // console.log('Редактирование: Загруженные данные:', locationToEdit);

    // --- Убираем старую логику загрузки списка и поиска в нем ---
    // // Сначала загружаем/обновляем список локаций
    // await locationStore.fetchLocations();
    // if (locationStore.error) {
    //     throw new Error(locationStore.error); // Если при загрузке списка была ошибка
    // }
    // // Затем получаем нужную локацию из стора
    // const locationToEdit = locationStore.getLocationById(id);
    // console.log('Редактирование: Загруженные данные:', locationToEdit);
    // --- Конец убранной логики ---

    if (locationToEdit) {
      // Заполняем formData
      formData.name = locationToEdit.name;
      formData.location.lon = locationToEdit.lon;
      formData.location.lat = locationToEdit.lat;
      formData.location.alt = locationToEdit.alt;
      // console.log('Редактирование: formData обновлен:', JSON.parse(JSON.stringify(formData.location))); // Логируем копию

      // Используем watch или nextTick, чтобы убедиться, что computed свойство обновилось
      // и карта готова перед центрированием
      await nextTick(); // Дождемся обновления DOM/computed свойств
      // console.log('Редактирование: Вычисленные координаты маркера:', markerCoords.value);

      // Устанавливаем центр и зум карты (после того как formData обновлен)
       if (mapRef.value?.leafletObject) {
          // Увеличим зум для редактирования
          const editZoom = 10;
          const centerCoords = mapCenter.value; // Используем computed свойство
          // console.log('Редактирование: Центрирование карты на:', centerCoords, 'с зумом:', editZoom);
          mapRef.value.leafletObject.setView(centerCoords, editZoom);
       }

    } else {
      // Локация не найдена (или ошибка загрузки)
      formError.value = locationStore.error || `Локация с ID ${id} не найдена.`;
      console.error(formError.value);
      // Возможно, стоит перенаправить пользователя или показать сообщение
      // router.push('/locations'); // Например, перенаправить
    }
  } catch (error: any) {
    console.error('Ошибка при загрузке данных локации:', error);
    // Дополнительно проверяем ошибку из стора, если она там установлена
    formError.value = locationStore.error || error.message || 'Не удалось загрузить данные локации.';
  } finally {
    isFetchingLocation.value = false;
  }
}
// --- Конец новой функции ---

// Валидация формы
function validateForm(): boolean {
  Object.keys(validationErrors).forEach(key => delete validationErrors[key]); // Очистка старых ошибок
  let isValid = true;

  if (!formData.name.trim()) {
    validationErrors.name = 'Имя обязательно для заполнения.';
    isValid = false;
  }
  if (formData.location.lon < -180 || formData.location.lon > 180) {
      validationErrors.lon = 'Долгота должна быть между -180 и 180.';
      isValid = false;
  }
   if (formData.location.lat < -90 || formData.location.lat > 90) {
      validationErrors.lat = 'Широта должна быть между -90 и 90.';
      isValid = false;
  }
   if (formData.location.alt < 0) {
      validationErrors.alt = 'Высота не может быть отрицательной.';
      isValid = false;
  }
  // Добавить другие проверки при необходимости (isNan и т.д.)

  return isValid;
}

// Обработка отправки формы
async function handleSubmit() {
  formError.value = null;
  if (!validateForm()) {
    return; 
  }

  isLoading.value = true;
  try {
    if (mode.value === 'new') {
      const newId = await locationStore.addLocation(formData);
      if (newId) {
        router.push('/locations'); // Переход на список после успеха
      } else {
        formError.value = locationStore.error || 'Не удалось добавить локацию.';
      }
    } else if (mode.value === 'edit' && locationId.value !== null) {
      // TODO: Implement edit logic
      const payload: LocationUpdatePayload = {
        locationId: locationId.value,
        location: {
          name: formData.name,
          location: formData.location
        }
      };
      const success = await locationStore.updateLocation(payload);
      if (success) {
        router.push('/locations');
      } else {
        formError.value = locationStore.error || 'Не удалось обновить локацию.';
      }
    }
  } catch (err: any) {
    console.error('Submit error:', err);
    formError.value = err.message || 'Произошла неизвестная ошибка.';
  } finally {
    isLoading.value = false;
  }
}

</script>

<style scoped>
/* Override leaflet-search button icon to use SVG */
:deep(.leaflet-control-search .search-button) { 
  background-image: none !important; /* Remove default image */
  width: 30px; /* Adjust size as needed */
  height: 30px; /* Adjust size as needed */
  /* Use SVG as a mask */
  background-color: #6b7280; /* Tailwind gray-500 */
  -webkit-mask-image: url('data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z" /></svg>');
  mask-image: url('data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke-width="1.5" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" d="m21 21-5.197-5.197m0 0A7.5 7.5 0 1 0 5.196 5.196a7.5 7.5 0 0 0 10.607 10.607Z" /></svg>');
  background-repeat: no-repeat;
  background-position: center;
  background-size: 18px 18px; /* Adjust icon size */
  border-radius: 4px; /* Optional: match other leaflet controls */
}

/* Optional: Style on hover */
:deep(.leaflet-control-search .search-button:hover) {
  background-color: #374151; /* Tailwind gray-700 */
}

/* Ensure the search control itself is visible */
:deep(.leaflet-control-search) {
  /* Add any necessary styles if the whole control is hidden */
  /* Example: */
  /* box-shadow: 0 1px 5px rgba(0,0,0,0.65); */
  /* background: #fff; */
}

/* Adjust input field style if needed */
:deep(.leaflet-control-search .search-input) {
   height: 30px;
   padding: 0 10px;
   border: 1px solid #ccc;
   box-shadow: none;
}

</style>
