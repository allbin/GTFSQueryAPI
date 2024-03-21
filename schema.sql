create table public.spatial_ref_sys
(
    srid      integer not null
        primary key
        constraint spatial_ref_sys_srid_check
            check ((srid > 0) AND (srid <= 998999)),
    auth_name varchar(256),
    auth_srid integer,
    srtext    varchar(2048),
    proj4text varchar(2048)
);

alter table public.spatial_ref_sys
    owner to allbinary;

grant select on public.spatial_ref_sys to public;

create table public.agency
(
    agency_id       varchar(255) not null
        primary key,
    agency_name     varchar(255),
    agency_url      varchar(255),
    agency_timezone varchar(255),
    agency_lang     varchar(255)
);

alter table public.agency
    owner to allbinary;

create table public.calendar_dates
(
    service_id     varchar(255) not null,
    date           date         not null,
    exception_type smallint,
    primary key (service_id, date)
);

alter table public.calendar_dates
    owner to allbinary;

create index calendar_dates_date_idx
    on public.calendar_dates (date);

create table public.routes
(
    route_id         varchar(255) not null
        primary key,
    agency_id        varchar(255)
        references public.agency,
    route_short_name varchar(255),
    route_long_name  varchar(255),
    route_type       varchar(255),
    route_url        varchar(255)
);

alter table public.routes
    owner to allbinary;

create table public.stops
(
    stop_id       varchar(255)  not null
        primary key,
    stop_name     varchar(255)  not null,
    stop_lat      numeric(8, 6) not null,
    stop_lon      numeric(9, 6) not null,
    location_type smallint
        constraint stops_location_type_check
            check ((location_type >= 0) AND (location_type <= 1)),
    geom          geometry(Point, 4326)
);

alter table public.stops
    owner to allbinary;

create index stops_geom
    on public.stops using gist (geom);

create table public.trips
(
    route_id        varchar(255)
        references public.routes,
    service_id      varchar(255) not null,
    trip_id         varchar(255) not null
        primary key,
    trip_headsign   varchar(255),
    trip_short_name varchar(255)
);

alter table public.trips
    owner to allbinary;

create table public.stop_times
(
    trip_id        varchar(255),
    arrival_time   varchar(8) not null,
    departure_time varchar(8) not null,
    stop_id        varchar(255)
        references public.stops,
    stop_sequence  integer    not null,
    pickup_type    smallint,
    drop_off_type  smallint
);

alter table public.stop_times
    owner to allbinary;

create index stop_times_stop_id_idx
    on public.stop_times (stop_id);

