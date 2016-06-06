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

        function register(twoFAMethod, login, email, password, totpcode, sms) {
            var url = '/register';
            var data = {
                twofamethod: twoFAMethod,
                login: login,
                email: email,
                password: password,
                totpcode: totpcode,
                phonenumber: sms
            };
            return $http.post('/register', data);
        }
    }
})();
