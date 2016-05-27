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
        vm.requestingorganization = queryParams['client_id'];
        vm.requestedScopes = queryParams['scope'];
        vm.requested = {
            address: [],
            bank: [],
            email: [],
            phone: [],
            organizations: {},
            facebook: false,
            github: false
        };
        vm.requestedorganizations = [];
        vm.username = $rootScope.user;

        vm.user = {};
        vm.authorizations = {
            organizations: {}
        };

        vm.authorize = authorize;


        activate();

        function activate() {
            parseScopes();
            fetch();
        }

//
        function fetch() {

            UserService
                .get(vm.username)
                .then(
                    function(data) {
                        vm.user = data;
                        var properties = ['address', 'email', 'phone', 'bank'];
                        angular.forEach(vm.requested, function (value, property) {
                            if (properties.indexOf(property) === -1) {
                                return;
                            }
                            // loop over requests
                            var prop = vm.user[property];
                            if (!vm.authorizations[property]) {
                                vm.authorizations[property] = {};
                            }
                            // select first by default
                            angular.forEach(value, function (requestedLabel) {
                                // Empty label -> "main"
                                if (!requestedLabel) {
                                    vm.requested[property].splice(vm.requested[property].indexOf(requestedLabel), 1);
                                    requestedLabel = 'main';
                                    vm.requested[property].push(requestedLabel);
                                }
                            });
                            angular.forEach(value, function (requestedLabel) {
                                vm.authorizations[property][requestedLabel] = Object.keys(prop)[0];
                            });
                        });
                    },
                    function(reason) {
                        $window.location.href = 'error' + reason.status;
                    }
                );
        }

        function parseScopes() {
            if (vm.requestedScopes) {
                var scopes = vm.requestedScopes.split(',');
                // Filter duplicated scopes
                scopes = scopes.filter(function (item, pos, self) {
                    return self.indexOf(item) === pos;
                });
                angular.forEach(scopes, function (scope, i) {
                    var splitPermission = scope.split(':');
                    var permissionLabel = splitPermission[splitPermission.length - 1];
                    if (scope.startsWith('user:memberof:')) {
                        vm.requested.organizations[permissionLabel] = true;
                    }
                    else if (scope.startsWith('user:address:')) {
                        vm.requested.address.push(permissionLabel);
                    }
                    else if (scope.startsWith('user:email:')) {
                        vm.requested.email.push(permissionLabel);
                    }
                    else if (scope.startsWith('user:phone:')) {
                        vm.requested.phone.push(permissionLabel);
                    }
                    else if (scope.startsWith('user:bankaccount:')) {
                        vm.requested.bank.push(permissionLabel);
                    }
                    else if (scope === 'user:github') {
                        vm.requested.github = true;
                        vm.authorizations.github = true;
                    }
                    else if (scope === 'user:facebook') {
                        vm.requested.facebook = true;
                        vm.authorizations.facebook = true;
                    }
                });
            }
        }

        function authorize() {
            vm.authorizations.organizations = [];
            angular.forEach(vm.requested.organizations, function (allowed, organization) {
                if (allowed) {
                    vm.authorizations.organizations.push(organization);
                }
            });
            vm.authorizations.username = vm.username;
            vm.authorizations.grantedTo = vm.requestingorganization;
            UserService
                .saveAuthorization(vm.authorizations)
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
