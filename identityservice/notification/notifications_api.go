package notification

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/itsyouonline/identityserver/identityservice/contract"
	"github.com/itsyouonline/identityserver/identityservice/organization"
)

type NotificationsAPI struct {
}

// Get the list of notifications, these are pending invitations or approvals
// It is handler for GET /users/{username}/notifications
func (api NotificationsAPI) Get(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	var notification Notification

	orgReqMgr := organization.NewOrganizationRequestManager(r)

	userOrgRequests, err := orgReqMgr.GetByUser(username)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	notification.Invitations = userOrgRequests

	// TODO: Get Approvals and Contract requests
	notification.Approvals = []organization.JoinOrganizationRequest{}
	notification.ContractRequests = []contract.ContractSigningRequest{}

	w.Header().Set("Content-type", "application/json")
	json.NewEncoder(w).Encode(&notification)
}
