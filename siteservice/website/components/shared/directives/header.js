(function () {
    'use strict';
    angular.module('itsyouonline.header', [])
        .directive('itsYouOnlineHeader', function () {
            return {
                restrict: 'E',
                replace: true,
                templateUrl: 'components/shared/directives/header.html',
                link: function (scope, element, attr) {
                    scope.header_login = attr.register !== undefined;
                }
            };
        });
})();