(function() {
    'use strict';


    angular
        .module("itsyouonlineApp")
        .controller("AuthorizeController", AuthorizeController);


    AuthorizeController.$inject = ['$scope', '$rootScope', '$location', '$window', 'UserService', 'UserDialogService'];

    function AuthorizeController($scope, $rootScope, $location, $window, UserService, UserDialogService) {
        var vm = this;

        var queryParams = $location.search();
        vm.requestingorganization = queryParams['client_id'];
        vm.requestedScopes = queryParams['scope'];
        vm.requestedorganizations = [];
        vm.username = $rootScope.user;

        vm.user = {};

        UserDialogService.init(vm);
        vm.showAddEmailDialog = addEmail;
        vm.showAddPhonenumberDialog = addPhone;
        vm.showAddAddressDialog = addAddress;
        vm.showBankAccountDialog = bank;


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
        $scope.isNew = true;

        activate();

        function activate() {
            fetch();
        }

        function fetch() {

            UserService
                .get(vm.username)
                .then(
                    function(data) {
                        vm.user = data;
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
                    if (scope === 'user:name') {
                        $scope.requested.name = true;
                        $scope.authorizations.name = true;
                    }
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
                var properties = ['address', 'email', 'phone', 'bank'];
                angular.forEach($scope.requested, function (value, property) {
                    if (properties.indexOf(property) === -1) {
                        return;
                    }
                    // loop over requests
                    if (!$scope.authorizations[property]) {
                        $scope.authorizations[property] = {};
                    }
                    // Empty label -> "main"
                    angular.forEach(value, function (requestedLabel) {
                        if (!requestedLabel) {
                            $scope.requested[property].splice($scope.requested[property].indexOf(requestedLabel), 1);
                            requestedLabel = 'main';
                            $scope.requested[property].push(requestedLabel);
                        }
                    });
                    var prop = vm.user[property];
                    angular.forEach(value, function (requestedLabel) {
                        // select first by default, None if the user did not configure this property yet
                        if (prop) {
                            $scope.authorizations[property][requestedLabel] = Object.keys(prop)[0] || '';
                        } else {
                            $scope.authorizations[property][requestedLabel] = '';
                        }
                    });
                });
            }
        }

        function update() {
            // called by the authorizationDetailsDirective
            $scope.authorizations.username = vm.username;
            $scope.authorizations.grantedTo = vm.requestingorganization;
            UserService
                .saveAuthorization($scope.authorizations)
                .then(
                    function (data) {
                        var u = URI($location.absUrl());
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

        function addEmail(event, label) {
            selectDefault(UserDialogService.addEmail, event, label, 'email');
        }

        function addPhone(event, label) {
            selectDefault(UserDialogService.addPhonenumber, event, label, 'phone');
        }

        function addAddress(event, label) {
            selectDefault(UserDialogService.addAddress, event, label, 'address');
        }

        function bank(event, label) {
            selectDefault(UserDialogService.bankAccount, event, label, 'bank');
        }

        function selectDefault(fx, event, label, property) {
            fx(event).then(function (data) {
                $scope.authorizations[property][label] = data.newLabel;

            }, function () {
                // Select first possible value, else 'None'
                $scope.authorizations[property][label] = Object.keys(vm.user[property])[0] || '';
            });
        }
    }
})();
