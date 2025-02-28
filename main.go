package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/koriebruh/gateway-simply/handlers"
	"github.com/koriebruh/gateway-simply/test"
	"log"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Yyws")
}

func main() {

	r := mux.NewRouter()
	//r.PathPrefix("/").HandlerFunc(handlers.ProxyRequest)
	//r.HandleFunc("/hai", Home).Methods("GET")
	r.PathPrefix("/products").HandlerFunc(handlers.ProxyRequest)
	r.PathPrefix("/orders").HandlerFunc(handlers.ProxyRequest)

	log.Print("RUNNING IN PORT 8080")

	go test.DummyClient("products", "8081")
	go test.DummyClient("orders", "8082")

	log.Fatal(http.ListenAndServe(":8080", r))

}
