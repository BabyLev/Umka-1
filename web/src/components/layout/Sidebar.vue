<template>
  <aside 
    :class="[
      'transition-all duration-300 bg-gray-800 text-gray-100 flex-shrink-0 h-screen',
      isCollapsed ? 'w-16' : 'w-64'
    ]"
  >
    <!-- Заголовок и кнопка сворачивания -->
    <div class="p-4 flex items-center justify-between border-b border-gray-700">
      <h2 
        :class="[
          'text-white font-semibold transition-all duration-300',
          isCollapsed ? 'text-sm' : 'text-xl'
        ]"
      >
        {{ isCollapsed ? 'UMKA' : 'UMKA Server' }}
      </h2>
      <button 
        @click="toggleCollapse" 
        class="text-gray-400 hover:text-white focus:outline-none" 
        title="Свернуть/развернуть меню"
      >
        <svg 
          xmlns="http://www.w3.org/2000/svg" 
          :class="[
            'h-5 w-5 transition-transform',
            isCollapsed ? 'transform rotate-180' : ''
          ]" 
          viewBox="0 0 20 20" 
          fill="currentColor"
        >
          <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
        </svg>
      </button>
    </div>

    <!-- Основное навигационное меню -->
    <nav class="mt-4">
      <ul class="space-y-1">
        <!-- Спутники -->
        <li>
          <router-link 
            to="/satellites" 
            :class="[
              'flex items-center px-4 py-3 transition-colors duration-200',
              'hover:bg-gray-700 focus:outline-none',
              isActiveRoute('/satellites') ? 'bg-gray-700 border-l-4 border-blue-500' : 'border-l-4 border-transparent'
            ]"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M5 2a1 1 0 011 1v1h1a1 1 0 010 2H6v1a1 1 0 01-2 0V6H3a1 1 0 010-2h1V3a1 1 0 011-1zm0 10a1 1 0 011 1v1h1a1 1 0 110 2H6v1a1 1 0 11-2 0v-1H3a1 1 0 110-2h1v-1a1 1 0 011-1zm7-10a1 1 0 01.707.293l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414-1.414L14.586 8l-3.293-3.293A1 1 0 0112 4zm-7 4a1 1 0 011 1v1h1a1 1 0 110 2H6v1a1 1 0 11-2 0v-1H3a1 1 0 110-2h1V7a1 1 0 011-1z" clip-rule="evenodd" />
            </svg>
            <span :class="['ml-3', isCollapsed ? 'sr-only' : '']">Спутники</span>
          </router-link>
        </li>

        <!-- Локации -->
        <li>
          <router-link 
            to="/locations" 
            :class="[
              'flex items-center px-4 py-3 transition-colors duration-200',
              'hover:bg-gray-700 focus:outline-none',
              isActiveRoute('/locations') ? 'bg-gray-700 border-l-4 border-blue-500' : 'border-l-4 border-transparent'
            ]"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M5.05 4.05a7 7 0 119.9 9.9L10 18.9l-4.95-4.95a7 7 0 010-9.9zM10 11a2 2 0 100-4 2 2 0 000 4z" clip-rule="evenodd" />
            </svg>
            <span :class="['ml-3', isCollapsed ? 'sr-only' : '']">Локации</span>
          </router-link>
        </li>

        <!-- Разделитель для секции расчетов -->
        <li class="px-4 py-2">
          <div class="border-t border-gray-700"></div>
        </li>

        <!-- Расчеты: Координаты -->
        <li>
          <router-link 
            to="/calculations/coordinates" 
            :class="[
              'flex items-center px-4 py-3 transition-colors duration-200',
              'hover:bg-gray-700 focus:outline-none',
              isActiveRoute('/calculations/coordinates') ? 'bg-gray-700 border-l-4 border-blue-500' : 'border-l-4 border-transparent'
            ]"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zm0 6a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1v-2zm0 6a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1v-2z" clip-rule="evenodd" />
            </svg>
            <span :class="['ml-3', isCollapsed ? 'sr-only' : '']">Координаты</span>
          </router-link>
        </li>

        <!-- Расчеты: Углы -->
        <li>
          <router-link 
            to="/calculations/look-angles" 
            :class="[
              'flex items-center px-4 py-3 transition-colors duration-200',
              'hover:bg-gray-700 focus:outline-none',
              isActiveRoute('/calculations/look-angles') ? 'bg-gray-700 border-l-4 border-blue-500' : 'border-l-4 border-transparent'
            ]"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path d="M10 12a2 2 0 100-4 2 2 0 000 4z" />
              <path fill-rule="evenodd" d="M.458 10C1.732 5.943 5.522 3 10 3s8.268 2.943 9.542 7c-1.274 4.057-5.064 7-9.542 7S1.732 14.057.458 10zM14 10a4 4 0 11-8 0 4 4 0 018 0z" clip-rule="evenodd" />
            </svg>
            <span :class="['ml-3', isCollapsed ? 'sr-only' : '']">Углы</span>
          </router-link>
        </li>

        <!-- Расчеты: Видимость -->
        <li>
          <router-link 
            to="/calculations/visibility" 
            :class="[
              'flex items-center px-4 py-3 transition-colors duration-200',
              'hover:bg-gray-700 focus:outline-none',
              isActiveRoute('/calculations/visibility') ? 'bg-gray-700 border-l-4 border-blue-500' : 'border-l-4 border-transparent'
            ]"
          >
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M3 5a2 2 0 012-2h10a2 2 0 012 2v10a2 2 0 01-2 2H5a2 2 0 01-2-2V5zm11 1H6v8l4-2 4 2V6z" clip-rule="evenodd" />
            </svg>
            <span :class="['ml-3', isCollapsed ? 'sr-only' : '']">Видимость</span>
          </router-link>
        </li>
      </ul>
    </nav>
  </aside>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useRoute } from 'vue-router';

// Состояние сворачивания меню
const isCollapsed = ref(false);

// Функция переключения состояния
const toggleCollapse = () => {
  isCollapsed.value = !isCollapsed.value;
};

// Получение текущего маршрута для выделения активных пунктов
const route = useRoute();

// Функция для проверки активного маршрута
const isActiveRoute = (path: string) => {
  return route.path.startsWith(path);
};
</script>

<style scoped>
/* Добавим плавные переходы для всех элементов меню */
.router-link-active {
  font-weight: 500;
}

/* Дополнительные стили для анимации при наведении */
@keyframes pulse {
  0% { transform: scale(1); }
  50% { transform: scale(1.05); }
  100% { transform: scale(1); }
}

.hover\:animate-pulse:hover {
  animation: pulse 1s infinite;
}
</style> 