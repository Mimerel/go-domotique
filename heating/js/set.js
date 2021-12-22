const URL = "http://192.168.222.55:9998";
const URLAction = URL + "/runAction";
const URLUpdate = URL + "/heating/updateValues";
const queryString = window.location.search;
var urlParams = new URLSearchParams(queryString);

window.onload = function exampleFunction() {
    var myParam = location.search.split('tab=')[1];
    console.log("active tab",myParam);

    if (myParam === undefined) {
        myParam = "lights";
    }
    console.log("active tab",myParam);
    changeActiveTabTo(myParam);
    updateValues();
};

setInterval(function(){
    //refresh();
    updateValues();
    }, 10000);


function updateValues() {
    $.get(URLUpdate, function (dataCollected, status) {
        var data = Array();
        var total = 0;
        data = dataCollected;
        data.forEach( device => {
            var roundPower = Math.round(device.Power * 100) / 100;
            total += device.Power;
            var theId = "power_"+device.DomotiqueId;
            if (document.getElementById(theId) !== null) {
                document.getElementById(theId).innerText = roundPower === 0 ? "-" : roundPower + " W" ;
            }
            theId = "temperature_"+device.DomotiqueId;
            if (document.getElementById(theId) !== null) {
                document.getElementById(theId).innerText = device.Temperature === 0 ? "-" : device.Temperature + " °" ;
            }
            theId = "status_"+device.DomotiqueId;
            if (document.getElementById(theId) !== null) {
                if (device.Status === "on") {
                    document.getElementById(theId).style.backgroundColor =  "#ADFF2F";
                } else {
                    document.getElementById(theId).style.backgroundColor =  "red";
                }
            }

        });
        total = Math.round(total * 100) / 100;
        document.getElementById("totalpower").innerText = total + " Watts";
    });
    snackbar("updated ...");
}

function refresh() {
    document.location.reload();
}

function changeActiveTabTo(newTab) {
    var tabs = document.getElementsByClassName("is-active");
    Array.from(tabs).forEach(tab => {
        console.log("removing tag",tab);
        tab.classList.remove("is-active");
    });
    tabs = document.getElementsByClassName(newTab);
    Array.from(tabs).forEach(tab => {
        console.log("adding tag",tab);
        tab.classList.add("is-active");
    });
    tabs = document.getElementsByClassName("is-active");
    Array.from(tabs).forEach(tab => {
        console.log("Show tag",tab);
    });
    var url = window.location.href;
    console.log("active url", url);
    if (location.search.includes("tab=")) {


    var myParam = location.search.split('tab=')[1];

    url = url.replace("tab="+myParam, "tab="+newTab);
    } else {
        url =  url + "?tab="+newTab ;
    }
    window.history.replaceState(null, null, url);
    console.log("active url", url);
}

function setTemporary(type) {
    var valueDay = document.getElementById('day').value !== "" ? parseFloat(document.getElementById('day').value) : 0;
    var valueHour = document.getElementById('hour').value !== "" ? parseFloat(document.getElementById('hour').value) : 0;
    valueHour = valueHour + valueDay * 24;
    ipPort = document.getElementById('ipPort').value;
    url = 'http://' + ipPort + '/heating/temporary/' + type + '/' + valueHour;
    console.log(url);
    refresh();
    //window.location.assign(url);
}

function toggleDevice(id, url) {
    console.log(url);
    $.get(url, function (data, status) {
        console.log(data);
    });
    snackbar("Done");
}

function runAction(id, action, payload) {
    const url = URLAction+ "?id="+id+"&action="+action+"&payload="+payload;
    console.log(url);
    $.get(url, function (data, status) {
        console.log(data);
    });
    snackbar("Done");
}

function runReconnect() {
    const url = URL+ "/reconnect";
    console.log(url);
    $.get(url, function (data, status) {
        console.log(data);
    });
    snackbar("Done");
}



function slideDevice(id) {
    var slider = document.getElementById("slider" + id).value;
    const action = "/roller/0/command/pos";
    const url = URLAction+ "?id="+id+"&action="+action+"&payload="+slider;
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
