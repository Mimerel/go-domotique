const URL = "http://192.168.222.55:9998";
const URLAction = URL + "/runAction";
const URLUpdate = URL + "/heating/updateValues";
const queryString = window.location.search;
var urlParams = new URLSearchParams(queryString);

window.onload = function exampleFunction() {
    var myParam = location.search.split('tab=')[1];
    //console.log("active tab",myParam);

    if (myParam === undefined) {
        myParam = "lights";
    }
    //console.log("active tab",myParam);
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
        var subTotals = new Map();
        data = dataCollected;
        //console.log(data);
        data.forEach( device => {

            var roundPower = Math.round(device.Power * 100) / 100;
            if (device.Room !== "") {
                temp = subTotals.get(device.Room);
                if (temp === undefined) {
                    temp = 0;
                }
                subTotals.set(device.Room, Math.round((temp+roundPower) * 100) / 100);
                if (document.getElementById(device.Room+"W") !== null) {
                    document.getElementById(device.Room+"W").innerText = subTotals.get(device.Room) + " W";
                }
            }
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
            var theId = "status_"+device.DomotiqueId;
            if (document.getElementById(theId) !== null) {
                if (device.Status === "on") {
                    document.getElementById(theId).style.backgroundColor =  "red";
                    document.getElementById(theId).innerText = "Enabled";
                } else if (device.Status === "off") {
                    document.getElementById(theId).style.backgroundColor =  "#ADFF2F";
                    document.getElementById(theId).innerText = "Disabled";
                }
            }

            theId = "temperature_"+device.DomotiqueId;
            if (document.getElementById(theId) !== null) {
                document.getElementById(theId).innerText = device.Temperature === 0 ? "-" : device.Temperature + " °C" ;
            }
            theId = "virbration_"+device.DomotiqueId;
            if (document.getElementById(theId) !== null) {
                document.getElementById(theId).innerText = device.Vibration === 0 ? "-" : device.Vibration + "" ;
            }
            theId = "tilt_"+device.DomotiqueId;
            if (document.getElementById(theId) !== null) {
                document.getElementById(theId).innerText = device.Tilt === -1 ? "-" : device.Tilt + "" ;
            }
            theId = "lux_"+device.DomotiqueId;
            if (document.getElementById(theId) !== null) {
                document.getElementById(theId).innerText = device.Lux === 0 ? "-" : device.Lux + "" ;
            }
            theId = "illumination_"+device.DomotiqueId;
            if (document.getElementById(theId) !== null) {
                document.getElementById(theId).innerText = device.Illumination === 0 ? "-" : device.Illumination + "" ;
            }
            theId = "state_"+device.DomotiqueId;
            if (document.getElementById(theId) !== null) {
                document.getElementById(theId).innerText = device.State === 0 ? "-" : device.State + "" ;
                if (device.State === "close") {
                    document.getElementById(theId).style.backgroundColor =  "#ADFF2F";
                } else if (device.State === "open") {
                    document.getElementById(theId).style.backgroundColor =  "red";
                }
            }
            theId = "humidity_"+device.DomotiqueId;
            if (document.getElementById(theId) !== null) {
                document.getElementById(theId).innerText = device.Humidity === 0 ? "-" : device.Humidity + " %" ;
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
                document.getElementById(theId).innerText = "Temp. Status : "+ device.TemperatureStatus  ;
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
            }
            theId = "lastDirection_"+device.DomotiqueId;
            if (document.getElementById(theId) !== null) {
                document.getElementById(theId).innerText = device.LastDirection === 0 ? "??" : device.LastDirection + "" ;
            }
            theId = "battery_"+device.DomotiqueId;
            if (document.getElementById(theId) !== null) {
                document.getElementById(theId).innerText = "Battery : " +device.Battery + " %" ;
            }
            theId = "active_"+device.DomotiqueId;
            if (document.getElementById(theId) !== null) {
                document.getElementById(theId).innerText = device.Active === false ? "Device is not active" : " Device is Active" ;
            }
            theId = "online_"+device.DomotiqueId;
            if (document.getElementById(theId) !== null) {
                document.getElementById(theId).innerText = device.Online === true ? "Device Online" : " Device OffLine" ;
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

    var tabtrs = document.getElementsByClassName(" room");
    Array.from(tabtrs).forEach(tab => {
        tab.style.display = "none";
        if (tab.classList.contains("room"+newTab)) {
            tab.style.display = "";
        }
    });

    var tabbuttons = document.getElementsByClassName(" tabButton");
    Array.from(tabbuttons).forEach(tab => {
       tab.style.backgroundColor = "gray";
       tab.style.color = "white";
       if (tab.classList.contains(newTab)) {
           tab.style.backgroundColor = "greenyellow";
           tab.style.color = "black";
       }
    });

    var tabs = document.getElementsByClassName("is-active");
    Array.from(tabs).forEach(tab => {
        //console.log("removing tag",tab);
        tab.classList.remove("is-active");
    });
    tabs = document.getElementsByClassName(newTab);
    Array.from(tabs).forEach(tab => {
        // console.log("adding tag",tab);
        tab.classList.add("is-active");
    });
    tabs = document.getElementsByClassName("is-active");
    Array.from(tabs).forEach(tab => {
        //console.log("Show tag",tab);
    });
    var url = window.location.href;
    //console.log("active url", url);
    if (location.search.includes("tab=")) {


    var myParam = location.search.split('tab=')[1];

    url = url.replace("tab="+myParam, "tab="+newTab);
    } else {
        url =  url + "?tab="+newTab ;
    }
    window.history.replaceState(null, null, url);
    //console.log("active url", url);
}

function setTemporary(type) {
    var valueDay = document.getElementById('day').value !== "" ? parseFloat(document.getElementById('day').value) : 0;
    var valueHour = document.getElementById('hour').value !== "" ? parseFloat(document.getElementById('hour').value) : 0;
    valueHour = valueHour + valueDay * 24;
    ipPort = document.getElementById('ipPort').value;
    url = 'http://' + ipPort + '/heating/temporary/' + type + '/' + valueHour;
    //console.log(url);
    $.get(url, function (data, status) {
        //console.log(data);
    });
    refresh();
    //window.location.assign(url);
}

function toggleDevice(id, url) {
    //console.log(url);
    $.get(url, function (data, status) {
        //console.log(data);
    });
    snackbar("Done");
}

function runAction(id, action, payload) {
    const url = URLAction+ "?id="+id+"&action="+action+"&payload="+payload;
    //console.log(url);
    $.get(url, function (data, status) {
        //console.log(data);
    });
    snackbar("Done");
}
function runActionShelly4PM(id, action, instance, payload) {
    const url = URLAction+ "?id="+id+"&action="+action+"&payload="+'{"id": 1, "src":"Mimerel", "method": "Switch.Set", "params": {"id": '+instance+', "on": '+payload+'}}';
    //console.log(url);
    $.get(url, function (data, status) {
        //console.log(data);
    });
    snackbar("Done");
}

function runActionValueChange(id, action,  payload) {

    var value = document.getElementById('temperatureTarget_'+id).innerText !== "" ? parseFloat(document.getElementById('temperatureTarget_'+id).innerText) : 0;
    const newvalue = value+parseFloat(payload);

    //console.log("Change temp", value, payload);

    const url = URLAction+ "?id="+id+"&action="+action+"&payload="+newvalue;
    //console.log(url);
    $.get(url, function (data, status) {
        //console.log(data);
    });
    snackbar("Done");
}


function runReconnect() {
    const url = URL+ "/reconnect";
    //console.log(url);
    $.get(url, function (data, status) {
        //console.log(data);
    });
    snackbar("Done");
}



function slideDevice(id) {
    var slider = document.getElementById("slider" + id).value;
    const action = "/roller/0/command/pos";
    const url = URLAction+ "?id="+id+"&action="+action+"&payload="+slider;
    $.get(url, function (data, status) {
        //console.log(data);
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
