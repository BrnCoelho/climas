package main

import "time"

type ClimaAtual struct {
	ID          int       `json:"id"`
	Cidade      string    `json:"cidade"`
	Data        time.Time `json:"data"`
	Hora        string    `json:"hora"`
	Temperatura int       `json:"temperatura"`
	Umidade     int       `json:"umidade"`
	Vento       string    `json:"vento"`
	Condicao    string    `json:"condicao"`
}

type PrevisaoDias struct {
	ID           int    `json:"id"`
	DiaSemana    string `json:"dia_semana"`
	Maxima       int    `json:"maxima"`
	Minima       int    `json:"minima"`
	Condicao     string `json:"condicao"`
	IDClimaAtual int    `json:"id_clima_atual"`
}
