(function() {
    'use strict';


    angular.module("itsyouonlineApp").service("OrganizationService",OrganizationService);


    OrganizationService.$inject = ['$http','$q'];

    function OrganizationService($http, $q) {
        var apiURL =  'api/organizations';

        var service = {
            create: create
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

    }


})();
