package organization

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type UsersusernameorganizationsAPI struct {
}

func exists(value string, list []string) bool {
	for _, val := range list {
		if val == value {
			return true
		}
	}

	return false
}

// Get the list organizations a user is owner or member of
// It is handler for GET /users/{username}/organizations
func (api UsersusernameorganizationsAPI) Get(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	orgMgr := NewManager(r)

	orgs, err := orgMgr.AllByUser(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	userOrgs := UserOrganizations{
		Member: []string{},
		Owner:  []string{},
	}

	for _, org := range orgs {
		if exists(username, org.Owners) {
			userOrgs.Owner = append(userOrgs.Owner, org.Globalid)
		} else {
			userOrgs.Member = append(userOrgs.Member, org.Globalid)
		}
	}
	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(&userOrgs)
}

// Accept membership in organization
// It is handler for POST /users/{username}/organizations/{globalid}/roles/{role}
func (api UsersusernameorganizationsAPI) globalidrolesrolePost(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	role := mux.Vars(r)["role"]
	organization := mux.Vars(r)["globalid"]

	var j JoinOrganizationRequest

	if err := json.NewDecoder(r.Body).Decode(&j); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	orgReqMgr := NewOrganizationRequestManager(r)

	orgRequest, err := orgReqMgr.Get(username, organization, role)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// TODO: Save member
	orgMgr := NewManager(r)

	if org, err := orgMgr.GetByName(organization); err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	} else {
		if exists(RoleOwner, orgRequest.Role) {
			// Accepted Owner role
			if err := orgMgr.SaveOwner(org, username); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		} else {
			// Accepted member role
			if err := orgMgr.SaveMember(org, username); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}
	}

	orgRequest.Status = RequestAccepted

	if err := orgReqMgr.Save(orgRequest); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(orgRequest)

	w.WriteHeader(http.StatusCreated)
}

// Reject membership invitation in an organization.

// It is handler for DELETE /users/{username}/organizations/{globalid}/roles/{role}
func (api UsersusernameorganizationsAPI) globalidrolesroleDelete(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	role := mux.Vars(r)["role"]
	organization := mux.Vars(r)["globalid"]

	orgReqMgr := NewOrganizationRequestManager(r)

	orgRequest, err := orgReqMgr.Get(username, organization, role)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	orgMgr := NewManager(r)

	if _, err := orgMgr.GetByName(organization); err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	orgRequest.Status = RequestRejected

	if err := orgReqMgr.Save(orgRequest); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
