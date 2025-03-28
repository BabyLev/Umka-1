--- схема таблицы для локаций

create table locations (
    id bigserial primary key, --- первичный ключ, идентификаторы локаций
    loc_name text not null, --- имя локации
    lat numeric(9,6) not null,
    lon numeric(9,6) not null,
    alt numeric(9,6) not null
)