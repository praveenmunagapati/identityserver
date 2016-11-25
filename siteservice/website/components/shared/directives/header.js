(function () {
    'use strict';
    angular.module('itsyouonline.header', ['pascalprecht.translate'])
        .directive('itsYouOnlineHeader', ['$translate', function ($translate) {
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
                        if (!localStorage.getItem('langKey')) {
                            localStorage.setItem('langKey', "en");
                        }
                        var language = localStorage.getItem('langKey');
                        $translate.use(language);
                        scope.langKey = language;
                    }

                    function hideCookieWarning(){
                        localStorage.setItem('cookiewarning-dismissed', true);
                        scope.showCookieWarning = false;
                    }

                    function updateLanguage(){
                        localStorage.setItem("langKey", scope.langKey);
                        $translate.use(scope.langKey);
                    }
                }
            };
        }]);
})();
