--- схема таблицы для спутников

create table satellites (
    id bigserial primary key, --- первичный ключ, идентификаторы спутников
    sat_name text not null, --- имя спутика
    norad_id int,
    line1 text not null,
    line2 text not null
)