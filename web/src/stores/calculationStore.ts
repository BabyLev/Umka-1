import { defineStore } from 'pinia';
import api from '@/api'; // Используем axios
import type {
  CoordinatesResult,
  LookAnglesResult,
  TimeRange,
  CalculateParams,
  LookAnglesParams,
  TimeRangesParams,
} from '@/types'; // Типы уже обновлены

interface CalculationState {
  coordinatesResult: CoordinatesResult | null;
  lookAnglesResult: LookAnglesResult | null;
  timeRangesResult: TimeRange[] | null;
  isLoading: boolean;
  error: string | null;
}

export const useCalculationStore = defineStore('calculation', {
  state: (): CalculationState => ({
    coordinatesResult: null,
    lookAnglesResult: null,
    timeRangesResult: null,
    isLoading: false,
    error: null,
  }),
  actions: {
    // Очистка конкретного результата или всех результатов
    clearResult(type: 'coordinates' | 'lookAngles' | 'timeRanges' | 'all' = 'all') {
      if (type === 'coordinates' || type === 'all') {
        this.coordinatesResult = null;
      }
      if (type === 'lookAngles' || type === 'all') {
        this.lookAnglesResult = null;
      }
      if (type === 'timeRanges' || type === 'all') {
        this.timeRangesResult = null;
      }
       if (type === 'all') {
         this.error = null; // Сбрасываем ошибку при полной очистке
       }
    },

    // Расчет координат
    async calculateCoordinates(params: CalculateParams) {
      this.isLoading = true;
      this.error = null;
      this.coordinatesResult = null; // Очищаем предыдущий результат
      try {
        // Используем обновленные типы параметров и результата
        const response = await api.post<CoordinatesResult>('/calculate/', params);
        this.coordinatesResult = response.data;
      } catch (err: any) {
        this.error = err.response?.data?.message || err.message || 'Failed to calculate coordinates';
        console.error("Error calculating coordinates:", err);
      } finally {
        this.isLoading = false;
      }
    },

    // Расчет углов (азимут, элевация)
    async calculateLookAngles(params: LookAnglesParams) {
      this.isLoading = true;
      this.error = null;
      this.lookAnglesResult = null; // Очищаем предыдущий результат
      try {
        // Используем обновленные типы параметров и результата
        const response = await api.post('/look-angles/', params);
        this.lookAnglesResult = response.data;
      } catch (err: any) {
        this.error = err.response?.data?.message || err.message || 'Failed to calculate look angles';
        console.error("Error calculating look angles:", err);
      } finally {
        this.isLoading = false;
      }
    },

    // Расчет видимости
    async calculateTimeRanges(params: TimeRangesParams): Promise<TimeRange[]> {
      this.isLoading = true;
      this.error = null;
      this.timeRangesResult = null; // Очищаем предыдущий результат
      try {
        // Используем обновленные типы параметров и результата (TimeRange[])
        const response = await api.post<TimeRange[]>('/time-ranges/', params);
        this.timeRangesResult = response.data;
        return response.data;
      } catch (err: any) {
        const errorMessage = err.response?.data?.message || err.message || 'Failed to calculate time ranges';
        this.error = errorMessage;
        console.error("Error calculating time ranges:", err);
        throw new Error(errorMessage);
      } finally {
        this.isLoading = false;
      }
    },
  },
}); 