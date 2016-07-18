# Registry per user/organization

## Concept

A simple key-value store per user/organization. Only the user/organization can list/add/modify/delete the entries, everyone can read, given the key is known.


## List/add/modify/delete entries

### Required scopes

To modify registry entries for a user, the `user:admin` scope is required, to modify the entries for an organization, you need the `organization:owner` scope.

### URL's

- Users:

    List entries:
    `GET /api/users/{username}/registry`

    Create entry:
    `POST /api/users/{username}/registry`

    Update Entry:
    `PUT /api/users/{username}/registry/{key}`

    Delete Entry:
    `DELETE /api/users/{username}/registry/{key}`

- Organizations:

    List entries:
    `GET /api/organizations/{globalid}/registry`

    Create entry:
    `POST /api/organizations/{globalid}/registry`

    Update Entry:
    `PUT /api/organizations/{globalid}/registry/{key}`

    Delete Entry:
    `DELETE /api/organizations/{globalid}/registry/{key}`

## Read an entry

No oauth scope is required

### URL's

- Users:
    `GET /api/users/{username}/registry/{key}`

- Organizations:
    `GET /api/organizations/{globalid}/registry/{key}`
