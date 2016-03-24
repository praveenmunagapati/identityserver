(function() {
    'use strict';


    angular
        .module("itsyouonlineApp")
        .controller("UserHomeController", UserHomeController);


    UserHomeController.$inject = [
        '$q', '$rootScope', '$location', '$window', '$mdToast', 'NotificationService', 'OrganizationService', 'UserService'];

    function UserHomeController($q, $rootScope, $location, $window, $mdToast, NotificationService, OrganizationService, UserService) {
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

        vm.checkSelected = checkSelected;
        vm.accept = accept;
        vm.reject = reject;
        vm.goToOrganization = goToOrganization;
        vm.getPendingCount = getPendingCount;

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

                    },
                    function(reason) {
                        $window.location.href = "error" + reason.status;
                    }
                );

            OrganizationService
                .getUserOrganizations(vm.username)
                .then(
                    function(data) {
                        vm.owner = data.owner;
                        vm.member = data.member;
                    },
                    function(reason) {
                        $window.location.href = "error" + reason.status;
                    }
                );


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
    }


})();
