package organization

import (
	"encoding/json"
	"net/http"
)

type OrganizationsAPI struct {
}

// Get organizations. Authorization limits are applied to requesting user.

// It is handler for GET /organizations
func (api OrganizationsAPI) Get(w http.ResponseWriter, r *http.Request) {
	var respBody Organizations
	json.NewEncoder(w).Encode(&respBody)
	// uncomment below line to add header
	// w.Header().Set("key","value")
}

// Create new organization
// It is handler for POST /organizations
func (api OrganizationsAPI) Post(w http.ResponseWriter, r *http.Request) {
	var reqBody Organization

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(400)
		return
	}

	var respBody Organization
	json.NewEncoder(w).Encode(&respBody)
	// uncomment below line to add header
	// w.Header().Set("key","value")
}

// Get organization info
// It is handler for GET /organizations/{globalid}
func (api OrganizationsAPI) globalidGet(w http.ResponseWriter, r *http.Request) {
	var respBody Organization
	json.NewEncoder(w).Encode(&respBody)
	// uncomment below line to add header
	// w.Header().Set("key","value")
}

// Update organization info
// It is handler for PUT /organizations/{globalid}
func (api OrganizationsAPI) globalidPut(w http.ResponseWriter, r *http.Request) {
	var reqBody Organization

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(400)
		return
	}
	var respBody Organization
	json.NewEncoder(w).Encode(&respBody)
	// uncomment below line to add header
	// w.Header().Set("key","value")
}

// Assign a member to organization
// It is handler for POST /organizations/{globalid}/members
func (api OrganizationsAPI) globalidmembersPost(w http.ResponseWriter, r *http.Request) {
	var reqBody member

	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(400)
		return
	}
	var respBody member
	json.NewEncoder(w).Encode(&respBody)
	// uncomment below line to add header
	// w.Header().Set("key","value")
}

// Remove a member from organization
// It is handler for DELETE /organizations/{globalid}/members/{username}
func (api OrganizationsAPI) globalidmembersusernameDelete(w http.ResponseWriter, r *http.Request) {
	// uncomment below line to add header
	// w.Header().Set("key","value")
}
