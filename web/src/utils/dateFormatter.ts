export function formatDateTime(dateString: string | Date): string {
  if (!dateString) return '-';

  try {
    const date = new Date(dateString);
    // Проверка на валидность даты
    if (isNaN(date.getTime())) {
        console.warn("Invalid date string passed to formatDateTime:", dateString);
        return 'Invalid Date';
    }

    // Форматируем в локальный формат даты и времени
    // Например: 15.07.2024, 14:35:10
    // Вы можете настроить опции для другого формата
    const options: Intl.DateTimeFormatOptions = {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false // 24-часовой формат
    };

    return new Intl.DateTimeFormat(navigator.language || 'ru-RU', options).format(date);

  } catch (error) {
    console.error("Error formatting date:", dateString, error);
    return 'Error'; // Возвращаем строку ошибки, если что-то пошло не так
  }
}

// Хелпер для форматирования только времени
export function formatTime(dateString: string | Date): string {
  if (!dateString) return '';
  try {
    const date = new Date(dateString);
    if (isNaN(date.getTime())) {
        console.warn("Invalid date string passed to formatTime:", dateString);
        return 'Invalid Time';
    }
    const options: Intl.DateTimeFormatOptions = {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false,
    };
    return new Intl.DateTimeFormat(navigator.language || 'ru-RU', options).format(date);
  } catch (error) {
    console.error("Error formatting time:", dateString, error);
    return 'Error';
  }
} 