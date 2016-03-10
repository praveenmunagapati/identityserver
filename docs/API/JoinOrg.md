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
    "id": "56e168de47ee2036ac81e3a7",
    "organization": "greenitglobe.com",
    "role": [
        "member"
    ],
    "user": "bob",
    "status": 0
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
            "id": "56e168de47ee2036ac81e3a7",
            "organization": "greenitglobe.com",
            "role": [
                "member"
            ],
            "user": "bob",
            "status": 0
        }
    ]
}
```

`status` indicates that request is `PENDING`. Status values are as the following:

- 0: Pending request
- 1: Accepted request
- 2: Rejected request

## Accepting invitation

A user can accept invitation via an API call

```
POST /users/bob/organizations/greenitglobe.com/roles/member
{
    "id": "56e168de47ee2036ac81e3a7",
    "organization": "greenitglobe.com",
    "role": [
        "member"
    ],
    "user": "bob",
    "status": 0
}

response: 201 CREATED
{
    "id": "56e168de47ee2036ac81e3a7",
    "organization": "greenitglobe.com",
    "role": [
        "member"
    ],
    "user": "bob",
    "status": 1
}
```

## Rejecting invitation

A user can reject an invitation via an API call

```
DELETE /users/bob/organizations/greenitglobe.com/roles/member

response: 204 NO CONTENT
```

