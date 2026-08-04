<script setup lang="ts">
/**
 * Extracted admin API key settings card shell from SettingsView
 * (Sprint 4 incremental split). Parent keeps state and actions.
 */
</script>

<template>
  <div class="card">
    <div class="border-b border-gray-100 px-6 py-4 dark:border-dark-700">
      <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
        <slot name="title" />
      </h2>
      <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
        <slot name="description" />
      </p>
    </div>
    <div class="space-y-4 p-6">
      <slot />
    </div>
  </div>
</template>
