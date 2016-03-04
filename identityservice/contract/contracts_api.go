package contract

import (
	"encoding/json"
	"net/http"
)

type ContractsAPI struct {
}

// Create a new contract.
// It is handler for POST /contracts
func (api ContractsAPI) Post(w http.ResponseWriter, r *http.Request) {
	var respBody Contract
	json.NewEncoder(w).Encode(&respBody)
	// uncomment below line to add header
	// w.Header().Set("key","value")
}

// Sign a contract
// It is handler for POST /contracts/{contractId}/signatures
func (api ContractsAPI) contractIdsignaturesPost(w http.ResponseWriter, r *http.Request) {
	var reqBody Signature

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(400)
		return
	}
	// uncomment below line to add header
	// w.Header().Set("key","value")
}

// Get a contract
// It is handler for GET /contracts/{contractId}
func (api ContractsAPI) contractIdGet(w http.ResponseWriter, r *http.Request) {
	var respBody Contract
	json.NewEncoder(w).Encode(&respBody)
	// uncomment below line to add header
	// w.Header().Set("key","value")
}
