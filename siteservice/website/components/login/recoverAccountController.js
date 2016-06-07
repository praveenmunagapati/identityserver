(function () {
    'use strict';
    angular.module('loginApp')
        .controller('resetPasswordController', ['$http', '$window', '$routeParams', resetPasswordController]);

    function resetPasswordController($http, $window, $routeParams) {
        var vm = this;
        vm.submit = submit;
        var code = $routeParams.code;

        function submit() {
            var data = {
                password: vm.password,
                code: code
            };
            $http
                .post('/login/resetpassword', data)
                .then(function (response) {
                        // redirect to login
                        $window.location.hash = '';
                    },
                    function (response) {
                        switch (response.status) {
                            default:
                                $window.location.href = 'error' + response.status;
                        }
                    }
                );
        }
    }
})();