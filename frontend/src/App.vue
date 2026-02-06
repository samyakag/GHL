<script setup lang="ts">
import { ref } from 'vue'
import TodoList from './components/TodoList.vue'
import EventList from './components/EventList.vue'
import './App.css'

const activeTab = ref<'todos' | 'events'>('events')
</script>

<template>
  <div class="container">
    <h1>Event & Todo App</h1>
    <p class="subtitle">ConnectRPC + Vue.js + Go</p>

    <div class="tabs">
      <button class="tab-btn" :class="{ active: activeTab === 'events' }" @click="activeTab = 'events'">Events</button>
      <button class="tab-btn" :class="{ active: activeTab === 'todos' }" @click="activeTab = 'todos'">Todos</button>
    </div>

    <EventList v-if="activeTab === 'events'" />
    <TodoList v-if="activeTab === 'todos'" />
  </div>
</template>
