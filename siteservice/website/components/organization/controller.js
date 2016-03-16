(function() {
    'use strict';


    angular
        .module("itsyouonlineApp")
        .controller("OrganizationController", OrganizationController)
        .controller("OrganizationDetailController", OrganizationDetailController)
        .controller("InvitationDialogController", InvitationDialogController);


    OrganizationController.$inject = ['$location','OrganizationService','$window'];
    OrganizationDetailController.$inject = ['$location', '$routeParams', '$window', 'OrganizationService', '$mdDialog', '$mdMedia'];

    function OrganizationController($location, OrganizationService, $window) {
        var vm = this;
        vm.create = create;

        vm.validationerrors = {};

        activate();

        function activate() {

        }

        function create(){
            var dns = []

            if (vm.dns) {
                dns.push(vm.dns)
            }

            vm.validationerrors = {};

            OrganizationService
                .create(vm.name, dns, "rob")
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

    function OrganizationDetailController($location, $routeParams, $window, OrganizationService, $mdDialog, $mdMedia) {
        var vm = this,
            globalid = $routeParams.globalid;
        vm.invitations = [];
        vm.organization = {};

        vm.showInvitationDialog = showInvitationDialog;

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


        function showInvitationDialog(ev) {
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'))
            $mdDialog.show({
                controller: InvitationDialogController,
                templateUrl: 'components/organization/views/invitationdialog.html',
                targetEvent: ev,
                fullscreen: useFullScreen,
                locals:
                    {
                        OrganizationService: OrganizationService,
                        organization : vm.organization.globalid,
                        $window: $window
                    }
            })
            .then(
                function(answer) {
                    vm.invitations.push(invitation);
                });
        }
    }

    function InvitationDialogController($scope, $mdDialog, organization, OrganizationService, $window) {

        $scope.role = "member";

        $scope.cancel = cancel;
        $scope.invite = invite;

        function cancel(){
            $mdDialog.cancel();
        }

        function invite(username, role){
            OrganizationService.invite(organization, username, role).then(
                function(data){
                    $mdDialog.hide(data);
                },
                function(reason){
                    if (reason.status == 404 && reason.data.errors && (reason.data.errors.contains("NoSuchUser") || reason.data.errors.contains("Duplicate")){
                        //Indicate error
                    }
                    else
                    {
                        $window.location.href = "error" + reason.status;
                    }
                }
            );

        }
    }



})();
