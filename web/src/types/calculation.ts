// Параметры для расчета координат спутника
export interface CalculateParams { // Имя типа оставляем для соответствия store
  satelliteId: number; // Имя поля satelliteId согласно API
  timestamp?: number; // Unix timestamp (секунды), опционально
}

// Результат расчета координат
export interface CoordinatesResult {
  lat: number;  // Имя поля lat согласно API
  lon: number;  // Имя поля lon согласно API
  alt: number;  // Имя поля alt согласно API
  gmapslink: string; // Ссылка на Google Maps
}

// Параметры для расчета углов (азимут, элевация)
export interface LookAnglesParams {
  satelliteId: number; // Имя поля satelliteId согласно API
  observerPositionId: number; // Имя поля observerPositionId согласно API
  timestamp?: number; // Unix timestamp (секунды), опционально
}

// Результат расчета углов
export interface LookAnglesResult {
  azimuth: number;
  elevation: number;
  range: number; // Имя поля range согласно API
}

// Параметры для расчета интервалов видимости
export interface TimeRangesParams {
  satelliteId: number; // Имя поля satelliteId согласно API
  timestamp?: number; // Имя поля timestamp согласно API
  lon: number;       // Координаты наблюдателя (lon, lat, alt согласно API)
  lat: number;
  alt: number;
  countOfTimeRanges?: number; // Имя поля countOfTimeRanges согласно API
}

// Представление одного интервала видимости
export interface TimeRange { // Имя типа оставляем
  // start: string; // Время начала в формате RFC3339 - Старое поле
  // end: string;   // Время конца в формате RFC3339 - Старое поле
  from: string; // Новое поле из API
  to: string;   // Новое поле из API
  difference?: string; // Необязательное поле из API
}

// Добавляем интерфейс для ответа расчета координат
export interface SatellitePositionResponse {
  lat: number;
  lon: number;
  alt: number;
  gmapsLink: string;
} 