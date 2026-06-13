import { createRouter, createWebHistory } from 'vue-router'
import TextList from '@/views/TextList.vue'
import TextEditor from '@/views/TextEditor.vue'
import TextReader from '@/views/TextReader.vue'

const routes = [
  {
    path: '/',
    name: 'TextList',
    component: TextList
  },
  {
    path: '/texts/new',
    name: 'TextNew',
    component: TextEditor
  },
  {
    path: '/texts/:id/edit',
    name: 'TextEdit',
    component: TextEditor,
    props: true
  },
  {
    path: '/texts/:id',
    name: 'TextReader',
    component: TextReader,
    props: true
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
