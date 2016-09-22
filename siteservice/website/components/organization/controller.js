(function() {
    'use strict';
    angular
        .module("itsyouonlineApp")
        .controller("OrganizationDetailController", OrganizationDetailController)
        .controller("InvitationDialogController", InvitationDialogController)
        .directive('customOnChange', customOnChange);

    InvitationDialogController.$inject = ['$scope', '$mdDialog', 'organization', 'OrganizationService', 'UserDialogService'];
    OrganizationDetailController.$inject = ['$routeParams', '$window', 'OrganizationService', '$mdDialog', '$mdMedia',
        '$rootScope', 'UserDialogService', 'UserService'];

    function OrganizationDetailController($routeParams, $window, OrganizationService, $mdDialog, $mdMedia, $rootScope,
                                          UserDialogService, UserService) {
        var vm = this,
            globalid = $routeParams.globalid;
        vm.invitations = [];
        vm.apikeylabels = [];
        vm.organization = {};
        vm.organizationRoot = {};
        vm.childOrganizationNames = [];
        vm.logo = "";

        vm.initSettings = initSettings;
        vm.showInvitationDialog = showInvitationDialog;
        vm.showAPIKeyCreationDialog = showAPIKeyCreationDialog;
        vm.showAPIKeyDialog = showAPIKeyDialog;
        vm.showDNSDialog = showDNSDialog;
        vm.getOrganizationDisplayname = getOrganizationDisplayname;
        vm.fetchInvitations = fetchInvitations;
        vm.fetchAPIKeyLabels = fetchAPIKeyLabels;
        vm.showCreateOrganizationDialog = UserDialogService.createOrganization;
        vm.showDeleteOrganizationDialog = showDeleteOrganizationDialog;
        vm.editMember = editMember;
        vm.canEditRole = canEditRole;
        vm.showLeaveOrganization = showLeaveOrganization;
        vm.showLogoDialog = showLogoDialog;

        activate();

        function activate() {
            fetch();
        }

        function fetch(){
            OrganizationService
                .get(globalid)
                .then(
                    function(data) {
                        vm.organization = data;
                        vm.childOrganizationNames = getChildOrganizations(vm.organization.globalid);
                        vm.hasEditPermission = vm.organization.owners.indexOf($rootScope.user) !== -1;
                        fetchInvitations();
                    }
                );

            OrganizationService.getOrganizationTree(globalid)
                .then(function (data) {
                    vm.organizationRoot.children = [];
                    vm.organizationRoot.children.push(data);
                });

            OrganizationService.getLogo(globalid).then(
                function(data) {
                    vm.logo = data.logo;
                }
            );
        }

        function renderLogo() {
            if (vm.logo !== "") {
                var img = new Image();
                img.src = vm.logo;

                var c = document.getElementById("logoview");
                if (!c) {
                    console.log("aborting logo render - canvas not loaded");
                    return;
                }
                var ctx = c.getContext("2d");
                ctx.clearRect(0, 0, c.width, c.height);
                ctx.drawImage(img, 0, 0);
            }
        }

        function fetchInvitations() {
            if (!vm.hasEditPermission || vm.invitations.length) {
                return;
            }
            OrganizationService
                .getInvitations(globalid)
                .then(
                    function (data) {
                        vm.invitations = data;
                    }
                );
        }

        function fetchAPIKeyLabels(){
            if (!vm.hasEditPermission || vm.apikeylabels.length) {
                return;
            }
            OrganizationService
                .getAPIKeyLabels(globalid)
                .then(
                    function(data) {
                        vm.apikeylabels = data;
                    }
                );
        }

        function initSettings() {
            fetchAPIKeyLabels();
            renderLogo();
        }

        function showInvitationDialog(ev) {
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'));
            $mdDialog.show({
                controller: InvitationDialogController,
                templateUrl: 'components/organization/views/invitationdialog.html',
                targetEvent: ev,
                fullscreen: useFullScreen,
                locals:
                    {
                        OrganizationService: OrganizationService,
                        organization : vm.organization.globalid,
                        $window: $window
                    }
            })
            .then(
                function(invitation) {
                    vm.invitations.push(invitation);
                });
        }

        function showAPIKeyCreationDialog(ev) {
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'));
            $mdDialog.show({
                controller: ['$scope', '$mdDialog', 'organization', 'OrganizationService', 'label', APIKeyDialogController],
                templateUrl: 'components/organization/views/apikeydialog.html',
                targetEvent: ev,
                fullscreen: useFullScreen,
                locals:
                    {
                        OrganizationService: OrganizationService,
                        organization : vm.organization.globalid,
                        $window: $window,
                        label: ""
                    }
            })
            .then(
                function(data) {
                    if (data.newLabel) {
                        vm.apikeylabels.push(data.newLabel);
                    }
                });
        }

        function showAPIKeyDialog(ev, label) {
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'));
            $mdDialog.show({
                controller: ['$scope', '$mdDialog', 'organization', 'OrganizationService', 'label', APIKeyDialogController],
                templateUrl: 'components/organization/views/apikeydialog.html',
                targetEvent: ev,
                fullscreen: useFullScreen,
                locals:
                    {
                        OrganizationService: OrganizationService,
                        organization : vm.organization.globalid,
                        $window: $window,
                        label: label
                    }
            })
            .then(
                function(data) {
                    if (data.originalLabel != data.newLabel){
                        if (data.originalLabel) {
                            if (data.newLabel){
                                //replace
                                vm.apikeylabels[vm.apikeylabels.indexOf(data.originalLabel)] = data.newLabel;
                            }
                            else {
                                //remove
                                vm.apikeylabels.splice(vm.apikeylabels.indexOf(data.originalLabel),1);
                            }
                        }
                        else {
                            //add
                            vm.apikeylabels.push(data.newLabel);
                        }
                    }
                });
        }

        function showDNSDialog(ev, dnsName) {
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'));
            $mdDialog.show({
                controller: ['$scope', '$mdDialog', 'organization', 'OrganizationService', 'dnsName', DNSDialogController],
                templateUrl: 'components/organization/views/dnsDialog.html',
                targetEvent: ev,
                fullscreen: useFullScreen,
                locals: {
                    OrganizationService: OrganizationService,
                    organization: vm.organization.globalid,
                    $window: $window,
                    dnsName: dnsName
                }
            })
                .then(
                    function (data) {
                        if (data.originalDns) {
                            vm.organization.dns.splice(vm.organization.dns.indexOf(data.originalDns), 1);
                        }
                        if (data.newDns) {
                            vm.organization.dns.push(data.newDns);
                        }
                    });
        }

        function showLogoDialog(ev) {
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'));
            $mdDialog.show({
                controller: ['$scope', '$mdDialog', 'organization', 'OrganizationService', logoDialogController],
                templateUrl: 'components/organization/views/logoDialog.html',
                targetEvent: ev,
                fullscreen: useFullScreen,
                locals: {
                    OrganizationService: OrganizationService,
                    organization: vm.organization.globalid,
                    $window: $window
                }
            }).then(
                function() {
                    OrganizationService.getLogo(vm.organization.globalid).then(
                        function(data) {
                            vm.logo = data.logo;
                        }
                    ).then(
                        function() {
                            renderLogo();
                        }
                    );
                }
            );
        }

        function showDeleteOrganizationDialog(event) {
            var text = 'Are you sure you want to delete the organization "' + globalid + '"?';
            var confirm = $mdDialog.confirm()
                .title('Delete organization')
                .textContent(text)
                .ariaLabel('Delete organization ' + globalid)
                .targetEvent(event)
                .ok('Yes')
                .cancel('No');
            $mdDialog.show(confirm).then(function () {
                OrganizationService
                    .deleteOrganization(globalid)
                    .then(function () {
                        $window.location.hash = '#/';
                    }, function (response) {
                        if (response.status === 422) {
                            var msg = 'This organization cannot be deleted because it still has child organizations.';
                            UserDialogService.showSimpleDialog(msg, 'Error', 'Ok', event);
                        }
                    });
            });
        }

        function canEditRole(member) {
            return vm.organization.owners.indexOf($rootScope.user) > -1 && member !== $rootScope.user;
        }

        function editMember(event, user) {
            var role = 'members';
            angular.forEach(['members', 'owners'], function (r) {
                if (vm.organization[r].indexOf(user) !== -1) {
                    role = r;
                }
            });
            var changeRoleDialog = {
                controller: ['$mdDialog', 'OrganizationService', 'UserDialogService', 'organization', 'user', 'initialRole', EditOrganizationMemberController],
                controllerAs: 'ctrl',
                templateUrl: 'components/organization/views/changeRoleDialog.html',
                targetEvent: event,
                fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                locals: {
                    organization: vm.organization,
                    user: user,
                    initialRole: role
                }
            };

            $mdDialog
                .show(changeRoleDialog)
                .then(function (data) {
                    if (data.action === 'edit') {
                        vm.organization = data;
                    } else if (data.action === 'remove') {
                        var people = vm.organization[data.data.role];
                        people.splice(people.indexOf(data.data.username), 1);
                    }
                });
        }

        function showLeaveOrganization(event) {
            var text = 'Are you sure you want to leave the organization "' + globalid + '"?';
            var confirm = $mdDialog.confirm()
                .title('Leave organization')
                .textContent(text)
                .ariaLabel('Leave organization ' + globalid)
                .targetEvent(event)
                .ok('Yes')
                .cancel('No');
            $mdDialog
                .show(confirm)
                .then(function () {
                    UserService
                        .leaveOrganization($rootScope.user, globalid)
                        .then(function () {
                            $window.location.hash = '#/';
                        }, function (response) {
                            if (response.status === 404) {
                                UserDialogService.showSimpleDialog('User or organization not found', 'Error', null, event);
                            }
                        });
                });
        }
    }

    function getOrganizationDisplayname(globalid) {
        if (globalid) {
            var split = globalid.split('.');
            return split[split.length - 1];
        }
    }

    function InvitationDialogController($scope, $mdDialog, organization, OrganizationService, UserDialogService) {

        $scope.role = "member";

        $scope.cancel = cancel;
        $scope.invite = invite;
        $scope.validationerrors = {};


        function cancel(){
            $mdDialog.cancel();
        }

        function invite(searchString, role){
            $scope.validationerrors = {};
            OrganizationService.invite(organization, searchString, role).then(
                function(data){
                    $mdDialog.hide(data);
                },
                function(reason){
                    if (reason.status == 409){
                        $scope.validationerrors.duplicate = true;
                    }
                    else if (reason.status == 404){
                        $scope.validationerrors.nosuchuser = true;
                    } else if (reason.status === 422) {
                        cancel();
                        var msg = 'Organization ' + organization + ' has reached the maximum amount of invitations.';
                        UserDialogService.showSimpleDialog(msg, 'Error');
                    }
                }
            );

        }
    }

    function APIKeyDialogController($scope, $mdDialog, organization, OrganizationService, label) {
        //If there is a key, it is already saved, if not, this means that a new secret is being created.

        $scope.apikey = {secret: ""};

        if (label) {
            $scope.secret = "-- Loading --";
            OrganizationService.getAPIKey(organization, label).then(
                function(data){
                    $scope.apikey = data;
                }
            );
        }

        $scope.originalLabel = label;
        $scope.savedLabel = label;
        $scope.label = label;
        $scope.organization = organization;

        $scope.cancel = cancel;
        $scope.validationerrors = {};
        $scope.create = create;
        $scope.update = update;
        $scope.deleteAPIKey = deleteAPIKey;

        $scope.modified = false;


        function cancel(){
            if ($scope.modified) {
                $mdDialog.hide({originalLabel: label, newLabel: $scope.label});
            }
            else {
                $mdDialog.cancel();
            }
        }

        function create(label, apiKey){
            $scope.validationerrors = {};
            apiKey.label = label;
            OrganizationService.createAPIKey(organization, apiKey).then(
                function(data){
                    $scope.modified = true;
                    $scope.apikey = data;
                    $scope.savedLabel = data.label;
                },
                function(reason){
                    if (reason.status === 409) {
                        $scope.validationerrors.duplicate = true;
                    }
                }
            );
        }

        function update(oldLabel, newLabel){
            $scope.validationerrors = {};
            OrganizationService.updateAPIKey(organization, oldLabel, newLabel, $scope.apikey).then(
                function () {
                    $mdDialog.hide({originalLabel: oldLabel, newLabel: newLabel});
                },
                function(reason){
                    if (reason.status === 409) {
                        $scope.validationerrors.duplicate = true;
                    }
                }
            );
        }


        function deleteAPIKey(label){
            $scope.validationerrors = {};
            OrganizationService.deleteAPIKey(organization, label).then(
                function () {
                    $mdDialog.hide({originalLabel: label, newLabel: ""});
                }
            );
        }

    }

    function DNSDialogController($scope, $mdDialog, organization, OrganizationService, dnsName) {
        $scope.organization = organization;
        $scope.dnsName = dnsName;
        $scope.newDnsName = dnsName;

        $scope.cancel = cancel;
        $scope.validationerrors = {};
        $scope.create = create;
        $scope.update = update;
        $scope.remove = remove;

        function cancel() {
            $mdDialog.cancel();
        }

        function create(dnsName) {
            if (!$scope.form.$valid) {
                return;
            }
            $scope.validationerrors = {};
            OrganizationService.createDNS(organization, dnsName).then(
                function (data) {
                    $mdDialog.hide({originalDns: "", newDns: data.name});
                },
                function (reason) {
                    if (reason.status === 409) {
                        $scope.validationerrors.duplicate = true;
                    }
                }
            );
        }

        function update(oldDns, newDns) {
            if (!$scope.form.$valid) {
                return;
            }
            $scope.validationerrors = {};
            OrganizationService.updateDNS(organization, oldDns, newDns).then(
                function (data) {
                    $mdDialog.hide({originalDns: oldDns, newDns: data.name});
                },
                function (reason) {
                    if (reason.status === 409) {
                        $scope.validationerrors.duplicate = true;
                    }
                }
            );
        }


        function remove(dnsName) {
            $scope.validationerrors = {};
            OrganizationService.deleteDNS(organization, dnsName)
                .then(function () {
                    $mdDialog.hide({originalDns: dnsName, newDns: ""});
                });
        }
    }

    function getChildOrganizations(organization) {
        var children = [];
        if (organization) {
            for (var i = 0, splitted = organization.split('.'); i < splitted.length; i++) {
                var parents = splitted.slice(0, i + 1);
                children.push({
                    name: splitted[i],
                    url: '#/organization/' + parents.join('.')
                });
            }
        }
        return children;

    }

    function EditOrganizationMemberController($mdDialog, OrganizationService, UserDialogService, organization, user, initialRole) {
        var ctrl = this;
        ctrl.role = initialRole;
        ctrl.user = user;
        ctrl.organization = organization;
        ctrl.cancel = cancel;
        ctrl.submit = submit;
        ctrl.remove = remove;

        function cancel() {
            $mdDialog.cancel();
        }

        function submit() {
            OrganizationService
                .updateMembership(organization.globalid, ctrl.user, ctrl.role)
                .then(function (data) {
                    $mdDialog.hide({action: 'edit', data: data});
                }, function () {
                    UserDialogService.showSimpleDialog('Could not change role, please try again later', 'Error', 'ok', event);
                });
        }

        function remove() {
            OrganizationService
                .removeMember(organization.globalid, user, initialRole)
                .then(function () {
                    $mdDialog.hide({action: 'remove', data: {role: initialRole, user: user}});
                }, function (response) {
                    $mdDialog.cancel(response);
                });
        }
    }

    function logoDialogController($scope, $mdDialog, organization, OrganizationService) {
        $scope.organization = organization;
        $scope.setFile = setFile;
        $scope.cancel = cancel;
        $scope.validationerrors = {};
        $scope.update = update;
        $scope.remove = remove;

        OrganizationService.getLogo(organization).then(
            function(data) {
                $scope.logo = data.logo;
                if ($scope.logo !== "") {
                    var img = new Image()
                    img.src = $scope.logo;

                    var c = document.getElementById("preview");
                    var ctx = c.getContext("2d");
                    ctx.clearRect(0, 0, c.width, c.height);
                    ctx.drawImage(img, 0, 0);
                }
            }
        );

        $scope.uploadFile = function(event){
                var files = event.target.files;
                setFile(files);
            };

        function setFile(element) {
            //$scope.logo = element[0];
            var c = document.getElementById("preview");
            var ctx = c.getContext("2d");
            ctx.clearRect(0, 0, c.width, c.height);
            var upload = document.getElementById("fileToUpload");
            var url = URL.createObjectURL(element[0]);
            var img = new Image();
            img.src = url;

            img.onload = function() {
                var wscale = 1;
                if (img.width > c.width) {
                    wscale = c.width / img.width;
                }
                var hscale = 1;
                if (img.height > c.height) {
                    hscale = c.height / img.height;
                }
                var canvasCopy = document.createElement("canvas");
                var copyContext = canvasCopy.getContext("2d");

                canvasCopy.width = img.width;
                canvasCopy.height = img.height;
                copyContext.drawImage(img, 0, 0);

                var ratio = (wscale <= hscale ? wscale : hscale);

                var widthOffset = (c.width - img.width * ratio) / 2;
                var heightOffset = (c.height - img.height * ratio) / 2;
                ctx.drawImage(canvasCopy, widthOffset, heightOffset, canvasCopy.width * ratio, canvasCopy.height * ratio);

                $scope.dataurl = c.toDataURL();
            }

        }

        function cancel() {
            $mdDialog.cancel();
        }

        function update(logo) {
            if (!$scope.form.$valid) {
                return;
            }
            $scope.validationerrors = {};
            OrganizationService.setLogo(organization, logo).then(
                function (data) {
                    $mdDialog.hide({logo: logo});
                },
                function (reason) {
                    if (reason.status === 413) {
                        $scope.validationerrors.filesize = true;
                    }
                }
            );
        }


        function remove() {
            $scope.validationerrors = {};
            OrganizationService.deleteLogo(organization)
                .then(function () {
                    $mdDialog.hide({logo: ""});
                });
        }
    }

    function customOnChange() {
        return {
            restrict: 'A',
            link: function (scope, element, attrs) {
                var onChangeHandler = scope.$eval(attrs.customOnChange);
                element.bind('change', onChangeHandler);
            }
        };
    }

})();
