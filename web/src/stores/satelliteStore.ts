import { defineStore } from 'pinia';
import api from '@/api'; // Теперь импортируем настроенный axios
import type {
    Satellite,
    SatelliteCreate, // Используется для входных данных формы, трансформируем в ApiSatelliteInfo перед отправкой
    // SatelliteUpdate, // Заменяем на SatelliteUpdatePayload для запроса
    SatelliteUpdatePayload,
    SatelliteFetchParams,
    SatelliteListApiResponse,
    ApiSatelliteInfo
} from '@/types';

// Интерфейс состояния остается прежним, хранит массив Satellite
interface SatelliteState {
  satellites: Satellite[];
  loading: boolean;
  error: string | null;
  currentSatellite: Satellite | null; // Добавляем поле для текущего/выбранного спутника
  // Убираем satelliteInfo, т.к. fetchSatelliteInfoByNorad убран
  // satelliteInfo: ApiSatelliteInfo | null;
}

export const useSatelliteStore = defineStore('satellite', {
  state: (): SatelliteState => ({
    satellites: [],
    loading: false,
    error: null,
    currentSatellite: null, // Инициализируем null
    // satelliteInfo: null,
  }),
  actions: {
    // Получение списка спутников
    async fetchSatellites(params: SatelliteFetchParams = {}) {
      this.loading = true;
      this.error = null;
      try {
        // Ожидаем ответ в формате SatelliteListApiResponse
        const response = await api.post<SatelliteListApiResponse>('/satellite/', params);
        // Трансформируем карту { id: info } в массив Satellite[]
        this.satellites = Object.entries(response.data.satellites || {}).map(([id, info]) => ({
            id: parseInt(id, 10),
            ...info,
        }));
      } catch (err: any) {
        this.error = err.response?.data?.message || err.message || 'Failed to fetch satellites';
        this.satellites = [];
        console.error("Error fetching satellites:", err);
      } finally {
        this.loading = false;
      }
    },

    // Убираем fetchSatelliteInfoByNorad, т.к. нет такого эндпоинта в API
    // и логика получения данных с r4uab относится скорее к моменту добавления/обновления
    // async fetchSatelliteInfoByNorad(noradId: number) { ... }

    // Добавление спутника
    async addSatellite(satelliteData: SatelliteCreate) {
        this.loading = true;
        this.error = null;
        try {
            // satelliteData соответствует ApiSatelliteInfo, отправляем его
            // Ожидаем ответ { satelliteId: number }
            const response = await api.put<{ satelliteId: number }>('/satellite/', satelliteData);
            const newSatelliteId = response.data.satelliteId;
            // TODO: Либо рефетчим спутники (проще), либо конструируем объект Satellite
            // Для конструирования нужны данные TLE/имя, которые могли прийти с r4uab на бэке
            // Просто добавим заглушку с ID, потом обновим через fetchSatellites
            // this.satellites.push({ id: newSatelliteId, ...satelliteData }); // Некорректно, т.к. satelliteData может быть неполным
            await this.fetchSatellites(); // Перезапрашиваем список
            return newSatelliteId; // Возвращаем ID добавленного спутника
        } catch (err: any) {
            this.error = err.response?.data?.message || err.message || 'Failed to add satellite';
            console.error("Error adding satellite:", err);
            return null;
        } finally {
            this.loading = false;
        }
    },

    // Обновление спутника
    async updateSatellite(payload: SatelliteUpdatePayload) {
        this.loading = true;
        this.error = null;
        // Лог перед вызовом API
        console.log(
          '[SatelliteStore] Attempting updateSatellite with payload:', 
          JSON.stringify(payload, null, 2) // Используем JSON.stringify для читаемого вывода объекта
        );
        try {
            // Отправляем данные в формате { satelliteId, satellite: {...} }
            const response = await api.patch<void>('/satellite/', payload);
            // Лог после успешного вызова
            console.log('[SatelliteStore] api.patch call successful.', response); 

            // Рефетч списка после успешного обновления
            await this.fetchSatellites(); 
            return true; // Успех
        } catch (err: any) {
             // Лог при ошибке
            console.error("[SatelliteStore] Error during api.patch:", err.response || err.message || err);
            this.error = err.response?.data?.message || err.message || 'Failed to update satellite';
            return false; // Неудача
        } finally {
            this.loading = false;
        }
    },

    // Получение информации о конкретном спутнике по ID
    async fetchSatelliteById(id: number) {
        this.loading = true;
        this.error = null;
        this.currentSatellite = null; // Сбрасываем перед загрузкой
        try {
            // API GET /satellite/{id} возвращает ApiSatelliteInfo
            const response = await api.get<ApiSatelliteInfo>(`/satellite/${id}`);
            // Сохраняем в currentSatellite, добавляя ID
            this.currentSatellite = { id, ...response.data };
            return this.currentSatellite;
        } catch (err: any) {
            this.error = err.response?.data?.message || err.message || `Failed to fetch satellite with id ${id}`;
            console.error(`Error fetching satellite with id ${id}:`, err);
            return null;
        } finally {
            this.loading = false;
        }
    },

    // Удаление спутника
    async deleteSatellite(id: number) {
        this.loading = true;
        this.error = null;
        try {
            await api.delete<string>(`/satellite/${id}`);
            this.satellites = this.satellites.filter(sat => sat.id !== id);
            // Если удалили текущий выбранный спутник, сбрасываем его
            if (this.currentSatellite?.id === id) {
                this.currentSatellite = null;
            }
        } catch (err: any) {
            this.error = err.response?.data || err.message || 'Failed to delete satellite';
            console.error("Error deleting satellite:", err);
        } finally {
            this.loading = false;
        }
    },
  },
  getters: {
      getSatelliteById: (state) => (id: number): Satellite | undefined => {
          return state.satellites.find((sat: Satellite) => sat.id === id);
      },
      // Геттер для текущего загруженного спутника
      getCurrentSatellite: (state): Satellite | null => {
          return state.currentSatellite;
      },
      // Геттер для получения списка спутников в формате для выбора (например, { label: name, value: id })
      satellitesForSelect: (state) => {
          return state.satellites.map(sat => ({ label: `${sat.name} (ID: ${sat.id})`, value: sat.id }));
      }
  }
}); 