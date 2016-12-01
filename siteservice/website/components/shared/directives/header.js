(function () {
    'use strict';
    angular.module('itsyouonline.header', ['pascalprecht.translate'])
        .directive('itsYouOnlineHeader', ['$location', '$translate', function ($location, $translate) {
            return {
                restrict: 'E',
                replace: true,
                templateUrl: 'components/shared/directives/header.html',
                link: function (scope, element, attr) {
                    scope.header_login = attr.register !== undefined;
                    scope.showCookieWarning = !localStorage.getItem('cookiewarning-dismissed');
                    scope.hideCookieWarning  = hideCookieWarning;
                    scope.updateLanguage = updateLanguage;
                    var supportedLangs = ["en", "nl"];
                    var defaultLang = "en";
                    init();

                    function init() {
                        // selectedLangKey is the language key that has explicitly been selected by the user
                        scope.langKey = localStorage.getItem('selectedLangKey');
                        // set the langKey, this is the sites language, to the selected language. if its null, it'll be overriden anyway
                        localStorage.setItem('langKey', scope.langKey);
                        // it the user hasn't set a language yet
                        if (!scope.langKey) {
                            var urlParams = $location.search();
                            var lang = urlParams["lang"];
                            // if a queryvalue 'lang' is set and within the supported languages use that
                            if (supportedLangs.indexOf(lang) > -1) {
                                localStorage.setItem('langKey', lang)
                                scope.langKey = lang;
                            } else {
                                localStorage.setItem('langKey', defaultLang)
                                scope.langKey = defaultLang
                            }
                        }
                        $translate.use(scope.langKey);
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
