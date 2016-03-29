package company

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/itsyouonline/identityserver/db"
)

type CompaniesAPI struct {
}

// Register a new company
// It is handler for POST /companies
func (api CompaniesAPI) Post(w http.ResponseWriter, r *http.Request) {

	var company Company

	if err := json.NewDecoder(r.Body).Decode(&company); err != nil {
		log.Debug("Error decoding the company:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !company.IsValid() {
		log.Debug("Invalid organization")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	companyMgr := NewCompanyManager(r)
	err := companyMgr.Create(&company)
	if err != nil && err != db.ErrDuplicate {
		log.Error("Error saving company:", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err == db.ErrDuplicate {
		log.Debug("Duplicate company:", company)
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

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

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(&respBody)
}

// It is handler for GET /companies/{globalid}/validate
func (api CompaniesAPI) globalIdvalidateGet(w http.ResponseWriter, r *http.Request) {
	log.Error("globalIdvalidateGet is not implemented")
}

// Get the contracts where the organization is 1 of the parties. Order descending by
// date.
// It is handler for GET /companies/{globalId}/contracts
func (api CompaniesAPI) globalIdcontractsGet(w http.ResponseWriter, r *http.Request) {
	log.Error("globalIdcontractsGet is not implemented")
}

// GetCompanyList is the handler for GET /companies
// Get companies. Authorization limits are applied to requesting user.
func (api CompaniesAPI) GetCompanyList(w http.ResponseWriter, r *http.Request) {
	log.Error("GetCompanyList is not implemented")
}

// globalIdGet is the handler for GET /companies/{globalId}
// Get organization info
func (api CompaniesAPI) globalIdGet(w http.ResponseWriter, r *http.Request) {
	log.Error("globalIdGet is not implemented")
}
