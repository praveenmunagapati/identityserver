(function() {
    'use strict';


    angular
        .module("itsyouonlineApp")
        .controller("AuthorizeController", AuthorizeController);


    AuthorizeController.$inject = ['$scope', '$rootScope', '$location', '$window', 'UserService'];

    function AuthorizeController($scope, $rootScope, $location, $window, UserService) {
        var vm = this;

        var queryParams = URI($location.absUrl()).search(true);
        vm.requestingorganization = queryParams['client_id'];
        vm.requestedScopes = queryParams['scope'];
        vm.requestedorganizations = [];
        vm.username = $rootScope.user;

        $scope.user = {};

        $scope.requested = {
            address: [],
            bank: [],
            email: [],
            phone: [],
            organizations: {},
            facebook: false,
            github: false
        };
        $scope.authorizations = {
            organizations: {}
        };

        $scope.update = update;


        activate();

        function activate() {
            fetch();
        }

        function fetch() {

            UserService
                .get(vm.username)
                .then(
                    function(data) {
                        $scope.user = data;
                        parseScopes();
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
                        $scope.requested.organizations[permissionLabel] = true;
                    }
                    else if (scope.startsWith('user:address:')) {
                        $scope.requested.address.push(permissionLabel);
                    }
                    else if (scope.startsWith('user:email:')) {
                        $scope.requested.email.push(permissionLabel);
                    }
                    else if (scope.startsWith('user:phone:')) {
                        $scope.requested.phone.push(permissionLabel);
                    }
                    else if (scope.startsWith('user:bankaccount:')) {
                        $scope.requested.bank.push(permissionLabel);
                    }
                    else if (scope === 'user:github') {
                        $scope.requested.github = true;
                        $scope.authorizations.github = true;
                    }
                    else if (scope === 'user:facebook') {
                        $scope.requested.facebook = true;
                        $scope.authorizations.facebook = true;
                    }
                });
                $scope.init();
            }
        }

        function update() {
            // called by the authorizationDetailsDirective
            $scope.authorizations.username = $scope.username;
            $scope.authorizations.grantedTo = $scope.requestingorganization;
            UserService
                .saveAuthorization($scope.authorizations)
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
