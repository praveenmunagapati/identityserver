package organization

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"github.com/itsyouonline/identityserver/identityservice/user"
)

type OrganizationsAPI struct {
}

// Get organizations. Authorization limits are applied to requesting user.

// It is handler for GET /organizations
func (api OrganizationsAPI) Get(w http.ResponseWriter, r *http.Request) {
	orgMgr := NewManager(r)

	// TODO: extract user to apply auth filters.
	respBody, err := orgMgr.All()
	if err != nil {
		log.Error("Error retrieving organizations,", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respBody)
}

// Create new organization
// It is handler for POST /organizations
func (api OrganizationsAPI) Post(w http.ResponseWriter, r *http.Request) {
	var org Organization

	if err := json.NewDecoder(r.Body).Decode(&org); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	orgMgr := NewManager(r)

	if err := orgMgr.Save(&org); err != nil {
		log.Error("Error saving organizations,", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&org)

	w.WriteHeader(http.StatusCreated)
}

// Get organization info
// It is handler for GET /organizations/{globalid}
func (api OrganizationsAPI) globalidGet(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]
	orgMgr := NewManager(r)

	org, err := orgMgr.GetByName(globalid)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(org)
}

// Update organization info
// It is handler for PUT /organizations/{globalid}
func (api OrganizationsAPI) globalidPut(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]

	var org Organization

	if err := json.NewDecoder(r.Body).Decode(&org); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	orgMgr := NewManager(r)

	oldOrg, err := orgMgr.GetByName(globalid)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if org.Globalid != globalid || org.GetId() != oldOrg.GetId() {
		http.Error(w, "Changing globalid or id is Forbidden!", http.StatusForbidden)
		return
	}

	// Update only certain fields
	oldOrg.PublicKeys = org.PublicKeys
	oldOrg.Dns = org.Dns

	if err := orgMgr.Save(oldOrg); err != nil {
		log.Error("Error while saving organization: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(oldOrg)
}

// Assign a member to organization
// It is handler for POST /organizations/{globalid}/members
func (api OrganizationsAPI) globalidmembersPost(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]

	var m member

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	orgMgr := NewManager(r)

	_, err := orgMgr.GetByName(globalid)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Check if user exists
	userMgr := user.NewManager(r)

	if ok, err := userMgr.Exists(m.Username); err != nil || !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Create JoinRequest
	orgReqMgr := NewOrganizationRequestManager(r)

	orgReq := &JoinOrganizationRequest{
		Role:         []string{RoleMember},
		Organization: globalid,
		User:         m.Username,
	}

	if err := orgReqMgr.Save(orgReq); err != nil {
		log.Error("Error inviting member: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orgReq)

	w.WriteHeader(http.StatusCreated)
}

// Remove a member from organization
// It is handler for DELETE /organizations/{globalid}/members/{username}
func (api OrganizationsAPI) globalidmembersusernameDelete(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]
	username := mux.Vars(r)["username"]

	orgMgr := NewManager(r)

	org, err := orgMgr.GetByName(globalid)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err := orgMgr.RemoveMember(org, username); err != nil {
		log.Error("Error adding member: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}