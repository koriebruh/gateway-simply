package handlers

import (
	"github.com/koriebruh/gateway-simply/config"
	"github.com/koriebruh/gateway-simply/utils"
	"io"
	"net/http"
	"strings"
)

func ProxyRequest(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/")
	segments := strings.SplitN(path, "/", 2)

	if len(segments) < 2 {
		utils.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid API path"})
		return
	}

	serviceName, endpoint := segments[0], segments[1]
	serviceURL, exists := config.Services[serviceName]
	if !exists {
		utils.WriteJSONResponse(w, http.StatusNotFound, map[string]string{"error": "Service not found"})
		return
	}

	targetURL := serviceURL + "/" + endpoint
	req, err := http.NewRequest(r.Method, targetURL, r.Body)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create request"})
		return
	}

	// Copy headers
	for key, values := range r.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to reach backend service"})
		return
	}
	defer resp.Body.Close()

	// Copy response
	body, _ := io.ReadAll(resp.Body)
	w.WriteHeader(resp.StatusCode)
	w.Write(body)
}
