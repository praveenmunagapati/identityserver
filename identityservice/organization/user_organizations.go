package organization

type UserOrganizations struct {
	Member []string `json:"member"`
	Owner  []string `json:"owner"`
}
