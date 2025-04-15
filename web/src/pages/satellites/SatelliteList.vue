<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-2xl font-semibold mb-4">Список спутников</h1>

    <div class="mb-4 flex justify-between items-center">
      <!-- Поиск -->
      <div class="relative w-1/3">
        <input
          type="text"
          v-model="searchName"
          @input="debouncedSearch"
          placeholder="Поиск по имени..."
          class="w-full px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
        <span v-if="isLoading && isSearching" class="absolute right-3 top-1/2 transform -translate-y-1/2">
           <!-- Иконка загрузки (можно заменить на spinner) -->
           <svg class="animate-spin h-5 w-5 text-gray-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
             <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
             <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
           </svg>
        </span>
      </div>

      <!-- Кнопка добавления -->
      <router-link
        to="/satellites/new"
        class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded transition duration-150 ease-in-out"
      >
        Добавить спутник
      </router-link>
    </div>

    <!-- Индикатор загрузки -->
    <div v-if="isLoading && !isSearching" class="text-center py-10">
      <p>Загрузка спутников...</p>
      <!-- Можно добавить spinner -->
    </div>

    <!-- Сообщение об ошибке -->
    <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
      <strong class="font-bold">Ошибка!</strong>
      <span class="block sm:inline"> {{ error }}</span>
    </div>

    <!-- Таблица со спутниками -->
    <div v-else-if="satellites.length > 0" class="bg-white shadow-md rounded my-6 overflow-x-auto">
      <table class="min-w-full leading-normal">
        <thead>
          <tr>
            <th class="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
              ID
            </th>
            <th class="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
              Имя
            </th>
            <th class="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">
              NORAD ID
            </th>
            <th class="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider max-w-md">
              TLE
            </th>
            <th class="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-center text-xs font-semibold text-gray-600 uppercase tracking-wider">
              Действия
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="satellite in satellites" :key="satellite.id">
            <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm">
              {{ satellite.id }}
            </td>
            <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm">
              {{ satellite.name }}
            </td>
            <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm">
              {{ satellite.noradId || '-' }}
            </td>
            <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm max-w-md" >
              <pre
                v-if="satellite.line1 && satellite.line2"
                class="font-mono bg-gray-100 p-2 rounded text-xs whitespace-pre overflow-x-auto"
                :title="`${satellite.line1}\n${satellite.line2}`"
                >{{ satellite.line1 }}
{{ satellite.line2 }}</pre>
              <span v-else class="text-gray-400">-</span>
            </td>
            <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm text-center">
              <router-link
                :to="`/satellites/${satellite.id}`"
                class="text-indigo-600 hover:text-indigo-900 mr-3"
                title="Редактировать"
              >
                 ✏️
              </router-link>
              <button
                @click="confirmDelete(satellite.id)"
                class="text-red-600 hover:text-red-900"
                title="Удалить"
              >
                🗑️
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Сообщение, если спутники не найдены -->
    <div v-else class="text-center py-10 text-gray-500">
      Спутники не найдены.
    </div>

    <!-- Модальное окно подтверждения удаления -->
    <BaseModal :show="showDeleteConfirm" @close="cancelDelete">
      <template #header>Подтверждение удаления</template>
      <template #body>
        <p>Вы уверены, что хотите удалить спутник с ID {{ satelliteToDeleteId }}?</p>
      </template>
      <template #footer>
        <button
           @click="deleteSatelliteConfirmed"
           class="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 mr-2 disabled:opacity-50"
           :disabled="isDeleting"
        >
          {{ isDeleting ? 'Удаление...' : 'Удалить' }}
        </button>
        <button
           @click="cancelDelete"
           class="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300"
           :disabled="isDeleting"
         >
          Отмена
        </button>
      </template>
    </BaseModal>

  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { storeToRefs } from 'pinia';
import { useSatelliteStore } from '@/stores/satelliteStore';
import { useRouter } from 'vue-router';
import BaseModal from '@/components/BaseModal.vue'; // Импортируем модальное окно

// Простая реализация debounce
function debounce<T extends (...args: any[]) => any>(func: T, wait: number): (...args: Parameters<T>) => void {
  let timeout: number | undefined;
  return function executedFunction(...args: Parameters<T>) {
    const later = () => {
      clearTimeout(timeout);
      func(...args);
    };
    clearTimeout(timeout);
    timeout = window.setTimeout(later, wait);
  };
}

const satelliteStore = useSatelliteStore();
const router = useRouter();
const { satellites, isLoading, error } = storeToRefs(satelliteStore); // Делаем state реактивным

const searchName = ref('');
const isSearching = ref(false); // Флаг для индикатора загрузки при поиске

// Загружаем спутники при монтировании
onMounted(() => {
  satelliteStore.fetchSatellites();
});

// Функция поиска с debounce
const performSearch = async () => {
  isSearching.value = true;
  await satelliteStore.fetchSatellites({ name: searchName.value || undefined });
  isSearching.value = false;
}

const debouncedSearch = debounce(performSearch, 500); // Задержка 500ms

// --- Логика удаления --- 
const showDeleteConfirm = ref(false);
const satelliteToDeleteId = ref<number | null>(null);
const isDeleting = ref(false); // Флаг для индикации процесса удаления

function confirmDelete(id: number) {
  console.log("Confirm delete for satellite ID:", id);
  satelliteToDeleteId.value = id;
  showDeleteConfirm.value = true; // Открываем модальное окно

  // Убираем прямой вызов удаления отсюда
  // satelliteStore.deleteSatellite(id)
  //   .then(() => {
  //     console.log(`Спутник ${id} успешно удален (предположительно)`);
  //   })
  //   .catch(err => {
  //     console.error("Ошибка при удалении спутника:", err);
  //     alert(`Ошибка при удалении спутника ${id}: ${err.message || err}`);
  //   });
}


function cancelDelete() {
  showDeleteConfirm.value = false;
  satelliteToDeleteId.value = null;
}

async function deleteSatelliteConfirmed() {
  if (satelliteToDeleteId.value !== null) {
    isDeleting.value = true;
    try {
      await satelliteStore.deleteSatellite(satelliteToDeleteId.value);
      // Уведомление об успехе (можно заменить на более красивое)
      console.log(`Спутник ${satelliteToDeleteId.value} успешно удален.`);
      // Список должен обновиться автоматически через store
    } catch (err: any) {
       // Уведомление об ошибке (можно заменить)
       console.error("Failed to delete satellite:", err);
       alert(`Ошибка при удалении спутника: ${err.message || 'Неизвестная ошибка'}`);
    } finally {
      isDeleting.value = false;
      cancelDelete(); // Закрываем модальное окно в любом случае
    }
  }
}

</script>

<style scoped>
/* Стили специфичные для компонента, если нужны */
</style>
