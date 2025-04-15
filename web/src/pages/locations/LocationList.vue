<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-2xl font-semibold mb-4">Список локаций</h1>

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
           <svg class="animate-spin h-5 w-5 text-gray-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
             <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
             <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
           </svg>
        </span>
      </div>

      <!-- Кнопка добавления -->
      <router-link
        to="/locations/new" 
        class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded transition duration-150 ease-in-out"
      >
        Добавить локацию
      </router-link>
    </div>

    <!-- Индикатор загрузки -->
    <div v-if="isLoading && !isSearching" class="text-center py-10">
      <p>Загрузка локаций...</p>
    </div>

    <!-- Сообщение об ошибке -->
    <div v-else-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative" role="alert">
      <strong class="font-bold">Ошибка!</strong>
      <span class="block sm:inline"> {{ error }}</span>
    </div>

    <!-- Таблица с локациями -->
    <div v-else-if="locations.length > 0" class="bg-white shadow-md rounded my-6 overflow-x-auto">
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
              Координаты (lon, lat, alt)
            </th>
            <th class="px-5 py-3 border-b-2 border-gray-200 bg-gray-100 text-center text-xs font-semibold text-gray-600 uppercase tracking-wider">
              Действия
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="location in locations" :key="location.id">
            <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm">
              {{ location.id }}
            </td>
            <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm">
              {{ location.name }}
            </td>
            <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm">
              {{ location.lon.toFixed(4) }}, {{ location.lat.toFixed(4) }}, {{ location.alt.toFixed(2) }} км
            </td>
            <td class="px-5 py-5 border-b border-gray-200 bg-white text-sm text-center">
              <router-link
                :to="`/locations/${location.id}`" 
                class="text-indigo-600 hover:text-indigo-900 mr-3"
                title="Редактировать"
              >
                 ✏️
              </router-link>
              <button
                @click="confirmDelete(location.id)"
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

    <!-- Сообщение, если локации не найдены -->
    <div v-else class="text-center py-10 text-gray-500">
      Локации не найдены.
    </div>

    <!-- Модальное окно подтверждения удаления -->
    <BaseModal :show="showDeleteConfirm" @close="cancelDelete">
      <template #header>Подтверждение удаления</template>
      <template #body>
        <p>Вы уверены, что хотите удалить локацию с ID {{ itemToDeleteId }}?</p>
      </template>
      <template #footer>
        <button
           @click="deleteItemConfirmed"
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
import { useLocationStore } from '@/stores/locationStore'; 
import { useRouter } from 'vue-router';
import BaseModal from '@/components/BaseModal.vue'; 

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

const locationStore = useLocationStore();
const router = useRouter();
const { locations, loading: isLoading, error } = storeToRefs(locationStore); 

const searchName = ref('');
const isSearching = ref(false); 

// Загружаем локации при монтировании
onMounted(() => {
  locationStore.fetchLocations();
});

// Функция поиска с debounce
const performSearch = async () => {
  isSearching.value = true;
  await locationStore.fetchLocations({ name: searchName.value || undefined });
  isSearching.value = false;
}

const debouncedSearch = debounce(performSearch, 500); 

// --- Логика удаления --- 
const showDeleteConfirm = ref(false);
const itemToDeleteId = ref<number | null>(null);
const isDeleting = ref(false); 

function confirmDelete(id: number) {
  console.log("Confirm delete for location ID:", id);
  itemToDeleteId.value = id;
  showDeleteConfirm.value = true; 
}


function cancelDelete() {
  showDeleteConfirm.value = false;
  itemToDeleteId.value = null;
}

async function deleteItemConfirmed() {
  if (itemToDeleteId.value !== null) {
    isDeleting.value = true;
    try {
      await locationStore.deleteLocation(itemToDeleteId.value);
      console.log(`Локация ${itemToDeleteId.value} успешно удалена.`);
    } catch (err: any) {
       console.error("Failed to delete location:", err);
       alert(`Ошибка при удалении локации: ${err.message || 'Неизвестная ошибка'}`);
    } finally {
      isDeleting.value = false;
      cancelDelete(); 
    }
  }
}

</script>

<style scoped>
/* Стили, если нужны */
</style>
