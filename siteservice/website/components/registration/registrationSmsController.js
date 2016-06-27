(function () {
    'use strict';
    angular
        .module('itsyouonline.registration')
        .controller('smsController', ['$http', '$timeout', '$window', '$scope', '$cookies', smsController]);

    function smsController($http, $timeout, $window, $scope, $cookies) {
        var vm = this;
        vm.submit = submit;
        vm.smsconfirmation = {confirmed: false};

        $timeout(checkconfirmation, 1000);
        function checkconfirmation() {
            $http.get('register/smsconfirmed').then(
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
                .post('register/smsconfirmation', data)
                .then(function (response) {
                    $cookies.remove('registrationdetails');
                    $window.location.href = response.data.redirecturl;
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