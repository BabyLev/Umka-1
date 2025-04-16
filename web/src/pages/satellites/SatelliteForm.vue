<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-2xl font-semibold mb-6">{{ formTitle }}</h1>

    <form @submit.prevent="handleSubmit" class="bg-white shadow-md rounded px-8 pt-6 pb-8 mb-4">
      <!-- Сообщение об ошибке -->
      <div v-if="error" class="mb-4 bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
        <strong class="font-bold">Ошибка!</strong>
        <span class="block sm:inline"> {{ error }}</span>
      </div>

      <!-- Имя -->
      <div class="mb-4">
        <label class="block text-gray-700 text-sm font-bold mb-2" for="name">
          Имя спутника *
        </label>
        <input
          id="name"
          type="text"
          v-model="formData.name"
          required
          class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
          :class="{ 'border-red-500': !isNameValid && formSubmitted }"
        />
         <p v-if="!isNameValid && formSubmitted" class="text-red-500 text-xs italic">Имя обязательно.</p>
      </div>

      <!-- NORAD ID -->
      <div class="mb-4">
        <label class="block text-gray-700 text-sm font-bold mb-2" for="noradId">
          NORAD ID (если указан, TLE необязательны)
        </label>
        <input
          id="noradId"
          type="number"
          v-model.number="formData.noradId"
          placeholder="Можно получить TLE с r4uab по этому ID"
          class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
        />
      </div>

      <!-- Блок TLE (отображается условно) -->
      <div v-if="showTleFields">
        <!-- Line 1 TLE -->
        <div class="mb-4">
          <label class="block text-gray-700 text-sm font-bold mb-2" for="line1">
            TLE Line 1 *
          </label>
          <input
            id="line1"
            type="text"
            v-model="formData.line1"
            maxlength="69"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline font-mono"
            :class="{ 'border-red-500': !isLine1Valid && formSubmitted }"
            :required="showTleFields"
          />
          <p v-if="showTleFields && !isLine1Valid && formSubmitted" class="text-red-500 text-xs italic">TLE Line 1 обязательна (69 символов), если NORAD ID не указан.</p>
        </div>

        <!-- Line 2 TLE -->
        <div class="mb-6">
          <label class="block text-gray-700 text-sm font-bold mb-2" for="line2">
            TLE Line 2 *
          </label>
          <input
            id="line2"
            type="text"
            v-model="formData.line2"
            maxlength="69"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline font-mono"
            :class="{ 'border-red-500': !isLine2Valid && formSubmitted }"
            :required="showTleFields"
          />
          <p v-if="showTleFields && !isLine2Valid && formSubmitted" class="text-red-500 text-xs italic">TLE Line 2 обязательна (69 символов), если NORAD ID не указан.</p>
        </div>
      </div>

      <!-- Кнопки -->
      <div class="flex items-center justify-between">
        <button
          type="submit"
          :disabled="isLoading"
          class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline transition duration-150 ease-in-out disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <span v-if="isLoading">Сохранение...</span>
          <span v-else>{{ submitButtonText }}</span>
        </button>
        <router-link
          to="/satellites"
          class="inline-block align-baseline font-bold text-sm text-blue-500 hover:text-blue-800"
        >
          Отмена
        </router-link>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { storeToRefs } from 'pinia';
import { useSatelliteStore } from '@/stores/satelliteStore';
import type { SatelliteCreate } from '@/types'; // Импортируем нужный тип

const route = useRoute();
const router = useRouter();
const satelliteStore = useSatelliteStore();
const { isLoading, error } = storeToRefs(satelliteStore);

// Определяем режим работы (создание или редактирование)
const isEditMode = computed(() => route.meta.mode === 'edit');
const satelliteId = computed(() => (isEditMode.value ? Number(route.params.id) : null));

const formTitle = computed(() => (isEditMode.value ? 'Редактировать спутник' : 'Добавить новый спутник'));
const submitButtonText = computed(() => (isEditMode.value ? 'Сохранить изменения' : 'Добавить спутник'));

// Реактивная модель формы
const formData = ref<SatelliteCreate>({
  name: '',
  noradId: undefined, // Используем undefined для необязательного числового поля
  line1: '',
  line2: '',
});

const formSubmitted = ref(false); // Флаг для отображения ошибок валидации после попытки отправки

// Определяем, нужно ли показывать поля TLE
const showTleFields = computed(() => !formData.value.noradId);

// Сбрасываем TLE, если введен NORAD ID
watch(() => formData.value.noradId, (newNoradId) => {
  if (newNoradId) {
    formData.value.line1 = '';
    formData.value.line2 = '';
  }
});

// --- Валидация ---
const isNameValid = computed(() => !!formData.value.name?.trim());
// TLE валидны, если NORAD ID указан ИЛИ если TLE заполнены правильно (69 символов)
const isLine1Valid = computed(() => !!formData.value.noradId || formData.value.line1?.trim().length === 69);
const isLine2Valid = computed(() => !!formData.value.noradId || formData.value.line2?.trim().length === 69);
// Форма валидна, если имя валидно И (указан NORAD ID ИЛИ обе строки TLE валидны)
const isFormValid = computed(() => {
    const tleValid = isLine1Valid.value && isLine2Valid.value;
    // Форма валидна, если:
    // 1. Имя указано
    // 2. ИЛИ NORAD ID указан (и тогда TLE не важны)
    // 3. ИЛИ NORAD ID НЕ указан, НО TLE валидны
    return isNameValid.value && (!!formData.value.noradId || tleValid);
});

// Загрузка данных для редактирования
onMounted(async () => {
  if (isEditMode.value && satelliteId.value !== null) {
    let satellite = satelliteStore.getSatelliteById(satelliteId.value);
    if (!satellite) {
      satellite = await satelliteStore.fetchSatelliteById(satelliteId.value);
    }
    if (satellite) {
      formData.value = {
        name: satellite.name,
        noradId: satellite.noradId ?? undefined,
        line1: satellite.line1,
        line2: satellite.line2,
      };
    } else {
       console.error("[SatelliteForm] Satellite not found for editing, redirecting...");
       router.push('/satellites');
    }
  }
});

// Обработчик отправки формы
const handleSubmit = async () => {
  // Логируем состояние ПЕРЕД принятием решения
  console.log(
    '[SatelliteForm] handleSubmit triggered. Mode:', 
    { 
      isEditMode: isEditMode.value, 
      satelliteId: satelliteId.value 
    }
  );

  formSubmitted.value = true;
  if (!isFormValid.value) {
      console.log("[SatelliteForm] Form is invalid", {
          name: isNameValid.value,
          norad: !!formData.value.noradId,
          line1: isLine1Valid.value,
          line2: isLine2Valid.value
      });
      return;
  }

  // Создаем копию данных для отправки
  const payload: SatelliteCreate = { ...formData.value };
  if (payload.noradId) {
      payload.line1 = '';
      payload.line2 = '';
  }

  let success = false;
  // Проверяем условие еще раз прямо перед вызовом
  const shouldUpdate = isEditMode.value && satelliteId.value !== null;
  console.log('[SatelliteForm] Decision:', shouldUpdate ? 'UPDATE' : 'ADD');

  if (shouldUpdate) { // Используем переменную для ясности
    // Режим редактирования
    console.log('[SatelliteForm] Calling satelliteStore.updateSatellite...');
    success = await satelliteStore.updateSatellite({
        satelliteId: satelliteId.value!, // Уверенность в non-null из-за shouldUpdate
        satellite: payload
    });
  } else {
    // Режим создания
    console.log('[SatelliteForm] Calling satelliteStore.addSatellite...');
    const newSatelliteId = await satelliteStore.addSatellite(payload);
    success = newSatelliteId !== null;
  }

  if (success && !error.value) {
    console.log('[SatelliteForm] Operation successful, navigating to /satellites');
    router.push('/satellites');
  } else {
      console.error("[SatelliteForm] Operation failed.");
  }
};
</script>

<style scoped>
/* Стили для формы, если нужны */
.font-mono {
  font-family: monospace;
}
</style>
