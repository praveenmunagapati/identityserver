package organization

import "github.com/itsyouonline/identityserver/oauthservice"

type APIKey struct {
	CallbackURL                string `json:"callbackURL,omitempty" validate:"min=5,max=250,nonzero"`
	ClientCredentialsGrantType bool   `json:"clientCredentialsGrantType,omitempty" validate:"nonzero"`
	Label                      string `json:"label" validate:"min=2,max=50"`
	Secret                     string `json:"secret,omitempty" validate:"max=250,nonzero"`
}

//FromOAuthClient creates an APIKey instance from an oauthservice.Oauth2Client
func FromOAuthClient(client *oauthservice.Oauth2Client) APIKey {
	apiKey := APIKey{
		CallbackURL:                client.CallbackURL,
		ClientCredentialsGrantType: client.ClientCredentialsGrantType,
		Label:  client.Label,
		Secret: client.Secret,
	}
	return apiKey
}
