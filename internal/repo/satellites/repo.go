package satellites

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

// CRUD Satellites
func (r *Repo) CreateSatellite(ctx context.Context, sat Satellite) (int, error) {
	query := `
	insert into satellites
	 (sat_name, norad_id, line1, line2) 
	 values ($1, $2, $3, $4) returning id;
	 `

	row := r.conn.QueryRow(ctx, query, sat.SatName, sat.NoradID, sat.Line1, sat.Line2)
	var id int
	err := row.Scan(&id)

	return id, err
}

func (r *Repo) GetSatellite(ctx context.Context, id int) (Satellite, error) {
	sat := Satellite{}

	err := r.conn.QueryRow(ctx, "select id, sat_name, norad_id, line1, line2 from satellites where id=$1", id).
		Scan(&sat.ID, &sat.SatName, &sat.NoradID, &sat.Line1, &sat.Line2)
	if err != nil {
		return Satellite{}, err
	}

	return sat, nil
}

func (r *Repo) FindSatellite(ctx context.Context, filter FilterSatellite) ([]Satellite, error) {
	query := `
	select * from satellites 
	where 1=1
	AND CASE 
		WHEN array_length($1::int[], 1) > 0 THEN id = ANY($1::int[])
		ELSE true
	END
	AND CASE
		WHEN array_length($2::int[], 1) > 0 THEN norad_id = ANY($2::int[])
		ELSE true
	END
	AND CASE
		WHEN $3::text IS NOT NULL THEN sat_name ilike '%' || '$3::text' || '%'
		ELSE true
	END
	AND CASE
		WHEN $4::boolean IS NOT NULL THEN 
			WHEN $4::boolean = TRUE THEN norad_id is not null
			ELSE norad_id is null
		ELSE true
	END
`

	rows, err := r.conn.Query(ctx, query, filter.IDs, filter.NoradIDs, filter.SatName, filter.NoradIDNotNull)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sats []Satellite

	for rows.Next() {
		var sat Satellite

		err := rows.Scan(&sat.ID, &sat.SatName, &sat.NoradID, &sat.Line1, &sat.Line2)
		if err != nil {
			return nil, fmt.Errorf("не удалось вернуть спутник %w", err)
		}

		sats = append(sats, sat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка итерации по результату из бд: %w", err)
	}

	return sats, nil
}

func (r *Repo) UpdateSatellite(ctx context.Context, sat Satellite) error {
	query := `
		update satellites 
		set sat_name = $1, norad_id = $2, line1 = $3, line2 = $4 
		where id=$5
	`

	_, err := r.conn.Exec(ctx, query, sat.SatName, sat.NoradID, sat.Line1, sat.Line2, sat.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repo) DeleteSatellite(ctx context.Context, id int) error {
	_, err := r.conn.Exec(ctx, "delete from satellites where id=$1", id)
	if err != nil {
		return err
	}

	return nil
}
