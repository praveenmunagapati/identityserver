(function() {
    'use strict';


    angular.module("itsyouonlineApp").controller("OrganizationController",OrganizationController);


    OrganizationController.$inject = ['$location','OrganizationService','$window'];

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
                    $location.path("/organization/" + vm.name);
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


})();
