package test

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/koriebruh/gateway-simply/utils"
	"log"
	"net/http"
)

func dummyService(serviceName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		welcome := fmt.Sprintf("WELCOME TO %s", serviceName)
		utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"msg": welcome})
	}
}

func DummyClient(serviceName string, dummyHost string) {
	r := mux.NewRouter()
	r.HandleFunc("/", dummyService(serviceName)).Methods("GET")

	port := fmt.Sprintf(":%s", dummyHost)
	log.Printf("OK SERVICE %s RUN IN %s", serviceName, port)
	log.Fatal(http.ListenAndServe(port, r))
}
