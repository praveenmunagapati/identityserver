(function () {
    'use strict';

    angular
        .module("itsyouonline.registration")
        .service("registrationService", ['$http', RegistrationService]);

    function RegistrationService($http) {
        return {
            validateUsername: validateUsername,
            register: register,
            getLogo: getLogo
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

        function getLogo(globalid) {
            var url = '/api/organizations/' + encodeURIComponent(globalid) + '/logo';
            console.log(url);
            return $http.get(url).then(
                function (response) {
                    return response.data;
                },
                function (reason) {
                    return $q.reject(reason);
                }
            );
        }
    }
})();
