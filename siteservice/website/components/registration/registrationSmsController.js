(function () {
    'use strict';
    angular
        .module('registrationApp')
        .controller('smsController', ['$http', '$timeout', '$window', smsController]);

    function smsController($http, $timeout, $window) {
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
                    console.log(response.data)
                    $window.location.href = response.data.redirecturl;
                }, function (response) {
                    switch (response.status) {
                        case 422:
                            if (response.data.error === 'invalidsmscode') {
                                $scope.phoneconfirmationform.phonenumber.$setValidity("invalidphonenumber", false);
                            }
                            break;
                        case 401:
                            // Session expired. Go back to registration page.
                            $window.location.hash = '';
                            break;
                        default:
                            $window.location.href = '/error' + response.status;
                            break;
                    }
                });
        }

    }
})();