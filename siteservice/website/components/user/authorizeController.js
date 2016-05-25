(function() {
    'use strict';


    angular
        .module("itsyouonlineApp")
        .controller("AuthorizeController", AuthorizeController);


    AuthorizeController.$inject = [
        '$q', '$rootScope', '$location', '$window', 'UserService'];

    function AuthorizeController($q, $rootScope, $location, $window, UserService) {
        var vm = this;

        var queryParams = URI($location.absUrl()).search(true);
        vm.requestingorganization = queryParams["client_id"];
        vm.requestedScopes = queryParams["scope"];
        vm.requestedorganizations = [];
        vm.requestedGithub = false
        vm.requestedFacebook = false
        vm.requestedAddresses = [];
        vm.requestedPhones = [];
        vm.requestedEmails = [];
        vm.requestedBankaccounts = [];
        vm.username = $rootScope.user;

        vm.user = {};

        vm.authorize = authorize;


        activate();

        function activate() {
            fetch();
            parseScopes();
        }

        function fetch() {

            UserService
                .get(vm.username)
                .then(
                    function(data) {
                        vm.user = data;
                    },
                    function(reason) {
                        $window.location.href = "error" + reason.status;
                    }
                );
        }

        function parseScopes() {
            if (vm.requestedScopes) {
                var splitScopes = vm.requestedScopes.split(",");
                for (var i = 0; i < splitScopes.length; i++) {
                    //TODO: make sure if the same scope is requested multiple times, it is only added to the lists ones
                    if (splitScopes[i].startsWith("user:memberof:")) {
                        var a = splitScopes[i].split(":");
                        vm.requestedorganizations.push(a[a.length - 1]);
                    }
                    if (splitScopes[i].startsWith("user:address:")) {
                        var a = splitScopes[i].split(":");
                        vm.requestedAddresses.push(a[a.length - 1]);
                    }
                    if (splitScopes[i].startsWith("user:email:")) {
                        var a = splitScopes[i].split(":");
                        vm.requestedEmails.push(a[a.length - 1]);
                    }
                    if (splitScopes[i].startsWith("user:phone:")) {
                        var a = splitScopes[i].split(":");
                        vm.requestedPhones.push(a[a.length - 1]);
                    }
                    if (splitScopes[i].startsWith("user:bankaccount:")) {
                        var a = splitScopes[i].split(":");
                        vm.requestedBankaccounts.push(a[a.length - 1]);
                    }
                    if (splitScopes[i] == "user:github") {
                        vm.requestedGithub = true;
                    }
                    if (splitScopes[i] == "user:facebook") {
                        vm.requestedFacebook = true;
                    }
                }
            }
        }

        function authorize() {
            UserService
                .saveAuthorization({username:vm.username, grantedTo:vm.requestingorganization, organizations:vm.requestedorganizations})
                .then(
                    function(data) {
                        var u = URI($location.absUrl());
                        var queryParams = u.search(true);
                        var endpoint = queryParams["endpoint"];
                        delete queryParams.endpoint;
                        u.pathname(endpoint);
                        u.search(queryParams);
                        $window.location.href = u.toString();
                    },
                    function(reason) {
                        $window.location.href = "error" + reason.status;
                    }
                );
        }


    }


})();
