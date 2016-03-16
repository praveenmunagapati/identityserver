(function() {
    'use strict';

    angular
        .module("itsyouonlineApp")
        .controller("ProfileController", ProfileController)
        .controller("ProfileEditController", ProfileEditController);

    ProfileController.$inject = ['$rootScope', '$location', '$window', 'UserService'];
    ProfileEditController.$inject = ['$rootScope', '$location', '$window', 'UserService'];

    function ProfileController($rootScope, $location, $window, UserService) {
        var vm = this;

        vm.username = $rootScope.user;

        vm.user = {};

        vm.goToEdit = goToEdit;

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
                    },
                    function(reason) {
                        $window.location.href = "error" + reason.status;
                    }
                );
        }

        function goToEdit() {
            $location.path("/profile/edit");
        }
    }

    function ProfileEditController($rootScope, $location, $window, UserService) {
        var vm = this;

        vm.username = $rootScope.user;

        vm.user = {};

        vm.goToProfile = goToProfile;

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
                    },
                    function(reason) {
                        $window.location.href = "error" + reason.status;
                    }
                );
        }

        function goToProfile() {
            $location.path("/profile");
        }
    }


})();
