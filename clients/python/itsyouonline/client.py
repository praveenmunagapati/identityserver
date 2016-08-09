import requests

BASE_URI = "https://itsyou.online/api"


class Client:
    def __init__(self):
        self.url = BASE_URI
        self.session = requests.Session()
        self.auth_header = ''
    
    def set_auth_header(self, val):
        ''' set authorization header value'''
        self.auth_header = val


    def GetCompanyList(self, headers=None, query_params=None):
        """
        Get companies. Authorization limits are applied to requesting user.
        It is method for GET /companies
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/companies"
        return self.session.get(uri, headers=headers, params=query_params)


    def CreateCompany(self, data, headers=None, query_params=None):
        """
        Register a new company
        It is method for POST /companies
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/companies"
        return self.session.post(uri, data, headers=headers, params=query_params)


    def GetCompany(self, globalId, headers=None, query_params=None):
        """
        Get organization info
        It is method for GET /companies/{globalId}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/companies/"+globalId
        return self.session.get(uri, headers=headers, params=query_params)


    def UpdateCompany(self, data, globalId, headers=None, query_params=None):
        """
        Update existing company. Updating ``globalId`` is not allowed.
        It is method for PUT /companies/{globalId}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/companies/"+globalId
        return self.session.put(uri, data, headers=headers, params=query_params)


    def GetCompanyInfo(self, globalId, headers=None, query_params=None):
        """
        It is method for GET /companies/{globalId}/info
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/companies/"+globalId+"/info"
        return self.session.get(uri, headers=headers, params=query_params)


    def companies_byGlobalId_validate_get(self, globalId, headers=None, query_params=None):
        """
        It is method for GET /companies/{globalId}/validate
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/companies/"+globalId+"/validate"
        return self.session.get(uri, headers=headers, params=query_params)


    def GetCompanyContracts(self, globalId, headers=None, query_params=None):
        """
        Get the contracts where the organization is 1 of the parties. Order descending by date.
        It is method for GET /companies/{globalId}/contracts
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/companies/"+globalId+"/contracts"
        return self.session.get(uri, headers=headers, params=query_params)


    def CreateCompanyContract(self, data, globalId, headers=None, query_params=None):
        """
        Create a new contract.
        It is method for POST /companies/{globalId}/contracts
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/companies/"+globalId+"/contracts"
        return self.session.post(uri, data, headers=headers, params=query_params)


    def GetContract(self, contractId, headers=None, query_params=None):
        """
        Get a contract
        It is method for GET /contracts/{contractId}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/contracts/"+contractId
        return self.session.get(uri, headers=headers, params=query_params)


    def SignContract(self, data, contractId, headers=None, query_params=None):
        """
        Sign a contract
        It is method for POST /contracts/{contractId}/signatures
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/contracts/"+contractId+"/signatures"
        return self.session.post(uri, data, headers=headers, params=query_params)


    def CreateUser(self, data, headers=None, query_params=None):
        """
        Create a new user
        It is method for POST /users
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users"
        return self.session.post(uri, data, headers=headers, params=query_params)


    def GetUserPhoneNumbers(self, username, headers=None, query_params=None):
        """
        It is method for GET /users/{username}/phonenumbers
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/phonenumbers"
        return self.session.get(uri, headers=headers, params=query_params)


    def RegisterNewUserPhonenumber(self, data, username, headers=None, query_params=None):
        """
        Register a new phonenumber
        It is method for POST /users/{username}/phonenumbers
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/phonenumbers"
        return self.session.post(uri, data, headers=headers, params=query_params)


    def GetUserPhonenumberByLabel(self, label, username, headers=None, query_params=None):
        """
        It is method for GET /users/{username}/phonenumbers/{label}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/phonenumbers/"+label
        return self.session.get(uri, headers=headers, params=query_params)


    def UpdateUserPhonenumber(self, data, label, username, headers=None, query_params=None):
        """
        Update the label and/or value of an existing phonenumber.
        It is method for PUT /users/{username}/phonenumbers/{label}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/phonenumbers/"+label
        return self.session.put(uri, data, headers=headers, params=query_params)


    def DeleteUserPhonenumber(self, label, username, headers=None, query_params=None):
        """
        Removes a phonenumber
        It is method for DELETE /users/{username}/phonenumbers/{label}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/phonenumbers/"+label
        return self.session.delete(uri, headers=headers, params=query_params)


    def ValidatePhonenumber(self, data, label, username, headers=None, query_params=None):
        """
        Sends validation text to phone numbers
        It is method for POST /users/{username}/phonenumbers/{label}/activate
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/phonenumbers/"+label+"/activate"
        return self.session.post(uri, data, headers=headers, params=query_params)


    def VerifyPhoneNumber(self, data, label, username, headers=None, query_params=None):
        """
        Verifies a phone number
        It is method for PUT /users/{username}/phonenumbers/{label}/activate
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/phonenumbers/"+label+"/activate"
        return self.session.put(uri, data, headers=headers, params=query_params)


    def GetNotifications(self, username, headers=None, query_params=None):
        """
        Get the list of notifications, these are pending invitations or approvals
        It is method for GET /users/{username}/notifications
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/notifications"
        return self.session.get(uri, headers=headers, params=query_params)


    def GetAllAuthorizations(self, username, headers=None, query_params=None):
        """
        Get the list of authorizations.
        It is method for GET /users/{username}/authorizations
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/authorizations"
        return self.session.get(uri, headers=headers, params=query_params)


    def GetAuthorization(self, grantedTo, username, headers=None, query_params=None):
        """
        Get the authorization for a specific organization.
        It is method for GET /users/{username}/authorizations/{grantedTo}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/authorizations/"+grantedTo
        return self.session.get(uri, headers=headers, params=query_params)


    def UpdateAuthorization(self, data, grantedTo, username, headers=None, query_params=None):
        """
        Modify which information an organization is able to see.
        It is method for PUT /users/{username}/authorizations/{grantedTo}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/authorizations/"+grantedTo
        return self.session.put(uri, data, headers=headers, params=query_params)


    def DeleteAuthorization(self, grantedTo, username, headers=None, query_params=None):
        """
        Remove the authorization for an organization, the granted organization will no longer have access the user's information.
        It is method for DELETE /users/{username}/authorizations/{grantedTo}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/authorizations/"+grantedTo
        return self.session.delete(uri, headers=headers, params=query_params)


    def GetUserOrganizations(self, username, headers=None, query_params=None):
        """
        Get the list organizations a user is owner or member of
        It is method for GET /users/{username}/organizations
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/organizations"
        return self.session.get(uri, headers=headers, params=query_params)


    def LeaveOrganization(self, globalid, username, headers=None, query_params=None):
        """
        Removes the user from an organization
        It is method for DELETE /users/{username}/organizations/{globalid}/leave
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/organizations/"+globalid+"/leave"
        return self.session.delete(uri, headers=headers, params=query_params)


    def AcceptMembership(self, data, role, globalid, username, headers=None, query_params=None):
        """
        Accept membership in organization
        It is method for POST /users/{username}/organizations/{globalid}/roles/{role}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/organizations/"+globalid+"/roles/"+role
        return self.session.post(uri, data, headers=headers, params=query_params)


    def users_byUsernameorganizations_byGlobalid_rolesrole_delete(self, role, globalid, username, headers=None, query_params=None):
        """
        Reject membership invitation in an organization.
        It is method for DELETE /users/{username}/organizations/{globalid}/roles/{role}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/organizations/"+globalid+"/roles/"+role
        return self.session.delete(uri, headers=headers, params=query_params)


    def GetUserInformation(self, username, headers=None, query_params=None):
        """
        It is method for GET /users/{username}/info
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/info"
        return self.session.get(uri, headers=headers, params=query_params)


    def ValidateUsername(self, username, headers=None, query_params=None):
        """
        It is method for GET /users/{username}/validate
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/validate"
        return self.session.get(uri, headers=headers, params=query_params)


    def GetDigitalWallet(self, username, headers=None, query_params=None):
        """
        It is method for GET /users/{username}/digitalwallet
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/digitalwallet"
        return self.session.get(uri, headers=headers, params=query_params)


    def RegisterNewDigitalAssetAddress(self, data, username, headers=None, query_params=None):
        """
        Register a new digital asset address
        It is method for POST /users/{username}/digitalwallet
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/digitalwallet"
        return self.session.post(uri, data, headers=headers, params=query_params)


    def GetDigitalAssetAddressByLabel(self, label, username, headers=None, query_params=None):
        """
        It is method for GET /users/{username}/digitalwallet/{label}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/digitalwallet/"+label
        return self.session.get(uri, headers=headers, params=query_params)


    def UpdateDigitalAssetAddress(self, data, label, username, headers=None, query_params=None):
        """
        Update the label and/or value of an existing address.
        It is method for PUT /users/{username}/digitalwallet/{label}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/digitalwallet/"+label
        return self.session.put(uri, data, headers=headers, params=query_params)


    def DeleteDigitalAssetAddress(self, label, username, headers=None, query_params=None):
        """
        Removes an address
        It is method for DELETE /users/{username}/digitalwallet/{label}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/digitalwallet/"+label
        return self.session.delete(uri, headers=headers, params=query_params)


    def GetUserBankAccounts(self, username, headers=None, query_params=None):
        """
        It is method for GET /users/{username}/banks
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/banks"
        return self.session.get(uri, headers=headers, params=query_params)


    def CreateUserBankAccount(self, data, username, headers=None, query_params=None):
        """
        Create new bank account
        It is method for POST /users/{username}/banks
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/banks"
        return self.session.post(uri, data, headers=headers, params=query_params)


    def GetUserBankAccountByLabel(self, username, label, headers=None, query_params=None):
        """
        It is method for GET /users/{username}/banks/{label}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/banks/"+label
        return self.session.get(uri, headers=headers, params=query_params)


    def UpdateUserBankAccount(self, data, username, label, headers=None, query_params=None):
        """
        Update an existing bankaccount and label.
        It is method for PUT /users/{username}/banks/{label}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/banks/"+label
        return self.session.put(uri, data, headers=headers, params=query_params)


    def DeleteUserBankAccount(self, username, label, headers=None, query_params=None):
        """
        Delete a BankAccount
        It is method for DELETE /users/{username}/banks/{label}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/banks/"+label
        return self.session.delete(uri, headers=headers, params=query_params)


    def GetUserContracts(self, username, headers=None, query_params=None):
        """
        Get the contracts where the user is 1 of the parties. Order descending by date.
        It is method for GET /users/{username}/contracts
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/contracts"
        return self.session.get(uri, headers=headers, params=query_params)


    def CreateUserContract(self, data, username, headers=None, query_params=None):
        """
        Create a new contract.
        It is method for POST /users/{username}/contracts
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/contracts"
        return self.session.post(uri, data, headers=headers, params=query_params)


    def GetUser(self, username, headers=None, query_params=None):
        """
        It is method for GET /users/{username}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username
        return self.session.get(uri, headers=headers, params=query_params)


    def GetTwoFAMethods(self, username, headers=None, query_params=None):
        """
        Get the possible two factor authentication methods
        It is method for GET /users/{username}/twofamethods
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/twofamethods"
        return self.session.get(uri, headers=headers, params=query_params)


    def GetTOTPSecret(self, username, headers=None, query_params=None):
        """
        It is method for GET /users/{username}/totp
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/totp"
        return self.session.get(uri, headers=headers, params=query_params)


    def SetupTOTP(self, data, username, headers=None, query_params=None):
        """
        It is method for POST /users/{username}/totp
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/totp"
        return self.session.post(uri, data, headers=headers, params=query_params)


    def RemoveTOTP(self, username, headers=None, query_params=None):
        """
        It is method for DELETE /users/{username}/totp
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/totp"
        return self.session.delete(uri, headers=headers, params=query_params)


    def GetEmailAddresses(self, username, headers=None, query_params=None):
        """
        Get a list of the user his email addresses.
        It is method for GET /users/{username}/emailaddresses
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/emailaddresses"
        return self.session.get(uri, headers=headers, params=query_params)


    def RegisterNewEmailAddress(self, data, username, headers=None, query_params=None):
        """
        Register a new email address
        It is method for POST /users/{username}/emailaddresses
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/emailaddresses"
        return self.session.post(uri, data, headers=headers, params=query_params)


    def UpdateEmailAddress(self, data, label, username, headers=None, query_params=None):
        """
        Updates the label and/or value of an email address
        It is method for PUT /users/{username}/emailaddresses/{label}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/emailaddresses/"+label
        return self.session.put(uri, data, headers=headers, params=query_params)


    def DeleteEmailAddress(self, label, username, headers=None, query_params=None):
        """
        Removes an email address
        It is method for DELETE /users/{username}/emailaddresses/{label}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/emailaddresses/"+label
        return self.session.delete(uri, headers=headers, params=query_params)


    def ValidateEmailAddress(self, data, label, username, headers=None, query_params=None):
        """
        Sends validation email to email address
        It is method for POST /users/{username}/emailaddresses/{label}/validate
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/emailaddresses/"+label+"/validate"
        return self.session.post(uri, data, headers=headers, params=query_params)


    def ListUserRegistry(self, username, headers=None, query_params=None):
        """
        Lists the Registry entries
        It is method for GET /users/{username}/registry
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/registry"
        return self.session.get(uri, headers=headers, params=query_params)


    def AddUserRegistryEntry(self, data, username, headers=None, query_params=None):
        """
        Adds a RegistryEntry to the user's registry, if the key is already used, it is overwritten.
        It is method for POST /users/{username}/registry
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/registry"
        return self.session.post(uri, data, headers=headers, params=query_params)


    def GetUserRegistryEntry(self, key, username, headers=None, query_params=None):
        """
        Get a RegistryEntry from the user's registry.
        It is method for GET /users/{username}/registry/{key}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/registry/"+key
        return self.session.get(uri, headers=headers, params=query_params)


    def DeleteUserRegistryEntry(self, key, username, headers=None, query_params=None):
        """
        Removes a RegistryEntry from the user's registry
        It is method for DELETE /users/{username}/registry/{key}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/registry/"+key
        return self.session.delete(uri, headers=headers, params=query_params)


    def DeleteFacebookAccount(self, username, headers=None, query_params=None):
        """
        Delete the associated facebook account
        It is method for DELETE /users/{username}/facebook
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/facebook"
        return self.session.delete(uri, headers=headers, params=query_params)


    def UpdateUserName(self, data, username, headers=None, query_params=None):
        """
        Update the user his firstname and lastname
        It is method for PUT /users/{username}/name
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/name"
        return self.session.put(uri, data, headers=headers, params=query_params)


    def UpdatePassword(self, data, username, headers=None, query_params=None):
        """
        Update the user his password
        It is method for PUT /users/{username}/password
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/password"
        return self.session.put(uri, data, headers=headers, params=query_params)


    def ListAPIKeys(self, username, headers=None, query_params=None):
        """
        Lists the API keys
        It is method for GET /users/{username}/apikeys
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/apikeys"
        return self.session.get(uri, headers=headers, params=query_params)


    def AddApiKey(self, data, username, headers=None, query_params=None):
        """
        Adds an APIKey to the user
        It is method for POST /users/{username}/apikeys
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/apikeys"
        return self.session.post(uri, data, headers=headers, params=query_params)


    def GetAPIkey(self, label, username, headers=None, query_params=None):
        """
        Get an API key by label
        It is method for GET /users/{username}/apikeys/{label}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/apikeys/"+label
        return self.session.get(uri, headers=headers, params=query_params)


    def UpdateAPIkey(self, data, label, username, headers=None, query_params=None):
        """
        Updates the label for the api key
        It is method for PUT /users/{username}/apikeys/{label}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/apikeys/"+label
        return self.session.put(uri, data, headers=headers, params=query_params)


    def DeleteAPIkey(self, label, username, headers=None, query_params=None):
        """
        Removes an API key
        It is method for DELETE /users/{username}/apikeys/{label}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/apikeys/"+label
        return self.session.delete(uri, headers=headers, params=query_params)


    def DeleteGithubAccount(self, username, headers=None, query_params=None):
        """
        Unlink Github Account
        It is method for DELETE /users/{username}/github
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/github"
        return self.session.delete(uri, headers=headers, params=query_params)


    def GetUserAddresses(self, username, headers=None, query_params=None):
        """
        It is method for GET /users/{username}/addresses
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/addresses"
        return self.session.get(uri, headers=headers, params=query_params)


    def RegisterNewUserAddress(self, data, username, headers=None, query_params=None):
        """
        Register a new address
        It is method for POST /users/{username}/addresses
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/addresses"
        return self.session.post(uri, data, headers=headers, params=query_params)


    def GetUserAddressByLabel(self, label, username, headers=None, query_params=None):
        """
        It is method for GET /users/{username}/addresses/{label}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/addresses/"+label
        return self.session.get(uri, headers=headers, params=query_params)


    def UpdateUserAddress(self, data, label, username, headers=None, query_params=None):
        """
        Update the label and/or value of an existing address.
        It is method for PUT /users/{username}/addresses/{label}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/addresses/"+label
        return self.session.put(uri, data, headers=headers, params=query_params)


    def DeleteUserAddress(self, label, username, headers=None, query_params=None):
        """
        Removes an address
        It is method for DELETE /users/{username}/addresses/{label}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/users/"+username+"/addresses/"+label
        return self.session.delete(uri, headers=headers, params=query_params)


    def CreateNewOrganization(self, data, headers=None, query_params=None):
        """
        Create a new organization. 1 user should be in the owners list. Validation is performed to check if the securityScheme allows management on this user.
        It is method for POST /organizations
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations"
        return self.session.post(uri, data, headers=headers, params=query_params)


    def GetOrganization(self, globalid, headers=None, query_params=None):
        """
        Get organization info
        It is method for GET /organizations/{globalid}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid
        return self.session.get(uri, headers=headers, params=query_params)


    def CreateNewSubOrganization(self, data, globalid, headers=None, query_params=None):
        """
        Create a new suborganization.
        It is method for POST /organizations/{globalid}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid
        return self.session.post(uri, data, headers=headers, params=query_params)


    def UpdateOrganization(self, data, globalid, headers=None, query_params=None):
        """
        Update organization info
        It is method for PUT /organizations/{globalid}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid
        return self.session.put(uri, data, headers=headers, params=query_params)


    def DeleteOrganization(self, globalid, headers=None, query_params=None):
        """
        Deletes an organization and all data linked to it (join-organization-invitations, oauth_access_tokens, oauth_clients)
        It is method for DELETE /organizations/{globalid}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid
        return self.session.delete(uri, headers=headers, params=query_params)


    def AddOrganizationOwner(self, data, globalid, headers=None, query_params=None):
        """
        Invite a user to become owner of an organization.
        It is method for POST /organizations/{globalid}/owners
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid+"/owners"
        return self.session.post(uri, data, headers=headers, params=query_params)


    def RemoveOrganizationOwner(self, username, globalid, headers=None, query_params=None):
        """
        Remove an owner from organization
        It is method for DELETE /organizations/{globalid}/owners/{username}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid+"/owners/"+username
        return self.session.delete(uri, headers=headers, params=query_params)


    def GetOrganizationContracts(self, globalid, headers=None, query_params=None):
        """
        Get the contracts where the organization is 1 of the parties. Order descending by date.
        It is method for GET /organizations/{globalid}/contracts
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid+"/contracts"
        return self.session.get(uri, headers=headers, params=query_params)


    def CreateOrganizationContracty(self, data, globalid, headers=None, query_params=None):
        """
        Create a new contract.
        It is method for POST /organizations/{globalid}/contracts
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid+"/contracts"
        return self.session.post(uri, data, headers=headers, params=query_params)


    def GetOrganizationAPIKeyLabels(self, globalid, headers=None, query_params=None):
        """
        Get the list of active api keys.
        It is method for GET /organizations/{globalid}/apikeys
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid+"/apikeys"
        return self.session.get(uri, headers=headers, params=query_params)


    def CreateNewOrganizationAPIKey(self, data, globalid, headers=None, query_params=None):
        """
        Create a new API Key, a secret itself should not be provided, it will be generated serverside.
        It is method for POST /organizations/{globalid}/apikeys
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid+"/apikeys"
        return self.session.post(uri, data, headers=headers, params=query_params)


    def GetOrganizationAPIKey(self, label, globalid, headers=None, query_params=None):
        """
        It is method for GET /organizations/{globalid}/apikeys/{label}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid+"/apikeys/"+label
        return self.session.get(uri, headers=headers, params=query_params)


    def UpdateOrganizationAPIKey(self, data, label, globalid, headers=None, query_params=None):
        """
        Updates the label or other properties of a key.
        It is method for PUT /organizations/{globalid}/apikeys/{label}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid+"/apikeys/"+label
        return self.session.put(uri, data, headers=headers, params=query_params)


    def DeleteOrganizationAPIKey(self, label, globalid, headers=None, query_params=None):
        """
        Removes an API key
        It is method for DELETE /organizations/{globalid}/apikeys/{label}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid+"/apikeys/"+label
        return self.session.delete(uri, headers=headers, params=query_params)


    def ListOrganizationRegistry(self, globalid, headers=None, query_params=None):
        """
        Lists the RegistryEntries in an organization's registry.
        It is method for GET /organizations/{globalid}/registry
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid+"/registry"
        return self.session.get(uri, headers=headers, params=query_params)


    def AddOrganizationRegistryEntry(self, data, globalid, headers=None, query_params=None):
        """
        Adds a RegistryEntry to the organization's registry, if the key is already used, it is overwritten.
        It is method for POST /organizations/{globalid}/registry
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid+"/registry"
        return self.session.post(uri, data, headers=headers, params=query_params)


    def GetOrganizationTree(self, globalid, headers=None, query_params=None):
        """
        It is method for GET /organizations/{globalid}/tree
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid+"/tree"
        return self.session.get(uri, headers=headers, params=query_params)


    def AddOrganizationMember(self, data, globalid, headers=None, query_params=None):
        """
        Assign a member to organization.
        It is method for POST /organizations/{globalid}/members
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid+"/members"
        return self.session.post(uri, data, headers=headers, params=query_params)


    def UpdateOrganizationMemberShip(self, data, globalid, headers=None, query_params=None):
        """
        Update an organization membership
        It is method for PUT /organizations/{globalid}/members
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid+"/members"
        return self.session.put(uri, data, headers=headers, params=query_params)


    def RemoveOrganizationMember(self, username, globalid, headers=None, query_params=None):
        """
        Remove a member from an organization.
        It is method for DELETE /organizations/{globalid}/members/{username}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid+"/members/"+username
        return self.session.delete(uri, headers=headers, params=query_params)


    def GetOrganizationRegistryEntry(self, key, globalid, headers=None, query_params=None):
        """
        Get a RegistryEntry from the organization's registry.
        It is method for GET /organizations/{globalid}/registry/{key}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid+"/registry/"+key
        return self.session.get(uri, headers=headers, params=query_params)


    def DeleteOrganizationRegistryEntry(self, key, globalid, headers=None, query_params=None):
        """
        Removes a RegistryEntry from the organization's registry
        It is method for DELETE /organizations/{globalid}/registry/{key}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid+"/registry/"+key
        return self.session.delete(uri, headers=headers, params=query_params)


    def CreateOrganizationDNS(self, data, dnsname, globalid, headers=None, query_params=None):
        """
        Creates a new DNS name associated with an organization
        It is method for POST /organizations/{globalid}/dns/{dnsname}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid+"/dns/"+dnsname
        return self.session.post(uri, data, headers=headers, params=query_params)


    def UpdateOrganizationDNS(self, data, dnsname, globalid, headers=None, query_params=None):
        """
        Updates an existing DNS name associated with an organization
        It is method for PUT /organizations/{globalid}/dns/{dnsname}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid+"/dns/"+dnsname
        return self.session.put(uri, data, headers=headers, params=query_params)


    def DeleteOrganizaitonDNS(self, dnsname, globalid, headers=None, query_params=None):
        """
        Removes a DNS name
        It is method for DELETE /organizations/{globalid}/dns/{dnsname}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid+"/dns/"+dnsname
        return self.session.delete(uri, headers=headers, params=query_params)


    def GetPendingOrganizationInvitations(self, globalid, headers=None, query_params=None):
        """
        Get the list of pending invitations for users to join this organization.
        It is method for GET /organizations/{globalid}/invitations
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid+"/invitations"
        return self.session.get(uri, headers=headers, params=query_params)


    def RemovePendingOrganizationInvitation(self, username, globalid, headers=None, query_params=None):
        """
        Cancel a pending invitation.
        It is method for DELETE /organizations/{globalid}/invitations/{username}
        """
        if self.auth_header:
            self.session.headers.update({"Authorization":self.auth_header})

        uri = self.url + "/organizations/"+globalid+"/invitations/"+username
        return self.session.delete(uri, headers=headers, params=query_params)
