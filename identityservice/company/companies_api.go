package company

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

type CompaniesAPI struct {
}

// Register a new company
// It is handler for POST /companies
func (api CompaniesAPI) Post(w http.ResponseWriter, r *http.Request) {

	var company Company

	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	companyMgr := NewCompanyManager(r)
	err := companyMgr.Save(&company)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(&company)
}

// Update existing company. Updating ``globalId`` is not allowed.
// It is handler for PUT /companies/{globalId}
func (api CompaniesAPI) globalIdPut(w http.ResponseWriter, r *http.Request) {

	globalID := mux.Vars(r)["globalId"]

	var company Company

	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	companyMgr := NewCompanyManager(r)

	oldCompany, cerr := companyMgr.GetByName(globalID)
	if cerr != nil {
		log.Debug(cerr)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if company.Globalid != globalID || company.GetId() != oldCompany.GetId() {
		http.Error(w, "Changing globalId or id is Forbidden!", http.StatusForbidden)
		return
	}

	if err := companyMgr.Save(&company); err != nil {
		log.Error("Error saving company:\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

}

// It is handler for GET /companies/{globalid}/info
func (api CompaniesAPI) globalIdinfoGet(w http.ResponseWriter, r *http.Request) {
	companyMgr := NewCompanyManager(r)

	globalID := mux.Vars(r)["globalId"]

	company, err := companyMgr.GetByName(globalID)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	respBody := company
	json.NewEncoder(w).Encode(&respBody)

}

// It is handler for GET /companies/{globalid}/validate
func (api CompaniesAPI) globalIdvalidateGet(w http.ResponseWriter, r *http.Request) {

	// token := req.FormValue("token")

	// uncomment below line to add header
	// w.Header.Set("key","value")
}

// Get the contracts where the organization is 1 of the parties. Order descending by
// date.
// It is handler for GET /companies/{globalId}/contracts
func (api CompaniesAPI) globalIdcontractsGet(w http.ResponseWriter, r *http.Request) {}
