var app = angular.module('myApp', []);
app.controller('myCtrl', function($scope, $http) {
    $scope.updateInterval = 30000;

    $scope.getStatus = function() {
        $http.get("/api/status")
            .then(function (response) {
                $scope.recentStatus = response.data;
            });
    };

    $scope.setUpdateInterval = function() {
        $http.get("/api/interval/" + $scope.updateInterval);
    };

    $scope.turnLight = function(flag) {
        flag = flag ? "on" : "off";
        $http.get("/api/light/" + flag);
    };

    $scope.conn = new WebSocket("ws://192.168.1.103:8080/ws/one");
    $scope.conn.onmessage = function (evt) {
        $scope.getStatus();
    };

    $scope.getStatus();
});
