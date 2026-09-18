<script setup lang="ts">
import { NConfigProvider, NMessageProvider, darkTheme, zhCN, dateZhCN } from 'naive-ui'
import { darkMode } from './theme'
// Design decisions DS-4A (dark default) + DS-5A (teal accent, Naive defaults).
</script>

<template>
  <NConfigProvider
    :theme="darkMode.isDark ? darkTheme : null"
    :locale="zhCN"
    :date-locale="dateZhCN"
    :theme-overrides="{
      common: {
        primaryColor: '#18a058',
        primaryColorHover: '#36ad6a',
        primaryColorPressed: '#0c7a43',
      },
    }"
  >
    <NMessageProvider>
      <router-view />
    </NMessageProvider>
  </NConfigProvider>
</template>
