(function () {
    'use strict';
    angular
        .module('loginApp')
        .controller('smsConfirmationController', ['$http', '$timeout', '$window', '$scope', smsConfirmationController]);

    function smsConfirmationController($http, $timeout, $window, $scope) {
        var vm = this;
        vm.submit = submit;
        vm.smsconfirmation = {confirmed: false};

        $timeout(checkconfirmation, 1000);

        function checkconfirmation() {
            $http.get('login/smsconfirmed').then(
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
                .post('login/smsconfirmation', data)
                .then(function (response) {
                    $window.location.href = response.data.redirecturl;
                    $cookies.remove('registrationdetails');
                }, function (response) {
                    switch (response.status) {
                        case 422:
                            $scope.phoneconfirmationform.smscode.$setValidity("invalid_sms_code", false);
                            break;
                    }
                });
        }

    }
})();