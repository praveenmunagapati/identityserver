package organization

type OrganizationUser struct {
	User          MemberView `json:"user"`
	Role          string     `json:"role"`
	MissingScopes []string   `json:"missingscopes"`
}
