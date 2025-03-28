package locations

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Repo struct {
	conn *pgx.Conn
}

func New(conn *pgx.Conn) *Repo {
	return &Repo{
		conn: conn,
	}
}

// CRUD Locations

func (r *Repo) CreateLocation(ctx context.Context, loc Location) (int, error) {
	query := `
	insert into locations
	 (loc_name, lon, lat, alt) 
	 values ($1, $2, $3, $4) returning id;
	 `

	row := r.conn.QueryRow(ctx, query, loc.Name, loc.Point.Lon, loc.Point.Lat, loc.Point.Alt)
	var id int
	err := row.Scan(&id)

	return id, err
}

func (r *Repo) GetLocation(ctx context.Context, id int) (Location, error) {
	loc := Location{}

	err := r.conn.QueryRow(ctx, "select id, loc_name, lon, lat, alt from locations where id=$1", id).
		Scan(&loc.ID, &loc.Name, &loc.Point.Lon, &loc.Point.Lat, &loc.Point.Alt)
	if err != nil {
		return Location{}, err
	}

	return loc, nil
}

func (r *Repo) UpdateLocation(ctx context.Context, loc Location) error {
	query := `
		update locations 
		set loc_name = $1, lon = $2, lat = $3, alt = $4 
		where id=$5
	`

	_, err := r.conn.Exec(ctx, query, loc.Name, loc.Point.Lon, loc.Point.Lat, loc.Point.Alt, loc.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) DeleteLocation(ctx context.Context, id int) error {
	_, err := r.conn.Exec(ctx, "delete from locations where id=$1", id)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) FindLocation(ctx context.Context, filter FilterLocation) ([]Location, error) {
	query := `
	select * from locations 
	where 1=1
	AND CASE
		WHEN $1::text IS NOT NULL THEN loc_name ilike '%' || $1 || '%'
		ELSE true
	END
`

	rows, err := r.conn.Query(ctx, query, filter.Name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var locs []Location

	for rows.Next() {
		var loc Location

		err := rows.Scan(&loc.ID, &loc.Name, &loc.Point.Lon, &loc.Point.Lat, &loc.Point.Alt)
		if err != nil {
			return nil, fmt.Errorf("не удалось вернуть локацию %w", err)
		}

		locs = append(locs, loc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации по результату из бд: %w", err)
	}

	return locs, nil
}
