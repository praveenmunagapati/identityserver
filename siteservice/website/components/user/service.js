(function() {
    'use strict';

    angular
        .module("itsyouonlineApp")
        .service("UserService", UserService)
        .service("NotificationService", NotificationService);

    UserService.$inject = ['$http','$q'];
    NotificationService.$inject = ['$http','$q'];

    function NotificationService($http, $q) {
        var apiURL = 'api/users';

        var service = {
            get: get,
            accept: accept,
            reject: reject
        }

        return service;

        function get(username) {
            var url = apiURL + '/' + username + '/notifications';

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

        function accept(invitation) {
            var url = apiURL + '/' + invitation.user + '/organizations/' + invitation.organization + '/roles/' + invitation.role ;

            return $http
                .post(url, invitation)
                .then(
                    function(response) {
                        return response.data;
                    },
                    function(reason) {
                        return $q.reject(reason);
                    }
                );
        }

        function reject(invitation) {
            var url = apiURL + '/' + invitation.user + '/organizations/' + invitation.organization + '/roles/' + invitation.role ;

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
    }

    function UserService($http, $q) {
        var apiURL = 'api/users';

        var service = {
            get: get,
            registerNewEmailAddress: registerNewEmailAddress,
            updateEmailAddress: updateEmailAddress,
            deleteEmailAddress: deleteEmailAddress,
            registerNewPhonenumber: registerNewPhonenumber,
            updatePhonenumber: updatePhonenumber,
            deletePhonenumber: deletePhonenumber,
            registerNewAddress: registerNewAddress,
            updateAddress: updateAddress,
            deleteAddress: deleteAddress

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

        function registerNewEmailAddress(username, label, emailaddress) {
            var url = apiURL + '/' + username + '/emailaddresses';

            return $http
                .post(url, {label: label, emailaddress: emailaddress})
                .then(
                    function(response) {
                        return response.data;
                    },
                    function(reason) {
                        return $q.reject(reason);
                    }
                );
        }

        function updateEmailAddress(username, oldlabel, newlabel, emailaddress) {
            var url = apiURL + '/' + username + '/emailaddresses/' + oldlabel ;

            return $http
                .put(url, {label: newlabel, emailaddress: emailaddress})
                .then(
                    function(response) {
                        return response.data;
                    },
                    function(reason) {
                        return $q.reject(reason);
                    }
                );
        }

        function deleteEmailAddress(username, label) {
            var url = apiURL + '/' + username + '/emailaddresses/' + label ;

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


        function registerNewPhonenumber(username, label, phonenumber) {
            var url = apiURL + '/' + username + '/phonenumbers';

            return $http
                .post(url, {label: label, phonenumber: phonenumber})
                .then(
                    function(response) {
                        return response.data;
                    },
                    function(reason) {
                        return $q.reject(reason);
                    }
                );
        }

        function updatePhonenumber(username, oldlabel, newlabel, phonenumber) {
            var url = apiURL + '/' + username + '/phonenumbers/' + oldlabel ;

            return $http
                .put(url, {label: newlabel, phonenumber: phonenumber})
                .then(
                    function(response) {
                        return response.data;
                    },
                    function(reason) {
                        return $q.reject(reason);
                    }
                );
        }

        function deletePhonenumber(username, label) {
            var url = apiURL + '/' + username + '/phonenumbers/' + label ;

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


        function registerNewAddress(username, label, address) {
            var url = apiURL + '/' + username + '/addresses';

            return $http
                .post(url, {label: label, address: address})
                .then(
                    function(response) {
                        return response.data;
                    },
                    function(reason) {
                        return $q.reject(reason);
                    }
                );
        }

        function updateAddress(username, oldlabel, newlabel, address) {
            var url = apiURL + '/' + username + '/addresses/' + oldlabel ;

            return $http
                .put(url, {label: newlabel, address: address})
                .then(
                    function(response) {
                        return response.data;
                    },
                    function(reason) {
                        return $q.reject(reason);
                    }
                );
        }

        function deleteAddress(username, label) {
            var url = apiURL + '/' + username + '/addresses/' + label ;

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





    }
})();
