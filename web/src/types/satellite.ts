// Тип для представления объекта спутника в API
export interface ApiSatelliteInfo {
  line1: string;
  line2: string;
  name: string;
  noradId?: number | null; // В API может быть null или отсутствовать? Судя по доке, 0, но ?/null безопаснее
}

// Тип для объекта спутника, используемый во фронтенде (с нашим ID)
export interface Satellite extends ApiSatelliteInfo {
  id: number;
}

// Тип для создания нового спутника (передаем объект, как в API PUT /satellite/)
export interface SatelliteCreate extends ApiSatelliteInfo {}

// Тип для обновления спутника (согласно API PATCH /satellite/)
export interface SatelliteUpdatePayload {
    satelliteId: number;
    satellite: Partial<ApiSatelliteInfo>; // Обновляемые поля
}

// Параметры для поиска/фильтрации спутников через POST /satellite/
export interface SatelliteFetchParams {
  ids?: number[];
  noradIds?: number[]; // Используем noradIds согласно API
  name?: string;
}

// Тип ответа для запроса списка спутников (POST /satellite/)
export interface SatelliteListApiResponse {
    satellites: { // Ключ - ID спутника (строка), значение - информация о нем
        [id: string]: ApiSatelliteInfo;
    };
}
