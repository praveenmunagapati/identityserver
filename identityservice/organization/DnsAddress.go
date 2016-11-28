package organization

type DnsAddress struct {
	Name string `json:"name,omitempty" validate:"min=4,max=250,nonzero"`
}
