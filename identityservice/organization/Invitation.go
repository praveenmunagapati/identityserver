package organization

type Invitation struct {
	Created Date   `json:"created"`
	Role    string `json:"role"`
	User    string `json:"user"`
}
