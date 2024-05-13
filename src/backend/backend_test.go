package backend_test

import (
	"bytes"
	"example.com/gymshop/backend"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a backend.App

const tablePacksCreationQuery = `CREATE TABLE IF NOT EXISTS packs
(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    order_id INT,
	quantity INT NOT NULL,
	FOREIGN KEY (order_id) REFERENCES orders (id)
)`

const tableOrdersCreationQuery = `CREATE TABLE IF NOT EXISTS orders
(
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    customerName VARCHAR(256) NOT NULL,
    total INT NOT NULL
)`

func TestMain(m *testing.M) {
	a = backend.App{}
	a.Initialize()
	ensureTableExists()
	code := m.Run()

	clearPacksTable()
	clearOrdersTable()
	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tablePacksCreationQuery); err != nil {
		log.Fatal(err)
	}

	if _, err := a.DB.Exec(tableOrdersCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearPacksTable() {
	a.DB.Exec("DELETE FROM packs")
	a.DB.Exec("DELETE FROM sqlite_sequence WHERE name = 'packs'")
}

func clearOrdersTable() {
	a.DB.Exec("DELETE FROM orders")
	a.DB.Exec("DELETE FROM sqlite_sequence WHERE name = 'orders'")
}

func TestCreateOrder(t *testing.T) {
	clearOrdersTable()

	payload := []byte(`{"customerName":"TestCustomer", "quantity":1}`)

	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
