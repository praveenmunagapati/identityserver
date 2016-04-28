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
                for (let scope of splitScopes) {
                    if (scope.startsWith("user:memberof:")) {
                        var a = scope.split(":")
                        vm.requestedorganizations.push(a[a.length - 1])
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
