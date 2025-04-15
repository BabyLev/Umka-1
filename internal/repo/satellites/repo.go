package satellites

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repo struct {
	conn *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repo {
	return &Repo{
		conn: pool,
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
	var args []interface{}
	query := "select id, sat_name, norad_id, line1, line2 from satellites where 1=1"

	argId := 1

	if len(filter.IDs) > 0 {
		query += fmt.Sprintf(" AND id = ANY($%d::int[])", argId)
		args = append(args, filter.IDs)
		argId++
	}

	if len(filter.NoradIDs) > 0 {
		var noradIDs64 []int64
		for _, id := range filter.NoradIDs {
			noradIDs64 = append(noradIDs64, int64(id))
		}
		query += fmt.Sprintf(" AND norad_id = ANY($%d::bigint[])", argId)
		args = append(args, noradIDs64)
		argId++
	}

	if filter.SatName != nil && *filter.SatName != "" {
		query += fmt.Sprintf(" AND sat_name ilike $%d", argId)
		args = append(args, "%"+*filter.SatName+"%")
		argId++
	}

	if filter.NoradIDNotNull != nil {
		if *filter.NoradIDNotNull {
			query += " AND norad_id IS NOT NULL"
		} else {
			query += " AND norad_id IS NULL"
		}
	}

	rows, err := r.conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса FindSatellite: %w", err) // Добавим контекст ошибке
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
