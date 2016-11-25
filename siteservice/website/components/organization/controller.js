(function() {
    'use strict';
    angular
        .module("itsyouonlineApp")
        .controller("OrganizationDetailController", OrganizationDetailController)
        .controller("InvitationDialogController", InvitationDialogController)
        .directive('customOnChange', customOnChange);

    InvitationDialogController.$inject = ['$scope', '$mdDialog', '$translate', 'organization', 'OrganizationService', 'UserDialogService'];
    OrganizationDetailController.$inject = ['$routeParams', '$window', '$translate', 'OrganizationService', '$mdDialog', '$mdMedia',
        '$rootScope', 'UserDialogService', 'UserService'];

    function OrganizationDetailController($routeParams, $window, $translate, OrganizationService, $mdDialog, $mdMedia, $rootScope,
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
        vm.showAddOrganizationDialog = showAddOrganizationDialog;
        vm.showAPIKeyCreationDialog = showAPIKeyCreationDialog;
        vm.showAPIKeyDialog = showAPIKeyDialog;
        vm.showDNSDialog = showDNSDialog;
        vm.getOrganizationDisplayname = getOrganizationDisplayname;
        vm.fetchInvitations = fetchInvitations;
        vm.fetchAPIKeyLabels = fetchAPIKeyLabels;
        vm.showCreateOrganizationDialog = UserDialogService.createOrganization;
        vm.showDeleteOrganizationDialog = showDeleteOrganizationDialog;
        vm.editMember = editMember;
        vm.editOrgMember = editOrgMember;
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
                    var pixelWidth = 200 + getBranchWidth(vm.organizationRoot.children[0]);
                    document.getElementById('treegraph').style.width = pixelWidth + 'px';
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

        function getBranchWidth(branch, rootDepth) {
            var splitted = branch.globalid.split(".")
            var length = splitted[splitted.length - 1].length * 6;
            var spacing = 0;
            if (branch.children.length > 1) {
                spacing = (branch.children.length - 1) * 80;
            }
            if (branch.children.length === 0) {
                return length;
            }
            else {
                var childWidth = spacing;
                for (var i = 0; i < branch.children.length; i++) {
                    childWidth += getBranchWidth(branch.children[i]);
                }
                return childWidth > length ? childWidth : length;
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

        function showAddOrganizationDialog(ev) {
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'));
            $mdDialog.show({
                controller: AddOrganizationDialogController,
                templateUrl: 'components/organization/views/addOrganizationMemberDialog.html',
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
                function(data) {
                  if (data.role === "member") {
                      if (vm.organization.orgmembers == null) {
                          vm.organization.orgmembers = [];
                      }
                      vm.organization.orgmembers.push(data.organization);
                  } else {
                      if (vm.organization.orgowners == null) {
                          vm.organization.orgowners = []
                      }
                      vm.organization.orgowners.push(data.organization);
                  }
                });
        }

        function showAPIKeyCreationDialog(ev) {
            var useFullScreen = ($mdMedia('sm') || $mdMedia('xs'));
            $mdDialog.show({
                controller: ['$scope', '$mdDialog', '$translate', 'organization', 'OrganizationService', 'label', APIKeyDialogController],
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
                controller: ['$scope', '$mdDialog', '$translate', 'organization', 'OrganizationService', 'label', APIKeyDialogController],
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
            $translate(['organization.controller.confirmdelete', 'organization.controller.deleteorg', 'organization.controller.deleteorganization',
                'organization.controller.haschildren', 'organization.controller.yes', 'organization.controller.no'], {organization: globalid}).then(function(translations){
                    var text = translations['organization.controller.confirmdelete'];
                    var confirm = $mdDialog.confirm()
                        .title(translations['organization.controller.deleteorg'])
                        .textContent(text)
                        .ariaLabel(translations['organization.controller.deleteorganization'])
                        .targetEvent(event)
                        .ok(translations['organization.controller.yes'])
                        .cancel(translations['organization.controller.no']);
                    $mdDialog.show(confirm).then(function () {
                        OrganizationService
                            .deleteOrganization(globalid)
                            .then(function () {
                                $window.location.hash = '#/';
                            }, function (response) {
                                if (response.status === 422) {
                                    var msg = translations['organization.controller.haschildren'];
                                    UserDialogService.showSimpleDialog(msg, 'Error', 'Ok', event);
                                }
                            });
                    });
                })
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
                controller: ['$mdDialog', '$translate', 'OrganizationService', 'UserDialogService', 'organization', 'user', 'initialRole', EditOrganizationMemberController],
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
                        vm.organization = data.data;
                    } else if (data.action === 'remove') {
                        var people = vm.organization[data.data.role];
                        people.splice(people.indexOf(data.data.username), 1);
                    }
                });
        }

        function editOrgMember(event, org) {
            var role = 'orgmembers';
            angular.forEach(['orgmembers', 'orgowners'], function (r) {
                if (vm.organization[r].indexOf(org) !== -1) {
                    role = r;
                }
            });
            var changeOrgRoleDialog = {
                controller: ['$mdDialog', '$translate', 'OrganizationService', 'UserDialogService', 'organization', 'org', 'initialRole', EditOrganizationMemberOrgController],
                controllerAs: 'ctrl',
                templateUrl: 'components/organization/views/changeOrganizationRoleDialog.html',
                targetEvent: event,
                fullscreen: $mdMedia('sm') || $mdMedia('xs'),
                locals: {
                    organization: vm.organization,
                    org: org,
                    initialRole: role
                }
            };

            $mdDialog
                .show(changeOrgRoleDialog)
                .then(function (data) {
                    if (data.action === 'edit') {
                        vm.organization = data.data;
                    } else if (data.action === 'remove') {
                        var collection;
                        if (data.data.role === 'orgmembers') {
                            collection = vm.organization.orgmembers;
                        } else {
                            collection = vm.organization.orgowners;
                        }
                        collection.splice(collection.indexOf(data.data.org), 1);
                    }
                });
        }

        function showLeaveOrganization(event) {
            $translate(['organization.controller.confirmleave', 'organization.controller.leaveorg', 'organization.controller.leaveorganization', 'organization.controller.yes',
                'organization.controller.no', 'organization.controller.notfound'], {organization: globalid}).then(function(translations){
                    var text = translations['organization.controller.confirmleave'];
                    var confirm = $mdDialog.confirm()
                        .title(translations['organization.controller.leaveorg'])
                        .textContent(text)
                        .ariaLabel(translations['organization.controller.leaveorganization'])
                        .targetEvent(event)
                        .ok(translations['organization.controller.yes'])
                        .cancel(translations['organization.controller.no']);
                    $mdDialog
                        .show(confirm)
                        .then(function () {
                            UserService
                                .leaveOrganization($rootScope.user, globalid)
                                .then(function () {
                                    $window.location.hash = '#/';
                                }, function (response) {
                                    if (response.status === 404) {
                                        UserDialogService.showSimpleDialog(translations['organization.controller.notfound'], 'Error', null, event);
                                    }
                                });
                        });
                })
        }
    }

    function getOrganizationDisplayname(globalid) {
        if (globalid) {
            var split = globalid.split('.');
            return split[split.length - 1];
        }
    }

    function InvitationDialogController($scope, $mdDialog, $translate, organization, OrganizationService, UserDialogService) {

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

    function AddOrganizationDialogController($scope, $mdDialog, organization, OrganizationService, UserDialogService) {

        $scope.role = "member";

        $scope.organization = organization;

        $scope.cancel = cancel;
        $scope.addOrganization = addOrganization;
        $scope.validationerrors = {};


        function cancel(){
            $mdDialog.cancel();
        }

        function addOrganization(searchString, role){
            $scope.validationerrors = {};
            OrganizationService.addOrganization(organization, searchString, role).then(
                function(data){
                    $mdDialog.hide({organization: searchString,role: role});
                },
                function(reason){
                    if (reason.status == 409){
                        $scope.validationerrors.duplicate = true;
                    }
                    else if (reason.status == 404){
                        $scope.validationerrors.nosuchorganization = true;
                    }
                }
            );

        }
    }

    function APIKeyDialogController($scope, $mdDialog, $translate, organization, OrganizationService, label) {
        //If there is a key, it is already saved, if not, this means that a new secret is being created.

        $scope.apikey = {secret: ""};

        if (label) {
            $translate(['organization.controller.loadingkey']).then(function(translations){
                $scope.secret = translations['organization.controller.loadingkey'];
                OrganizationService.getAPIKey(organization, label).then(
                    function(data){
                        $scope.apikey = data;
                    }
                );
            })
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

    function EditOrganizationMemberController($mdDialog, $translate, OrganizationService, UserDialogService, organization, user, initialRole) {
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
                    $translate(['organization.controller.cantchangerole']).then(function(translations){
                        UserDialogService.showSimpleDialog(translations['organization.controller.cantchangerole'], 'Error', 'ok', event);
                    })
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

    function EditOrganizationMemberOrgController($mdDialog, $translate, OrganizationService, UserDialogService, organization, org, initialRole) {
        var ctrl = this;
        ctrl.role = initialRole;
        ctrl.org = org;
        ctrl.organization = organization;
        ctrl.cancel = cancel;
        ctrl.submit = submit;
        ctrl.remove = remove;

        function cancel() {
            $mdDialog.cancel();
        }

        function submit() {
            OrganizationService
                .updateOrgMembership(organization.globalid, ctrl.org, ctrl.role)
                .then(function (data) {
                    $mdDialog.hide({action: 'edit', data: data});
                }, function () {
                    $translate(['organization.controller.cantchangerole']).then(function(translations){
                        UserDialogService.showSimpleDialog(translations['organization.controller.cantchangerole'], 'Error', 'ok', event);
                    })
                });
        }

        function remove() {
            OrganizationService
                .removeOrgMember(organization.globalid, org, initialRole)
                .then(function () {
                    $mdDialog.hide({action: 'remove', data: {role: initialRole, org: org}});
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
                makeFileDrop();
                if ($scope.logo !== "") {
                    var img = new Image()
                    img.src = $scope.logo;

                    var c = document.getElementById("logo-upload-preview");
                    var ctx = c.getContext("2d");
                    ctx.clearRect(0, 0, c.width, c.height);
                    ctx.drawImage(img, 0, 0);
                } else {
                    var c = document.getElementById("logo-upload-preview");
                    c.className += " dirty-background";
                }
            }
        );

        function makeFileDrop() {
            var target = document.getElementById("logo-upload-preview");
            target.addEventListener("dragover", function(e){e.preventDefault();}, true);
            target.addEventListener("drop", function(src){
	              src.preventDefault();
                //only allow image files, ignore others
                if (!src.dataTransfer.files[0].type.match(/image.*/)) {
                    return;
                }
                var reader = new FileReader();
	              reader.onload = function(e){
		                setFile(e.target.result);
	              };
	              reader.readAsDataURL(src.dataTransfer.files[0]);
            }, true);
        }

        $scope.uploadFile = function(event){
                var files = event.target.files;
                var url = URL.createObjectURL(files[0]);
                setFile(url);
            };


        function setFile(url) {
            var c = document.getElementById("logo-upload-preview");
            var ctx = c.getContext("2d");
            ctx.clearRect(0, 0, c.width, c.height);
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

                //check if the dirty-background css class is applied to the canvas
                if (c.className.match(/(?:^|\s)dirty-background(?!\S)/) ) {
                    //remove the dirty-background css class from the canvas
                    c.className = c.className.replace( /(?:^|\s)dirty-background(?!\S)/g , '' );
                }

                $scope.dataurl = c.toDataURL();
                // forces the update button after a file drop, migh fix safari issues?
                $scope.$digest();
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
