(function() {
    'use strict';


    angular
        .module("itsyouonlineApp")
        .controller("UserHomeController", UserHomeController);


    UserHomeController.$inject = [
        '$q', '$rootScope', '$location', '$window', '$mdToast', '$mdMedia', '$mdDialog', 'NotificationService', 'OrganizationService', 'UserService'];

    function UserHomeController($q, $rootScope, $location, $window, $mdToast, $mdMedia, $mdDialog, NotificationService, OrganizationService, UserService) {
        var vm = this;

        vm.username = $rootScope.user;
        vm.notifications = {
            invitations: [],
            approvals: [],
            contractRequests: []
        };
        vm.notificationMessage = '';

        vm.owner = [];
        vm.member = [];

        vm.user = {};

        vm.checkSelected = checkSelected;
        vm.accept = accept;
        vm.reject = reject;
        vm.goToOrganization = goToOrganization;
        vm.getPendingCount = getPendingCount;
        vm.showEmailDetailDialog = showEmailDetailDialog;
        vm.showAddEmailDialog = showAddEmailDialog;
        vm.showPhonenumberDetailDialog = showPhonenumberDetailDialog;
        vm.showAddPhonenumberDialog = showAddPhonenumberDialog;
        vm.showAddressDetailDialog = showAddressDetailDialog;
        vm.showAddAddressDialog = showAddAddressDialog;

        activate();

        function activate() {
            fetch();
        }

        function fetch() {
            NotificationService
                .get(vm.username)
                .then(
                    function(data) {
                        vm.notifications = data;
                        var count = getPendingCount(data.invitations);

                        if (count === 0) {
                            vm.notificationMessage = 'No unhandled notifications';
                        } else {
                            vm.notificationMessage = '';
                        }

                    }
                );

            OrganizationService
                .getUserOrganizations(vm.username)
                .then(
                    function(data) {
                        vm.owner = data.owner;
                        vm.member = data.member;
                    }
                );


            UserService
                .get(vm.username)
                .then(
                    function(data) {
                        vm.user = data;
                    }
                );
        }

        function getPendingCount(invitations) {
            var count = 0;
            invitations.forEach(function(invitation) {
                if (invitation.status === 'pending') {
                    count += 1;
                }
            });

            return count;
        }

        function checkSelected() {
            var selected = false;

            vm.notifications.invitations.forEach(function(invitation) {
                if (invitation.selected === true) {
                    selected = true;
                }
            });

            return selected;
        }

        function accept() {
            var requests = [];

            vm.notifications.invitations.forEach(function(invitation) {
                if (invitation.selected === true) {
                    requests.push(NotificationService.accept(invitation));
                }
            });

            $q
                .all(requests)
                .then(
                    function(responses) {
                        toast('Accepted ' + responses.length + ' invitations!');
                        fetch();
                    },
                    function(reason) {
                        $window.location.href = "error" + reason.status;
                    }
                );
        }

        function reject() {
            var requests = [];

            vm.notifications.invitations.forEach(function(invitation) {
                if (invitation.selected === true) {
                    requests.push(NotificationService.reject(invitation));
                }
            });

            $q
                .all(requests)
                .then(
                    function(responses) {
                        toast('Rejected ' + responses.length + ' invitations!');
                        fetch();
                    },
                    function(reason) {
                        $window.location.href = "error" + reason.status;
                    }
                );
        }

        function toast(message) {
            var toast = $mdToast
                .simple()
                .textContent(message)
                .hideDelay(2500)
                .position('top right');

            // Show toast!
            $mdToast.show(toast);
        }

        function goToOrganization(organization) {
            $location.path("/organizations/" + organization);
        }

        function showEmailDetailDialog(ev, label, emailaddress){
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'))
            $mdDialog.show({
                controller: EmailDetailDialogController,
                templateUrl: 'components/user/views/emailaddressdialog.html',
                targetEvent: ev,
                fullscreen: useFullScreen,
                locals:
                    {
                        UserService: UserService,
                        username : vm.username,
                        $window: $window,
                        label: label,
                        emailaddress : emailaddress,
                        deleteIsPossible: (Object.keys(vm.user.email).length > 1)
                    }
            })
            .then(
                function(data) {
                    if (data.newLabel) {
                        vm.user.email[data.newLabel] = data.emailaddress;
                    }
                    if (!data.newLabel || data.newLabel != data.originalLabel){
                        delete vm.user.email[data.originalLabel];
                    }
                });
        }

        function showAddEmailDialog(ev) {
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'))
            $mdDialog.show({
                controller: EmailDetailDialogController,
                templateUrl: 'components/user/views/emailaddressdialog.html',
                targetEvent: ev,
                fullscreen: useFullScreen,
                locals:
                    {
                        UserService: UserService,
                        username : vm.username,
                        $window: $window,
                        label: "",
                        emailaddress: "",
                        deleteIsPossible: false
                    }
            })
            .then(
                function(data) {
                    vm.user.email[data.newLabel] = data.emailaddress;
                });
        }

        function showPhonenumberDetailDialog(ev, label, phonenumber){
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'))
            $mdDialog.show({
                controller: GenericDetailDialogController,
                templateUrl: 'components/user/views/phonenumberdialog.html',
                targetEvent: ev,
                fullscreen: useFullScreen,
                locals:
                    {
                        username : vm.username,
                        $window: $window,
                        label: label,
                        data : phonenumber,
                        createFunction: UserService.registerNewPhonenumber,
                        updateFunction: UserService.updatePhonenumber,
                        deleteFunction: UserService.deletePhonenumber
                    }
            })
            .then(
                function(data) {
                    if (data.newLabel) {
                        vm.user.phone[data.newLabel] = data.data;
                    }
                    if (!data.newLabel || data.newLabel != data.originalLabel){
                        delete vm.user.phone[data.originalLabel];
                    }
                });
        }

        function showAddPhonenumberDialog(ev) {
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'))
            $mdDialog.show({
                controller: GenericDetailDialogController,
                templateUrl: 'components/user/views/phonenumberdialog.html',
                targetEvent: ev,
                fullscreen: useFullScreen,
                locals:
                    {
                        username : vm.username,
                        $window: $window,
                        label: "",
                        data: "",
                        createFunction: UserService.registerNewPhonenumber,
                        updateFunction: UserService.updatePhonenumber,
                        deleteFunction: UserService.deletePhonenumber
                    }
            })
            .then(
                function(data) {
                    vm.user.phone[data.newLabel] = data.data;
                });
        }


        function showAddressDetailDialog(ev, label, address){
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'))
            $mdDialog.show({
                controller: GenericDetailDialogController,
                templateUrl: 'components/user/views/addressdialog.html',
                targetEvent: ev,
                fullscreen: useFullScreen,
                locals:
                    {
                        username : vm.username,
                        $window: $window,
                        label: label,
                        data : address,
                        createFunction: UserService.registerNewAddress,
                        updateFunction: UserService.updateAddress,
                        deleteFunction: UserService.deleteAddress
                    }
            })
            .then(
                function(data) {
                    if (data.newLabel) {
                        vm.user.address[data.newLabel] = data.data;
                    }
                    if (!data.newLabel || data.newLabel != data.originalLabel){
                        delete vm.user.address[data.originalLabel];
                    }
                });
        }

        function showAddAddressDialog(ev) {
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'))
            $mdDialog.show({
                controller: GenericDetailDialogController,
                templateUrl: 'components/user/views/addressdialog.html',
                targetEvent: ev,
                fullscreen: useFullScreen,
                locals:
                    {
                        username : vm.username,
                        $window: $window,
                        label: "",
                        data: {},
                        createFunction: UserService.registerNewAddress,
                        updateFunction: UserService.updateAddress,
                        deleteFunction: UserService.deleteAddress
                    }
            })
            .then(
                function(data) {
                    vm.user.address[data.newLabel] = data.data;
                });
        }

    }


    function EmailDetailDialogController($scope, $mdDialog, username, UserService, $window, label, emailaddress, deleteIsPossible) {
        //If there is an emailaddress, it is already saved, if not, this means that a new one is being registered.

        $scope.emailaddress = emailaddress;
        $scope.deleteIsPossible = deleteIsPossible;

        $scope.originalLabel = label;
        $scope.label = label;
        $scope.username = username;

        $scope.cancel = cancel;
        $scope.validationerrors = {}
        $scope.create = create;
        $scope.update = update;
        $scope.deleteEmailAddress = deleteEmailAddress;


        function cancel(){
            $mdDialog.cancel();
        }

        function create(label, emailaddress){
            if (Object.keys($scope.emailaddressform.$error).length > 0 ){return;}
            $scope.validationerrors = {};
            UserService.registerNewEmailAddress(username, label, emailaddress).then(
                function(data){
                    $mdDialog.hide({originalLabel: "", newLabel: label, emailaddress: emailaddress});
                },
                function(reason){
                    if (reason.status == 409){
                        $scope.validationerrors.duplicate = true;
                    }
                    else
                    {
                        $window.location.href = "error" + reason.status;
                    }
                }
            );
        }

        function update(oldLabel, newLabel, emailaddress){
            if (Object.keys($scope.emailaddressform.$error).length > 0 ){return;}
            $scope.validationerrors = {};
            UserService.updateEmailAddress(username, oldLabel, newLabel, emailaddress).then(
                function(data){
                    $mdDialog.hide({originalLabel: oldLabel, newLabel: newLabel, emailaddress: emailaddress});
                },
                function(reason){
                    if (reason.status == 409){
                        $scope.validationerrors.duplicate = true;
                    }
                    else
                    {
                        $window.location.href = "error" + reason.status;
                    }
                }
            );
        }

        function deleteEmailAddress(label){
            $scope.validationerrors = {};
            UserService.deleteEmailAddress(username, label).then(
                function(data){
                    $mdDialog.hide({originalLabel: label, newLabel: ""});
                },
                function(reason){
                    $window.location.href = "error" + reason.status;
                }
            );
        }

    }


    function GenericDetailDialogController($scope, $mdDialog, username, $window, label, data, createFunction, updateFunction, deleteFunction) {

        $scope.data = data;

        $scope.originalLabel = label;
        $scope.label = label;
        $scope.username = username;

        $scope.cancel = cancel;
        $scope.validationerrors = {}
        $scope.create = create;
        $scope.update = update;
        $scope.remove = remove;

        function cancel(){
            $mdDialog.cancel();
        }

        function create(label, data){
            if (Object.keys($scope.dataform.$error).length > 0 ){return;}
            $scope.validationerrors = {};
            createFunction(username, label, data).then(
                function(response){
                    $mdDialog.hide({originalLabel: "", newLabel: label, data: data});
                },
                function(reason){
                    if (reason.status == 409){
                        $scope.validationerrors.duplicate = true;
                    }
                    else
                    {
                        $window.location.href = "error" + reason.status;
                    }
                }
            );
        }

        function update(oldLabel, newLabel, data){
            if (Object.keys($scope.dataform.$error).length > 0 ){return;}
            $scope.validationerrors = {};
            updateFunction(username, oldLabel, newLabel, data).then(
                function(response){
                    $mdDialog.hide({originalLabel: oldLabel, newLabel: newLabel, data: data});
                },
                function(reason){
                    if (reason.status == 409){
                        $scope.validationerrors.duplicate = true;
                    }
                    else
                    {
                        $window.location.href = "error" + reason.status;
                    }
                }
            );
        }

        function remove(label){
            $scope.validationerrors = {};
            deleteFunction(username, label).then(
                function(response){
                    $mdDialog.hide({originalLabel: label, newLabel: ""});
                },
                function(reason){
                    $window.location.href = "error" + reason.status;
                }
            );
        }

    }


})();
