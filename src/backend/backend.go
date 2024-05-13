package backend

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"io"
	"log"
	"net/http"
)

type App struct {
	DB     *sql.DB
	Port   string
	Router *mux.Router
}

func (a *App) Initialize() {
	DB, err := sql.Open("sqlite3", "../gymshop.db")
	if err != nil {
		log.Fatal(err.Error())
	}

	a.DB = DB
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

func (a *App) Run() {
	fmt.Println("Server started listening on port ", a.Port)
	log.Fatal(http.ListenAndServe(a.Port, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/orders", a.allOrders).Methods("GET")
	a.Router.HandleFunc("/orders", a.newOrder).Methods("POST")
}

func (a *App) allOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := getOrders(a.DB)
	if err != nil {
		fmt.Printf("getOrders error: %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, orders)
}

func (a *App) newOrder(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := io.ReadAll(r.Body)

	var o order
	json.Unmarshal(reqBody, &o)

	err := o.createOrder(a.DB)
	if err != nil {
		fmt.Printf("newOrder error: %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
