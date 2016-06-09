(function () {
    'use strict';
    angular.module('loginApp')
        .controller('twoFactorAuthenticationController', ['$scope', '$window', '$interval', '$mdDialog', 'LoginService',
            twoFactorAuthenticationController]);

    function twoFactorAuthenticationController($scope, $window, $interval, $mdDialog, LoginService) {
        var vm = this;
        vm.resetValidation = resetValidation;
        vm.shouldShowSendButton = shouldShowSendButton;
        vm.sendSmsCode = sendSmsCode;
        vm.login = login;
        vm.getHelpText = getHelpText;
        vm.smsSend = false;
        vm.selectedTwoFaMethod = null;
        vm.hasMoreThanOneTwoFaMethod = false;
        init();
        var interval;

        function init() {
            LoginService
                .getTwoFactorAuthenticationMethods()
                .then(function (data) {
                    vm.possibleTwoFaMethods = {};
                    if (data['totp']) {
                        vm.possibleTwoFaMethods['totp'] = 'Authenticator application';
                    }
                    if (data['sms'] && Object.keys(data['sms']).length) {
                        angular.forEach(data['sms'], function (sms, label) {
                            vm.possibleTwoFaMethods['sms-' + label] = 'SMS - ' + sms + ' (' + label + ')';
                        });
                    }
                    var methods = Object.keys(vm.possibleTwoFaMethods);
                    vm.hasMoreThanOneTwoFaMethod = methods.length > 1;
                    if (!methods.length) {
                        $mdDialog.show(
                            $mdDialog.alert()
                                .clickOutsideToClose(true)
                                .title('Error')
                                .htmlContent('You do not have any two factor authentication methods available. <br /> Please contact an administrator to recover your account.')
                                .ariaLabel('Error')
                                .ok('Ok')
                        );
                        return;
                    }
                    if (!vm.hasMoreThanOneTwoFaMethod && methods[0].indexOf('sms-') === 0) {
                        vm.selectedTwoFaMethod = methods[0];
                        sendSmsCode();
                    } else {
                        // Preselect based on what was selected when previously logging in
                        var method = localStorage.getItem('itsyouonline.last2falabel');
                        if (!method || methods.indexOf(method) === -1) {
                            method = methods[0];
                        }
                        vm.selectedTwoFaMethod = method;
                    }
                }, function (response) {
                    $window.location.href = 'error' + response.status;
                });
        }

        function getHelpText() {
            if (vm.selectedTwoFaMethod && vm.selectedTwoFaMethod.indexOf('sms-') === 0 && vm.smsSend) {
                return 'Click the link in the sms sent to your phone or enter the code from the sms here to continue.';
            }
            if (vm.selectedTwoFaMethod === 'totp') {
                return 'Fill in the 6 digit code from the authenticator application on your phone.';
            }
        }

        function shouldShowSendButton() {
            return vm.selectedTwoFaMethod && vm.selectedTwoFaMethod.indexOf('sms-') === 0;
        }

        function resetValidation() {
            $scope.twoFaForm.code.$setValidity("invalid_code", true);
        }

        function sendSmsCode() {
            if (interval) {
                $interval.cancel(interval);
            }
            var phoneLabel = vm.selectedTwoFaMethod.replace('sms-', '');
            vm.smsSend = true;
            LoginService
                .sendSmsCode(phoneLabel)
                .then(function () {
                    interval = $interval(checkSmsConfirmation, 1000);
                }, function (response) {
                    switch (response.status) {
                        case 401:
                            // Login session expired. Go back to username/password screen.
                            goToPage('');
                            break;
                        default:
                            goToPage('/error' + response.status);
                            break;
                    }
                });
        }

        function login() {
            var method;
            if (vm.selectedTwoFaMethod === 'totp') {
                method = LoginService.submitTotpCode;
            } else if (vm.selectedTwoFaMethod.indexOf('sms-') === 0) {
                method = LoginService.submitSmsCode;
            }
            method(vm.code)
                .then(
                    function (data) {
                        localStorage.setItem('itsyouonline.last2falabel', vm.selectedTwoFaMethod);
                        goToPage(data.redirecturl);
                    },
                    function (response) {
                        switch (response.status) {
                            case 422:
                                $scope.twoFaForm.code.$setValidity("invalid_code", false);
                                break;
                            case 401:
                                // Login session expired. Go back to username/password screen.
                                goToPage('');
                                break;
                            default:
                                goToPage('/error' + response.status);
                                break;
                        }
                    });
        }

        function checkSmsConfirmation() {
            LoginService.checkSmsConfirmation()
                .then(function (data) {
                    if (data.confirmed) {
                        login();
                    }
                });
        }

        function goToPage(url) {
            if (interval) {
                $interval.cancel(interval);
            }
            $window.location.href = url;
        }

    }
})();