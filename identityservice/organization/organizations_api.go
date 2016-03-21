package organization

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"github.com/itsyouonline/identityserver/db"
	"github.com/itsyouonline/identityserver/identityservice/invitations"
	"github.com/itsyouonline/identityserver/identityservice/user"
	"github.com/itsyouonline/identityserver/oauthservice"
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
		log.Debug("Error decoding the organization:", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if !org.IsValid() {
		log.Debug("Invalid organization")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	orgMgr := NewManager(r)

	err := orgMgr.Create(&org)

	if err != nil && err != db.ErrDuplicate {
		log.Error("Error saving organizations:", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err == db.ErrDuplicate {
		log.Debug("Duplicate organization")
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(&org)
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

	org, err := orgMgr.GetByName(globalid)
	if err != nil { //TODO: make a distinction with an internal server error
		log.Debug("Error while getting the organization: ", err)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// Check if user exists
	userMgr := user.NewManager(r)

	if ok, err := userMgr.Exists(m.Username); err != nil || !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	for _, membername := range org.Members {
		if membername == m.Username {
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}
	}

	for _, membername := range org.Owners {
		if membername == m.Username {
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}
	}

	// Create JoinRequest
	invitationMgr := invitations.NewInvitationManager(r)

	orgReq := &invitations.JoinOrganizationInvitation{
		Role:         invitations.RoleMember,
		Organization: globalid,
		User:         m.Username,
		Status:       invitations.RequestPending,
	}

	if err := invitationMgr.Save(orgReq); err != nil {
		log.Error("Error inviting member: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(orgReq)
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

// It is handler for POST /organizations/{globalid}/members
func (api OrganizationsAPI) globalidownersPost(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]

	var m member

	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	orgMgr := NewManager(r)

	org, err := orgMgr.GetByName(globalid)
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

	for _, membername := range org.Owners {
		if membername == m.Username {
			http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
			return
		}
	}
	// Create JoinRequest
	invitationMgr := invitations.NewInvitationManager(r)

	orgReq := &invitations.JoinOrganizationInvitation{
		Role:         invitations.RoleOwner,
		Organization: globalid,
		User:         m.Username,
		Status:       invitations.RequestPending,
	}

	if err := invitationMgr.Save(orgReq); err != nil {
		log.Error("Error inviting owner: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(orgReq)
}

// Remove a member from organization
// It is handler for DELETE /organizations/{globalid}/members/{username}
func (api OrganizationsAPI) globalidownersusernameDelete(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]
	username := mux.Vars(r)["username"]

	orgMgr := NewManager(r)

	org, err := orgMgr.GetByName(globalid)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err := orgMgr.RemoveOwner(org, username); err != nil {
		log.Error("Error removing owner: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Get the list of pending invitations for users to join this organization.
// It is handler for GET /organizations/{globalid}/invitations
func (api OrganizationsAPI) GetPendingInvitations(w http.ResponseWriter, r *http.Request) {
	globalid := mux.Vars(r)["globalid"]

	invitationMgr := invitations.NewInvitationManager(r)

	requests, err := invitationMgr.GetPendingByOrganization(globalid)

	if err != nil {
		log.Error("Error in GetPendingByOrganization: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	pendingInvites := make([]Invitation, len(requests), len(requests))
	for index, request := range requests {
		pendingInvites[index] = Invitation{
			Role: request.Role,
			User: request.User,
		}
	}
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(pendingInvites)
}

// Cancel a pending invitation.
// It is handler for DELETE /organizations/{globalid}/invitations/{username}
func (api OrganizationsAPI) RemovePendingInvitation(w http.ResponseWriter, r *http.Request) {
	log.Error("RemovePendingInvitation is not implemented")
}

// Get the contracts where the organization is 1 of the parties. Order descending by
// date.
// It is handler for GET /organizations/{globalid}/contracts
func (api OrganizationsAPI) GetContracts(w http.ResponseWriter, r *http.Request) {
	log.Error("GetContracts is not implemented")
}

// GetAPISecretLabels gets the list of labels that are defined for active api secrets. The secrets themselves
// are not included.
// It is handler for GET /organizations/{globalid}/apisecrets
func (api OrganizationsAPI) GetAPISecretLabels(w http.ResponseWriter, r *http.Request) {
	organization := mux.Vars(r)["globalid"]

	mgr := oauthservice.NewManager(r)
	labels, err := mgr.GetClientSecretLabels(organization)
	if err != nil {
		log.Error("Error getting a client secret labels: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(labels)
}

func isValidAPISecretLabel(label string) (valid bool) {
	valid = true
	labelLength := len(label)
	valid = valid && labelLength > 2 && labelLength < 51
	return valid
}

// GetAPISecret is the handler for GET /organizations/{globalid}/apisecrets/{label}
func (api OrganizationsAPI) GetAPISecret(w http.ResponseWriter, r *http.Request) {
	organization := mux.Vars(r)["globalid"]
	label := mux.Vars(r)["label"]

	mgr := oauthservice.NewManager(r)
	secret, err := mgr.GetClientSecret(organization, label)
	if err != nil {
		log.Error("Error getting a client secret: ", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if secret == "" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	response := struct {
		Label  string `json:"label"`
		Secret string `json:"secret"`
	}{
		Label:  label,
		Secret: secret,
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}

// CreateNewAPISecret creates a new API Secret
// It is handler for POST /organizations/{globalid}/apisecrets
func (api OrganizationsAPI) CreateNewAPISecret(w http.ResponseWriter, r *http.Request) {
	organization := mux.Vars(r)["globalid"]

	body := struct{ label string }{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if !isValidAPISecretLabel(body.label) {
		log.Debug("Invalid label: ", body.label)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	c := oauthservice.NewOauth2Client(organization, body.label)

	mgr := oauthservice.NewManager(r)
	err := mgr.CreateClientSecret(c)

	if err != nil && err != db.ErrDuplicate {
		log.Error("Error creating api secret label", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err == db.ErrDuplicate {
		log.Debug("Duplicate label")
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	response := struct {
		Label  string `json:"label"`
		Secret string `json:"secret"`
	}{
		Label:  c.Label,
		Secret: c.Secret,
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}

// UpdateAPISecretLabel updates the label of the secret
// It is handler for PUT /organizations/{globalid}/apisecrets/{label}
func (api OrganizationsAPI) UpdateAPISecretLabel(w http.ResponseWriter, r *http.Request) {
	organization := mux.Vars(r)["globalid"]
	oldlabel := mux.Vars(r)["label"]

	body := struct{ label string }{}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	if !isValidAPISecretLabel(body.label) {
		log.Debug("Invalid label: ", body.label)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	mgr := oauthservice.NewManager(r)
	err := mgr.RenameClientSecret(organization, oldlabel, body.label)

	if err != nil && err != db.ErrDuplicate {
		log.Error("Error renaming api secret label", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err == db.ErrDuplicate {
		log.Debug("Duplicate label")
		http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// DeleteAPISecret removes an API secret
// It is handler for DELETE /organizations/{globalid}/apisecrets/{label}
func (api OrganizationsAPI) DeleteAPISecret(w http.ResponseWriter, r *http.Request) {
	organization := mux.Vars(r)["globalid"]
	label := mux.Vars(r)["label"]

	mgr := oauthservice.NewManager(r)
	mgr.DeleteClientSecret(organization, label)

	w.WriteHeader(http.StatusNoContent)
}
