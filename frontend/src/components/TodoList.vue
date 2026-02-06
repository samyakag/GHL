<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { todoClient } from '../client'
import type { Todo } from '../gen/todo/v1/todo_pb'

const todos = ref<Todo[]>([])
const newTodoTitle = ref('')
const loading = ref(false)
const error = ref('')

async function fetchTodos() {
  try {
    loading.value = true
    error.value = ''
    const response = await todoClient.listTodos({})
    todos.value = response.todos
  } catch (e) {
    error.value = `Failed to fetch todos: ${e}`
  } finally {
    loading.value = false
  }
}

async function createTodo() {
  if (!newTodoTitle.value.trim()) return
  try {
    loading.value = true
    error.value = ''
    await todoClient.createTodo({ title: newTodoTitle.value })
    newTodoTitle.value = ''
    await fetchTodos()
  } catch (e) {
    error.value = `Failed to create todo: ${e}`
  } finally {
    loading.value = false
  }
}

async function toggleTodo(todo: Todo) {
  try {
    error.value = ''
    await todoClient.updateTodo({
      id: todo.id,
      title: todo.title,
      completed: !todo.completed,
    })
    await fetchTodos()
  } catch (e) {
    error.value = `Failed to update todo: ${e}`
  }
}

async function deleteTodo(id: string) {
  try {
    error.value = ''
    await todoClient.deleteTodo({ id })
    await fetchTodos()
  } catch (e) {
    error.value = `Failed to delete todo: ${e}`
  }
}

onMounted(() => {
  fetchTodos()
})
</script>

<template>
  <div v-if="error" class="error">{{ error }}</div>

  <form @submit.prevent="createTodo" class="add-form">
    <input v-model="newTodoTitle" type="text" placeholder="What needs to be done?" :disabled="loading" />
    <button type="submit" :disabled="loading || !newTodoTitle.trim()">Add</button>
  </form>

  <div v-if="loading && todos.length === 0" class="loading">Loading...</div>

  <ul v-else class="todo-list">
    <li v-for="todo in todos" :key="todo.id" :class="{ completed: todo.completed }">
      <input type="checkbox" :checked="todo.completed" @change="toggleTodo(todo)" />
      <span class="title">{{ todo.title }}</span>
      <button class="delete-btn" @click="deleteTodo(todo.id)">Delete</button>
    </li>
  </ul>

  <p v-if="!loading && todos.length === 0" class="empty">No todos yet. Add one above!</p>
</template>
