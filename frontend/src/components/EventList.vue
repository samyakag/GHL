<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { eventClient } from '../client'
import type { Event } from '../gen/event/v1/event_pb'
import { ConnectError, Code } from '@connectrpc/connect'

const events = ref<Event[]>([])
const newEventTitle = ref('')
const newEventDescription = ref('')
const newEventDate = ref('')
const newEventCapacity = ref(10)

const selectedEventId = ref<string | null>(null)
const regName = ref('')
const regEmail = ref('')

const loading = ref(false)
const error = ref('')

// Retry helper with exponential backoff
async function withRetry<T>(
  fn: () => Promise<T>,
  maxRetries = 3,
  initialDelayMs = 1000
): Promise<T> {
  let lastError: unknown

  for (let attempt = 0; attempt <= maxRetries; attempt++) {
    try {
      return await fn()
    } catch (e) {
      lastError = e

      // Don't retry on business logic errors (validation, not found, capacity full, etc.)
      if (e instanceof ConnectError) {
        const nonRetryableCodes = [
          Code.InvalidArgument,
          Code.NotFound,
          Code.AlreadyExists,
          Code.PermissionDenied,
          Code.FailedPrecondition,
          Code.Unauthenticated,
        ]
        if (nonRetryableCodes.includes(e.code)) {
          throw e
        }
      }

      // If this was the last attempt, throw the error
      if (attempt === maxRetries) {
        throw e
      }

      // Exponential backoff: wait 1s, 2s, 4s...
      const delayMs = initialDelayMs * Math.pow(2, attempt)
      await new Promise(resolve => setTimeout(resolve, delayMs))
    }
  }

  throw lastError
}

async function fetchEvents() {
  try {
    loading.value = true
    error.value = ''
    const response = await withRetry(() => eventClient.listEvents({}))
    events.value = response.events
  } catch (e) {
    error.value = `Failed to fetch events: ${e}`
  } finally {
    loading.value = false
  }
}

async function createEvent() {
  if (!newEventTitle.value.trim() || !newEventDate.value) return
  try {
    loading.value = true
    error.value = ''
    await withRetry(() =>
      eventClient.createEvent({
        title: newEventTitle.value,
        description: newEventDescription.value,
        date: newEventDate.value,
        capacity: newEventCapacity.value,
      })
    )
    newEventTitle.value = ''
    newEventDescription.value = ''
    newEventDate.value = ''
    newEventCapacity.value = 10
    await fetchEvents()
  } catch (e) {
    error.value = `Failed to create event: ${e}`
  } finally {
    loading.value = false
  }
}

function openRegForm(eventId: string) {
  selectedEventId.value = selectedEventId.value === eventId ? null : eventId
  regName.value = ''
  regEmail.value = ''
}

async function registerForEvent(eventId: string) {
  if (!regName.value.trim() || !regEmail.value.trim()) return
  try {
    loading.value = true
    error.value = ''
    await withRetry(() =>
      eventClient.registerForEvent({
        eventId,
        name: regName.value,
        email: regEmail.value,
      })
    )
    regName.value = ''
    regEmail.value = ''
    selectedEventId.value = null
    await fetchEvents()
  } catch (e) {
    error.value = `Failed to register: ${e}`
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchEvents()
})
</script>

<template>
  <div v-if="error" class="error">{{ error }}</div>

  <form @submit.prevent="createEvent" class="event-form">
    <input v-model="newEventTitle" type="text" placeholder="Event title" :disabled="loading" />
    <input v-model="newEventDescription" type="text" placeholder="Description" :disabled="loading" />
    <div class="event-form-row">
      <input v-model="newEventDate" type="date" :disabled="loading" />
      <input v-model.number="newEventCapacity" type="number" min="1" placeholder="Capacity" :disabled="loading" />
      <button type="submit" :disabled="loading || !newEventTitle.trim() || !newEventDate">Create Event</button>
    </div>
  </form>

  <div v-if="loading && events.length === 0" class="loading">Loading...</div>

  <div v-else class="event-list">
    <div v-for="event in events" :key="event.id" class="event-card">
      <div class="event-header">
        <div class="event-info">
          <h3>{{ event.title }}</h3>
          <p class="event-description">{{ event.description }}</p>
          <p class="event-date">{{ event.date }}</p>
        </div>
        <span class="capacity-badge" :class="{ 'sold-out': event.registeredCount >= event.capacity }">
          {{ event.registeredCount >= event.capacity ? 'Sold Out' : `${event.registeredCount}/${event.capacity}` }}
        </span>
      </div>

      <div v-if="event.registeredCount < event.capacity">
        <button class="register-btn" @click="openRegForm(event.id)" :disabled="loading">
          {{ selectedEventId === event.id ? 'Cancel' : 'Register' }}
        </button>
      </div>

      <form v-if="selectedEventId === event.id" @submit.prevent="registerForEvent(event.id)" class="register-form">
        <input v-model="regName" type="text" placeholder="Your name" :disabled="loading" />
        <input v-model="regEmail" type="email" placeholder="Your email" :disabled="loading" />
        <button type="submit" :disabled="loading || !regName.trim() || !regEmail.trim()">Submit</button>
      </form>
    </div>
  </div>

  <p v-if="!loading && events.length === 0" class="empty">No events yet. Create one above!</p>
</template>
