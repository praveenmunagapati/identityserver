package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	rootURL = "/"
)

type Itsyouonline struct {
	client http.Client
}

func NewItsyouonline() *Itsyouonline {
	c := new(Itsyouonline)
	c.client = http.Client{}
	return c
}

// Register a new company
func (c *Itsyouonline) CreateCompany(company Company, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)

	resp, err := doReqWithBody("POST", rootURL+"/companies", &company, c.client, headers, qsParam)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

// Get companies. Authorization limits are applied to requesting user.
func (c *Itsyouonline) GetCompanyList(headers, queryParams map[string]interface{}) ([]Company, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u []Company

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/companies"+qsParam, nil)
	if err != nil {
		return u, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Get organization info
func (c *Itsyouonline) GetCompany(globalId string, headers, queryParams map[string]interface{}) (Company, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u Company

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/companies/"+globalId+qsParam, nil)
	if err != nil {
		return u, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Update existing company. Updating ``globalId`` is not allowed.
func (c *Itsyouonline) UpdateCompany(globalId string, headers, queryParams map[string]interface{}) (Company, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u Company

	resp, err := doReqWithBody("PUT", rootURL+"/companies/"+globalId, nil, c.client, headers, qsParam)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Get the contracts where the organization is 1 of the parties. Order descending by date.
func (c *Itsyouonline) GetCompanyContracts(globalId string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/companies/"+globalId+"/contracts"+qsParam, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

// Create a new contract.
func (c *Itsyouonline) CreateCompanyContract(globalId string, contract Contract, headers, queryParams map[string]interface{}) (Contract, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u Contract

	resp, err := doReqWithBody("POST", rootURL+"/companies/"+globalId+"/contracts", &contract, c.client, headers, qsParam)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

func (c *Itsyouonline) GetCompanyInfo(globalId string, headers, queryParams map[string]interface{}) (companyview, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u companyview

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/companies/"+globalId+"/info"+qsParam, nil)
	if err != nil {
		return u, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

func (c *Itsyouonline) CompaniesGlobalIdValidateGet(globalId string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/companies/"+globalId+"/validate"+qsParam, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

// Get a contract
func (c *Itsyouonline) GetContract(contractId string, headers, queryParams map[string]interface{}) (Contract, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u Contract

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/contracts/"+contractId+qsParam, nil)
	if err != nil {
		return u, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Sign a contract
func (c *Itsyouonline) SignContract(contractId string, signature Signature, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)

	resp, err := doReqWithBody("POST", rootURL+"/contracts/"+contractId+"/signatures", &signature, c.client, headers, qsParam)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

// Create a new organization. 1 user should be in the owners list. Validation is performed to check if the securityScheme allows management on this user.
func (c *Itsyouonline) CreateNewOrganization(organization Organization, headers, queryParams map[string]interface{}) (Organization, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u Organization

	resp, err := doReqWithBody("POST", rootURL+"/organizations", &organization, c.client, headers, qsParam)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Create a new suborganization.
func (c *Itsyouonline) CreateNewSubOrganization(globalid string, organization Organization, headers, queryParams map[string]interface{}) (Organization, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u Organization

	resp, err := doReqWithBody("POST", rootURL+"/organizations/"+globalid, &organization, c.client, headers, qsParam)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Get organization info
func (c *Itsyouonline) GetOrganization(globalid string, headers, queryParams map[string]interface{}) (Organization, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u Organization

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/organizations/"+globalid+qsParam, nil)
	if err != nil {
		return u, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Update organization info
func (c *Itsyouonline) UpdateOrganization(globalid string, organization Organization, headers, queryParams map[string]interface{}) (Organization, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u Organization

	resp, err := doReqWithBody("PUT", rootURL+"/organizations/"+globalid, &organization, c.client, headers, qsParam)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Create a new API Key, a secret itself should not be provided, it will be generated serverside.
func (c *Itsyouonline) CreateNewOrganizationAPIKey(globalid string, organizationapikey OrganizationAPIKey, headers, queryParams map[string]interface{}) (OrganizationAPIKey, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u OrganizationAPIKey

	resp, err := doReqWithBody("POST", rootURL+"/organizations/"+globalid+"/apikeys", &organizationapikey, c.client, headers, qsParam)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Get the list of active api keys.
func (c *Itsyouonline) GetOrganizationAPIKeyLabels(globalid string, headers, queryParams map[string]interface{}) ([]string, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u []string

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/organizations/"+globalid+"/apikeys"+qsParam, nil)
	if err != nil {
		return u, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Updates the label or other properties of a key.
func (c *Itsyouonline) UpdateOrganizationAPIKey(label, globalid string, organizationsglobalidapikeyslabelputreqbody OrganizationsGlobalidApikeysLabelPutReqBody, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)

	resp, err := doReqWithBody("PUT", rootURL+"/organizations/"+globalid+"/apikeys/"+label, &organizationsglobalidapikeyslabelputreqbody, c.client, headers, qsParam)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

func (c *Itsyouonline) GetOrganizationAPIKey(label, globalid string, headers, queryParams map[string]interface{}) (OrganizationAPIKey, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u OrganizationAPIKey

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/organizations/"+globalid+"/apikeys/"+label+qsParam, nil)
	if err != nil {
		return u, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Removes an API key
func (c *Itsyouonline) DeleteOrganizationAPIKey(label, globalid string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)
	// create request object
	req, err := http.NewRequest("DELETE", rootURL+"/organizations/"+globalid+"/apikeys/"+label+qsParam, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	return c.client.Do(req)
}

// Get the contracts where the organization is 1 of the parties. Order descending by date.
func (c *Itsyouonline) GetOrganizationContracts(globalid string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/organizations/"+globalid+"/contracts"+qsParam, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

// Create a new contract.
func (c *Itsyouonline) CreateOrganizationContracty(globalid string, contract Contract, headers, queryParams map[string]interface{}) (Contract, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u Contract

	resp, err := doReqWithBody("POST", rootURL+"/organizations/"+globalid+"/contracts", &contract, c.client, headers, qsParam)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Creates a new DNS name associated with an organization
func (c *Itsyouonline) CreateOrganizationDNS(dnsname, globalid string, organizationsglobaliddnsdnsnamepostreqbody OrganizationsGlobalidDnsDnsnamePostReqBody, headers, queryParams map[string]interface{}) (OrganizationsGlobalidDnsDnsnamePostRespBody, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u OrganizationsGlobalidDnsDnsnamePostRespBody

	resp, err := doReqWithBody("POST", rootURL+"/organizations/"+globalid+"/dns/"+dnsname, &organizationsglobaliddnsdnsnamepostreqbody, c.client, headers, qsParam)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Removes a DNS name
func (c *Itsyouonline) DeleteOrganizaitonDNS(dnsname, globalid string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)
	// create request object
	req, err := http.NewRequest("DELETE", rootURL+"/organizations/"+globalid+"/dns/"+dnsname+qsParam, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	return c.client.Do(req)
}

// Updates an existing DNS name associated with an organization
func (c *Itsyouonline) UpdateOrganizationDNS(dnsname, globalid string, organizationsglobaliddnsdnsnameputreqbody OrganizationsGlobalidDnsDnsnamePutReqBody, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)

	resp, err := doReqWithBody("PUT", rootURL+"/organizations/"+globalid+"/dns/"+dnsname, &organizationsglobaliddnsdnsnameputreqbody, c.client, headers, qsParam)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

// Get the list of pending invitations for users to join this organization.
func (c *Itsyouonline) GetPendingOrganizationInvitations(globalid string, headers, queryParams map[string]interface{}) ([]JoinOrganizationInvitation, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u []JoinOrganizationInvitation

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/organizations/"+globalid+"/invitations"+qsParam, nil)
	if err != nil {
		return u, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Cancel a pending invitation.
func (c *Itsyouonline) RemovePendingOrganizationInvitation(username, globalid string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)
	// create request object
	req, err := http.NewRequest("DELETE", rootURL+"/organizations/"+globalid+"/invitations/"+username+qsParam, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	return c.client.Do(req)
}

// Assign a member to organization.
func (c *Itsyouonline) AddOrganizationMember(globalid string, member Member, headers, queryParams map[string]interface{}) (Member, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u Member

	resp, err := doReqWithBody("POST", rootURL+"/organizations/"+globalid+"/members", &member, c.client, headers, qsParam)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Remove a member from organization
func (c *Itsyouonline) RemoveOrganizationMember(username, globalid string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)
	// create request object
	req, err := http.NewRequest("DELETE", rootURL+"/organizations/"+globalid+"/members/"+username+qsParam, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	return c.client.Do(req)
}

// Invite a user to become owner of an organization.
func (c *Itsyouonline) AddOrganizationOwner(globalid string, member Member, headers, queryParams map[string]interface{}) (Member, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u Member

	resp, err := doReqWithBody("POST", rootURL+"/organizations/"+globalid+"/owners", &member, c.client, headers, qsParam)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Remove an owner from organization
func (c *Itsyouonline) RemoveOrganizationOwner(username, globalid string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)
	// create request object
	req, err := http.NewRequest("DELETE", rootURL+"/organizations/"+globalid+"/owners/"+username+qsParam, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	return c.client.Do(req)
}

func (c *Itsyouonline) GetOrganizationTree(globalid string, headers, queryParams map[string]interface{}) ([]OrganizationTreeItem, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u []OrganizationTreeItem

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/organizations/"+globalid+"/tree"+qsParam, nil)
	if err != nil {
		return u, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Create a new user
func (c *Itsyouonline) CreateUser(user User, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)

	resp, err := doReqWithBody("POST", rootURL+"/users", &user, c.client, headers, qsParam)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

func (c *Itsyouonline) GetUser(username string, headers, queryParams map[string]interface{}) (User, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u User

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/users/"+username+qsParam, nil)
	if err != nil {
		return u, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Register a new address
func (c *Itsyouonline) RegisterNewUserAddress(username string, usersusernameaddressespostreqbody UsersUsernameAddressesPostReqBody, headers, queryParams map[string]interface{}) (UsersUsernameAddressesPostRespBody, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u UsersUsernameAddressesPostRespBody

	resp, err := doReqWithBody("POST", rootURL+"/users/"+username+"/addresses", &usersusernameaddressespostreqbody, c.client, headers, qsParam)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

func (c *Itsyouonline) GetUserAddresses(username string, headers, queryParams map[string]interface{}) (UsersUsernameAddressesGetRespBody, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u UsersUsernameAddressesGetRespBody

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/users/"+username+"/addresses"+qsParam, nil)
	if err != nil {
		return u, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Removes an address
func (c *Itsyouonline) DeleteUserAddress(label, username string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)
	// create request object
	req, err := http.NewRequest("DELETE", rootURL+"/users/"+username+"/addresses/"+label+qsParam, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	return c.client.Do(req)
}

func (c *Itsyouonline) GetUserAddressByLabel(label, username string, headers, queryParams map[string]interface{}) (Address, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u Address

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/users/"+username+"/addresses/"+label+qsParam, nil)
	if err != nil {
		return u, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Update the label and/or value of an existing address.
func (c *Itsyouonline) UpdateUserAddress(label, username string, usersusernameaddresseslabelputreqbody UsersUsernameAddressesLabelPutReqBody, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)

	resp, err := doReqWithBody("PUT", rootURL+"/users/"+username+"/addresses/"+label, &usersusernameaddresseslabelputreqbody, c.client, headers, qsParam)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

// Adds an APIKey to the user
func (c *Itsyouonline) AddApiKey(username string, usersusernameapikeyspostreqbody UsersUsernameApikeysPostReqBody, headers, queryParams map[string]interface{}) (UserAPIKey, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u UserAPIKey

	resp, err := doReqWithBody("POST", rootURL+"/users/"+username+"/apikeys", &usersusernameapikeyspostreqbody, c.client, headers, qsParam)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Lists the API keys
func (c *Itsyouonline) ListAPIKeys(username string, headers, queryParams map[string]interface{}) ([]UserAPIKey, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u []UserAPIKey

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/users/"+username+"/apikeys"+qsParam, nil)
	if err != nil {
		return u, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Updates the label for the api key
func (c *Itsyouonline) UpdateAPIkey(label, username string, usersusernameapikeyslabelputreqbody UsersUsernameApikeysLabelPutReqBody, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)

	resp, err := doReqWithBody("PUT", rootURL+"/users/"+username+"/apikeys/"+label, &usersusernameapikeyslabelputreqbody, c.client, headers, qsParam)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

// Removes an API key
func (c *Itsyouonline) DeleteAPIkey(label, username string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)
	// create request object
	req, err := http.NewRequest("DELETE", rootURL+"/users/"+username+"/apikeys/"+label+qsParam, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	return c.client.Do(req)
}

// Get an API key by label
func (c *Itsyouonline) GetAPIkey(label, username string, headers, queryParams map[string]interface{}) (UserAPIKey, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u UserAPIKey

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/users/"+username+"/apikeys/"+label+qsParam, nil)
	if err != nil {
		return u, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Get the list of authorizations.
func (c *Itsyouonline) GetAllAuthorizations(username string, headers, queryParams map[string]interface{}) ([]Authorization, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u []Authorization

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/users/"+username+"/authorizations"+qsParam, nil)
	if err != nil {
		return u, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Remove the authorization for an organization, the granted organization will no longer have access the user's information.
func (c *Itsyouonline) DeleteAuthorization(grantedTo, username string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)
	// create request object
	req, err := http.NewRequest("DELETE", rootURL+"/users/"+username+"/authorizations/"+grantedTo+qsParam, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	return c.client.Do(req)
}

// Modify which information an organization is able to see.
func (c *Itsyouonline) UpdateAuthorization(grantedTo, username string, authorization Authorization, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)

	resp, err := doReqWithBody("PUT", rootURL+"/users/"+username+"/authorizations/"+grantedTo, &authorization, c.client, headers, qsParam)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

// Get the authorization for a specific organization.
func (c *Itsyouonline) GetAuthorization(grantedTo, username string, headers, queryParams map[string]interface{}) (Authorization, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u Authorization

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/users/"+username+"/authorizations/"+grantedTo+qsParam, nil)
	if err != nil {
		return u, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

func (c *Itsyouonline) GetUserBankAccounts(username string, headers, queryParams map[string]interface{}) (UsersUsernameBanksGetRespBody, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u UsersUsernameBanksGetRespBody

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/users/"+username+"/banks"+qsParam, nil)
	if err != nil {
		return u, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Create new bank account
func (c *Itsyouonline) CreateUserBankAccount(username string, usersusernamebankspostreqbody UsersUsernameBanksPostReqBody, headers, queryParams map[string]interface{}) (UsersUsernameBanksPostRespBody, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u UsersUsernameBanksPostRespBody

	resp, err := doReqWithBody("POST", rootURL+"/users/"+username+"/banks", &usersusernamebankspostreqbody, c.client, headers, qsParam)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Update an existing bankaccount and label.
func (c *Itsyouonline) UpdateUserBankAccount(username, label string, bankaccount BankAccount, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)

	resp, err := doReqWithBody("PUT", rootURL+"/users/"+username+"/banks/"+label, &bankaccount, c.client, headers, qsParam)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

// Delete a BankAccount
func (c *Itsyouonline) UsersUsernameBanksLabelDelete(username, label string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)
	// create request object
	req, err := http.NewRequest("DELETE", rootURL+"/users/"+username+"/banks/"+label+qsParam, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	return c.client.Do(req)
}

func (c *Itsyouonline) GetUserBankAccountByLabel(username, label string, headers, queryParams map[string]interface{}) (BankAccount, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u BankAccount

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/users/"+username+"/banks/"+label+qsParam, nil)
	if err != nil {
		return u, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Get the contracts where the user is 1 of the parties. Order descending by date.
func (c *Itsyouonline) GetUserContracts(username string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/users/"+username+"/contracts"+qsParam, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

// Create a new contract.
func (c *Itsyouonline) CreateUserContract(username string, contract Contract, headers, queryParams map[string]interface{}) (Contract, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u Contract

	resp, err := doReqWithBody("POST", rootURL+"/users/"+username+"/contracts", &contract, c.client, headers, qsParam)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Get a list of the user his email addresses.
func (c *Itsyouonline) GetEmailAddresses(username string, headers, queryParams map[string]interface{}) ([]EmailAddress, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u []EmailAddress

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/users/"+username+"/emailaddresses"+qsParam, nil)
	if err != nil {
		return u, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Register a new email address
func (c *Itsyouonline) RegisterNewEmailAddress(username string, usersusernameemailaddressespostreqbody UsersUsernameEmailaddressesPostReqBody, headers, queryParams map[string]interface{}) (UsersUsernameEmailaddressesPostRespBody, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u UsersUsernameEmailaddressesPostRespBody

	resp, err := doReqWithBody("POST", rootURL+"/users/"+username+"/emailaddresses", &usersusernameemailaddressespostreqbody, c.client, headers, qsParam)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Removes an email address
func (c *Itsyouonline) DeleteEmailAddress(label, username string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)
	// create request object
	req, err := http.NewRequest("DELETE", rootURL+"/users/"+username+"/emailaddresses/"+label+qsParam, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	return c.client.Do(req)
}

// Updates the label and/or value of an email address
func (c *Itsyouonline) UpdateEmailAddress(label, username string, usersusernameemailaddresseslabelputreqbody UsersUsernameEmailaddressesLabelPutReqBody, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)

	resp, err := doReqWithBody("PUT", rootURL+"/users/"+username+"/emailaddresses/"+label, &usersusernameemailaddresseslabelputreqbody, c.client, headers, qsParam)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

// Sends validation email to email address
func (c *Itsyouonline) ValidateEmailAddress(label, username string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)

	resp, err := doReqWithBody("POST", rootURL+"/users/"+username+"/emailaddresses/"+label+"/validate", nil, c.client, headers, qsParam)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

// Delete the associated facebook account
func (c *Itsyouonline) DeleteFacebookAccount(username string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)
	// create request object
	req, err := http.NewRequest("DELETE", rootURL+"/users/"+username+"/facebook"+qsParam, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	return c.client.Do(req)
}

// Unlink Github Account
func (c *Itsyouonline) DeleteGithubAccount(username string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)
	// create request object
	req, err := http.NewRequest("DELETE", rootURL+"/users/"+username+"/github"+qsParam, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	return c.client.Do(req)
}

func (c *Itsyouonline) GetUserInformation(username string, headers, queryParams map[string]interface{}) (userview, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u userview

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/users/"+username+"/info"+qsParam, nil)
	if err != nil {
		return u, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Update the user his firstname and lastname
func (c *Itsyouonline) UpdateUserName(username string, usersusernamenameputreqbody UsersUsernameNamePutReqBody, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)

	resp, err := doReqWithBody("PUT", rootURL+"/users/"+username+"/name", &usersusernamenameputreqbody, c.client, headers, qsParam)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

// Get the list of notifications, these are pending invitations or approvals
func (c *Itsyouonline) GetNotifications(username string, headers, queryParams map[string]interface{}) (UsersUsernameNotificationsGetRespBody, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u UsersUsernameNotificationsGetRespBody

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/users/"+username+"/notifications"+qsParam, nil)
	if err != nil {
		return u, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Get the list organizations a user is owner or member of
func (c *Itsyouonline) GetUserOrganizations(username string, headers, queryParams map[string]interface{}) (UsersUsernameOrganizationsGetRespBody, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u UsersUsernameOrganizationsGetRespBody

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/users/"+username+"/organizations"+qsParam, nil)
	if err != nil {
		return u, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Accept membership in organization
func (c *Itsyouonline) AcceptMembership(globalid, role, username string, joinorganizationinvitation JoinOrganizationInvitation, headers, queryParams map[string]interface{}) (JoinOrganizationInvitation, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u JoinOrganizationInvitation

	resp, err := doReqWithBody("POST", rootURL+"/users/"+username+"/organizations/"+globalid+"/roles/"+role, &joinorganizationinvitation, c.client, headers, qsParam)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Reject membership invitation in an organization.
func (c *Itsyouonline) UsersUsernameOrganizationsGlobalidRolesRoleDelete(globalid, role, username string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)
	// create request object
	req, err := http.NewRequest("DELETE", rootURL+"/users/"+username+"/organizations/"+globalid+"/roles/"+role+qsParam, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	return c.client.Do(req)
}

// Update the user his password
func (c *Itsyouonline) UpdatePassword(username string, usersusernamepasswordputreqbody UsersUsernamePasswordPutReqBody, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)

	resp, err := doReqWithBody("PUT", rootURL+"/users/"+username+"/password", &usersusernamepasswordputreqbody, c.client, headers, qsParam)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

// Register a new phonenumber
func (c *Itsyouonline) RegisterNewUserPhonenumber(username string, usersusernamephonenumberspostreqbody UsersUsernamePhonenumbersPostReqBody, headers, queryParams map[string]interface{}) (Phonenumber, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u Phonenumber

	resp, err := doReqWithBody("POST", rootURL+"/users/"+username+"/phonenumbers", &usersusernamephonenumberspostreqbody, c.client, headers, qsParam)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

func (c *Itsyouonline) GetUserPhoneNumbers(username string, headers, queryParams map[string]interface{}) (UsersUsernamePhonenumbersGetRespBody, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u UsersUsernamePhonenumbersGetRespBody

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/users/"+username+"/phonenumbers"+qsParam, nil)
	if err != nil {
		return u, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

func (c *Itsyouonline) GetUserPhonenumberByLabel(label, username string, headers, queryParams map[string]interface{}) (Phonenumber, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u Phonenumber

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/users/"+username+"/phonenumbers/"+label+qsParam, nil)
	if err != nil {
		return u, nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

// Removes a phonenumber
func (c *Itsyouonline) DeleteUserPhonenumber(label, username string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)
	// create request object
	req, err := http.NewRequest("DELETE", rootURL+"/users/"+username+"/phonenumbers/"+label+qsParam, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	return c.client.Do(req)
}

// Update the label and/or value of an existing phonenumber.
func (c *Itsyouonline) UpdateUserPhonenumber(label, username string, usersusernamephonenumberslabelputreqbody UsersUsernamePhonenumbersLabelPutReqBody, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)

	resp, err := doReqWithBody("PUT", rootURL+"/users/"+username+"/phonenumbers/"+label, &usersusernamephonenumberslabelputreqbody, c.client, headers, qsParam)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

// Verifies a phone number
func (c *Itsyouonline) VerifyPhoneNumber(label, username string, usersusernamephonenumberslabelactivateputreqbody UsersUsernamePhonenumbersLabelActivatePutReqBody, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)

	resp, err := doReqWithBody("PUT", rootURL+"/users/"+username+"/phonenumbers/"+label+"/activate", &usersusernamephonenumberslabelactivateputreqbody, c.client, headers, qsParam)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

// Sends validation text to phone numbers
func (c *Itsyouonline) ValidatePhonenumber(label, username string, headers, queryParams map[string]interface{}) (UsersUsernamePhonenumbersLabelActivatePostRespBody, *http.Response, error) {
	qsParam := buildQueryString(queryParams)
	var u UsersUsernamePhonenumbersLabelActivatePostRespBody

	resp, err := doReqWithBody("POST", rootURL+"/users/"+username+"/phonenumbers/"+label+"/activate", nil, c.client, headers, qsParam)
	if err != nil {
		return u, nil, err
	}
	defer resp.Body.Close()

	return u, resp, json.NewDecoder(resp.Body).Decode(&u)
}

func (c *Itsyouonline) UsersUsernameValidateGet(username string, headers, queryParams map[string]interface{}) (*http.Response, error) {
	qsParam := buildQueryString(queryParams)

	// create request object
	req, err := http.NewRequest("GET", rootURL+"/users/"+username+"/validate"+qsParam, nil)
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, fmt.Sprintf("%v", v))
	}

	//do the request
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}
