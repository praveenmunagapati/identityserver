(function () {
    'use strict';
    angular
        .module('itsyouonlineApp', ['ngCookies', 'ngMaterial', 'ngRoute', 'ngMessages', 'ngSanitize',
            'monospaced.qrcode',
            'itsyouonline.shared', 'itsyouonline.header', 'itsyouonline.footer', 'itsyouonline.user',
            'itsyouonline.validation', 'itsyouonline.telinput', 'pascalprecht.translate'])
        .config(['$mdThemingProvider', themingConfig])
        .config(['$httpProvider', httpConfig])
        .config(['$routeProvider', routeConfig])
        .config(['$translateProvider', translateConfig])
        .config([init])
        .factory('authenticationInterceptor', ['$q', '$window', authenticationInterceptor])
        .directive('pagetitle', ['$rootScope', '$timeout', 'footerService', pagetitle])
        .run(['$route', '$cookies', '$rootScope', '$location', runFunction]);

    function themingConfig($mdThemingProvider) {
        $mdThemingProvider.definePalette('blueish', {
            '50': '#f7fbfd',
            '100': '#badeed',
            '200': '#8ec8e2',
            '300': '#55add3',
            '400': '#3ca1cd',
            '500': '#4d738a',
            '600': '#2a7ea3',
            '700': '#236b8a',
            '800': '#1d5872',
            '900': '#17455a',
            'A100': '#f7fbfd',
            'A200': '#badeed',
            'A400': '#3ca1cd',
            'A700': '#236b8a',
            'contrastDefaultColor': 'light',
            'contrastDarkColors': '50 100 200 300 400 A100 A200 A400'
        });
        $mdThemingProvider.theme('default')
            .primaryPalette('blueish');
        $mdThemingProvider.enableBrowserColor({
            palette: 'primary', // Default is 'primary', any basic material palette and extended palettes are available
            hue: '800' // Default is '800'
        });
    }

    function httpConfig($httpProvider) {
        $httpProvider.interceptors.push('authenticationInterceptor');
        //initialize get if not there
        if (!$httpProvider.defaults.headers.get) {
            $httpProvider.defaults.headers.get = {};
        }
        //disable IE ajax request caching
        $httpProvider.defaults.headers.get['If-Modified-Since'] = '0';
    }

    function authenticationInterceptor($q, $window) {
        return {
            'request': function (config) {
                return config || $q.when(config);
            },
            'response': function (response) {
                return response || $q.when(response);
            },

            'responseError': function (rejection) {
                if (rejection.status === 401 || rejection.status === 403 || rejection.status === 419) {
                    $window.location.href = '/login';
                } else if (rejection.status.toString().startsWith('5')) {
                    $window.location.href = 'error' + rejection.status;
                }
                return $q.reject(rejection);
            }
        };
    }

    function routeConfig($routeProvider) {
        $routeProvider
            .when('/authorize', {
                templateUrl: 'components/user/views/authorize.html',
                controller: 'AuthorizeController',
                controllerAs: 'vm',
                data: {
                    pageTitle: 'Authorize',
                    showFooter: false
                }
            })
            .when('/company/new', {
                templateUrl: 'components/company/views/new.html',
                controller: 'CompanyController',
                controllerAs: 'vm',
                data: {
                    pageTitle: 'New company'
                }
            })
            .when('/organization/:globalid', {
                templateUrl: 'components/organization/views/detail.html',
                controller: 'OrganizationDetailController',
                controllerAs: 'vm',
                data: {
                    pageTitle: 'Organization detail'
                }
            })
            .when('/profile', {
                templateUrl: 'components/user/views/profile.html',
                controller: 'UserHomeController',
                controllerAs: 'vm',
                data: {
                    pageTitle: 'Profile'
                }
            })
            .when('/notifications', {
                templateUrl: 'components/user/views/notifications.html',
                controller: 'UserHomeController',
                controllerAs: 'vm',
                data: {
                    pageTitle: 'Notifications'
                }
            })
            .when('/organizations', {
                templateUrl: 'components/user/views/organizations.html',
                controller: 'UserHomeController',
                controllerAs: 'vm',
                data: {
                    pageTitle: 'Organizations'
                }
            })
            .when('/authorizations', {
                templateUrl: 'components/user/views/authorizations.html',
                controller: 'UserHomeController',
                controllerAs: 'vm',
                data: {
                    pageTitle: 'Authorizations'
                }
            })
            .when('/settings', {
                templateUrl: 'components/user/views/settings.html',
                controller: 'UserHomeController',
                controllerAs: 'vm',
                data: {
                    pageTitle: 'Settings'
                }
            })
            .otherwise('/profile');
    }

    function translateConfig($translateProvider) {
        $translateProvider.useStaticFilesLoader({
            prefix: 'assets/i18n/',
            suffix: '.json'
        });
        $translateProvider.useSanitizeValueStrategy('sanitize');
        $translateProvider.useMissingTranslationHandlerLog();

        var supportedLangs = ["en", "nl"];
        var defaultLang = "en";

        // selectedLangKey is the language key that has explicitly been selected by the user
        var langKey = localStorage.getItem('selectedLangKey');
        // set the langKey, this is the sites language, to the selected language. if its null, it'll be overriden anyway
        localStorage.setItem('langKey', langKey);
        // it the user hasn't set a language yet
        if (!langKey) {
            var langParam = getParameterByName("lang");
            var lang = langParam || URI(window.location.href).search(true).lang;
            // if a queryvalue 'lang' is set and within the supported languages use that
            if (supportedLangs.indexOf(lang) > -1) {
                localStorage.setItem('langKey', lang);
                // Store the langkey requested through the url params
                localStorage.setItem('requestedLangKey', lang);
                langKey = lang;
            } else {
                var previousLang = localStorage.getItem('requestedLangKey');
                // if a language was set thourgh an URL in a previous request use that
                if (previousLang) {
                    localStorage.setItem('langKey', previousLang);
                    langKey = previousLang;
                } else {
                    //if all else fails just use English
                    localStorage.setItem('langKey', defaultLang);
                    langKey = defaultLang;
                }
            }
        }
        $translateProvider.use(langKey);
    }

    function init() {
        localStorage.setItem('hasLoggedIn', true);
    }

    function pagetitle($rootScope, $timeout, footerService) {
        return {
            link: function (scope, element) {
                var listener = function (event, current) {
                    var pageTitle = 'It\'s You Online';
                    var routeData = current.$$route && current.$$route.data || {};
                    if (routeData.pageTitle) {
                        pageTitle = current.$$route.data.pageTitle + ' - ' + pageTitle;
                    }
                    footerService.setFooter(routeData.showFooter !== undefined ? routeData.showFooter : true);
                    $timeout(function () {
                        element.text(pageTitle);
                    }, 0, false);
                };

                $rootScope.$on('$routeChangeSuccess', listener);
            }
        };
    }

    function runFunction($route, $cookies, $rootScope, $location) {
        $rootScope.user = $cookies.get('itsyou.online.user');
        var original = $location.path;
        // prevent controller reload when changing route params in code because we aren't using states
        $location.path = function (path, reload) {
            if (reload === false) {
                var lastRoute = $route.current;
                var un = $rootScope.$on('$locationChangeSuccess', function () {
                    $route.current = lastRoute;
                    un();
                });
            }
            return original.apply($location, [path]);
        };
        if (window.location.hostname === 'dev.itsyou.online') {
            setTimeout(function () {
                window.location.reload();
            }, 9 * 60 * 1000);
        }
        initializePolyfills();
    }

    function getParameterByName(name, url) {
        if (!url) {
              url = window.location.href;
        }
        name = name.replace(/[\[\]]/g, "\\$&");
        var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
            results = regex.exec(url);
        if (!results) return null;
        if (!results[2]) return '';
        return decodeURIComponent(results[2].replace(/\+/g, " "));
    }

})();
