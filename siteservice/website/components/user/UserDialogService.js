/**
 * Created by lucas on 13/06/16.
 */
(function () {
    'use strict';

    angular
        .module('itsyouonline.user', [])
        .factory('UserDialogService', ['$window', '$q', '$interval', '$mdMedia', '$mdDialog', 'UserService', 'configService', UserDialogService]);

    function UserDialogService($window, $q, $interval, $mdMedia, $mdDialog, UserService, configService) {
        var vm;
        var genericDetailControllerParams = ['$scope', '$mdDialog', 'username', '$window', 'label', 'data',
            'createFunction', 'updateFunction', 'deleteFunction', GenericDetailDialogController];
        return {
            init: init,
            addEmail: addEmail,
            emailDetail: emailDetail,
            addPhonenumber: addPhonenumber,
            addAddress: addAddress,
            addressDetail: addressDetail,
            phonenumberDetail: phonenumberDetail,
            verifyPhone: verifyPhone,
            bankAccount: bankAccount,
            facebook: facebook,
            addFacebook: addFacebook,
            github: github,
            addGithub: addGithub,
            showSimpleDialog: showSimpleDialog
        };

        function init(scope) {
            vm = scope;
        }

        function doNothing() {
        }

        function addEmail(ev) {
            return $q(function (resolve, reject) {
                $mdDialog.show({
                    controller: ['$scope', '$mdDialog', '$window', 'UserService', 'username', 'label', 'emailaddress', 'deleteIsPossible', EmailDetailDialogController],
                    controllerAs: 'ctrl',
                    templateUrl: 'components/user/views/emailaddressdialog.html',
                    targetEvent: ev,
                    fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                    locals: {
                        username: vm.username,
                        label: "",
                        emailaddress: "",
                        deleteIsPossible: false
                    }
                })
                    .then(
                        function (data) {
                            vm.user.email[data.newLabel] = data.emailaddress;
                            resolve(data);
                        }, reject);
            });
        }

        function emailDetail(ev, label, emailaddress) {
            $mdDialog.show({
                controller: ['$scope', '$mdDialog', '$window', 'UserService', 'username', 'label', 'emailaddress', 'deleteIsPossible', EmailDetailDialogController],
                controllerAs: 'ctrl',
                templateUrl: 'components/user/views/emailaddressdialog.html',
                targetEvent: ev,
                fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                locals: {
                    username: vm.username,
                    label: label,
                    emailaddress: emailaddress,
                    deleteIsPossible: Object.keys(vm.user.email).length > 1
                }
            })
                .then(
                    function (data) {
                        if (data.newLabel) {
                            vm.user.email[data.newLabel] = data.emailaddress;
                        }
                        if (!data.newLabel || data.newLabel != data.originalLabel) {
                            delete vm.user.email[data.originalLabel];
                        }
                    });
        }

        function addPhonenumber(ev) {
            return $q(function (resolve, reject) {
                $mdDialog.show({
                    controller: genericDetailControllerParams,
                    templateUrl: 'components/user/views/phonenumberdialog.html',
                    targetEvent: ev,
                    fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                    locals: {
                        username: vm.username,
                        label: "",
                        data: "",
                        createFunction: UserService.registerNewPhonenumber,
                        updateFunction: UserService.updatePhonenumber,
                        deleteFunction: UserService.deletePhonenumber
                    }
                })
                    .then(
                        function (data) {
                            vm.user.phone[data.newLabel] = data.data;
                            // Verify a phonenumber if it's the same number as an already verified one.
                            var isVerified = false;
                            angular.forEach(vm.user.verifiedPhones, function (number, label) {
                                if (number === data.data) {
                                    vm.user.verifiedPhones[data.newLabel] = data.data;
                                    isVerified = true;
                                }
                            });
                            if (!isVerified) {
                                verifyPhone(ev, data.newLabel, data.data);
                            }
                            resolve(data);
                        }, reject);
            });
        }

        function phonenumberDetail(ev, label, phonenumber) {
            function deletePhoneNumber(username, label) {
                return $q(function (resolve, reject) {
                    UserService
                        .deletePhonenumber(username, label, false)
                        .then(resolve, function (response) {
                            if (response.status === 409) {
                                var errorMsg, dialog;
                                if (response.data.error === 'warning_delete_last_verified_phone_number') {
                                    errorMsg = 'Are you sure you want to delete this phone number? <br />' +
                                        'It is your last verified phone number, which means you will <br />' +
                                        'no longer be able to login using sms confirmations.';
                                    dialog = $mdDialog.confirm()
                                        .title('Confirm deletion')
                                        .ok('Confirm')
                                        .cancel('Cancel');
                                }
                                else if (response.data.error === 'cannot_delete_last_verified_phone_number') {
                                    errorMsg = 'You cannot delete your last verified phone number. <br />' +
                                        'Please change your 2 factor authentication method or add another verified phone number.';
                                    dialog = $mdDialog.alert()
                                        .title('Error')
                                        .ok('Close');
                                }
                                dialog = dialog.htmlContent(errorMsg)
                                    .ariaLabel('Delete phone number')
                                    .targetEvent(ev);
                                $mdDialog.show(dialog)
                                    .then(function () {
                                        UserService
                                            .deletePhonenumber(username, label, true)
                                            .then(function () {
                                                // Manually remove phone number since the dialog which executes the updatePhoneNumber promise callback had been closed before
                                                delete vm.user.phone[label];
                                            }, function () {
                                                showSimpleDialog('Could not delete phone number. Please try again later.');
                                            });
                                    });
                            } else {
                                showSimpleDialog('Could not delete phone number. Please try again later.');
                                reject();
                            }
                        });
                });
            }

            $mdDialog.show({
                controller: genericDetailControllerParams,
                templateUrl: 'components/user/views/phonenumberdialog.html',
                targetEvent: ev,
                fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                locals: {
                    username: vm.username,
                    label: label,
                    data: phonenumber,
                    createFunction: UserService.registerNewPhonenumber,
                    updateFunction: UserService.updatePhonenumber,
                    deleteFunction: deletePhoneNumber
                }
            })
                .then(updatePhoneNumber, function () {
                });

            function updatePhoneNumber(data) {
                // no data is provided when dialog is closed because another dialog opened (in case a confirmation is asked)
                if (data) {
                    if (data.newLabel) {
                        vm.user.phone[data.newLabel] = data.data;
                    }
                    if (!data.newLabel || data.newLabel != data.originalLabel) {
                        delete vm.user.phone[data.originalLabel];
                    }
                }
            }
        }

        function verifyPhone(event, label, phonenumber) {
            var interval;
            $mdDialog.show({
                controller: ['$scope', '$mdDialog', '$interval', 'user', 'label', 'phonenumber', verifyPhoneDialogController],
                controllerAs: 'ctrl',
                templateUrl: 'components/user/views/verifyPhoneDialog.html',
                targetEvent: event,
                fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                locals: {
                    user: vm.user,
                    label: label,
                    phonenumber: phonenumber
                }
            }).finally(function () {
                $interval.cancel(interval);
            });

            function verifyPhoneDialogController($scope, $mdDialog, $interval, user, label, phonenumber) {
                var ctrl = this;
                ctrl.label = label;
                ctrl.phonenumber = phonenumber;
                ctrl.close = close;
                ctrl.submit = submit;
                ctrl.validationKey = '';
                ctrl.resetValidation = resetValidation;

                init();

                function init() {
                    UserService
                        .sendPhoneVerificationCode(vm.username, label)
                        .then(function (responseData) {
                            ctrl.validationKey = responseData.validationkey;
                            interval = $interval(checkconfirmation, 1000);
                        }, function (response) {
                            $mdDialog.show(
                                $mdDialog.alert()
                                    .clickOutsideToClose(true)
                                    .title('Error')
                                    .textContent('Failed to send verification code. Please try again later.')
                                    .ariaLabel('Error while sending verification code')
                                    .ok('Close')
                                    .targetEvent(event)
                            );
                        });
                }

                function close() {
                    $mdDialog.cancel();
                }

                function checkconfirmation() {
                    UserService
                        .getVerifiedPhones(vm.username)
                        .then(function success(confirmedPhones) {
                            var isConfirmed = false;
                            angular.forEach(confirmedPhones, function (number, label) {
                                if (label === ctrl.label) {
                                    isConfirmed = true;
                                }
                            });
                            if (isConfirmed) {
                                if (vm.user.verifiedPhones) {
                                    vm.user.verifiedPhones[ctrl.label] = ctrl.phonenumber;
                                }
                                close();
                            }
                        });
                }

                function submit() {
                    UserService
                        .verifyPhone(user.username, ctrl.label, ctrl.validationKey, ctrl.smscode)
                        .then(function () {
                            if (vm.user.verifiedPhones) {
                                vm.user.verifiedPhones[ctrl.label] = ctrl.phonenumber;
                            }
                            close();
                        }, function (response) {
                            if (response.status === 422) {
                                $scope.form.smscode.$setValidity('invalid_code', false);
                            }
                        });
                }

                function resetValidation() {
                    $scope.form.smscode.$setValidity('invalid_code', true);
                }
            }
        }

        function addAddress(ev) {
            return $q(function (resolve, reject) {
                $mdDialog.show({
                    controller: genericDetailControllerParams,
                    templateUrl: 'components/user/views/addressdialog.html',
                    targetEvent: ev,
                    fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                    locals: {
                        username: vm.username,
                        label: "",
                        data: {},
                        createFunction: UserService.registerNewAddress,
                        updateFunction: UserService.updateAddress,
                        deleteFunction: UserService.deleteAddress
                    }
                })
                    .then(
                        function (data) {
                            vm.user.address[data.newLabel] = data.data;
                            resolve(data);
                        }, reject);
            });
        }

        function addressDetail(ev, label, address) {
            $mdDialog.show({
                controller: genericDetailControllerParams,
                templateUrl: 'components/user/views/addressdialog.html',
                targetEvent: ev,
                fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                locals: {
                    username: vm.username,
                    label: label,
                    data: address,
                    createFunction: UserService.registerNewAddress,
                    updateFunction: UserService.updateAddress,
                    deleteFunction: UserService.deleteAddress
                }
            })
                .then(
                    function (data) {
                        if (data.newLabel) {
                            vm.user.address[data.newLabel] = data.data;
                        }
                        if (!data.newLabel || data.newLabel != data.originalLabel) {
                            delete vm.user.address[data.originalLabel];
                        }
                    });
        }

        function bankAccount(ev, label, bank) {
            return $q(function (resolve, reject) {
                $mdDialog.show({
                    controller: genericDetailControllerParams,
                    templateUrl: 'components/user/views/bankAccountDialog.html',
                    targetEvent: ev,
                    fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                    locals: {
                        username: vm.username,
                        label: label,
                        data: bank,
                        createFunction: UserService.registerNewBankAccount,
                        updateFunction: UserService.updateBankAccount,
                        deleteFunction: UserService.deleteBankAccount
                    }
                })
                    .then(
                        function (data) {
                            if (data.originalLabel || !data.newLabel) {
                                delete vm.user.bank[data.originalLabel];
                            }
                            if (data.newLabel) {
                                vm.user.bank[data.newLabel] = data.data;
                            }
                            resolve(data);
                        }, reject);
            });
        }

        function addFacebook() {
            configService.getConfig(function (config) {
                $window.location.href = 'https://www.facebook.com/dialog/oauth?client_id='
                    + config.facebookclientid
                    + '&response_type=code&redirect_uri='
                    + $window.location.origin
                    + '/facebook_callback';
            });
        }

        function facebook(ev) {
            $mdDialog.show({
                controller: genericDetailControllerParams,
                templateUrl: 'components/user/views/facebookDialog.html',
                targetEvent: ev,
                fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                locals: {
                    username: vm.username,
                    label: "",
                    data: vm.user.facebook,
                    createFunction: doNothing,
                    updateFunction: doNothing,
                    deleteFunction: UserService.deleteFacebookAccount
                }
            })
                .then(
                    function () {
                        vm.user.facebook = {};
                    });
        }

        function github(ev) {
            $mdDialog.show({
                controller: genericDetailControllerParams,
                templateUrl: 'components/user/views/githubDialog.html',
                targetEvent: ev,
                fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                locals: {
                    username: vm.username,
                    label: "",
                    data: vm.user.github,
                    createFunction: doNothing,
                    updateFunction: doNothing,
                    deleteFunction: UserService.deleteGithubAccount
                }
            })
                .then(
                    function () {
                        vm.user.github = {};
                    });
        }

        function addGithub() {
            configService.getConfig(function (config) {
                $window.location.href = 'https://github.com/login/oauth/authorize/?client_id=' + config.githubclientid;
            });
        }

        /**
         *
         * @param message
         * @param title
         * @param closeText
         * @returns promise
         */
        function showSimpleDialog(message, title, closeText) {
            title = title || 'Alert';
            closeText = closeText || 'Close';
            return $mdDialog.show(
                $mdDialog.alert({
                    title: title,
                    htmlContent: message,
                    ok: closeText
                })
            );
        }

        function GenericDetailDialogController($scope, $mdDialog, username, $window, label, data, createFunction, updateFunction, deleteFunction) {

            $scope.data = data;

            $scope.originalLabel = label;
            $scope.label = label;
            $scope.username = username;

            $scope.cancel = cancel;
            $scope.validationerrors = {};
            $scope.create = create;
            $scope.update = update;
            $scope.remove = remove;

            function cancel() {
                $mdDialog.cancel();
            }

            function create(label, data) {
                if (Object.keys($scope.dataform.$error).length > 0) {
                    return;
                }
                $scope.validationerrors = {};
                createFunction(username, label, data).then(
                    function (response) {
                        $mdDialog.hide({originalLabel: "", newLabel: label, data: data});
                    },
                    function (reason) {
                        if (reason.status == 409) {
                            $scope.validationerrors.duplicate = true;
                        }
                        else {
                            $window.location.href = "error" + reason.status;
                        }
                    }
                );
            }

            function update(oldLabel, newLabel, data) {
                if (Object.keys($scope.dataform.$error).length > 0) {
                    return;
                }
                $scope.validationerrors = {};
                updateFunction(username, oldLabel, newLabel, data).then(
                    function (response) {
                        $mdDialog.hide({originalLabel: oldLabel, newLabel: newLabel, data: data});
                    },
                    function (response) {
                        if (response.data && response.data.error) {
                            $scope.validationerrors[response.data.error] = true;
                        }
                        if (response.status == 409) {
                            $scope.validationerrors.duplicate = true;
                        }
                        else {
                            $window.location.href = "error" + response.status;
                        }
                    }
                );
            }

            function remove(label) {
                $scope.validationerrors = {};
                deleteFunction(username, label).then(
                    function (response) {
                        $mdDialog.hide({originalLabel: label, newLabel: ""});
                    },
                    function (response) {
                        if (response) {
                            $window.location.href = "error" + response.status;
                        }
                    }
                );
            }

        }

        function EmailDetailDialogController($scope, $mdDialog, $window, UserService, username, label, emailaddress, deleteIsPossible) {
            //If there is an emailaddress, it is already saved, if not, this means that a new one is being registered.
            $scope.emailaddress = emailaddress;
            $scope.deleteIsPossible = deleteIsPossible;

            $scope.originalLabel = label;
            $scope.label = label;
            $scope.username = username;

            $scope.cancel = cancel;
            $scope.validationerrors = {};
            $scope.create = create;
            $scope.update = update;
            $scope.deleteEmailAddress = deleteEmailAddress;


            function cancel() {
                $mdDialog.cancel();
            }

            function create(label, emailaddress) {
                if (Object.keys($scope.emailaddressform.$error).length > 0) {
                    return;
                }
                $scope.validationerrors = {};
                UserService.registerNewEmailAddress(username, label, emailaddress).then(
                    function (data) {
                        $mdDialog.hide({originalLabel: "", newLabel: label, emailaddress: emailaddress});
                    },
                    function (reason) {
                        if (reason.status == 409) {
                            $scope.validationerrors.duplicate = true;
                        }
                        else {
                            $window.location.href = "error" + reason.status;
                        }
                    }
                );
            }

            function update(oldLabel, newLabel, emailaddress) {
                if (Object.keys($scope.emailaddressform.$error).length > 0) {
                    return;
                }
                $scope.validationerrors = {};
                UserService.updateEmailAddress(username, oldLabel, newLabel, emailaddress).then(
                    function (data) {
                        $mdDialog.hide({originalLabel: oldLabel, newLabel: newLabel, emailaddress: emailaddress});
                    },
                    function (reason) {
                        if (reason.status == 409) {
                            $scope.validationerrors.duplicate = true;
                        }
                        else {
                            $window.location.href = "error" + reason.status;
                        }
                    }
                );
            }

            function deleteEmailAddress(label) {
                $scope.validationerrors = {};
                UserService.deleteEmailAddress(username, label).then(
                    function (data) {
                        $mdDialog.hide({originalLabel: label, newLabel: ""});
                    },
                    function (reason) {
                        $window.location.href = "error" + reason.status;
                    }
                );
            }

        }

    }
})();