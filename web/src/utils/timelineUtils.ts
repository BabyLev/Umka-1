/**
 * Интерфейс для стилей полосы таймлайна.
 */
export interface TimelineBarStyle {
  left: string;
  width: string;
}

/**
 * Вычисляет смещение и ширину для полосы интервала на 24-часовой шкале.
 * @param from - Время начала интервала (ISO строка или объект Date).
 * @param to - Время конца интервала (ISO строка или объект Date).
 * @param referenceDate - Дата, относительно которой строятся 24ч.
 * @returns Объект со свойствами left и width в процентах или стили по умолчанию при ошибке.
 */
export function calculateBarStyle(from: string | Date, to: string | Date, referenceDate: Date): { left: string; width: string } {
  try {
    const startDate = new Date(from);
    const endDate = new Date(to);
    if (isNaN(startDate.getTime()) || isNaN(endDate.getTime()) || isNaN(referenceDate.getTime())) {
        console.warn("Invalid date provided to calculateBarStyle", { from, to, referenceDate });
      return { left: '0%', width: '0%' };
    }

    // Определяем начало и конец суток для referenceDate
    const startOfDay = new Date(referenceDate);
    startOfDay.setHours(0, 0, 0, 0);
    const startOfDayMs = startOfDay.getTime();
    const endOfDayMs = startOfDayMs + 24 * 60 * 60 * 1000; // Конец суток (начало следующих)

    const totalDayMilliseconds = endOfDayMs - startOfDayMs;
    if (totalDayMilliseconds <= 0) {
         console.error("Failed to calculate total milliseconds in a day.");
         return { left: '0%', width: '0%' };
    }

    // Определяем фактические начало и конец интервала в рамках референсных суток
    const effectiveStartMs = Math.max(startDate.getTime(), startOfDayMs);
    const effectiveEndMs = Math.min(endDate.getTime(), endOfDayMs);

    // Если интервал не пересекается с референсными сутками
    if (effectiveEndMs <= effectiveStartMs) {
        return { left: '0%', width: '0%' };
    }

    const startOffsetMilliseconds = effectiveStartMs - startOfDayMs;
    const endOffsetMilliseconds = effectiveEndMs - startOfDayMs;

    // Вычисляем позицию и ширину в процентах
    const leftPercent = (startOffsetMilliseconds / totalDayMilliseconds) * 100;
    const widthPercent = ((endOffsetMilliseconds - startOffsetMilliseconds) / totalDayMilliseconds) * 100;

    // Ограничиваем значения, чтобы избежать выхода за пределы 0-100%
    const finalLeft = Math.max(0, leftPercent);
    // Минимальная ширина 0.5% для видимости очень коротких интервалов
    const finalWidth = Math.max(0.5, Math.min(widthPercent, 100 - finalLeft));

    return {
      left: `${finalLeft.toFixed(2)}%`,
      width: `${finalWidth.toFixed(2)}%`,
    };
  } catch (error) {
    console.error("Error calculating bar style:", from, to, error);
    return { left: '0%', width: '0%' }; // Возвращаем стили по умолчанию при любой ошибке
  }
}

/**
 * Рассчитывает и форматирует длительность интервала.
 * @param from - Время начала интервала (ISO строка или объект Date).
 * @param to - Время конца интервала (ISO строка или объект Date).
 * @returns Строка с длительностью (напр., "10м 32с") или пустая строка при ошибке.
 */
export function calculateDuration(from: string | Date, to: string | Date): string {
    try {
        const startMs = new Date(from).getTime();
        const endMs = new Date(to).getTime();
        if (isNaN(startMs) || isNaN(endMs)) {
            console.warn("Invalid date provided to calculateDuration", { from, to });
            return '';
        }

        // Рассчитываем разницу в секундах
        const diffSeconds = Math.round((endMs - startMs) / 1000);
        if (diffSeconds < 0) {
             console.warn("End date is before start date in calculateDuration", { from, to });
             return '0м 0с'; // Или вернуть пустую строку/ошибку?
        }

        const minutes = Math.floor(diffSeconds / 60);
        const seconds = diffSeconds % 60;
        return `${minutes}м ${seconds}с`;
    } catch (error) {
        console.error("Error calculating duration:", from, to, error);
        return ''; // Возвращаем пустую строку при любой ошибке
    }
} 