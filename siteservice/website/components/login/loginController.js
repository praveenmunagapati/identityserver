(function () {
    'use strict';
    angular.module('loginApp')
        .controller('loginController', ['$http', '$window', '$scope', loginController]);

    function loginController($http, $window, $scope) {
        var vm = this;
        vm.submit = submit;
        vm.clearValidation = clearValidation;
        vm.externalSite = URI($window.location.href).search(true).client_id;

        function submit() {
            var data = {
                login: vm.login,
                password: vm.password
            };
            $http.post('/login', data).then(
                function (response) {
                    // Redirect 2 factor authentication page
                    $window.location.hash = '#/2fa';
                },
                function (response) {
                    if (response.status === 422) {
                        $scope.loginform.password.$setValidity("invalidcredentials", false);
                    }
                }
            );
        }

        function clearValidation() {
            $scope.loginform.password.$setValidity("invalidcredentials", true);
        }
    }
})();