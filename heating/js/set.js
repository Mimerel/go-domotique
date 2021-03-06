function setTemporary(type) {
    var valueDay = document.getElementById('day').value !== "" ? parseFloat(document.getElementById('day').value) : 0;
    var valueHour = document.getElementById('hour').value !== "" ? parseFloat(document.getElementById('hour').value) : 0;
    valueHour = valueHour + valueDay * 24;
    ipPort = document.getElementById('ipPort').value;
    url = 'http://' + ipPort + '/heating/temporary/' + type + '/' + valueHour;
    window.location.assign(url);
}

function toggleDevice(id, url) {
    console.log(url);
    $.get(url, function (data, status) {
        console.log(data);
    });
    snackbar("Done");
}

function slideDevice(id, url) {
    console.log(url);
    url = url.substr(0, url.length - 1);
    console.log(url);
    console.log("slider" + id);
    var slider = document.getElementById("slider" + id).value;
    console.log(slider);
    url += slider;
    $.get(url, function (data, status) {
        console.log(data);
    });
    snackbar("Done");
}

function snackbar(message, error) {
    'use strict';
    var snackbarContainer = document.querySelector('#demo-snackbar-example');
    var data = {
        message: message,
        timeout: 2000,
    };
    snackbarContainer.MaterialSnackbar.showSnackbar(data);
}
