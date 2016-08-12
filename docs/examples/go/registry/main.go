package main

import (
	"fmt"
	"os"

	"github.com/itsyouonline/identityserver/clients/go/itsyouonline"
)

func parseArguments() (clientID, secret, host string) {
	arguments := os.Args
	if len(arguments) < 3 {
		fmt.Println("Usage: " + arguments[0] + " client_id client_secret [host]")
		os.Exit(1)
	}
	clientID = arguments[1]
	secret = arguments[2]
	if len(arguments) > 3 {
		host = arguments[3]
	} else {
		host = "https://itsyou.online"
	}
	return
}

func main() {
	clientID, secret, host := parseArguments()

	//Step 1: Create an itsyou.online client, log in with client credentials and create a registry entry
	authenticatedClient := itsyouonline.NewItsyouonline()
	authenticatedClient.BaseURI = host + "/api"
	username, globalid, _, err := authenticatedClient.LoginWithClientCredentials(clientID, secret)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	entry := itsyouonline.RegistryEntry{Key: "testscriptkey", Value: "testscriptvalue"}
	var entries []itsyouonline.RegistryEntry

	if username != "" {
		fmt.Printf("Creating/updating a registry entry for user: %s\n", username)
		authenticatedClient.AddUserRegistryEntry(username, entry, nil, nil)

		entries, _, _ = authenticatedClient.ListUserRegistry(username, nil, nil)

	}
	if globalid != "" {
		fmt.Printf("Creating/updating a registry entry for organization: %s\n", globalid)
		authenticatedClient.AddOrganizationRegistryEntry(globalid, entry, nil, nil)

		entries, _, _ = authenticatedClient.ListOrganizationRegistry(globalid, nil, nil)
	}

	for _, e := range entries {
		fmt.Printf("Key: %s - Value: %s\n", e.Key, e.Value)
	}

	//Step 2: Get the registry entry using a anonymous client
	anonymousClient := itsyouonline.NewItsyouonline()
	anonymousClient.BaseURI = host + "/api"

	var requestedEntry itsyouonline.RegistryEntry
	fmt.Println("Getting the registry entry using an unauthenticated client")
	if username != "" {
		requestedEntry, _, err = anonymousClient.GetUserRegistryEntry("testscriptkey", username, nil, nil)
	}
	if globalid != "" {
		requestedEntry, _, err = anonymousClient.GetOrganizationRegistryEntry("testscriptkey", username, nil, nil)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("Key: %s - Value: %s\n", requestedEntry.Key, requestedEntry.Value)
}
