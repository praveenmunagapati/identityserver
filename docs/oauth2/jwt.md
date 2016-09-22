# JWT (JSON Web Token) support

Even though the OAuth token support works great for applications that need to access the information of a user, when passing on some of these authorizations to a third party service it is not a good idea to pass on your token itself.

The token you acquired might give access to a lot more information that you want to pass on to the third party service and it is required to invoke itsyou.online to verify that the authorization claim is valid.

For these use cases, itsyou.online supports JWT [RFC7519](https://tools.ietf.org/html/rfc7519).

Itsyou.online supports two way of obtaining JWTs:
1. Use an OAuth2 token for JWT creation where the JWT's claim set is a subset of the OAuth token's scopes.
2. Directly get a JWT instead of a normal OAuth token when following the OAut2 grant type flows.

## Case 1: Use an OAuth2 token for JWT creation where the JWT's claim set is a subset of the OAuth token's scopes

Suppose you have an OAuth token OAUTH-TOKEN with the following scopes:

- user:memberOf:org1
- user:memberOf:org2
- user:address:billing

and you want to call a third party service that only needs to know if the user is member of org1, there is no need to expose the billing address you are authorized to see.

You can create a JWT like this:
```
curl -H "Authorization: token OAUTH-TOKEN" https://itsyou.online/v1/oauth/jwt?scope=user:memberof:org1
```

The `scope` parameter can be a comma separated list of scopes. Instead of a query parameter, an http `POST` can also be submitted to this url with the scope parameter as a form value.

The response will be a JWT with:
* Header

    ```
    {
      "alg": "ES384",
      "typ": "JWT"
    }
    ```

* Data

    ```
    {
      "username": "bob",
      "scope": "user:memberof:org1",
      "iss": "itsyouonline",
      "aud": ["CLIENTID"],
      "exp": 1463554314
    }
    ```

    - iss: Issuer, in this case "itsyouonline"
    - exp: Expiration time in seconds since the epoch. This is set to the same time as the expiration time of the OAuth token used to acquire this JWT.
    - aud: An array with at least 1 entry: the `client_id` of the OAuth token used to acquire this JWT.

    If the OAuth token is not for a user but for an organization application that authenticated using the client credentials flow, the `username` field is replaced with a `globalid` field containing the globalid of the organization.

* Signature

    The JWT is signed by itsyou.online. The public key to verify if this JWT was really issued by itsyou.online is
    ```
    -----BEGIN PUBLIC KEY-----
    MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAES5X8XrfKdx9gYayFITc89wad4usrk0n2
    7MjiGYvqalizeSWTHEpnd7oea9IQ8T5oJjMVH5cc0H5tFSKilFFeh//wngxIyny6
    6+Vq5t5B0V0Ehy01+2ceEon2Y0XDkIKv
    -----END PUBLIC KEY-----
    ```

In case the requested scopes are not available for your OAuth token or the token has expired, an http 401 status code is returned.


### Creating JWT's for other audiences

In case you want to pass on a JWT and your authorization to a different audience, you can specify extra audiences to the call for creating a JWT.

```
curl -H "Authorization: token OAUTH-TOKEN" https://itsyou.online/v1/oauth/jwt??scope=user:memberof:org1&aud=external1,external2
```

In this case, this results in the following JWT data

    ```
    {
      "username": "bob",
      "scope": "user:memberOf:org1",
      "iss": "itsyouonline",
      "aud": [
            "CLIENTID",
            "external1",
            "external2"
            ]
      "exp": 1463554314
    }
    ```

The audience field is a list of audiences, the first audience is always the `client_id` of the OAuth token used to acquire this JWT followed by the audiences passed in the request. The extra audiences are not required to be valid globalid's of organizations in itsyou.online.


## Case 2: Directly get a JWT instead of a normal oauth2 token when following the oauth2 grant type flows

When using 1 of the authorization flows explained in the [Authorization grant types](oauth2.md) documentation, it is also possible to directly get a JWT returned instead of an OAuth2 token itself.
Add the `return_type=id_token` and a `scope` parameter with the desired scopes to the `/v1/oauth/access_token` call to do this.
For example:
```
https://itsyou.online/v1/oauth/access_token?grant_type=client_credentials&client_id=CLIENT_ID&client_secret=CLIENT_SECRET&reponse_type=id_token&scope=user:memberof:org1&aud=external1
```

In this case, the scope parameter needs to be given to prevent consumers to accidentally handing out `user:admin` or `organization:owner` scoped tokens to third party services

As shown in the example. it is also possible to specify additional audiences in the `/v1/oauth/access_token` call.
