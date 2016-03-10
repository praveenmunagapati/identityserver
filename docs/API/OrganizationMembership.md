# Organization membership

## Invite members

In order to invite a member, an API call can be done:

```
POST /organizations/<org>/members
{
    "username": "bob"    
}

response: 201 CREATED
{
    "organization": "mysoccerclub.com",
    "role": [
        "member"
    ],
    "user": "bob",
    "status": "pending"
}
```

## User notifications

A user can view different types of notifications, we are interested here in Join Organization requests (i.e `invitations`).

```
GET /users/bob/notifications
{
    "approvals": [],
    "contractRequests": [],
    "invitations": [
        {
            "organization": "mysoccerclub.com",
            "role": [
                "member"
            ],
            "user": "bob",
            "status": "pending"
        }
    ]
}
```


## Accepting invitation

A user can accept invitation via an API call

```
POST /users/bob/organizations/mysoccerclub.com/roles/member
{
    "organization": "mysoccerclub.com",
    "role": [
        "member"
    ],
    "user": "bob",
    "status": "pending"
}

response: 201 CREATED
{
    "organization": "mysoccerclub.com",
    "role": [
        "member"
    ],
    "user": "bob",
    "status": "accepted"
}
```

## Rejecting invitation

A user can reject an invitation via an API call

```
DELETE /users/bob/organizations/mysoccerclub.com/roles/member

response: 204 NO CONTENT
```
