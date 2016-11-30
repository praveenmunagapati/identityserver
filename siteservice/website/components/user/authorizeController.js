(function() {
    'use strict';


    angular
        .module("itsyouonlineApp")
        .controller("AuthorizeController", AuthorizeController);


    AuthorizeController.$inject = ['$scope', '$rootScope', '$location', '$window', '$q',
        'UserService', 'UserDialogService', 'NotificationService'];

    function AuthorizeController($scope, $rootScope, $location, $window, $q, UserService, UserDialogService, NotificationService) {
        var vm = this;

        var queryParams = $location.search();
        vm.requestingorganization = queryParams['client_id'];
        vm.requestedScopes = queryParams['scope'];
        vm.requestedorganizations = [];
        vm.username = $rootScope.user;

        vm.user = {};
        vm.pendingNotifications = [];
        vm.pendingOrganizationInvites = {};

        UserDialogService.init(vm);
        vm.showEmailDialog = addEmail;
        vm.showPhonenumberDialog = addPhone;
        vm.showAddressDialog = addAddress;
        vm.showBankAccountDialog = bank;
        vm.showPublicKeyDialog = addPublicKey;
        vm.submit = submit;
        vm.showDigitalWalletAddressDialog = digitalWalletAddress;
        var properties = ['addresses', 'emailaddresses', 'phonenumbers', 'bankaccounts', 'digitalwallet', 'publicKeys'];
        $scope.requested = {
            organizations: {}
        };
        $scope.authorizations = {};
        angular.forEach(properties, function (prop) {
            $scope.authorizations[prop] = [];
        });

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
                    });
            NotificationService
                .get(vm.username)
                .then(
                    function (data) {
                        vm.pendingNotifications = data.invitations.filter(function (invitation) {
                            return invitation.status === 'pending';
                        });
                        angular.forEach(vm.pendingNotifications, function (invite) {
                            vm.pendingOrganizationInvites[invite.organization] = true;
                        });
                    }
                );

        }

        function parseScopes() {
            if (vm.requestedScopes) {
                var listAuthorizations = {
                    'address': 'addresses',
                    'email': 'emailaddresses',
                    'phone': 'phonenumbers',
                    'bankaccount': 'bankaccounts',
                    'publickey': 'publicKeys' // why ???
                };
                var scopes = vm.requestedScopes.split(',');
                // Filter duplicated scopes
                scopes = scopes.filter(function (item, pos, self) {
                    return self.indexOf(item) === pos;
                });
                angular.forEach(scopes, function (scope) {
                    var splitPermission = scope.split(':');
                    if (!splitPermission.length > 1) {
                        return;
                    }
                    // Empty label -> 'main'
                    var userScope = splitPermission[1];
                    var permissionLabel = splitPermission.length > 2 && splitPermission[2] ? splitPermission[2] : 'main';
                    // last part is always the read or write permission
                    var readWrite = splitPermission[splitPermission.length - 1];
                    if (!['read', 'write'].includes(readWrite)) {
                        readWrite = null;
                    }
                    var auth = {
                        requestedlabel: permissionLabel,
                        reallabel: '',
                        scope: readWrite
                    };
                    var listScope = listAuthorizations[userScope];
                    if (listScope) {
                        auth.reallabel = vm.user[listScope].length ? vm.user[listScope][0].label : '';
                        $scope.authorizations[listScope].push(auth);
                    }
                    else if (scope === 'user:name') {
                        $scope.authorizations.name = true;
                    }
                    else if (scope.startsWith('user:memberof:')) {
                        $scope.requested.organizations[permissionLabel] = true;
                    }
                    else if (scope.startsWith('user:digitalwalletaddress:')) {
                        auth.reallabel = vm.user.digitalwallet.length ? vm.user.digitalwallet[0].label : '';
                        auth.currency = splitPermission.length === 4 ? splitPermission[3] : '';
                        $scope.authorizations.digitalwallet.push(auth);
                    }
                    else if (scope === 'user:github') {
                        $scope.authorizations.github = true;
                    }
                    else if (scope === 'user:facebook') {
                        $scope.authorizations.facebook = true;
                    }
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
                    function () {
                        var u = URI($location.absUrl());
                        var endpoint = queryParams["endpoint"];
                        delete queryParams.endpoint;
                        u.pathname(endpoint);
                        u.search(queryParams);
                        $window.location.href = u.toString();
                    }
                );
        }

        function addEmail(event, auth) {
            selectDefault(UserDialogService.emailDetail, event, auth, 'emailaddresses');
        }

        function addPhone(event, auth) {
            selectDefault(UserDialogService.phonenumberDetail, event, auth, 'phonenumbers');
        }

        function addAddress(event, auth) {
            selectDefault(UserDialogService.addressDetail, event, auth, 'addresses');
        }

        function bank(event, auth) {
            selectDefault(UserDialogService.bankAccount, event, auth, 'bankaccounts');
        }

        function addPublicKey(event, auth) {
            selectDefault(UserDialogService.publicKey, event, auth, 'publicKeys');
        }

        function submit() {
            // accept all the invites first
            var requests = [];

            vm.pendingNotifications.forEach(function (invitation) {
                requests.push(NotificationService.accept(invitation));
            });

            $q.all(requests)
                .then(function () {
                    $scope.save();
                });
        }

        function digitalWalletAddress(event, auth) {
            selectDefault(UserDialogService.digitalWalletAddressDetail, event, auth, 'digitalwallet');
        }

        function selectDefault(fx, event, auth, property) {
            fx(event).then(function (data) {
                auth.reallabel = data.data.label;

            }, function () {
                // Select first possible value, else 'None'
                auth.reallabel = vm.user[property][0] ? vm.user[property][0].label : '';
            });
        }
    }
})();
