<template>
  <div class="tree-node">
    <button
      v-if="isDirectory"
      type="button"
      class="node-item"
      :class="{ 'is-directory': true }"
      :aria-expanded="expanded.toString()"
      @click="expanded = !expanded"
    >
      <span class="node-icon">{{ expanded ? '📂' : '📁' }}</span>
      <span class="node-name">{{ node.name }}</span>
    </button>
    <router-link
      v-else
      :to="{ name: 'note', params: { path: node.path } }"
      class="node-item"
      :class="{ 'is-active': isActive }"
      @click="emit('note-click')"
    >
      <span class="node-icon">📄</span>
      <span class="node-name">{{ node.name }}</span>
    </router-link>
    
    <div v-if="isDirectory && expanded" class="node-children">
      <TreeNode 
        v-for="child in node.children" 
        :key="child.path" 
        :node="child"
        @note-click="$emit('note-click')"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'

const props = defineProps({
  node: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['note-click'])

const route = useRoute()

const expanded = ref(false)
const isDirectory = computed(() => props.node.type === 'directory')
const isActive = computed(() => {
  if (isDirectory.value) return false
  const currentPath = Array.isArray(route.params.path) ? route.params.path.join('/') : route.params.path
  return currentPath === props.node.path
})
</script>
