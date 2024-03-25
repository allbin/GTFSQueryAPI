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

-- name: GetStop :one
select * from stops where stop_id = @stop_id;

-- name: GetStopDepartures :many
SELECT st.stop_id::text                       AS id,
       cd.date::text                          AS date,
       st.arrival_time::text                  AS arrival,
       st.departure_time::text                AS departure,
       s.stop_name::text                      AS name,
       s.stop_lat::double precision           AS lat,
       s.stop_lon::double precision           AS lon,
       t.trip_headsign::text                  AS headsign,
       COALESCE(r.route_short_name, '')::text AS short_name,
       COALESCE(r.route_long_name, '')::text  AS long_name
FROM stop_times st
         INNER JOIN stops s ON st.stop_id = s.stop_id
         INNER JOIN trips t ON st.trip_id = t.trip_id
         INNER JOIN routes r ON t.route_id = r.route_id
         INNER JOIN calendar_dates cd ON t.service_id = cd.service_id
WHERE s.stop_id = @stop_id
  AND cd.date IN (current_date, current_date + interval '1 day')
ORDER BY date,
         departure
LIMIT @lim::bigint;

-- name: GetDeparturesForStops :many
SELECT st.stop_id::text                       AS id,
       cd.date::text                          AS date,
       st.arrival_time::text                  AS arrival,
       st.departure_time::text                AS departure,
       s.stop_name::text                      AS name,
       s.stop_lat::double precision           AS lat,
       s.stop_lon::double precision           AS lon,
       t.trip_headsign::text                  AS headsign,
       COALESCE(r.route_short_name, '')::text AS short_name,
       COALESCE(r.route_long_name, '')::text  AS long_name
FROM stop_times st
         INNER JOIN stops s ON st.stop_id = s.stop_id
         INNER JOIN trips t ON st.trip_id = t.trip_id
         INNER JOIN routes r ON t.route_id = r.route_id
         INNER JOIN calendar_dates cd ON t.service_id = cd.service_id
WHERE s.stop_id = ANY(@stop_id::text[])
  AND cd.date IN (current_date, current_date + interval '1 day')
ORDER BY date,
         departure
LIMIT @lim::bigint;
