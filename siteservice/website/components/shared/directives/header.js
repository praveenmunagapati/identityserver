(function () {
    'use strict';
    angular.module('itsyouonline.header', ['pascalprecht.translate'])
        .directive('itsYouOnlineHeader', ['$location', '$window', '$translate', function ($location, $window, $translate) {
            return {
                restrict: 'E',
                replace: true,
                templateUrl: 'components/shared/directives/header.html',
                link: function (scope, element, attr) {
                    scope.header_login = attr.register !== undefined;
                    scope.showCookieWarning = !localStorage.getItem('cookiewarning-dismissed');
                    scope.hideCookieWarning  = hideCookieWarning;
                    scope.updateLanguage = updateLanguage;
                    init();

                    function init() {
                        scope.langKey = localStorage.getItem('langKey');
                    }

                    function hideCookieWarning(){
                        localStorage.setItem('cookiewarning-dismissed', true);
                        scope.showCookieWarning = false;
                    }

                    function updateLanguage(){
                        localStorage.setItem('langKey', scope.langKey);
                        localStorage.setItem('selectedLangKey', scope.langKey)
                        $translate.use(scope.langKey);
                    }
                }
            };
        }]);
})();
