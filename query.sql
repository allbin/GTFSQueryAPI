-- name: GetStops :many
select * from stops limit 50;

-- name: GetStopsNearby :many
select *
from stops
where st_dwithin(
    geography(st_point(stop_lon, stop_lat)),
    geography(st_point(@lon::double precision, @lat::double precision)),
    @radius::double precision
)
order by st_distance(
    st_point(stop_lon, stop_lat),
    st_point(@lon::double precision, @lat::double precision)
)
limit @lim::bigint;
