package main

import (
	"database/sql"
)

type ClimaRepository struct {
	db *sql.DB
}

func NewClimaRepository(db *sql.DB) *ClimaRepository {
	return &ClimaRepository{db: db}
}

func (r *ClimaRepository) Create(c ClimaAtual) (int, error) {
	var id int
	query := `INSERT INTO clima_atual (cidade, data, hora, temperatura, umidade, vento, condicao) 
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := r.db.QueryRow(query, c.Cidade, c.Data, c.Hora, c.Temperatura, c.Umidade, c.Vento, c.Condicao).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *ClimaRepository) GetByCidade(cidade string) ([]ClimaAtual, error) {
	query := `SELECT id, cidade, data, hora, temperatura, umidade, vento, condicao FROM clima_atual WHERE cidade ILIKE $1`
	rows, err := r.db.Query(query, cidade)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var climas []ClimaAtual
	for rows.Next() {
		var c ClimaAtual
		if err := rows.Scan(&c.ID, &c.Cidade, &c.Data, &c.Hora, &c.Temperatura, &c.Umidade, &c.Vento, &c.Condicao); err != nil {
			return nil, err
		}
		climas = append(climas, c)
	}
	return climas, nil
}

func (r *ClimaRepository) GetAll() ([]ClimaAtual, error) {
	query := `SELECT id, cidade, data, hora, temperatura, umidade, vento, condicao FROM clima_atual`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var climas []ClimaAtual
	for rows.Next() {
		var c ClimaAtual
		if err := rows.Scan(&c.ID, &c.Cidade, &c.Data, &c.Hora, &c.Temperatura, &c.Umidade, &c.Vento, &c.Condicao); err != nil {
			return nil, err
		}
		climas = append(climas, c)
	}
	return climas, nil
}

func (r *ClimaRepository) Update(c ClimaAtual) error {
	query := `UPDATE clima_atual SET cidade=$1, data=$2, hora=$3, temperatura=$4, umidade=$5, vento=$6, condicao=$7 WHERE id=$8`
	_, err := r.db.Exec(query, c.Cidade, c.Data, c.Hora, c.Temperatura, c.Umidade, c.Vento, c.Condicao, c.ID)
	return err
}

func (r *ClimaRepository) Delete(id int) error {
	query := `DELETE FROM clima_atual WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

// Nova função para buscar previsões usando id_clima_atual
func (r *PrevisaoRepository) GetByClimaID(id int) ([]PrevisaoDias, error) {
	query := `SELECT id, dia_semana, maxima, minima, condicao, id_clima_atual 
              FROM previsao_dias 
              WHERE id_clima_atual = $1`
	rows, err := r.db.Query(query, id)
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
