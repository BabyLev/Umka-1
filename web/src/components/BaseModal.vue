<template>
  <Transition name="modal">
    <div
      v-if="show"
      class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
      @click.self="closeModal"
    >
      <div class="bg-white rounded-lg shadow-xl max-w-md w-full mx-4 p-6">
        <!-- Header -->
        <div class="flex justify-between items-center border-b pb-3 mb-4">
          <h3 class="text-lg font-medium">
            <slot name="header">Заголовок по умолчанию</slot>
          </h3>
          <button
            @click="closeModal"
            class="text-gray-400 hover:text-gray-600 text-2xl"
            aria-label="Закрыть"
          >
            &times;
          </button>
        </div>

        <!-- Body -->
        <div class="mb-4">
          <slot name="body">Тело модального окна по умолчанию.</slot>
        </div>

        <!-- Footer -->
        <div class="flex justify-end space-x-3 border-t pt-3">
          <slot name="footer">
            <button
              @click="closeModal"
              class="px-4 py-2 bg-gray-200 text-gray-800 rounded hover:bg-gray-300"
            >
              Закрыть
            </button>
          </slot>
        </div>
      </div>
    </div>
  </Transition>
</template>

<script setup lang="ts">
interface Props {
  show: boolean;
}

const props = defineProps<Props>();
const emit = defineEmits(['close']);

function closeModal() {
  emit('close');
}

// Добавляем немного стилей для анимации
</script>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .bg-white,
.modal-leave-active .bg-white {
  transition: transform 0.3s ease;
}

.modal-enter-from .bg-white,
.modal-leave-to .bg-white {
  transform: scale(0.95);
}
</style> 