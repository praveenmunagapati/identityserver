(function () {
    'use strict';
    angular
        .module('itsyouonline.registration')
        .controller('registrationController', [
            '$scope', '$window', '$mdUtil', 'configService', 'registrationService',
            registrationController]);

    function registrationController($scope, $window, $mdUtil, configService, registrationService) {
        var vm = this;
        configService.getConfig(function (config) {
            vm.totpsecret = config.totpsecret;
        });
        vm.register = register;
        vm.resetValidation = resetValidation;
        vm.basicInfoValid = basicInfoValid;
        vm.twoFAMethod = 'sms';
        vm.validateUsername = $mdUtil.debounce(function () {
            registrationService
                .validateUsername(vm.login)
                .then(function (response) {
                    $scope.signupform['login'].$setValidity('duplicateusername', response.data.valid);
                });
        }, 500, true);

        function register() {
            registrationService
                .register(vm.twoFAMethod, vm.login, vm.email, vm.password, vm.totpcode, vm.sms)
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

        function basicInfoValid() {
            return $scope.signupform.login
                && $scope.signupform.login.$valid
                && $scope.signupform.email.$valid
                && $scope.signupform.password.$valid
                && $scope.signupform.passwordvalidation.$valid;
        }
    }
})();