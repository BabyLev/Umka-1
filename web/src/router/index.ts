import { createRouter, createWebHistory } from 'vue-router'

// Импортируем страницы (пока заглушки)
// Можно использовать динамический импорт для lazy loading:
// const SatelliteList = () => import('@/pages/satellites/SatelliteList.vue');

import SatelliteList from '@/pages/satellites/SatelliteList.vue'
import SatelliteForm from '@/pages/satellites/SatelliteForm.vue'
import LocationList from '@/pages/locations/LocationList.vue'
import LocationForm from '@/pages/locations/LocationForm.vue'
import CalculationCoordinates from '@/pages/calculations/CalculationCoordinates.vue'
import CalculationLookAngles from '@/pages/calculations/CalculationLookAngles.vue'
import CalculationVisibility from '@/pages/calculations/CalculationVisibility.vue'

const routes = [
  {
    path: '/',
    redirect: '/satellites', // Редирект с главной на список спутников по умолчанию
  },
  {
    path: '/satellites',
    name: 'SatelliteList',
    component: SatelliteList,
  },
  {
    path: '/satellites/new',
    name: 'SatelliteNew',
    component: SatelliteForm, // Используем ту же форму для создания
    // Можно передать props или meta для различения new/edit
    meta: { mode: 'new' },
  },
  {
    path: '/satellites/:id',
    name: 'SatelliteEdit',
    component: SatelliteForm, // Используем ту же форму для редактирования
    props: true, // Передаем :id как props в компонент
    meta: { mode: 'edit' },
  },
  {
    path: '/locations',
    name: 'LocationList',
    component: LocationList,
  },
  {
    path: '/locations/new',
    name: 'LocationNew',
    component: LocationForm,
    meta: { mode: 'new' },
  },
  {
    path: '/locations/:id',
    name: 'LocationEdit',
    component: LocationForm,
    props: true,
    meta: { mode: 'edit' },
  },
  {
    path: '/calculations/coordinates',
    name: 'CalculationCoordinates',
    component: CalculationCoordinates,
  },
  {
    path: '/calculations/look-angles',
    name: 'CalculationLookAngles',
    component: CalculationLookAngles,
  },
  {
    path: '/calculations/visibility',
    name: 'CalculationVisibility',
    component: CalculationVisibility,
  },
  // Добавить маршрут 404, если нужно
  // { path: '/:pathMatch(.*)*', name: 'NotFound', component: NotFoundComponent },
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
})

export default router 