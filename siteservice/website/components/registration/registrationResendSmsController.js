(function () {
    'use strict';
    angular
        .module('itsyouonline.registration')
        .controller('resendSmsController', ['$scope', '$window', '$http', resendSmsController]);

    function resendSmsController($scope, $window, $http) {
        var vm = this;
        vm.submit = submit;
        vm.resetValidation = resetValidation;

        function submit() {
            var data = {
                phonenumber: vm.phonenumber.replace(' ', ''),
                langkey: localStorage.getItem('langKey')
            };
            $http
                .post('/register/resendsms', data)
                .then(function (response) {
                    $window.location.href = response.data.redirecturl;
                }, function (response) {
                    switch (response.status) {
                        case 422:
                            if (response.data.error === 'invalid_phonenumber') {
                                $scope.phoneconfirmationform.phonenumber.$setValidity("invalid_phonenumber", false);
                            }
                            break;
                    }
                });
        }

        function resetValidation() {
            $scope.phoneconfirmationform.phonenumber.$setValidity("invalid_phonenumber", true);
        }
    }
})();
