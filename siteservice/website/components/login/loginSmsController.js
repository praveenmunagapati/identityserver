(function () {
    'use strict';
    angular
        .module('loginApp')
        .controller('smsController', ['$scope', '$http', '$timeout', '$window', smsController]);

    function smsController($scope, $http, $timeout, $window) {
        var vm = this;
        vm.submit = submit;
        vm.smsconfirmation = {confirmed: false};
        vm.resetValidation = resetValidation;

        $timeout(checkconfirmation, 1000);

        function checkconfirmation() {
            $http.get("login/smsconfirmed").then(
                function success(response) {
                    vm.smsconfirmation = response.data;
                    if (!response.data.confirmed) {
                        $timeout(checkconfirmation, 1000);
                    } else {
                        submit();
                    }
                },
                function failed(response) {
                    $timeout(checkconfirmation, 1000);
                }
            );
        }

        function submit() {
            var data = {
                smscode: vm.smscode
            };
            $http
                .post('/login/smsconfirmation', data)
                .then(function (response) {
                    // success, redirect to the specified redirect URL.
                    $window.location.href = response.data.redirecturl;
                }, function (response) {
                    switch (response.status) {
                        case 422:
                            $scope.smsform.smscode.$setValidity("invalid_code", false);
                            break;
                        case 401:
                            // Login session expired. Go back to username/password screen.
                            $window.location.hash = '#/';
                            break;
                        default:
                            $window.location.href = '/error' + response.status;
                            break;
                    }
                });
        }

        function resetValidation() {
            $scope.smsform.smscode.$setValidity("invalid_code", true);
        }
    }
})();