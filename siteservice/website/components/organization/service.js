(function() {
    'use strict';


    angular.module("itsyouonlineApp").service("OrganizationService", OrganizationService);


    OrganizationService.$inject = ['$http','$q'];

    function OrganizationService($http, $q) {
        var apiURL =  'api/organizations';

        var service = {
            create: create,
            get: get,
            invite: invite,
            getUserOrganizations: getUserOrganizations,
            getInvitations: getInvitations,
            createAPISecret: createAPISecret,
            deleteAPISecret: deleteAPISecret,
            updateAPISecretLabel: updateAPISecretLabel,
            getAPISecretLabels: getAPISecretLabels,
            getAPISecret: getAPISecret
        }

        return service;

        function create(name, dns, owner){
            var url = apiURL;
            return $http.post(url, {globalid:name,dns:dns,owners:[owner]}).then(
                function(response) {
                    return response.data;
                },
                function(reason){
                    return $q.reject(reason);
                }
            );
        }

        function get(globalid) {
            var url = apiURL + '/' + globalid;

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

        function invite(globalid, member, role) {
            var url = apiURL + '/' + globalid + '/' + role + 's';

            return $http
                .post(url, {username: member})
                .then(
                    function(response) {
                        return response.data;
                    },
                    function(reason) {
                        return $q.reject(reason);
                    }
                );
        }

        function getUserOrganizations(username) {
            var url = '/api/users/' + username + '/organizations';

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


        function getInvitations(globalid){
            var url = apiURL + '/' + globalid + '/invitations';

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

        function getAPISecretLabels(globalid){
            var url = apiURL + '/' + globalid + '/apisecrets';

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

        function createAPISecret(globalid, label){
            var url = apiURL + '/' + globalid + '/apisecrets';

            return $http
                .post(url, {label: label})
                .then(
                    function(response) {
                        return response.data;
                    },
                    function(reason) {
                        return $q.reject(reason);
                    }
                );
        }

        function updateAPISecretLabel(globalid, oldLabel, newLabel){
            var url = apiURL + '/' + globalid + '/apisecrets/' + oldLabel;

            return $http
                .put(url, {label: newLabel})
                .then(
                    function(response) {
                        return response.data;
                    },
                    function(reason) {
                        return $q.reject(reason);
                    }
                );
        }

        function deleteAPISecret(globalid, label){
            var url = apiURL + '/' + globalid + '/apisecrets/' + label;

            return $http
                .delete(url)
                .then(
                    function(response) {
                        return response.data;
                    },
                    function(reason) {
                        return $q.reject(reason);
                    }
                );
        }


        function getAPISecret(globalid, label){
            var url = apiURL + '/' + globalid + '/apisecrets/' + label;

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


    }
})();
