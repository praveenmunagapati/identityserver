(function () {
    'use strict';
    angular.module('itsyouonline.footer', [])
        .directive('itsYouOnlineFooter', function () {
            return {
                restrict: 'E',
                replace: true,
                templateUrl: 'components/shared/directives/footer.html'
            };
        });
})();