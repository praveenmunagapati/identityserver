package userorganization

//The reason this api is not in the user package is because this would cause circular imports

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	organizationdb "github.com/itsyouonline/identityserver/db/organization"
	"github.com/itsyouonline/identityserver/identityservice/invitations"
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

	orgMgr := organizationdb.NewManager(r)

	orgs, err := orgMgr.AllByUser(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	type UserOrganizations struct {
		Member []string `json:"member"`
		Owner  []string `json:"owner"`
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

	var j invitations.JoinOrganizationInvitation

	if err := json.NewDecoder(r.Body).Decode(&j); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	orgReqMgr := invitations.NewInvitationManager(r)

	orgRequest, err := orgReqMgr.Get(username, organization, role, invitations.RequestPending)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// TODO: Save member
	orgMgr := organizationdb.NewManager(r)

	if org, err := orgMgr.GetByName(organization); err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	} else {
		if invitations.RoleOwner == orgRequest.Role {
			// Accepted Owner role
			if err := orgMgr.SaveOwner(org, username); err != nil {
				log.Error("Failed to save owner: ", username)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		} else {
			// Accepted member role
			if err := orgMgr.SaveMember(org, username); err != nil {
				log.Error("Failed to save member: ", username)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}
	}

	orgRequest.Status = invitations.RequestAccepted

	if err := orgReqMgr.Save(orgRequest); err != nil {
		log.Error("Failed to update org request status: ", orgRequest.Organization)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(orgRequest)
}

// Reject membership invitation in an organization.

// It is handler for DELETE /users/{username}/organizations/{globalid}/roles/{role}
func (api UsersusernameorganizationsAPI) globalidrolesroleDelete(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	role := mux.Vars(r)["role"]
	organization := mux.Vars(r)["globalid"]

	orgReqMgr := invitations.NewInvitationManager(r)

	orgRequest, err := orgReqMgr.Get(username, organization, role, invitations.RequestPending)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	orgMgr := organizationdb.NewManager(r)

	if org, err := orgMgr.GetByName(organization); err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	} else {
		if invitations.RoleOwner == orgRequest.Role {
			// Rejected Owner role
			if err := orgMgr.RemoveOwner(org, username); err != nil {
				log.Error("Failed to remove owner: ", username)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		} else {
			// Rejected member role
			if err := orgMgr.RemoveMember(org, username); err != nil {
				log.Error("Failed to reject member: ", username)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}
	}

	orgRequest.Status = invitations.RequestRejected

	if err := orgReqMgr.Save(orgRequest); err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
