<template>
  <div>
    <div class="mb-6">
      <router-link
        to="/sources"
        class="text-sm text-gray-500 hover:text-gray-700 inline-flex items-center mb-4"
      >
        <ArrowLeftIcon class="h-4 w-4 mr-1" />
        Back to Sources
      </router-link>
      <h2 class="text-2xl font-bold text-gray-900">
        {{ isEdit ? 'Edit Source' : 'New Source' }}
      </h2>
    </div>

    <div v-if="loading" class="text-center py-12">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
    </div>

    <form v-else @submit.prevent="handleSubmit" class="bg-white shadow-sm rounded-lg border border-gray-200 p-6">
      <div class="space-y-6">
        <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
          <div>
            <label for="name" class="block text-sm font-medium text-gray-700">
              Name <span class="text-red-500">*</span>
            </label>
            <input
              id="name"
              v-model="form.name"
              type="text"
              required
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
            />
          </div>

          <div>
            <label for="url" class="block text-sm font-medium text-gray-700">
              URL <span class="text-red-500">*</span>
            </label>
            <input
              id="url"
              v-model="form.url"
              type="url"
              required
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
            />
          </div>

          <div>
            <label for="article_index" class="block text-sm font-medium text-gray-700">
              Article Index <span class="text-red-500">*</span>
            </label>
            <input
              id="article_index"
              v-model="form.article_index"
              type="text"
              required
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
            />
          </div>

          <div>
            <label for="page_index" class="block text-sm font-medium text-gray-700">
              Page Index <span class="text-red-500">*</span>
            </label>
            <input
              id="page_index"
              v-model="form.page_index"
              type="text"
              required
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
            />
          </div>

          <div>
            <label for="rate_limit" class="block text-sm font-medium text-gray-700">
              Rate Limit
            </label>
            <input
              id="rate_limit"
              v-model="form.rate_limit"
              type="text"
              placeholder="1s"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
            />
          </div>

          <div>
            <label for="max_depth" class="block text-sm font-medium text-gray-700">
              Max Depth
            </label>
            <input
              id="max_depth"
              v-model.number="form.max_depth"
              type="number"
              min="1"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
            />
          </div>

          <div>
            <label for="city_name" class="block text-sm font-medium text-gray-700">
              City Name
            </label>
            <input
              id="city_name"
              v-model="form.city_name"
              type="text"
              placeholder="sudbury_com"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
            />
          </div>

          <div>
            <label for="group_id" class="block text-sm font-medium text-gray-700">
              Group ID (Drupal UUID)
            </label>
            <input
              id="group_id"
              v-model="form.group_id"
              type="text"
              placeholder="550e8400-e29b-41d4-a716-446655440000"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
            />
          </div>
        </div>

        <div>
          <label class="flex items-center">
            <input
              v-model="form.enabled"
              type="checkbox"
              class="rounded border-gray-300 text-blue-600 shadow-sm focus:border-blue-500 focus:ring-blue-500"
            />
            <span class="ml-2 text-sm text-gray-700">Enabled</span>
          </label>
        </div>

        <div v-if="error" class="rounded-md bg-red-50 p-4">
          <div class="flex">
            <ExclamationCircleIcon class="h-5 w-5 text-red-400" />
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">Error</h3>
              <div class="mt-2 text-sm text-red-700">{{ error }}</div>
            </div>
          </div>
        </div>
      </div>

      <div class="mt-6 flex justify-end space-x-3">
        <router-link
          to="/sources"
          class="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        >
          Cancel
        </router-link>
        <button
          type="submit"
          :disabled="submitting"
          class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <span v-if="submitting">Saving...</span>
          <span v-else>{{ isEdit ? 'Update' : 'Create' }} Source</span>
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { sourcesApi } from '../api/client'
import { ArrowLeftIcon, ExclamationCircleIcon } from '@heroicons/vue/24/outline'

const router = useRouter()
const route = useRoute()

const isEdit = computed(() => !!route.params.id)

const form = ref({
  name: '',
  url: '',
  article_index: '',
  page_index: '',
  rate_limit: '1s',
  max_depth: 2,
  time: [],
  selectors: {
    article: {},
    list: {},
    page: {},
  },
  city_name: null,
  group_id: null,
  enabled: true,
})

const loading = ref(false)
const submitting = ref(false)
const error = ref(null)

const loadSource = async () => {
  if (!isEdit.value) return
  
  loading.value = true
  error.value = null
  try {
    const source = await sourcesApi.get(route.params.id)
    form.value = {
      ...source,
      city_name: source.city_name || null,
      group_id: source.group_id || null,
    }
  } catch (err) {
    error.value = err.response?.data?.error || err.message || 'Failed to load source'
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  submitting.value = true
  error.value = null
  
  try {
    const data = {
      ...form.value,
      city_name: form.value.city_name || null,
      group_id: form.value.group_id || null,
    }
    
    if (isEdit.value) {
      await sourcesApi.update(route.params.id, data)
    } else {
      await sourcesApi.create(data)
    }
    
    router.push('/sources')
  } catch (err) {
    error.value = err.response?.data?.error || err.response?.data?.details || err.message || 'Failed to save source'
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  loadSource()
})
</script>

