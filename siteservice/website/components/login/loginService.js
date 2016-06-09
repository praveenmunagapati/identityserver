(function () {
    'use strict';

    angular
        .module("loginApp")
        .service("LoginService", LoginService);

    LoginService.$inject = ['$http', '$q'];


    function LoginService($http, $q) {
        var apiURL = '/login';

        return {
            getTwoFactorAuthenticationMethods: getTwoFactorAuthenticationMethods,
            sendSmsCode: sendSmsCode,
            submitTotpCode: submitTotpCode,
            submitSmsCode: submitSmsCode,
            checkSmsConfirmation: checkSmsConfirmation
        };

        function genericHttpCall(httpFunction, url, data) {
            if (data) {
                return httpFunction(url, data)
                    .then(
                        function (response) {
                            return response.data;
                        },
                        function (reason) {
                            return $q.reject(reason);
                        }
                    );
            }
            else {
                return httpFunction(url)
                    .then(
                        function (response) {
                            return response.data;
                        },
                        function (reason) {
                            return $q.reject(reason);
                        }
                    );
            }
        }

        function getTwoFactorAuthenticationMethods() {
            var url = apiURL + '/twofamethods';
            return genericHttpCall($http.get, url);
        }

        function sendSmsCode(phoneLabel) {
            var url = apiURL + '/smscode/' + encodeURIComponent(phoneLabel);
            return genericHttpCall($http.post, url);
        }

        function submitTotpCode(code) {
            var url = apiURL + '/totpconfirmation';
            var data = {
                totpcode: code
            };
            return genericHttpCall($http.post, url, data);
        }

        function submitSmsCode(code) {
            var url = apiURL + '/smsconfirmation';
            var data = {
                smscode: code
            };
            return genericHttpCall($http.post, url, data);
        }

        function checkSmsConfirmation() {
            var url = apiURL + '/smsconfirmed';
            return genericHttpCall($http.get, url);
        }
    }
})();
