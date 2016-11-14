package itsyouonline

import (
	"encoding/json"
	"net/http"
)

type OrganizationsService service

// Create a new organization. 1 user should be in the owners list. Validation is performed to check if the securityScheme allows management on this user.
func (s *OrganizationsService) CreateNewOrganization(organization Organization, headers, queryParams map[string]interface{}) (Organization, *http.Response, error) {
	var u Organization

	resp, err := s.client.doReqWithBody("POST", s.client.BaseURI+"/organizations", &organization, headers, queryParams)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Get organization info
func (s *OrganizationsService) GetOrganization(globalid string, headers, queryParams map[string]interface{}) (Organization, *http.Response, error) {
	var u Organization

	resp, err := s.client.doReqNoBody("GET", s.client.BaseURI+"/organizations/"+globalid, headers, queryParams)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Create a new suborganization.
func (s *OrganizationsService) CreateNewSubOrganization(globalid string, organization Organization, headers, queryParams map[string]interface{}) (Organization, *http.Response, error) {
	var u Organization

	resp, err := s.client.doReqWithBody("POST", s.client.BaseURI+"/organizations/"+globalid, &organization, headers, queryParams)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Update organization info
func (s *OrganizationsService) UpdateOrganization(globalid string, organization Organization, headers, queryParams map[string]interface{}) (Organization, *http.Response, error) {
	var u Organization

	resp, err := s.client.doReqWithBody("PUT", s.client.BaseURI+"/organizations/"+globalid, &organization, headers, queryParams)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Deletes an organization and all data linked to it (join-organization-invitations, oauth_access_tokens, oauth_clients, logo)
func (s *OrganizationsService) DeleteOrganization(globalid string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	// create request object
	return s.client.doReqNoBody("DELETE", s.client.BaseURI+"/organizations/"+globalid, headers, queryParams)
}

// Update the 2FA validity time for the organization
func (s *OrganizationsService) Set2faValidityTime(globalid string, int int, headers, queryParams map[string]interface{}) (*http.Response, error) {

	resp, err := s.client.doReqWithBody("POST", s.client.BaseURI+"/organizations/"+globalid+"/2fa/validity", &int, headers, queryParams)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

// Get the 2FA validity time for the organization, in seconds
func (s *OrganizationsService) Get2faValidityTime(globalid string, headers, queryParams map[string]interface{}) (int, *http.Response, error) {
	var u int

	resp, err := s.client.doReqNoBody("GET", s.client.BaseURI+"/organizations/"+globalid+"/2fa/validity", headers, queryParams)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Create a new API Key, a secret itself should not be provided, it will be generated serverside.
func (s *OrganizationsService) CreateNewOrganizationAPIKey(globalid string, organizationapikey OrganizationAPIKey, headers, queryParams map[string]interface{}) (OrganizationAPIKey, *http.Response, error) {
	var u OrganizationAPIKey

	resp, err := s.client.doReqWithBody("POST", s.client.BaseURI+"/organizations/"+globalid+"/apikeys", &organizationapikey, headers, queryParams)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Get the list of active api keys.
func (s *OrganizationsService) GetOrganizationAPIKeyLabels(globalid string, headers, queryParams map[string]interface{}) ([]string, *http.Response, error) {
	var u []string

	resp, err := s.client.doReqNoBody("GET", s.client.BaseURI+"/organizations/"+globalid+"/apikeys", headers, queryParams)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Updates the label or other properties of a key.
func (s *OrganizationsService) UpdateOrganizationAPIKey(label, globalid string, organizationsglobalidapikeyslabelputreqbody OrganizationsGlobalidApikeysLabelPutReqBody, headers, queryParams map[string]interface{}) (*http.Response, error) {

	resp, err := s.client.doReqWithBody("PUT", s.client.BaseURI+"/organizations/"+globalid+"/apikeys/"+label, &organizationsglobalidapikeyslabelputreqbody, headers, queryParams)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

func (s *OrganizationsService) GetOrganizationAPIKey(label, globalid string, headers, queryParams map[string]interface{}) (OrganizationAPIKey, *http.Response, error) {
	var u OrganizationAPIKey

	resp, err := s.client.doReqNoBody("GET", s.client.BaseURI+"/organizations/"+globalid+"/apikeys/"+label, headers, queryParams)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Removes an API key
func (s *OrganizationsService) DeleteOrganizationAPIKey(label, globalid string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	// create request object
	return s.client.doReqNoBody("DELETE", s.client.BaseURI+"/organizations/"+globalid+"/apikeys/"+label, headers, queryParams)
}

// Create a new contract.
func (s *OrganizationsService) CreateOrganizationContracty(globalid string, contract Contract, headers, queryParams map[string]interface{}) (Contract, *http.Response, error) {
	var u Contract

	resp, err := s.client.doReqWithBody("POST", s.client.BaseURI+"/organizations/"+globalid+"/contracts", &contract, headers, queryParams)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Get the contracts where the organization is 1 of the parties. Order descending by date.
func (s *OrganizationsService) GetOrganizationContracts(globalid string, headers, queryParams map[string]interface{}) (*http.Response, error) {

	resp, err := s.client.doReqNoBody("GET", s.client.BaseURI+"/organizations/"+globalid+"/contracts", headers, queryParams)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

// Creates a new DNS name associated with an organization
func (s *OrganizationsService) CreateOrganizationDNS(dnsname, globalid string, organizationsglobaliddnsdnsnamepostreqbody OrganizationsGlobalidDnsDnsnamePostReqBody, headers, queryParams map[string]interface{}) (OrganizationsGlobalidDnsDnsnamePostRespBody, *http.Response, error) {
	var u OrganizationsGlobalidDnsDnsnamePostRespBody

	resp, err := s.client.doReqWithBody("POST", s.client.BaseURI+"/organizations/"+globalid+"/dns/"+dnsname, &organizationsglobaliddnsdnsnamepostreqbody, headers, queryParams)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Updates an existing DNS name associated with an organization
func (s *OrganizationsService) UpdateOrganizationDNS(dnsname, globalid string, organizationsglobaliddnsdnsnameputreqbody OrganizationsGlobalidDnsDnsnamePutReqBody, headers, queryParams map[string]interface{}) (*http.Response, error) {

	resp, err := s.client.doReqWithBody("PUT", s.client.BaseURI+"/organizations/"+globalid+"/dns/"+dnsname, &organizationsglobaliddnsdnsnameputreqbody, headers, queryParams)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

// Removes a DNS name
func (s *OrganizationsService) DeleteOrganizaitonDNS(dnsname, globalid string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	// create request object
	return s.client.doReqNoBody("DELETE", s.client.BaseURI+"/organizations/"+globalid+"/dns/"+dnsname, headers, queryParams)
}

// Get the list of pending invitations for users to join this organization.
func (s *OrganizationsService) GetPendingOrganizationInvitations(globalid string, headers, queryParams map[string]interface{}) ([]JoinOrganizationInvitation, *http.Response, error) {
	var u []JoinOrganizationInvitation

	resp, err := s.client.doReqNoBody("GET", s.client.BaseURI+"/organizations/"+globalid+"/invitations", headers, queryParams)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Cancel a pending invitation.
func (s *OrganizationsService) RemovePendingOrganizationInvitation(username, globalid string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	// create request object
	return s.client.doReqNoBody("DELETE", s.client.BaseURI+"/organizations/"+globalid+"/invitations/"+username, headers, queryParams)
}

// Removes the Logo from an organization
func (s *OrganizationsService) DeleteOrganizationLogo(globalid string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	// create request object
	return s.client.doReqNoBody("DELETE", s.client.BaseURI+"/organizations/"+globalid+"/logo", headers, queryParams)
}

// Set the organization Logo for the organization
func (s *OrganizationsService) SetOrganizationLogo(globalid string, organizationsglobalidlogoputreqbody OrganizationsGlobalidLogoPutReqBody, headers, queryParams map[string]interface{}) (string, *http.Response, error) {
	var u string

	resp, err := s.client.doReqWithBody("PUT", s.client.BaseURI+"/organizations/"+globalid+"/logo", &organizationsglobalidlogoputreqbody, headers, queryParams)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Get the Logo from an organization
func (s *OrganizationsService) GetOrganizationLogo(globalid string, headers, queryParams map[string]interface{}) (string, *http.Response, error) {
	var u string

	resp, err := s.client.doReqNoBody("GET", s.client.BaseURI+"/organizations/"+globalid+"/logo", headers, queryParams)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Assign a member to organization.
func (s *OrganizationsService) AddOrganizationMember(globalid string, member Member, headers, queryParams map[string]interface{}) (Member, *http.Response, error) {
	var u Member

	resp, err := s.client.doReqWithBody("POST", s.client.BaseURI+"/organizations/"+globalid+"/members", &member, headers, queryParams)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Update an organization membership
func (s *OrganizationsService) UpdateOrganizationMemberShip(globalid string, organizationsglobalidmembersputreqbody OrganizationsGlobalidMembersPutReqBody, headers, queryParams map[string]interface{}) (Organization, *http.Response, error) {
	var u Organization

	resp, err := s.client.doReqWithBody("PUT", s.client.BaseURI+"/organizations/"+globalid+"/members", &organizationsglobalidmembersputreqbody, headers, queryParams)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Remove a member from an organization.
func (s *OrganizationsService) RemoveOrganizationMember(username, globalid string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	// create request object
	return s.client.doReqNoBody("DELETE", s.client.BaseURI+"/organizations/"+globalid+"/members/"+username, headers, queryParams)
}

// Invite a user to become owner of an organization.
func (s *OrganizationsService) AddOrganizationOwner(globalid string, member Member, headers, queryParams map[string]interface{}) (Member, *http.Response, error) {
	var u Member

	resp, err := s.client.doReqWithBody("POST", s.client.BaseURI+"/organizations/"+globalid+"/owners", &member, headers, queryParams)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Remove an owner from organization
func (s *OrganizationsService) RemoveOrganizationOwner(username, globalid string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	// create request object
	return s.client.doReqNoBody("DELETE", s.client.BaseURI+"/organizations/"+globalid+"/owners/"+username, headers, queryParams)
}

// Lists the RegistryEntries in an organization's registry.
func (s *OrganizationsService) ListOrganizationRegistry(globalid string, headers, queryParams map[string]interface{}) ([]RegistryEntry, *http.Response, error) {
	var u []RegistryEntry

	resp, err := s.client.doReqNoBody("GET", s.client.BaseURI+"/organizations/"+globalid+"/registry", headers, queryParams)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Adds a RegistryEntry to the organization's registry, if the key is already used, it is overwritten.
func (s *OrganizationsService) AddOrganizationRegistryEntry(globalid string, registryentry RegistryEntry, headers, queryParams map[string]interface{}) (RegistryEntry, *http.Response, error) {
	var u RegistryEntry

	resp, err := s.client.doReqWithBody("POST", s.client.BaseURI+"/organizations/"+globalid+"/registry", &registryentry, headers, queryParams)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Removes a RegistryEntry from the organization's registry
func (s *OrganizationsService) DeleteOrganizationRegistryEntry(key, globalid string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	// create request object
	return s.client.doReqNoBody("DELETE", s.client.BaseURI+"/organizations/"+globalid+"/registry/"+key, headers, queryParams)
}

// Get a RegistryEntry from the organization's registry.
func (s *OrganizationsService) GetOrganizationRegistryEntry(key, globalid string, headers, queryParams map[string]interface{}) (RegistryEntry, *http.Response, error) {
	var u RegistryEntry

	resp, err := s.client.doReqNoBody("GET", s.client.BaseURI+"/organizations/"+globalid+"/registry/"+key, headers, queryParams)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

func (s *OrganizationsService) GetOrganizationTree(globalid string, headers, queryParams map[string]interface{}) ([]OrganizationTreeItem, *http.Response, error) {
	var u []OrganizationTreeItem

	resp, err := s.client.doReqNoBody("GET", s.client.BaseURI+"/organizations/"+globalid+"/tree", headers, queryParams)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}
