(function() {
    'use strict';


    angular
        .module("itsyouonlineApp")
        .controller("OrganizationController", OrganizationController)
        .controller("OrganizationDetailController", OrganizationDetailController)
        .controller("OrganizationInviteController", OrganizationInviteController);


    OrganizationController.$inject = ['$location','OrganizationService','$window'];
    OrganizationDetailController.$inject = ['$location', '$routeParams', '$window', 'OrganizationService'];
    OrganizationInviteController.$inject = [
        '$q', '$location', '$routeParams', '$window', '$mdToast', 'OrganizationService'];

    function OrganizationController($location, OrganizationService, $window) {
        var vm = this;
        vm.create = create;

        vm.validationerrors = {};

        activate();

        function activate() {

        }

        function create(){
            var dns = []
            if( vm.dns ){
                dns.push(vm.dns)
            }
            vm.validationerrors = {};
            OrganizationService.create(vm.name,dns,"rob")
            .then(
                function(data){
                    $location.path("/organizations/" + vm.name);
                },
                function(reason){
                    if (reason.status == 409) {
                         vm.validationerrors = {duplicate: true};
                    }
                    else{
                        $window.location.href = "error" + reason.status;
                    }
                }
            );
        }
    }

    function OrganizationDetailController($location, $routeParams, $window, OrganizationService) {
        var vm = this,
            globalid = $routeParams.globalid;

        vm.organization = {};

        vm.goToInvite = goToInvite;

        activate();

        function activate() {
            fetch();
        }

        function fetch(){
            OrganizationService
                .get(globalid)
                .then(
                    function(data) {
                        vm.organization = data;
                    },
                    function(reason) {
                        $window.location.href = "error" + reason.status;
                    }
                );
        }

        function goToInvite() {
            $location.path("/organizations/" + globalid + '/invite');
        }
    }


    function OrganizationInviteController($q, $location, $routeParams, $window, $mdToast, OrganizationService) {
        var vm = this,
            globalid = $routeParams.globalid;

        vm.organization = {};
        vm.members = [];

        vm.invite = invite;

        activate();

        function activate() {
            fetch();
        }

        function fetch(){
            OrganizationService
                .get(globalid)
                .then(
                    function(data) {
                        vm.organization = data;
                    },
                    function(reason) {
                        $window.location.href = "error" + reason.status;
                    }
                );
        }

        function invite() {
            var invitations = [];

            if (vm.members.length == 0) {
                return;
            }

            vm.members.forEach(function(member){
                invitations.push(inviteMember(member));
            });

            $q.all(invitations)
                .then(
                    function(responses) {
                        var toast = $mdToast
                            .simple()
                            .textContent('Invited ' + responses.length + ' members to ' + globalid)
                            .hideDelay(4000)
                            .position('top right')
                            .action('Go to organization')
                            .highlightAction(true);

                        // Show toast!
                        $mdToast
                            .show(toast)
                            .then(function(response) {
                                $location.path("/organizations/" + globalid);
                            });
                    },
                    function(reason) {
                        $window.location.href = "error" + reason.status;
                    }
                );
        }

        function inviteMember(member) {
            return OrganizationService.invite(globalid, member);
        }
    }

})();
