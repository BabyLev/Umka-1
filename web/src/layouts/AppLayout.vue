<template>
  <div class="flex h-screen bg-gray-100">
    <!-- Sidebar -->
    <Sidebar />

    <!-- Content Area -->
    <div class="flex-1 flex flex-col overflow-hidden">
      <!-- Header (Optional) -->
      <!-- <header class="bg-white shadow-md p-4">Header Content</header> -->

      <!-- Main Content -->
      <main class="flex-1 overflow-x-hidden overflow-y-auto bg-gray-100 p-6">
        <router-view /> <!-- Where routed components will be displayed -->
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import Sidebar from '@/components/layout/Sidebar.vue';
</script>

<style scoped>
/* Add any layout-specific styles here if needed */
</style> 