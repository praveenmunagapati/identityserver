# JWT (JSON Web Token) support

Even though the oauth token support works great for applications that need to access the information of a user, when passing on some of these authorizations to a third party service it is not a good idea to pass on your token itself.

The token you acquired might give access to a lot more information that you want to pass on to the third party service and it is required to invoke itsyou.online to verify that the authorization claim is valid.

For these use cases, itsyou.online supports JWT [RFC7519](https://tools.ietf.org/html/rfc7519) creation where the JWT's claimset is a subset of the oauth token's scopes.

Suppose you have an oauth token OAUTH_TOKEN with the following scopes:

- user:memberOf:org1
- user:memberOf:org2
- user:address:billing

and you want to call a third party service that only needs to know if the user is member of org1, there is no need to expose the billing address you are authorized to see.

You can create a JWT like this:
```
curl -H "Authorization: token OAUTH-TOKEN" https://itsyou.online/v1/oath/jwt?scope=user:memberOf:org1
```

The `scope` parameter can be a comma seperated list of scopes. Instead of a query parameter, an http `POST` can also be submitted to this url with the scope parameter as a form value.

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
      "scope": "user:memberOf:org1",
      "iss": "itsyouonline",
      "aud": "CLIENTID",
      "exp": 1463554314,
    }
    ```

    - iss: Issuer, in this case "itsyouonline"
    - exp: Expiration time in seconds since the epoch. This is set to the same time as the expiration time of the oauth token used to acquire this JWT.
    - aud: The `client_id` of the oauth token used to acquire this JWT

* Signature

    The JWT is signed by itsyou.online. The public key to verify if this JWT was really issued by itsyou.online is
    ```
    AAAAAA
    ```

In case the requested scopes are not available for your oauth token or the token has expired, an http 401 status code is returned.
