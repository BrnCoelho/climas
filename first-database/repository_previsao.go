package main

import (
	"database/sql"
)

type PrevisaoRepository struct {
	db *sql.DB
}

func NewPrevisaoRepository(db *sql.DB) *PrevisaoRepository {
	return &PrevisaoRepository{db: db}
}

func (r *PrevisaoRepository) Create(p PrevisaoDias) (int, error) {
	var id int
	query := `INSERT INTO previsao_dias (dia_semana, maxima, minima, condicao, id_clima_atual) 
              VALUES ($1, $2, $3, $4, $5) RETURNING id`
	err := r.db.QueryRow(query, p.DiaSemana, p.Maxima, p.Minima, p.Condicao, p.IDClimaAtual).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *PrevisaoRepository) GetByDiaSemana(dia string) ([]PrevisaoDias, error) {
	query := `SELECT id, dia_semana, maxima, minima, condicao, id_clima_atual FROM previsao_dias WHERE dia_semana ILIKE $1`
	rows, err := r.db.Query(query, dia)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var previsoes []PrevisaoDias
	for rows.Next() {
		var p PrevisaoDias
		if err := rows.Scan(&p.ID, &p.DiaSemana, &p.Maxima, &p.Minima, &p.Condicao, &p.IDClimaAtual); err != nil {
			return nil, err
		}
		previsoes = append(previsoes, p)
	}
	return previsoes, nil
}

func (r *PrevisaoRepository) GetAll() ([]PrevisaoDias, error) {
	query := `SELECT id, dia_semana, maxima, minima, condicao, id_clima_atual FROM previsao_dias`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var previsoes []PrevisaoDias
	for rows.Next() {
		var p PrevisaoDias
		if err := rows.Scan(&p.ID, &p.DiaSemana, &p.Maxima, &p.Minima, &p.Condicao, &p.IDClimaAtual); err != nil {
			return nil, err
		}
		previsoes = append(previsoes, p)
	}
	return previsoes, nil
}

func (r *PrevisaoRepository) Update(p PrevisaoDias) error {
	query := `UPDATE previsao_dias SET dia_semana=$1, maxima=$2, minima=$3, condicao=$4, id_clima_atual=$5 WHERE id=$6`
	_, err := r.db.Exec(query, p.DiaSemana, p.Maxima, p.Minima, p.Condicao, p.IDClimaAtual, p.ID)
	return err
}

func (r *PrevisaoRepository) Delete(id int) error {
	query := `DELETE FROM previsao_dias WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
