(function () {
    'use strict';

    angular
        .module("itsyouonline.registration")
        .service("registrationService", ['$http', '$q', RegistrationService]);

    function RegistrationService($http, $q) {
        return {
            validateUsername: validateUsername,
            register: register
        };

        function validateUsername(username) {
            var options = {
                params: {
                    username: username
                }
            };
            return $http.get('/validateusername', options);
        }

        function register(twoFAMethod, login, email, password, totpcode, sms, redirectparams) {
            var url = '/register';
            var data = {
                twofamethod: twoFAMethod,
                login: login.trim(),
                email: email.trim(),
                password: password,
                totpcode: totpcode,
                phonenumber: sms,
                redirectparams: redirectparams
            };
            return $http.post(url, data);
        }
    }
})();
