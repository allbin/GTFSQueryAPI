// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package storage

import (
	"context"
)

const getStops = `-- name: GetStops :many
select stop_id, stop_name, stop_lat, stop_lon, location_type from stops limit 50
`

func (q *Queries) GetStops(ctx context.Context) ([]Stop, error) {
	rows, err := q.db.Query(ctx, getStops)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Stop
	for rows.Next() {
		var i Stop
		if err := rows.Scan(
			&i.StopID,
			&i.StopName,
			&i.StopLat,
			&i.StopLon,
			&i.LocationType,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getStopsNearby = `-- name: GetStopsNearby :many
select stop_id, stop_name, stop_lat, stop_lon, location_type
from stops
where st_dwithin(
    geography(st_point(stop_lon, stop_lat)),
    geography(st_point($1::double precision, $2::double precision)),
    $3::double precision
)
order by st_distance(
    st_point(stop_lon, stop_lat),
    st_point($1::double precision, $2::double precision)
)
limit $4::bigint
`

type GetStopsNearbyParams struct {
	Lon    float64
	Lat    float64
	Radius float64
	Lim    int64
}

func (q *Queries) GetStopsNearby(ctx context.Context, arg GetStopsNearbyParams) ([]Stop, error) {
	rows, err := q.db.Query(ctx, getStopsNearby,
		arg.Lon,
		arg.Lat,
		arg.Radius,
		arg.Lim,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Stop
	for rows.Next() {
		var i Stop
		if err := rows.Scan(
			&i.StopID,
			&i.StopName,
			&i.StopLat,
			&i.StopLon,
			&i.LocationType,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}