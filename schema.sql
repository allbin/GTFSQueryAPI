create table stops
(
    stop_id       varchar(255)  not null primary key,
    stop_name     varchar(255)  not null,
    stop_lat      numeric(8, 6) not null,
    stop_lon      numeric(9, 6) not null,
    location_type smallint constraint stops_location_type_check check ((location_type >= 0) AND (location_type <= 1))
    -- geom          geometry(Point, 4326)
);
