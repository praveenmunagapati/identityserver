/**
 * Created by lucas on 27/05/16.
 */
(function () {
    'use strict';
    angular.module('itsyouonlineApp')
        .directive('authorizationDetails', function () {
            return {
                restrict: 'AE',
                templateUrl: 'components/user/directives/authorizationDetails.html',
                link: function (scope, element, attr) {
                    scope.revoke = revoke;
                    scope.save = save;

                    function revoke(property, label) {
                        delete scope.authorizations[property][label];
                    }

                    function save() {
                        scope.authorizations.organizations = [];
                        angular.forEach(scope.requested.organizations, function (allowed, organization) {
                            if (allowed) {
                                scope.authorizations.organizations.push(organization);
                            }
                        });
                        // Filter unauthorized permission labels
                        angular.forEach(scope.authorizations, function (value, key) {
                            if (Array.isArray(value)) {
                                angular.forEach(value, function (val, i) {
                                    if (!val) {
                                        value.splice(value.indexOf(val), 1);
                                    }
                                });
                            } else if (typeof value === 'object') {
                                angular.forEach(value, function (val, k) {
                                    if (val === '') {
                                        delete value[k];
                                    }
                                });
                            }
                        });
                        scope.update();
                    }
                }
            };
        });
})();