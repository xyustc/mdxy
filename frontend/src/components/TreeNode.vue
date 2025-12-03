<template>
  <div class="tree-node">
    <div 
      class="node-item" 
      :class="{ 'is-directory': isDirectory, 'is-active': isActive }"
      @click="handleClick"
    >
      <span class="node-icon">
        <template v-if="isDirectory">
          {{ expanded ? '📂' : '📁' }}
        </template>
        <template v-else>📄</template>
      </span>
      <span class="node-name">{{ node.name }}</span>
    </div>
    
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
import { useRouter, useRoute } from 'vue-router'

const props = defineProps({
  node: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['note-click'])

const router = useRouter()
const route = useRoute()

const expanded = ref(false)
const isDirectory = computed(() => props.node.type === 'directory')
const isActive = computed(() => {
  if (isDirectory.value) return false
  return route.params.path === props.node.path
})

const handleClick = () => {
  if (isDirectory.value) {
    expanded.value = !expanded.value
  } else {
    router.push({ name: 'note', params: { path: props.node.path } })
    emit('note-click')
  }
}
</script>
