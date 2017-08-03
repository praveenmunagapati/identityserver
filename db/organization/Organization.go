package organization

import (
	"regexp"

	"github.com/itsyouonline/identityserver/db/user"
	"github.com/itsyouonline/identityserver/db/validation"

	"gopkg.in/validator.v2"
)

type Organization struct {
	DNS              []string        `json:"dns"`
	Globalid         string          `json:"globalid"`
	Members          []string        `json:"members"`
	Owners           []string        `json:"owners"`
	PublicKeys       []string        `json:"publicKeys"`
	SecondsValidity  int             `json:"secondsvalidity"`
	OrgOwners        []string        `json:"orgowners"`  //OrgOwners are other organizations that are owner of this organization
	OrgMembers       []string        `json:"orgmembers"` //OrgMembers are other organizations that are member of this organization
	RequiredScopes   []RequiredScope `json:"requiredscopes"`
	IncludeSubOrgsOf []string        `json:"includesuborgsof"`
}

// IsValid performs basic validation on the content of an organizations fields
func (org *Organization) IsValid() bool {
	regex, _ := regexp.Compile(`^[a-z\d\-_\s]{3,150}$`)
	return validator.Validate(org) == nil && regex.MatchString(org.Globalid)
}

func (org *Organization) IsValidSubOrganization() bool {
	regex, _ := regexp.Compile(`^[a-z\d\-_\s\.]{3,150}$`)
	return validator.Validate(org) == nil && regex.MatchString(org.Globalid)
}

func (org *Organization) ConvertToView(usrMgr *user.Manager, valMgr *validation.Manager) (*OrganizationView, error) {
	view := &OrganizationView{}
	view.DNS = org.DNS
	view.Globalid = org.Globalid
	view.PublicKeys = org.PublicKeys
	view.SecondsValidity = org.SecondsValidity
	view.OrgOwners = org.OrgOwners
	view.OrgMembers = org.OrgMembers
	view.RequiredScopes = org.RequiredScopes
	view.IncludeSubOrgsOf = org.IncludeSubOrgsOf

	var err error
	view.Members, err = convertUsernameToUserview(org.Members, usrMgr, valMgr)
	if err != nil {
		return view, err
	}
	view.Owners, err = convertUsernameToUserview(org.Owners, usrMgr, valMgr)

	return view, err
}

func convertUsernameToUserview(usernames []string, usrMgr *user.Manager, valMgr *validation.Manager) ([]MemberView, error) {
	views := []MemberView{}
	for _, username := range usernames {
		mv, err := ConvertUserToUserView(username, usrMgr, valMgr)
		if err != nil {
			return views, err
		}
		views = append(views, mv)
	}
	return views, nil
}

func ConvertUserToUserView(username string, usrMgr *user.Manager, valMgr *validation.Manager) (MemberView, error) {
	mv := MemberView{}
	mv.Username = username
	usr, err := usrMgr.GetByName(username)
	if err != nil {
		return mv, err
	}
	// check if real name is filled in
	if usr.Firstname != "" || usr.Lastname != "" {
		var usrId string
		if usr.Firstname != "" {
			usrId += usr.Firstname + " "
		}
		usrId += usr.Lastname
		mv.UserIdentifier = usrId
		return mv, err
	}
	// check for a validated email address
	for _, email := range usr.EmailAddresses {
		validated, err := valMgr.IsEmailAddressValidated(username, email.EmailAddress)
		if err != nil {
			return mv, err
		}
		if validated {
			mv.UserIdentifier = email.EmailAddress
			return mv, err
		}
	}
	// try the phone numbers
	for _, phone := range usr.Phonenumbers {
		validated, err := valMgr.IsPhonenumberValidated(username, phone.Phonenumber)
		if err != nil {
			return mv, err
		}
		if validated {
			mv.UserIdentifier = phone.Phonenumber
			return mv, err
		}
	}
	return mv, err
}

type OrganizationView struct {
	DNS              []string        `json:"dns"`
	Globalid         string          `json:"globalid"`
	Members          []MemberView    `json:"members"`
	Owners           []MemberView    `json:"owners"`
	PublicKeys       []string        `json:"publicKeys"`
	SecondsValidity  int             `json:"secondsvalidity"`
	OrgOwners        []string        `json:"orgowners"`  //OrgOwners are other organizations that are owner of this organization
	OrgMembers       []string        `json:"orgmembers"` //OrgMembers are other organizations that are member of this organization
	RequiredScopes   []RequiredScope `json:"requiredscopes"`
	IncludeSubOrgsOf []string        `json:"includesuborgsof"`
}

type MemberView struct {
	Username       string `json:"username"`
	UserIdentifier string `json:"useridentifier"`
}
