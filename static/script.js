document.addEventListener('DOMContentLoaded', () => {
    const tabs = document.querySelectorAll('[role="tab"]');
    const tabContents = document.querySelectorAll('[role="tabpanel"]');

    // Устанавливаем активной первую вкладку по умолчанию
    if (tabs.length > 0) {
        activateTab(tabs[0], tabContents);
    }

    tabs.forEach(tab => {
        tab.addEventListener('click', () => {
            activateTab(tab, tabContents);
        });
    });

    // Инициализация форм
    setupCalculateForms();
    setupSatelliteForms();
    setupLocationForms();

    // Загрузка начальных данных
    loadInitialData();
});

let allSatellites = {}; // Кэш для спутников
let allLocations = {}; // Кэш для локаций

async function loadInitialData() {
    try {
        // Параллельная загрузка спутников и локаций
        await Promise.all([
            loadSatellites(),
            loadLocations()
        ]);
        // После загрузки обновляем выпадающие списки
        populateSatelliteDropdowns();
        populateLocationDropdowns();
    } catch (error) {
        console.error("Ошибка при начальной загрузке данных:", error);
        alert("Не удалось загрузить начальные данные. Проверьте консоль.");
    }
}

function populateSatelliteDropdowns(satellites = allSatellites) {
    const selects = document.querySelectorAll('select[id$="-sat-id"]'); // Выбираем все select для спутников
    selects.forEach(select => {
        const currentVal = select.value; // Сохраняем текущее значение, если есть
        select.innerHTML = '<option value="">-- Выберите спутник --</option>'; // Очищаем и добавляем плейсхолдер
        for (const id in satellites) {
            const sat = satellites[id];
            const option = document.createElement('option');
            option.value = id;
            option.textContent = `${escapeHTML(sat.name || 'Без имени')} (ID: ${id}, NORAD: ${sat.noradId || 'N/A'})`;
            select.appendChild(option);
        }
        select.value = currentVal; // Восстанавливаем значение, если оно было
    });
}

function populateLocationDropdowns(locations = allLocations) {
    const selects = document.querySelectorAll('select[id$="-loc-id"]'); // Выбираем все select для локаций
    selects.forEach(select => {
        const currentVal = select.value;
        select.innerHTML = '<option value="">-- Выберите локацию --</option>';
        for (const id in locations) {
            const loc = locations[id];
            const option = document.createElement('option');
            option.value = id;
            option.textContent = `${escapeHTML(loc.name || 'Без имени')} (ID: ${id})`;
            select.appendChild(option);
        }
        select.value = currentVal;
    });
}

function activateTab(selectedTab, tabContents) {
    const targetId = selectedTab.getAttribute('data-tabs-target');
    const targetContent = document.querySelector(targetId);

    // Сброс стилей для всех вкладок и скрытие контента
    document.querySelectorAll('[role="tab"]').forEach(t => {
        t.classList.remove('border-blue-500', 'text-blue-600');
        t.classList.add('hover:text-gray-600', 'hover:border-gray-300');
        t.setAttribute('aria-selected', 'false');
    });
    tabContents.forEach(content => {
        content.classList.add('hidden');
    });

    // Активация выбранной вкладки и ее контента
    selectedTab.classList.add('border-blue-500', 'text-blue-600');
    selectedTab.classList.remove('hover:text-gray-600', 'hover:border-gray-300');
    selectedTab.setAttribute('aria-selected', 'true');
    if (targetContent) {
        targetContent.classList.remove('hidden');
    }

    // Возможно, перезагружать данные при переключении вкладок, если они могли измениться
    // switch (targetId) {
    //     case '#satellites':
    //         loadSatellites();
    //         break;
    //     case '#locations':
    //         loadLocations();
    //         break;
    // }
}

// --- Функции для взаимодействия с API --- //

const API_BASE_URL = ''; // Установите базовый URL вашего API, если он отличается от корня сайта

async function fetchData(url, options = {}) {
    try {
        const response = await fetch(API_BASE_URL + url, {
            headers: {
                'Content-Type': 'application/json',
                ...options.headers,
            },
            ...options,
        });
        if (!response.ok) {
            const errorText = await response.text();
            // Попытка распарсить ошибку как JSON
            let errorJson = {};
            try {
                 errorJson = JSON.parse(errorText);
            } catch(e) {}

            // Отображаем более понятное сообщение, если возможно
            const displayError = errorJson.error || errorText || 'Неизвестная ошибка';
            console.error(`HTTP error! status: ${response.status}, message: ${errorText}`);
            throw new Error(`Ошибка API (${response.status}): ${displayError}`);
        }
        // Если ожидается JSON ответ
        if (response.headers.get('Content-Type')?.includes('application/json') && response.status !== 204) { // 204 No Content
           return await response.json();
        }
        // Если ответ пустой или текстовый
        return await response.text();
    } catch (error) {
        console.error('Ошибка при запросе к API:', error);
        // Отображение ошибки пользователю (можно улучшить)
        displayError(error.message);
        throw error; // Повторно выбрасываем ошибку для обработки выше
    }
}

// --- Функции для вкладки "Расчеты" --- //

function setupCalculateForms() {
    const posForm = document.getElementById('calculate-pos-form');
    const anglesForm = document.getElementById('calculate-angles-form');
    const trForm = document.getElementById('calculate-tr-form');

    posForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        clearResults('calculate-results');
        const formData = new FormData(posForm);
        const satelliteId = parseInt(formData.get('satelliteId'));
        const timestamp = formData.get('timestamp') ? parseInt(formData.get('timestamp')) : undefined;
        try {
            await calculateSatellitePosition(satelliteId, timestamp);
        } catch (err) {/* Ошибка уже отображена в fetchData */}
    });

    anglesForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        clearResults('calculate-results');
        const formData = new FormData(anglesForm);
        const satelliteId = parseInt(formData.get('satelliteId'));
        const observerPositionId = parseInt(formData.get('observerPositionId'));
        const timestamp = formData.get('timestamp') ? parseInt(formData.get('timestamp')) : undefined;
         try {
            await calculateLookAngles(satelliteId, observerPositionId, timestamp);
         } catch (err) {/* Ошибка уже отображена в fetchData */}
    });

    trForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        clearResults('calculate-results');
        const formData = new FormData(trForm);
        const satelliteId = parseInt(formData.get('satelliteId'));
        const lat = parseFloat(formData.get('lat'));
        const lon = parseFloat(formData.get('lon'));
        const alt = parseFloat(formData.get('alt'));
        const count = formData.get('countOfTimeRanges') ? parseInt(formData.get('countOfTimeRanges')) : undefined;
        const timestamp = formData.get('timestamp') ? parseInt(formData.get('timestamp')) : undefined;
         try {
            await calculateVisibleTimeRanges(satelliteId, lat, lon, alt, count, timestamp);
         } catch (err) {/* Ошибка уже отображена в fetchData */}
    });
}

async function calculateSatellitePosition(satelliteId, timestamp) {
    const body = { satelliteID: satelliteId }; // Обратите внимание на регистр ID
    if (timestamp) {
        body.timestamp = timestamp;
    }
    const result = await fetchData('/calculate/', { method: 'POST', body: JSON.stringify(body) });
    displayResults('calculate-results', result);
}

async function calculateLookAngles(satelliteId, observerPositionId, timestamp) {
     const body = { satelliteID: satelliteId, observerPositionID: observerPositionId }; // Обратите внимание на регистр ID
    if (timestamp) {
        body.timestamp = timestamp;
    }
    const result = await fetchData('/look-angles/', { method: 'POST', body: JSON.stringify(body) });
    displayResults('calculate-results', result);
}

async function calculateVisibleTimeRanges(satelliteId, lat, lon, alt, countOfTimeRanges, timestamp) {
     const body = { satelliteID: satelliteId, lat, lon, alt }; // Обратите внимание на регистр ID
    if (timestamp) {
        body.timestamp = timestamp;
    }
    if (countOfTimeRanges) {
        body.countOfTimeRanges = countOfTimeRanges;
    }
    const result = await fetchData('/time-ranges/', { method: 'POST', body: JSON.stringify(body) });
    displayResults('calculate-results', result);
}

// --- Функции для вкладки "Спутники" --- //
function setupSatelliteForms() {
    const addEditForm = document.getElementById('satellite-form');
    const findForm = document.getElementById('find-satellite-form');

    addEditForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const formData = new FormData(addEditForm);
        const editId = formData.get('satelliteId') ? parseInt(formData.get('satelliteId')) : null;

        const satelliteData = {
            name: formData.get('name'),
            noradId: formData.get('noradId') ? parseInt(formData.get('noradId')) : null,
            line1: formData.get('line1') || null,
            line2: formData.get('line2') || null,
        };

        // Удаляем null значения, чтобы API не получал их
        Object.keys(satelliteData).forEach(key => {
             if (satelliteData[key] === null || satelliteData[key] === '') delete satelliteData[key];
        });

        try {
            if (editId) {
                 // Формируем тело запроса для PATCH
                const patchData = { satelliteID: editId, satellite: satelliteData }; // Обратите внимание на регистр ID
                await updateSatellite(editId, patchData.satellite); // Передаем только данные спутника
                displaySuccess(`Спутник ID ${editId} успешно обновлен.`);
            } else {
                const result = await addSatellite(satelliteData);
                displaySuccess(`Спутник успешно добавлен с ID: ${result.satelliteId}`);
            }
            resetSatelliteForm(); // Сброс формы после успеха
        } catch (err) {/* Ошибка уже отображена */}
    });

     findForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const formData = new FormData(findForm);
        const filter = {
            ids: formData.get('ids') ? formData.get('ids').split(',').map(id => parseInt(id.trim())).filter(id => !isNaN(id)) : undefined,
            noradIds: formData.get('noradIds') ? formData.get('noradIds').split(',').map(id => parseInt(id.trim())).filter(id => !isNaN(id)) : undefined,
            name: formData.get('name') || undefined,
        };
        // Удаляем undefined значения
        Object.keys(filter).forEach(key => filter[key] === undefined && delete filter[key]);
         try {
             await findSatellites(filter);
         } catch (err) {/* Ошибка уже отображена */}
    });

    // Скрытие/отображение полей TLE в зависимости от NORAD ID
    const noradIdInput = document.getElementById('satellite-norad-id');
    const tleInputsDiv = document.getElementById('satellite-tle-inputs');
    noradIdInput.addEventListener('input', () => {
        if (noradIdInput.value) {
            tleInputsDiv.classList.add('hidden');
            document.getElementById('satellite-line1').value = '';
            document.getElementById('satellite-line2').value = '';
        } else {
            tleInputsDiv.classList.remove('hidden');
        }
    });
}

async function loadSatellites(filter = {}) {
     try {
        // Очищаем фильтр от пустых строк/массивов перед отправкой
        const cleanFilter = { ...filter };
        Object.keys(cleanFilter).forEach(key => {
            if (cleanFilter[key] === undefined || cleanFilter[key] === '' || (Array.isArray(cleanFilter[key]) && cleanFilter[key].length === 0)) {
                 delete cleanFilter[key];
             }
         });

        const result = await fetchData('/satellite/', { method: 'POST', body: JSON.stringify(cleanFilter) });
        allSatellites = result.satellites || {}; // Обновляем кэш
        populateTable('satellites-table-body', allSatellites, createSatelliteRow);
        populateSatelliteDropdowns(); // Обновляем выпадающие списки после загрузки/фильтрации
         return allSatellites;
     } catch (err) {
         populateTable('satellites-table-body', {}, createSatelliteRow); // Очистка таблицы при ошибке
         populateSatelliteDropdowns({});
         return {};
     }
}

async function addSatellite(satelliteData) {
    const result = await fetchData('/satellite/', { method: 'PUT', body: JSON.stringify(satelliteData) });
    await loadSatellites(); // Перезагрузить список
    return result; // Возвращаем ответ (должен содержать ID)
}

async function updateSatellite(satelliteId, updateData) {
    // API ожидает { "satelliteID": id, "satellite": { ... } }
    await fetchData('/satellite/', { method: 'PATCH', body: JSON.stringify({ satelliteID: satelliteId, satellite: updateData }) });
    await loadSatellites(); // Перезагрузить список
}

async function deleteSatellite(satelliteId) {
    await fetchData(`/satellite/${satelliteId}`, { method: 'DELETE' });
    displaySuccess(`Спутник ID ${satelliteId} удален.`);
    await loadSatellites(); // Перезагрузить список
}

async function findSatellites(filter) {
    // Используем loadSatellites, т.к. эндпоинт один
    await loadSatellites(filter);
}

async function getSatellite(satelliteId) {
    try {
        const result = await fetchData(`/satellite/${satelliteId}`);
        return result;
    } catch (err) {
        return null;
    }
}

async function editSatellite(id) {
    console.log(`Редактирование спутника с ID: ${id}`);
    const satData = await getSatellite(id);
    if (!satData) {
        displayError(`Не удалось загрузить данные для спутника ID ${id}`);
        return;
    }

    resetSatelliteForm(); // Сброс на всякий случай

    // Заполняем форму
    document.getElementById('satellite-edit-id').value = id;
    document.getElementById('satellite-name').value = satData.name || '';
    document.getElementById('satellite-norad-id').value = satData.noradId || '';
    document.getElementById('satellite-line1').value = satData.line1 || '';
    document.getElementById('satellite-line2').value = satData.line2 || '';

    // Обновляем заголовок и кнопку
    document.getElementById('satellite-form-title').textContent = `Редактировать спутник ID: ${id}`;
    document.getElementById('satellite-form-submit').textContent = 'Сохранить изменения';
    document.getElementById('satellite-form-submit').classList.remove('bg-green-600', 'hover:bg-green-700');
    document.getElementById('satellite-form-submit').classList.add('bg-yellow-500', 'hover:bg-yellow-600');
    document.getElementById('satellite-form-cancel').classList.remove('hidden');

    // Скрываем/показываем TLE
     const tleInputsDiv = document.getElementById('satellite-tle-inputs');
     if (satData.noradId) {
         tleInputsDiv.classList.add('hidden');
     } else {
         tleInputsDiv.classList.remove('hidden');
     }

    // Прокрутка к форме
    document.getElementById('satellite-form').scrollIntoView({ behavior: 'smooth' });
}

function resetSatelliteForm() {
    const form = document.getElementById('satellite-form');
    form.reset();
    document.getElementById('satellite-edit-id').value = '';
    document.getElementById('satellite-form-title').textContent = 'Добавить новый спутник';
    const submitButton = document.getElementById('satellite-form-submit');
    submitButton.textContent = 'Добавить спутник';
    submitButton.classList.remove('bg-yellow-500', 'hover:bg-yellow-600');
    submitButton.classList.add('bg-green-600', 'hover:bg-green-700');
    document.getElementById('satellite-form-cancel').classList.add('hidden');
    document.getElementById('satellite-tle-inputs').classList.remove('hidden');
    clearMessages(); // Очистка сообщений об успехе/ошибках
}

// --- Функции для вкладки "Локации" --- //
function setupLocationForms() {
    const addEditForm = document.getElementById('location-form');
    const findForm = document.getElementById('find-location-form');

    addEditForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const formData = new FormData(addEditForm);
        const editId = formData.get('locationId') ? parseInt(formData.get('locationId')) : null;

        const locationData = {
            name: formData.get('name'),
            location: {
                lat: parseFloat(formData.get('lat')),
                lon: parseFloat(formData.get('lon')),
                alt: parseFloat(formData.get('alt')),
            }
        };

        try {
            clearMessages(); // Очистим сообщения перед запросом
            if (editId) {
                 const patchData = { locationID: editId, location: locationData }; 
                 await updateLocation(editId, patchData.location); 
                 displaySuccess(`Локация ID ${editId} успешно обновлена.`);
            } else {
                // Передаем locationData напрямую, без обертки observerLocation
                const result = await addLocation(locationData); 
                displaySuccess(`Локация успешно добавлена с ID: ${result.id}`); // Используем result.id
            }
            resetLocationForm();
        } catch (err) {/* Ошибка уже отображена */}
    });

     findForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const formData = new FormData(findForm);
        const filter = {
            name: formData.get('name') || undefined,
        };
        Object.keys(filter).forEach(key => filter[key] === undefined && delete filter[key]);
         try {
             await findLocations(filter);
         } catch (err) {/* Ошибка уже отображена */}
    });
}


async function loadLocations(filter = {}) {
     try {
        const cleanFilter = { ...filter };
        Object.keys(cleanFilter).forEach(key => {
            if (cleanFilter[key] === undefined || cleanFilter[key] === '') {
                delete cleanFilter[key];
            }
        });
        const result = await fetchData('/location/', { method: 'POST', body: JSON.stringify(cleanFilter) });
        allLocations = result.locations || {}; // Обновляем кэш
        populateTable('locations-table-body', allLocations, createLocationRow);
        populateLocationDropdowns(); // Обновляем выпадающие списки
        return allLocations;
     } catch (err) {
         populateTable('locations-table-body', {}, createLocationRow);
         populateLocationDropdowns({});
         return {};
     }
}

async function addLocation(locationData) {
    // Отправляем данные как есть, без дополнительной обертки
    const result = await fetchData('/location/', { method: 'PUT', body: JSON.stringify(locationData) });
    await loadLocations(); // Перезагрузить список
    return result;
}

async function updateLocation(locationId, updateData) {
     // API ожидает { "locationID": id, "location": { "name": ..., "location": { ... } } }
    await fetchData('/location/', { method: 'PATCH', body: JSON.stringify({ locationID: locationId, location: updateData }) });
    await loadLocations(); // Перезагрузить список
}

async function deleteLocation(locationId) {
    await fetchData(`/location/${locationId}`, { method: 'DELETE' });
    displaySuccess(`Локация ID ${locationId} удалена.`);
    await loadLocations(); // Перезагрузить список
}

async function findLocations(filter) {
    await loadLocations(filter);
}

async function getLocation(locationId) {
    // В API нет прямого GET /location/{id}, будем искать по ID
    try {
        const result = await findLocations({ ids: [locationId] }); // Используем поиск
        return result && result[locationId] ? result[locationId] : null; // Возвращаем найденный объект или null
    } catch (err) {
        return null;
    }
}

async function editLocation(id) {
    console.log(`Редактирование локации с ID: ${id}`);
    // Используем данные из кэша, т.к. нет прямого GET
    const locData = allLocations[id];

    if (!locData) {
        displayError(`Не удалось найти данные для локации ID ${id}. Попробуйте обновить список.`);
        // Можно попробовать вызвать findLocations({ ids: [id] }) здесь, но это усложнит поток
        return;
    }

    resetLocationForm();

    // Заполняем форму
    document.getElementById('location-edit-id').value = id;
    document.getElementById('location-name').value = locData.name || '';
    if (locData.location) {
        document.getElementById('location-lat').value = locData.location.lat ?? '';
        document.getElementById('location-lon').value = locData.location.lon ?? '';
        document.getElementById('location-alt').value = locData.location.alt ?? '';
    }

    // Обновляем заголовок и кнопку
    document.getElementById('location-form-title').textContent = `Редактировать локацию ID: ${id}`;
    const submitButton = document.getElementById('location-form-submit');
    submitButton.textContent = 'Сохранить изменения';
    submitButton.classList.remove('bg-green-600', 'hover:bg-green-700');
    submitButton.classList.add('bg-yellow-500', 'hover:bg-yellow-600');
    document.getElementById('location-form-cancel').classList.remove('hidden');

    // Прокрутка к форме
    document.getElementById('location-form').scrollIntoView({ behavior: 'smooth' });
}

function resetLocationForm() {
    const form = document.getElementById('location-form');
    form.reset();
    document.getElementById('location-edit-id').value = '';
    document.getElementById('location-form-title').textContent = 'Добавить новую локацию';
    const submitButton = document.getElementById('location-form-submit');
    submitButton.textContent = 'Добавить локацию';
    submitButton.classList.remove('bg-yellow-500', 'hover:bg-yellow-600');
    submitButton.classList.add('bg-green-600', 'hover:bg-green-700');
    document.getElementById('location-form-cancel').classList.add('hidden');
    clearMessages();
}

// --- Вспомогательные функции --- //

// Отображение сообщений об успехе или ошибке
let messageTimeout;
function showMessage(text, isError = false) {
    const container = document.body; // Или более специфичный контейнер
    let messageDiv = document.getElementById('api-message');
    if (!messageDiv) {
        messageDiv = document.createElement('div');
        messageDiv.id = 'api-message';
        messageDiv.className = 'fixed bottom-4 right-4 p-4 rounded-md shadow-lg text-white z-50 max-w-sm';
        container.appendChild(messageDiv);
    }

    messageDiv.textContent = text;
    messageDiv.className = `fixed bottom-4 right-4 p-4 rounded-md shadow-lg text-white z-50 max-w-sm ${isError ? 'bg-red-600' : 'bg-green-600'}`;
    messageDiv.classList.remove('hidden');

    // Скрыть сообщение через 5 секунд
    clearTimeout(messageTimeout);
    messageTimeout = setTimeout(() => {
        messageDiv.classList.add('hidden');
    }, 5000);
}

function displaySuccess(message) {
    showMessage(message, false);
}

function displayError(message) {
    showMessage(message, true);
}

function clearMessages() {
     const messageDiv = document.getElementById('api-message');
     if (messageDiv) {
         messageDiv.classList.add('hidden');
         clearTimeout(messageTimeout);
     }
     // Также очищаем область результатов расчетов
     clearResults('calculate-results');
}

function clearResults(elementId) {
    const resultsElement = document.getElementById(elementId);
    if (resultsElement) {
        const preElement = resultsElement.querySelector('pre');
        if (preElement) {
            preElement.textContent = ''; // Очищаем результаты
        }
    }
}

function displayResults(elementId, data) {
    const resultsElement = document.getElementById(elementId);
    if (resultsElement) {
        const preElement = resultsElement.querySelector('pre');
        if (preElement) {
            preElement.textContent = JSON.stringify(data, null, 2); // Красивый вывод JSON
        }
    }
}

function populateTable(tbodyId, data, createRowFunction) {
    const tbody = document.getElementById(tbodyId);
    if (!tbody) return;

    tbody.innerHTML = ''; // Очистить таблицу перед заполнением

    if (!data || Object.keys(data).length === 0) {
        tbody.innerHTML = '<tr><td colspan="100%" class="text-center p-4 text-gray-500">Нет данных для отображения</td></tr>';
        return;
    }

    // Сортируем по ID для консистентности
    const sortedIds = Object.keys(data).sort((a, b) => parseInt(a) - parseInt(b));

    for (const id of sortedIds) {
        const item = data[id];
        const row = createRowFunction(id, item);
        if (row) {
            tbody.appendChild(row);
        }
    }
}

// --- Функции создания строк таблиц и подтверждения удаления (без изменений) --- //

function createSatelliteRow(id, satellite) {
    const tr = document.createElement('tr');
    // Добавляем TLE в title для удобства
    const tleTitle = `Line 1: ${satellite.line1 || ''}\nLine 2: ${satellite.line2 || ''}`;
    tr.innerHTML = `
        <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">${id}</td>
        <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500" title="${escapeHTML(tleTitle)}">${escapeHTML(satellite.name || '')}</td>
        <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">${satellite.noradId || 'N/A'}</td>
        <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
            <button onclick="editSatellite(${id})" class="text-indigo-600 hover:text-indigo-900 mr-2">Редактировать</button>
            <button onclick="confirmDelete('satellite', ${id})" class="text-red-600 hover:text-red-900">Удалить</button>
        </td>
    `;
    return tr;
}

function createLocationRow(id, location) {
    const tr = document.createElement('tr');
    const coords = location.location ? `Ш: ${location.location.lat?.toFixed(6)}, Д: ${location.location.lon?.toFixed(6)}, В: ${location.location.alt?.toFixed(3)}` : 'N/A';
    tr.innerHTML = `
        <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">${id}</td>
        <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">${escapeHTML(location.name || '')}</td>
        <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">${coords}</td>
        <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
            <button onclick="editLocation(${id})" class="text-indigo-600 hover:text-indigo-900 mr-2">Редактировать</button>
            <button onclick="confirmDelete('location', ${id})" class="text-red-600 hover:text-red-900">Удалить</button>
        </td>
    `;
    return tr;
}

function confirmDelete(type, id) {
    const typeName = type === 'satellite' ? 'спутник' : 'локацию';
    if (confirm(`Вы уверены, что хотите удалить ${typeName} с ID ${id}?`)) {
        clearMessages(); // Очищаем старые сообщения
        if (type === 'satellite') {
            deleteSatellite(id).catch(err => {}); // Ловим ошибку, т.к. она уже отображена
        } else if (type === 'location') {
            deleteLocation(id).catch(err => {}); // Ловим ошибку, т.к. она уже отображена
        }
    }
}

// Функция для экранирования HTML (базовая безопасность)
function escapeHTML(str) {
     if (typeof str !== 'string') return '';
    return str.replace(/[&<>'"/]/g, function (s) {
        const entityMap = {
            '&': '&amp;',
            '<': '&lt;',
            '>': '&gt;',
            '"': '&quot;',
            "'": '&#39;', // Используем &#39; для одинарной кавычки
            '/': '&#x2F;'
        };
        return entityMap[s];
    });
}

// TODO: Реализовать формы для ввода данных во вкладках
// TODO: Добавить обработчики событий для кнопок в формах
// TODO: Улучшить обработку ошибок и обратную связь с пользователем
// TODO: Вызывать loadSatellites() и loadLocations() при загрузке страницы или активации вкладок 