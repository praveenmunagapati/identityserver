(function () {
    'use strict';
    angular.module('loginApp')
        .controller('forgotPasswordController', ['$http', '$window', '$scope', forgotPasswordController]);

    function forgotPasswordController($http, $window, $scope) {
        var vm = this;
        vm.submit = submit;
        vm.clearValidation = clearValidation;
        vm.emailSend = false;
        function submit() {
            var data = {
                login: vm.login
            };
            $http.post('/login/forgotpassword', data).then(
                function (response) {
                    vm.emailSend = true;
                },
                function (response) {
                    switch (response.status) {
                        case 404:
                            $scope.form.login.$setValidity("invalid", false);
                            break;
                    }
                }
            );
        }

        function clearValidation() {
            $scope.form.login.$setValidity("invalid", true);
        }
    }
})();