(function () {
    'use strict';
    angular
        .module('registrationApp')
        .controller('registrationController', ['$scope', '$http', '$window', 'configService', registrationController]);

    function registrationController($scope, $http, $window, configService) {
        var vm = this;
        configService.getConfig(function (config) {
            vm.totpsecret = config.totpsecret;
        });
        vm.register = register;
        vm.resetValidation = resetValidation;

        function register() {
            var data = {
                twofamethod: vm.twoFAMethod,
                login: vm.login,
                email: vm.email,
                password: vm.password,
                totpcode: vm.totpcode,
                phonenumber: vm.sms

            };
            console.log(data, data.password)
            $http
                .post('/register', data)
                .then(function (response) {
                    $window.location.href = response.data.redirecturl;
                }, function (response) {
                    console.log(response.data);
                    switch (response.status) {
                        case 422:
                            switch (response.data.error) {
                                case 'invalidphonenumber':
                                    $scope.signupform.phonenumber.$setValidity("invalidphonenumber", false);
                                    break;
                                case 'invalidtotpcode':
                                    $scope.signupform.totpcode.$setValidity("invalidtotpcode", false);
                                    break;
                                case 'invalidpassword':
                                    $scope.signupform.password.$setValidity("invalidpassword", false);
                                    break;
                                default:
                                    console.error('Unconfigured error:', response.data.error);
                            }
                            break;
                        case 409:
                            $scope.signupform.login.$setValidity('duplicateusername', false);
                            break;
                        case 401:
                            // Session expired. Reload page.
                            $window.location.reload();
                            break;
                        default:
                            $window.location.href = '/error' + response.status;
                            break;
                    }
                });
        }

        function resetValidation(prop) {
            switch (prop) {
                case 'phonenumber':
                    $scope.signupform[prop].$setValidity("invalidphonenumber", true);
                    break;
                case 'totpcode':
                    $scope.signupform[prop].$setValidity("totpcode", true);
                    break;
                case 'login':
                    $scope.signupform[prop].$setValidity("duplicateusername", true);
                    break;
                case 'twoFAMethod':
                    $scope.signupform.totpcode.$setValidity("totpcode", true);
                    $scope.signupform.phonenumber.$setValidity("invalidphonenumber", true);
                    $scope.signupform.phonenumber.$setValidity("pattern", true);
                    break;
            }
        }
    }
})();