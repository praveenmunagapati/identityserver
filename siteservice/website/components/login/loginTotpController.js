(function () {
    'use strict';
    angular.module('loginApp')
        .controller('totpController', ['$scope', '$http', '$window', '$location', totpController]);

    function totpController($scope, $http, $window, $location) {
        var vm = this;
        vm.submit = submit;
        vm.resetValidation = resetValidation;

        function submit() {
            // todo
            var data = {
                totpcode: vm.totpcode
            };
            $http
                .post('/login/totpconfirmation', data)
                .then(function (response) {
                    // success, redirect to the specified redirect URL.
                    $window.location.href = response.data.redirecturl;

                }, function (response) {
                    switch (response.status) {
                        case 422:
                            $scope.totpform.totpcode.$setValidity("invalidcode", false);
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
            $scope.totpform.totpcode.$setValidity("invalidcode", true);
        }
    }
})();