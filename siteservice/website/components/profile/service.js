(function() {
    'use strict';

    angular
        .module("itsyouonlineApp")
        .service("UserService", UserService);

    UserService.$inject = ['$http','$q'];

    function UserService($http, $q) {
        var apiURL = 'api/users';

        var service = {
            get: get,
        }

        return service;

        function get(username) {
            var url = apiURL + '/' + username;

            return $http
                .get(url)
                .then(
                    function(response) {
                        return response.data;
                    },
                    function(reason) {
                        return $q.reject(reason);
                    }
                );
        }

        function update(username, user) {
            var url = apiURL + '/' + username;

            return $http
                .put(url, user)
                .then(
                    function(response) {
                        return response.data;
                    },
                    function(reason) {
                        return $q.reject(reason);
                    }
                );
        }
    }
})();
