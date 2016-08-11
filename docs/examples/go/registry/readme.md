# Access the registry examples

This sample takes a client_id and secret and adds an entry in the registry depending if the client_id/secret combination is a user or an organization api key.
- ** user ** : Add/update an entry in a user registry using client credentials
- ** organization ** : Add/update an entry in an organization registry using client credentials

After the update of the registry, an anonymous client is created and the registry entry created in the previous step is requested and printed.
