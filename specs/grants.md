#Grants

## Problem

Currently the only wat for an organization or application to handle fine grained authorization or classification of users is to use suborganizations or store this classification itself. Neither of them are practical or a one fits all solution.
Suborganizations are a handy concept bu not really lightweight sine a user needs to accept membership for a right or authorization an app wants to giveHence the concept of  ** grants**. These are scopes an organization can stick on a user (usrname/email/phone nmber... and when u user authenticates using an oauth flow, the applications receives these grants as ** grant:name ** scopes, the users can olso be listed on a per grant basis thtough the api allowing a lightweight grouping and listing

## Todo: detailed spec to be later moved to the documentation
