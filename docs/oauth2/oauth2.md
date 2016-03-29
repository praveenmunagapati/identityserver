## Authorization grant types

OAuth2 defines four grant types, each of which is useful in different cases:

1. Authorization Code: used with server-side Applications
2. Implicit: used with Mobile Apps or Web Applications (applications that run on the user's device)
3. Resource Owner Password Credentials: used with trusted Applications, such as those owned by the service itself
4. Client Credentials: used with Applications API access

Currently the **authorization code** and **implicit** grant types are supported.


## Authorization Code Flow
The authorization code grant type is the most commonly used because it is optimized for server-side applications where Client Secret confidentiality can be maintained. This is a redirection-based flow, which means that the application must be capable of interacting with the user-agent (i.e. the user's web browser) and receiving API authorization codes that are routed through the user-agent.

### Prerequisite: clientid and client secret

In order to acquire an oauth access token, a client id and client secret are required.

In itsyou.online, organizations map to clients in the oauth2 terminology and the organization's globalid is used as the clientid. Client secrets can be created through the UI or through the `organizations/{globalid}/apisecrets` api.

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
POST https://itsyou.online/v1/oauth/access_token?client_id=CLIENT_ID&client_secret=CLIENT_SECRET&code=AUTHORIZATION_CODE&redirect_uri=CALLBACK_URL
```

### Step 5: Application Receives Access Token

If the authorization is valid, the API will send a response containing the access token (and optionally, a refresh token) to the application. The entire response will look something like this:

```
{"access_token":"ACCESS_TOKEN","token_type":"bearer","expires_in":2592000,"refresh_token":"REFRESH_TOKEN","scope":"read","info":{"username":"bob"}}
```
Now the application is authorized! It may use the token to access the user's account via the service API, limited to the scope of access, until the token expires or is revoked. If a refresh token was issued, it may be used to request new access tokens if the original token has expired.


### Use the access token to access the API

The access token allows you to make requests to the API on a behalf of a user.

```
GET https://itsyou.online/users/bob/info?access_token=...
```
You can pass the token in the query params like shown above, but a cleaner approach is to include it in the Authorization header

```
Authorization: token OAUTH-TOKEN
```
For example, in curl you can set the Authorization header like this:

```
curl -H "Authorization: token OAUTH-TOKEN" https://itsyou.online/users/bob/info
```

## Implicit flow
The implicit grant type is used for mobile apps and web applications (i.e. applications that run in a web browser), where the client secret confidentiality is not guaranteed. The implicit grant type is also a redirection-based flow but the access token is given to the user-agent to forward to the application, so it may be exposed to the user and other applications on the user's device. Also, this flow does not authenticate the identity of the application, and relies on the redirect URI (that was registered with the service) to serve this purpose.

The implicit grant type does not support refresh tokens.

### Step 1: Implicit Authorization Link

With the implicit grant type, the user is presented with an authorization link, that requests a token from the API. This link looks just like the authorization code link, except it is requesting a token instead of a code (note the response type "token"):

First, the user is given an authorization code link that looks like the following:

```
https://itsyou.online/v1/oauth/authorize?response_type=token&client_id=CLIENT_ID&redirect_uri=CALLBACK_URL&scope=read
```

### Step 2: User Authorizes Application

When the user clicks the link, they must first log in to itsyou.online, to authenticate their identity (unless they are already logged in). Then they will be prompted by to authorize or deny the application access to their account.

### Step 3: User-agent Receives Access Token with Redirect URI

When the user authorizes the application, itsyou.online redirects the user-agent to the application redirect URI, and includes a URI fragment containing the access token. It would look something like this:
```
https://petshop.com/callback#token=ACCESS_TOKEN
```

### Step 4: User-agent Follows the Redirect URI

The user-agent follows the redirect URI but retains the access token (notice the `#` in the url).

### Step 5: Application Sends Access Token Extraction Script

The application returns a webpage that contains a script that can extract the access token from the full redirect URI that the user-agent has retained.

### Step 6: Access Token Passed to Application

The user-agent executes the provided script and passes the extracted access token to the application.

Now the application is authorized! It may use the token to access the user's account via the itsyou.online API, limited to the scope of access, until the token expires or is revoked.
