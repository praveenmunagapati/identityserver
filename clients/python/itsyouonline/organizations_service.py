class OrganizationsService:
    def __init__(self, client):
        self.client = client



    def CreateNewOrganization(self, data, headers=None, query_params=None):
        """
        Create a new organization. 1 user should be in the owners list. Validation is performed to check if the securityScheme allows management on this user.
        It is method for POST /organizations
        """
        uri = self.client.base_url + "/organizations"
        return self.client.post(uri, data, headers=headers, params=query_params)


    def GetOrganization(self, globalid, headers=None, query_params=None):
        """
        Get organization info
        It is method for GET /organizations/{globalid}
        """
        uri = self.client.base_url + "/organizations/"+globalid
        return self.client.session.get(uri, headers=headers, params=query_params)


    def UpdateOrganization(self, data, globalid, headers=None, query_params=None):
        """
        Update organization info
        It is method for PUT /organizations/{globalid}
        """
        uri = self.client.base_url + "/organizations/"+globalid
        return self.client.put(uri, data, headers=headers, params=query_params)


    def DeleteOrganization(self, globalid, headers=None, query_params=None):
        """
        Deletes an organization and all data linked to it (join-organization-invitations, oauth_access_tokens, oauth_clients, logo)
        It is method for DELETE /organizations/{globalid}
        """
        uri = self.client.base_url + "/organizations/"+globalid
        return self.client.session.delete(uri, headers=headers, params=query_params)


    def CreateNewSubOrganization(self, data, globalid, headers=None, query_params=None):
        """
        Create a new suborganization.
        It is method for POST /organizations/{globalid}
        """
        uri = self.client.base_url + "/organizations/"+globalid
        return self.client.post(uri, data, headers=headers, params=query_params)


    def Set2faValidityTime(self, data, globalid, headers=None, query_params=None):
        """
        Update the 2FA validity time for the organization
        It is method for POST /organizations/{globalid}/2fa/validity
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/2fa/validity"
        return self.client.post(uri, data, headers=headers, params=query_params)


    def Get2faValidityTime(self, globalid, headers=None, query_params=None):
        """
        Get the 2FA validity time for the organization, in seconds
        It is method for GET /organizations/{globalid}/2fa/validity
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/2fa/validity"
        return self.client.session.get(uri, headers=headers, params=query_params)


    def GetOrganizationAPIKeyLabels(self, globalid, headers=None, query_params=None):
        """
        Get the list of active api keys.
        It is method for GET /organizations/{globalid}/apikeys
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/apikeys"
        return self.client.session.get(uri, headers=headers, params=query_params)


    def CreateNewOrganizationAPIKey(self, data, globalid, headers=None, query_params=None):
        """
        Create a new API Key, a secret itself should not be provided, it will be generated serverside.
        It is method for POST /organizations/{globalid}/apikeys
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/apikeys"
        return self.client.post(uri, data, headers=headers, params=query_params)


    def UpdateOrganizationAPIKey(self, data, label, globalid, headers=None, query_params=None):
        """
        Updates the label or other properties of a key.
        It is method for PUT /organizations/{globalid}/apikeys/{label}
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/apikeys/"+label
        return self.client.put(uri, data, headers=headers, params=query_params)


    def DeleteOrganizationAPIKey(self, label, globalid, headers=None, query_params=None):
        """
        Removes an API key
        It is method for DELETE /organizations/{globalid}/apikeys/{label}
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/apikeys/"+label
        return self.client.session.delete(uri, headers=headers, params=query_params)


    def GetOrganizationAPIKey(self, label, globalid, headers=None, query_params=None):
        """
        It is method for GET /organizations/{globalid}/apikeys/{label}
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/apikeys/"+label
        return self.client.session.get(uri, headers=headers, params=query_params)


    def CreateOrganizationContracty(self, data, globalid, headers=None, query_params=None):
        """
        Create a new contract.
        It is method for POST /organizations/{globalid}/contracts
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/contracts"
        return self.client.post(uri, data, headers=headers, params=query_params)


    def GetOrganizationContracts(self, globalid, headers=None, query_params=None):
        """
        Get the contracts where the organization is 1 of the parties. Order descending by date.
        It is method for GET /organizations/{globalid}/contracts
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/contracts"
        return self.client.session.get(uri, headers=headers, params=query_params)


    def CreateOrganizationDNS(self, data, dnsname, globalid, headers=None, query_params=None):
        """
        Creates a new DNS name associated with an organization
        It is method for POST /organizations/{globalid}/dns/{dnsname}
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/dns/"+dnsname
        return self.client.post(uri, data, headers=headers, params=query_params)


    def UpdateOrganizationDNS(self, data, dnsname, globalid, headers=None, query_params=None):
        """
        Updates an existing DNS name associated with an organization
        It is method for PUT /organizations/{globalid}/dns/{dnsname}
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/dns/"+dnsname
        return self.client.put(uri, data, headers=headers, params=query_params)


    def DeleteOrganizaitonDNS(self, dnsname, globalid, headers=None, query_params=None):
        """
        Removes a DNS name
        It is method for DELETE /organizations/{globalid}/dns/{dnsname}
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/dns/"+dnsname
        return self.client.session.delete(uri, headers=headers, params=query_params)


    def GetPendingOrganizationInvitations(self, globalid, headers=None, query_params=None):
        """
        Get the list of pending invitations for users to join this organization.
        It is method for GET /organizations/{globalid}/invitations
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/invitations"
        return self.client.session.get(uri, headers=headers, params=query_params)


    def RemovePendingOrganizationInvitation(self, username, globalid, headers=None, query_params=None):
        """
        Cancel a pending invitation.
        It is method for DELETE /organizations/{globalid}/invitations/{username}
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/invitations/"+username
        return self.client.session.delete(uri, headers=headers, params=query_params)


    def SetOrganizationLogo(self, data, globalid, headers=None, query_params=None):
        """
        Set the organization Logo for the organization
        It is method for PUT /organizations/{globalid}/logo
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/logo"
        return self.client.put(uri, data, headers=headers, params=query_params)


    def DeleteOrganizationLogo(self, globalid, headers=None, query_params=None):
        """
        Removes the Logo from an organization
        It is method for DELETE /organizations/{globalid}/logo
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/logo"
        return self.client.session.delete(uri, headers=headers, params=query_params)


    def GetOrganizationLogo(self, globalid, headers=None, query_params=None):
        """
        Get the Logo from an organization
        It is method for GET /organizations/{globalid}/logo
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/logo"
        return self.client.session.get(uri, headers=headers, params=query_params)


    def UpdateOrganizationMemberShip(self, data, globalid, headers=None, query_params=None):
        """
        Update an organization membership
        It is method for PUT /organizations/{globalid}/members
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/members"
        return self.client.put(uri, data, headers=headers, params=query_params)


    def AddOrganizationMember(self, data, globalid, headers=None, query_params=None):
        """
        Assign a member to organization.
        It is method for POST /organizations/{globalid}/members
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/members"
        return self.client.post(uri, data, headers=headers, params=query_params)


    def RemoveOrganizationMember(self, username, globalid, headers=None, query_params=None):
        """
        Remove a member from an organization.
        It is method for DELETE /organizations/{globalid}/members/{username}
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/members/"+username
        return self.client.session.delete(uri, headers=headers, params=query_params)


    def AddOrganizationOwner(self, data, globalid, headers=None, query_params=None):
        """
        Invite a user to become owner of an organization.
        It is method for POST /organizations/{globalid}/owners
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/owners"
        return self.client.post(uri, data, headers=headers, params=query_params)


    def RemoveOrganizationOwner(self, username, globalid, headers=None, query_params=None):
        """
        Remove an owner from organization
        It is method for DELETE /organizations/{globalid}/owners/{username}
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/owners/"+username
        return self.client.session.delete(uri, headers=headers, params=query_params)


    def AddOrganizationRegistryEntry(self, data, globalid, headers=None, query_params=None):
        """
        Adds a RegistryEntry to the organization's registry, if the key is already used, it is overwritten.
        It is method for POST /organizations/{globalid}/registry
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/registry"
        return self.client.post(uri, data, headers=headers, params=query_params)


    def ListOrganizationRegistry(self, globalid, headers=None, query_params=None):
        """
        Lists the RegistryEntries in an organization's registry.
        It is method for GET /organizations/{globalid}/registry
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/registry"
        return self.client.session.get(uri, headers=headers, params=query_params)


    def GetOrganizationRegistryEntry(self, key, globalid, headers=None, query_params=None):
        """
        Get a RegistryEntry from the organization's registry.
        It is method for GET /organizations/{globalid}/registry/{key}
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/registry/"+key
        return self.client.session.get(uri, headers=headers, params=query_params)


    def DeleteOrganizationRegistryEntry(self, key, globalid, headers=None, query_params=None):
        """
        Removes a RegistryEntry from the organization's registry
        It is method for DELETE /organizations/{globalid}/registry/{key}
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/registry/"+key
        return self.client.session.delete(uri, headers=headers, params=query_params)


    def GetOrganizationTree(self, globalid, headers=None, query_params=None):
        """
        It is method for GET /organizations/{globalid}/tree
        """
        uri = self.client.base_url + "/organizations/"+globalid+"/tree"
        return self.client.session.get(uri, headers=headers, params=query_params)
