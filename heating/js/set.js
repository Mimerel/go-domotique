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
    }, 5000);


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

                if (device.Status === "on") {
                    document.getElementById(theId).style.backgroundColor =  "#ADFF2F";
                } else if (device.Status === "off") {
                    document.getElementById(theId).style.backgroundColor =  "red";
                }
            }
            theId = "temperature_"+device.DomotiqueId;
            if (document.getElementById(theId) !== null) {
                document.getElementById(theId).innerText = device.Temperature === 0 ? "-" : device.Temperature + " °C" ;
            }
            theId = "deviceTemperature_"+device.DomotiqueId;
            if (document.getElementById(theId) !== null) {
                document.getElementById(theId).innerText = device.DeviceTemperature === 0 ? "-" : device.DeviceTemperature + " °C" ;
            }
            theId = "deviceOverTemperature_"+device.DomotiqueId;
            if (document.getElementById(theId) !== null) {
                document.getElementById(theId).innerText = device.DeviceOverTemperature === 0 ? "-" : device.DeviceOverTemperature + " °C" ;
            }
            theId = "temperatureStatus_"+device.DomotiqueId;
            if (document.getElementById(theId) !== null) {
                document.getElementById(theId).innerText = device.TemperatureStatus === 0 ? "-" : device.TemperatureStatus + " °C" ;
            }
            theId = "voltage_"+device.DomotiqueId;
            if (document.getElementById(theId) !== null) {
                document.getElementById(theId).innerText = device.Voltage === 0 ? "-" : device.Voltage + " " ;
            }
            theId = "temperatureTarget_"+device.DomotiqueId;
            if (document.getElementById(theId) !== null) {
                document.getElementById(theId).innerText = device.TemperatureTarget === 0 ? "-" : device.TemperatureTarget + " °C" ;
            }
            theId = "position_"+device.DomotiqueId;
            if (document.getElementById(theId) !== null) {
                document.getElementById(theId).innerText = device.CurrentPos === 0 ? "??" : device.CurrentPos + "" ;
                console.log("device pos", device.CurrentPos);
            }
            theId = "lastDirection_"+device.DomotiqueId;
            if (document.getElementById(theId) !== null) {
                document.getElementById(theId).innerText = device.LastDirection === 0 ? "??" : device.LastDirection + "" ;
                console.log("last direction", device.LastDirection);
            }
            theId = "battery"+device.DomotiqueId;
            if (document.getElementById(theId) !== null) {
                document.getElementById(theId).innerText = device.Battery === 0 ? "??" : device.Battery + " %" ;
                console.log("Battery", device.Battery);
            }

        });
        total = Math.round(total * 100) / 100;
        document.getElementById("totalpower").innerText = total + " Watts";
    });
    //snackbar("updated ...");
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
    $.get(url, function (data, status) {
        console.log(data);
    });
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

function runActionValueChange(id, action, value,  payload) {
    const newvalue = Number(value)+Number(payload);
    const url = URLAction+ "?id="+id+"&action="+action+"&payload="+newvalue;
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
