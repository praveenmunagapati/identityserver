(function() {
    'use strict';


    angular.module("itsyouonlineApp").service("OrganizationService", OrganizationService);


    OrganizationService.$inject = ['$http','$q'];

    function OrganizationService($http, $q) {
        var apiURL =  'api/organizations';
        var GET = $http.get;
        var POST = $http.post;
        var PUT = $http.put;
        var DELETE = $http.delete;

        return {
            create: create,
            get: get,
            invite: invite,
            getUserOrganizations: getUserOrganizations,
            getInvitations: getInvitations,
            createAPIKey: createAPIKey,
            deleteAPIKey: deleteAPIKey,
            updateAPIKey: updateAPIKey,
            getAPIKeyLabels: getAPIKeyLabels,
            getAPIKey: getAPIKey,
            getOrganizationTree: getOrganizationTree,
            createDNS: createDNS,
            updateDNS: updateDNS,
            deleteDNS: deleteDNS,
            deleteOrganization: deleteOrganization

        };

        function genericHttpCall(httpFunction, url, data) {
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

        function create(name, dns, owner, parentOrganization) {
            var url = apiURL;
            if (parentOrganization){
                url += '/' + encodeURIComponent(parentOrganization) + '/suborganizations';
                name = parentOrganization + '.' + name;
            }
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
            var url = apiURL + '/' + encodeURIComponent(globalid);

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

        function invite(globalid, searchString, role) {
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/' + encodeURIComponent(role) + 's';

            return $http
                .post(url, {searchstring: searchString})
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
            var url = '/api/users/' + encodeURIComponent(username) + '/organizations';

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
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/invitations';

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

        function getAPIKeyLabels(globalid){
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/apikeys';

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

        function createAPIKey(globalid, apiKey){
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/apikeys';

            return $http
                .post(url, apiKey)
                .then(
                    function(response) {
                        return response.data;
                    },
                    function(reason) {
                        return $q.reject(reason);
                    }
                );
        }

        function updateAPIKey(globalid, oldLabel, newLabel, apikey){
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/apikeys/' + encodeURIComponent(oldLabel);
            apikey.label = newLabel;
            return $http
                .put(url, apikey)
                .then(
                    function(response) {
                        return response.data;
                    },
                    function(reason) {
                        return $q.reject(reason);
                    }
                );
        }

        function deleteAPIKey(globalid, label){
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/apikeys/' + encodeURIComponent(label);

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

        function getAPIKey(globalid, label){
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/apikeys/' + encodeURIComponent(label);

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

        function getOrganizationTree(globalid) {
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/tree';
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

        function createDNS(globalid, dnsName) {
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/dns/' + encodeURIComponent(dnsName);

            return $http
                .post(url)
                .then(
                    function (response) {
                        return response.data;
                    },
                    function (reason) {
                        return $q.reject(reason);
                    }
                );
        }

        function updateDNS(globalid, oldDnsName, newDnsName) {
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/dns/' + encodeURIComponent(oldDnsName);

            return $http
                .put(url, {name: newDnsName})
                .then(
                    function (response) {
                        return response.data;
                    },
                    function (reason) {
                        return $q.reject(reason);
                    }
                );
        }

        function deleteDNS(globalid, dnsName) {
            var url = apiURL + '/' + encodeURIComponent(globalid) + '/dns/' + encodeURIComponent(dnsName);
            return genericHttpCall(DELETE, url);
        }

        function deleteOrganization(globalid) {
            var url = apiURL + '/' + encodeURIComponent(globalid);
            return genericHttpCall(DELETE, url);
        }
    }
})();
