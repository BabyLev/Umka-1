import { defineStore } from 'pinia';
import api from '@/api'; // Используем axios
import type {
    Location, // Тип для фронтенда
    LocationCreate, // Тип для формы, соответствует ApiObserverLocation
    // LocationUpdate, // Заменяем на LocationUpdatePayload
    LocationUpdatePayload,
    LocationSearchParams,
    LocationListApiResponse,
    ApiObserverLocation
} from '@/types';

// Интерфейс состояния хранит массив Location
interface LocationState {
  locations: Location[];
  loading: boolean;
  error: string | null;
}

export const useLocationStore = defineStore('location', {
  state: (): LocationState => ({
    locations: [],
    loading: false,
    error: null,
  }),
  actions: {
    // Получение списка локаций
    async fetchLocations(params: LocationSearchParams = {}) {
      this.loading = true;
      this.error = null;
      try {
        // API POST /location/ принимает { name: string }, возвращает LocationListApiResponse
        const response = await api.post<LocationListApiResponse>('/location/', params);
        // Трансформируем карту { id: observerLocation } в массив Location[]
        this.locations = Object.entries(response.data.locations || {}).map(([id, obsLocation]) => ({
            id: parseInt(id, 10),
            name: obsLocation.name,
            // Раскладываем вложенный объект location
            lon: obsLocation.location.lon,
            lat: obsLocation.location.lat,
            alt: obsLocation.location.alt,
        }));
      } catch (err: any) {
        this.error = err.response?.data?.message || err.message || 'Failed to fetch locations';
        this.locations = [];
        console.error("Error fetching locations:", err);
      } finally {
        this.loading = false;
      }
    },

    // Добавление локации
    async addLocation(locationData: LocationCreate) {
      this.loading = true;
      this.error = null;
      try {
        // locationData соответствует ApiObserverLocation, отправляем его (PUT /location/)
        // Ожидаем ответ { observerLocationId: number }
        const response = await api.put<{ observerLocationId: number }>('/location/', locationData);
        const newLocationId = response.data.observerLocationId;
        // Рефетчим список
        await this.fetchLocations();
        return newLocationId;
      } catch (err: any) {
        this.error = err.response?.data?.message || err.message || 'Failed to add location';
        console.error("Error adding location:", err);
        return null;
      } finally {
        this.loading = false;
      }
    },

    // Обновление локации
    async updateLocation(payload: LocationUpdatePayload) {
      this.loading = true;
      this.error = null;
      try {
        // Отправляем { locationId, location: {...} } (PATCH /location/)
        await api.patch<void>('/location/', payload);
        // Рефетчим список
        await this.fetchLocations();
        return true;
      } catch (err: any) {
        this.error = err.response?.data?.message || err.message || 'Failed to update location';
         console.error("Error updating location:", err);
        return false;
      } finally {
        this.loading = false;
      }
    },

    // Удаление локации
    async deleteLocation(id: number) {
      this.loading = true;
      this.error = null;
      try {
        // API DELETE /location/{id} возвращает текст
        await api.delete<string>(`/location/${id}`);
        // Удаляем из локального состояния
        this.locations = this.locations.filter((loc: Location) => loc.id !== id);
      } catch (err: any) {
        this.error = err.response?.data || err.message || 'Failed to delete location';
        console.error("Error deleting location:", err);
      } finally {
        this.loading = false;
      }
    },

    // --- Новое действие для получения одной локации ---
    async fetchLocationById(id: number): Promise<Location | null> {
      this.loading = true;
      this.error = null;
      try {
        // GET /location/{id} возвращает ApiObserverLocation
        const response = await api.get<ApiObserverLocation>(`/location/${id}`);
        const obsLocation = response.data;

        // Трансформируем ApiObserverLocation в Location (плоскую структуру для фронта)
        const location: Location = {
            id: id, // ID берем из параметра запроса, т.к. API его не возвращает в теле
            name: obsLocation.name,
            lon: obsLocation.location.lon,
            lat: obsLocation.location.lat,
            alt: obsLocation.location.alt,
        };

        // Опционально: можно обновить эту локацию в общем списке this.locations,
        // если она там есть, или добавить, если нет. Пока просто вернем ее.
        // const index = this.locations.findIndex(l => l.id === id);
        // if (index !== -1) {
        //   this.locations[index] = location;
        // } else {
        //   this.locations.push(location);
        // }

        return location;
      } catch (err: any) {
        this.error = err.response?.data?.message || err.message || `Failed to fetch location with ID ${id}`;
        console.error(`Error fetching location ${id}:`, err);
        return null;
      } finally {
        this.loading = false;
      }
    }
    // --- Конец нового действия ---

  },
  getters: {
    getLocationById: (state) => (id: number): Location | undefined => {
      return state.locations.find((loc: Location) => loc.id === id);
    },
    // Геттер для списка локаций в формате для выбора
    locationsForSelect: (state) => {
        return state.locations.map(loc => ({ label: `${loc.name} (ID: ${loc.id})`, value: loc.id }));
    }
  }
}); 