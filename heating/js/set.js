function setTemporary(type) {
    var valueDay = parseFloat(document.getElementById('day').value);
    var valueHour = parseFloat(document.getElementById('hour').value);
    if (valueDay !=="" ) {
        valueHour = valueHour + valueDay * 24;
    }
    ipPort = document.getElementById('ipPort').value;
    url = 'http://' + ipPort + '/heating/temporary/' + type +'/' +valueHour;
    window.location.assign( url);
}