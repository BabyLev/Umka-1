// Тип для координат в API
export interface ApiLocation {
  lon: number; // Долгота
  lat: number; // Широта
  alt: number; // Высота (км)
}

// Тип для локации наблюдателя в API
export interface ApiObserverLocation {
  name: string;
  location: ApiLocation;
}

// Тип для объекта локации, используемый во фронтенде (с нашим ID)
export interface Location {
  id: number;
  name: string;
  lon: number;
  lat: number;
  alt: number;
}

// Тип для создания новой локации (передаем объект, как в API PUT /location/)
export interface LocationCreate extends ApiObserverLocation {}

// Тип для обновления локации (согласно API PATCH /location/)
export interface LocationUpdatePayload {
  locationId: number;
  location: ApiObserverLocation; // API требует полный объект
}

// Параметры для поиска локаций через POST /location/
export interface LocationSearchParams {
    name?: string; // Поиск только по имени согласно API
}

// Тип ответа для запроса списка локаций (POST /location/)
export interface LocationListApiResponse {
    locations: {
        [id: string]: ApiObserverLocation;
    };
}
