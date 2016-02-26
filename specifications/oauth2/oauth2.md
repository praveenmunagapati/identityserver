## Authorization grant types

OAuth2 defines four grant types, each of which is useful in different cases:

1. Authorization Code: used with server-side Applications
2. Implicit: used with Mobile Apps or Web Applications (applications that run on the user's device)
3. Resource Owner Password Credentials: used with trusted Applications, such as those owned by the service itself
4. Client Credentials: used with Applications API access

Currently only the authorization code grant type is supported.

**Grant Type: Authorization Code**
The authorization code grant type is the most commonly used because it is optimized for server-side applications where Client Secret confidentiality can be maintained. This is a redirection-based flow, which means that the application must be capable of interacting with the user-agent (i.e. the user's web browser) and receiving API authorization codes that are routed through the user-agent.

## Authorization Code Flow

### Step 1: Authorization Code Link

First, the user is given an authorization code link that looks like the following:

```
https://itsyou.online/v1/oauth/authorize?response_type=code&client_id=CLIENT_ID&redirect_uri=CALLBACK_URL&scope=read&state=STATE
```

* https://itsyou.online/v1/oauth/authorize: the API authorization endpoint
* client_id=client_id

    the application's client ID
* redirect_uri=CALLBACK_URL

    The redirect_uri parameter is optional. If left out, the users will redirected to the callback URL configured in the OAuth Application settings. If provided, the redirect URL's host and port must exactly match the callback URL. The redirect URL's path must reference a subdirectory of the callback URL.
* response_type=code

    specifies that your application is requesting an authorization code grant
* scope=read

    specifies the level of access that the application is requesting

* state=STATE

    A random string. It is used to protect against csrf attacks.

### Step 2: User Authorizes Application

When the user clicks the link, they must first log in to the service, to authenticate their identity (unless they are already logged in). Then they will be prompted by the service to authorize or deny the application access to the requested information.

### Step 3: Application Receives Authorization Code

After the the user authorizes the application some of it's information, itsyou.online redirects the user-agent to the application redirect URI, which was specified during the client registration, along with an authorization code and a state parameter passed in step 1. If the state parameters don't match, the reqeust has been created by a third party and the process should be aborted.
The redirect would look something like this (assuming the application is "petshop.com"):

```
https://petshop.com/callback?code=AUTHORIZATION_CODE&state=STATE
```

### Step 4: Application Requests Access Token

The application requests an access token from the API, by passing the authorization code along with authentication details, including the client secret, to the API token endpoint. Here is an example POST request to itsyou.online's token endpoint:

```
https://itsyou.online/v1/oauth/token?client_id=CLIENT_ID&client_secret=CLIENT_SECRET&grant_type=authorization_code&code=AUTHORIZATION_CODE&redirect_uri=CALLBACK_URL
```

### Step 5: Application Receives Access Token

If the authorization is valid, the API will send a response containing the access token (and optionally, a refresh token) to the application. The entire response will look something like this:

```
{"access_token":"ACCESS_TOKEN","token_type":"bearer","expires_in":2592000,"refresh_token":"REFRESH_TOKEN","scope":"read","uid":100101,"info":{"username":"bob"}}
```
Now the application is authorized! It may use the token to access the user's account via the service API, limited to the scope of access, until the token expires or is revoked. If a refresh token was issued, it may be used to request new access tokens if the original token has expired.
