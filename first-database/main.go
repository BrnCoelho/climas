package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

var climaRepo *ClimaRepository
var previsaoRepo *PrevisaoRepository

// Middleware para logar requisições HTTP
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("➡️  %s %s", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
		log.Printf("✅ %s %s concluído em %v", r.Method, r.RequestURI, time.Since(start))
	})
}

func main() {
	connStr := os.Getenv("DB_CONN")
	if connStr == "" {
		connStr = "user=postgres password=bruno2001 dbname=meudatabase sslmode=disable host=localhost port=5777"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(100)

	climaRepo = NewClimaRepository(db)
	previsaoRepo = NewPrevisaoRepository(db)

	router := mux.NewRouter()
	router.Use(logMiddleware) // Ativa o log

	// Rotas para clima_atual
	router.HandleFunc("/clima", createClima).Methods("POST")
	router.HandleFunc("/clima", getClimaByCidade).Methods("GET")
	router.HandleFunc("/clima/all", getAllClimas).Methods("GET")
	router.HandleFunc("/clima/{id}", updateClima).Methods("PUT")
	router.HandleFunc("/clima/{id}", deleteClima).Methods("DELETE")

	// Rotas para previsao_dias
	router.HandleFunc("/previsao", createPrevisao).Methods("POST")
	router.HandleFunc("/previsao", getPrevisaoByDia).Methods("GET")
	router.HandleFunc("/previsao/all", getAllPrevisoes).Methods("GET")
	router.HandleFunc("/previsao/{id}", updatePrevisao).Methods("PUT")
	router.HandleFunc("/previsao/{id}", deletePrevisao).Methods("DELETE")
	router.HandleFunc("/previsao/por-clima/{id}", getPrevisoesPorClimaID).Methods("GET")

	//CORS primeiro
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
	})

	// Envolve com CORS -> depois aplica o log
	handlerComCors := c.Handler(router)
	finalHandler := logMiddleware(handlerComCors) // AQUI É A MÁGICA

	log.Println("🟢 Log de requisições ativado!")
	log.Println("🚀 Servidor iniciado na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", finalHandler))
}

// Handlers clima_atual

func createClima(w http.ResponseWriter, r *http.Request) {
	var c ClimaAtual
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := climaRepo.Create(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	c.ID = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(c)
}

func getClimaByCidade(w http.ResponseWriter, r *http.Request) {
	cidade := r.URL.Query().Get("cidade")
	if cidade == "" {
		http.Error(w, "Parâmetro 'cidade' é obrigatório", http.StatusBadRequest)
		return
	}
	climas, err := climaRepo.GetByCidade(cidade)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(climas) == 0 {
		http.Error(w, "Nenhum registro encontrado para essa cidade", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(climas)
}

func getAllClimas(w http.ResponseWriter, r *http.Request) {
	climas, err := climaRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(climas)
}

func updateClima(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	var c ClimaAtual
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	c.ID = id
	err = climaRepo.Update(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func deleteClima(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	err = climaRepo.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Handlers previsao_dias

func createPrevisao(w http.ResponseWriter, r *http.Request) {
	var p PrevisaoDias
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := previsaoRepo.Create(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p.ID = id
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func getPrevisaoByDia(w http.ResponseWriter, r *http.Request) {
	dia := r.URL.Query().Get("dia")
	if dia == "" {
		http.Error(w, "Parâmetro 'dia' é obrigatório", http.StatusBadRequest)
		return
	}
	previsoes, err := previsaoRepo.GetByDiaSemana(dia)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(previsoes) == 0 {
		http.Error(w, "Nenhuma previsão encontrada para esse dia", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(previsoes)
}

func getAllPrevisoes(w http.ResponseWriter, r *http.Request) {
	previsoes, err := previsaoRepo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(previsoes)
}

func updatePrevisao(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	var p PrevisaoDias
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	p.ID = id
	err = previsaoRepo.Update(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func deletePrevisao(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	err = previsaoRepo.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func getPrevisoesPorClimaID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	previsoes, err := previsaoRepo.GetByClimaID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(previsoes) == 0 {
		http.Error(w, "Nenhuma previsão encontrada para este clima", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(previsoes)
}
