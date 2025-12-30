// Modern ES6+ JavaScript - No jQuery dependency
'use strict';

const BASE_URL = window.location.origin;
const WS_URL = `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}/ws`;
const URLAction = `${BASE_URL}/runAction`;

let websocket = null;
let reconnectAttempts = 0;
const MAX_RECONNECT_ATTEMPTS = 10;
const RECONNECT_DELAY = 3000;

// Initialize on page load
document.addEventListener('DOMContentLoaded', () => {
    const urlTab = new URLSearchParams(window.location.search).get('tab');
    // Get first available room tab as default
    const firstTabBtn = document.querySelector('.tabButton');
    const firstRoom = firstTabBtn ? firstTabBtn.classList[1] : 'logs'; // Second class is room name
    const tab = urlTab || firstRoom;
    changeActiveTabTo(tab);
    initWebSocket();
    initDarkMode();
});

// WebSocket connection
function initWebSocket() {
    if (websocket && websocket.readyState === WebSocket.OPEN) {
        return;
    }

    websocket = new WebSocket(WS_URL);

    websocket.onopen = () => {
        console.log('WebSocket connected');
        reconnectAttempts = 0;
        snackbar('Connected to server');
    };

    websocket.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            updateDeviceValues(data);
        } catch (e) {
            console.error('Failed to parse WebSocket message:', e);
        }
    };

    websocket.onclose = () => {
        console.log('WebSocket disconnected');
        if (reconnectAttempts < MAX_RECONNECT_ATTEMPTS) {
            reconnectAttempts++;
            setTimeout(initWebSocket, RECONNECT_DELAY);
        } else {
            snackbar('Connection lost. Please refresh the page.', true);
        }
    };

    websocket.onerror = (error) => {
        console.error('WebSocket error:', error);
    };
}

// Update all device values from received data
function updateDeviceValues(devices) {
    if (!Array.isArray(devices)) return;

    let total = 0;
    const subTotals = new Map();

    devices.forEach(device => {
        const roundPower = Math.round(device.Power * 100) / 100;

        // Calculate room subtotals
        if (device.Room) {
            const current = subTotals.get(device.Room) || 0;
            subTotals.set(device.Room, Math.round((current + roundPower) * 100) / 100);
            updateElement(`${device.Room}W`, `${subTotals.get(device.Room)} W`);
        }

        total += device.Power;

        // Update device-specific elements
        updatePower(device.DomotiqueId, roundPower, device.Status);
        updateStatus(device.DomotiqueId, device.Status);
        updateElement(`temperature_${device.DomotiqueId}`, device.Temperature ? `${device.Temperature} °C` : '-');
        updateElement(`humidity_${device.DomotiqueId}`, device.Humidity ? `${device.Humidity} %` : '-');
        updateElement(`battery_${device.DomotiqueId}`, device.Battery ? `Battery : ${device.Battery} %` : '');
        updateElement(`position_${device.DomotiqueId}`, `${device.CurrentPos || 0} %`);
        updateElement(`temperatureTarget_${device.DomotiqueId}`, device.TemperatureTarget ? `${device.TemperatureTarget} °C` : '-');
        updateElement(`deviceTemperature_${device.DomotiqueId}`, device.DeviceTemperature ? `${device.DeviceTemperature} °C` : '-');
        updateElement(`voltage_${device.DomotiqueId}`, device.Voltage ? `${device.Voltage}` : '-');
        updateElement(`tilt_${device.DomotiqueId}`, device.Tilt !== -1 ? `${device.Tilt}` : '-');
        updateElement(`lux_${device.DomotiqueId}`, device.Lux ? `${device.Lux}` : '-');
        updateElement(`illumination_${device.DomotiqueId}`, device.Illumination ? `${device.Illumination}` : '-');
        updateState(device.DomotiqueId, device.State);
        updateElement(`online_${device.DomotiqueId}`, device.Online ? 'Device Online' : 'Device OffLine');
        updateElement(`active_${device.DomotiqueId}`, device.Active ? 'Device is Active' : 'Device is not active');
    });

    updateElement('totalpower', `${Math.round(total * 100) / 100} Watts`);
}

function updateElement(id, value) {
    const el = document.getElementById(id);
    if (el) el.textContent = value;
}

function updatePower(id, power, status) {
    const el = document.getElementById(`power_${id}`);
    if (!el) return;
    el.textContent = power === 0 ? '-' : `${power} W`;
    el.style.backgroundColor = status === 'on' ? 'var(--color-success)' : status === 'off' ? 'var(--color-danger)' : '';
}

function updateStatus(id, status) {
    const el = document.getElementById(`status_${id}`);
    if (!el) return;
    if (status === 'on') {
        el.style.backgroundColor = 'var(--color-danger)';
        el.textContent = 'Enabled';
    } else if (status === 'off') {
        el.style.backgroundColor = 'var(--color-success)';
        el.textContent = 'Disabled';
    }
}

function updateState(id, state) {
    const el = document.getElementById(`state_${id}`);
    if (!el) return;
    el.textContent = state || '-';
    el.style.backgroundColor = state === 'close' ? 'var(--color-success)' : state === 'open' ? 'var(--color-danger)' : '';
}

// Tab navigation
function changeActiveTabTo(newTab) {
    // Hide all room rows, show selected
    document.querySelectorAll('.room').forEach(el => {
        el.style.display = el.classList.contains(`room${newTab}`) ? '' : 'none';
    });

    // Update tab button styles
    document.querySelectorAll('.tabButton').forEach(btn => {
        const isActive = btn.classList.contains(newTab);
        btn.style.backgroundColor = isActive ? 'var(--color-primary)' : 'var(--color-tab-inactive)';
        btn.style.color = isActive ? 'var(--color-text-inverse)' : 'var(--color-text)';
    });

    // Update URL
    const url = new URL(window.location);
    url.searchParams.set('tab', newTab);
    window.history.replaceState(null, null, url);
}

// Dark mode toggle
function initDarkMode() {
    const savedTheme = localStorage.getItem('theme');
    if (savedTheme === 'dark' || (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
        document.documentElement.setAttribute('data-theme', 'dark');
    }
}

function toggleDarkMode() {
    const isDark = document.documentElement.getAttribute('data-theme') === 'dark';
    document.documentElement.setAttribute('data-theme', isDark ? 'light' : 'dark');
    localStorage.setItem('theme', isDark ? 'light' : 'dark');
    snackbar(`${isDark ? 'Light' : 'Dark'} mode enabled`);
}

// Device actions
function toggleDevice(id, url) {
    fetch(url)
        .then(response => {
            if (!response.ok) throw new Error('Network response was not ok');
            snackbar('Device toggled');
        })
        .catch(() => snackbar('Failed to toggle device', true));
}

function runAction(id, action, payload) {
    const url = `${URLAction}?id=${id}&action=${encodeURIComponent(action)}&payload=${encodeURIComponent(payload)}`;
    fetch(url)
        .then(response => {
            if (!response.ok) throw new Error('Network response was not ok');
            snackbar('Action executed');
        })
        .catch(() => snackbar('Failed to execute action', true));
}

function runActionShelly4PM(id, action, instance, payload) {
    const jsonPayload = JSON.stringify({
        id: 1,
        src: "Mimerel",
        method: "Switch.Set",
        params: { id: instance, on: payload }
    });
    const url = `${URLAction}?id=${id}&action=${encodeURIComponent(action)}&payload=${encodeURIComponent(jsonPayload)}`;
    fetch(url)
        .then(response => {
            if (!response.ok) throw new Error('Network response was not ok');
            snackbar('Action executed');
        })
        .catch(() => snackbar('Failed to execute action', true));
}

function runActionValueChange(id, action, payload) {
    const posEl = document.getElementById(`position_${id}`);
    const currentValue = posEl ? parseFloat(posEl.textContent) || 0 : 0;
    const newValue = currentValue + parseFloat(payload);
    const url = `${URLAction}?id=${id}&action=${encodeURIComponent(action)}&payload=${newValue}`;
    fetch(url)
        .then(response => {
            if (!response.ok) throw new Error('Network response was not ok');
            snackbar('Value updated');
        })
        .catch(() => snackbar('Failed to update value', true));
}

function slideDevice(id) {
    const slider = document.getElementById(`slider${id}`);
    if (!slider) return;
    const value = slider.value;
    const url = `${URLAction}?id=${id}&action=${encodeURIComponent('/roller/0/command/pos')}&payload=${value}`;
    fetch(url)
        .then(response => {
            if (!response.ok) throw new Error('Network response was not ok');
            snackbar(`Position set to ${value}%`);
        })
        .catch(() => snackbar('Failed to set position', true));
}

function runReconnect() {
    fetch(`${BASE_URL}/reconnect`)
        .then(response => {
            if (!response.ok) throw new Error('Network response was not ok');
            snackbar('Reconnect initiated');
        })
        .catch(() => snackbar('Failed to reconnect', true));
}

function setTemporary(type) {
    const dayEl = document.getElementById('day');
    const hourEl = document.getElementById('hour');
    const valueDay = dayEl && dayEl.value ? parseFloat(dayEl.value) : 0;
    const valueHour = hourEl && hourEl.value ? parseFloat(hourEl.value) : 0;
    const totalHours = valueHour + valueDay * 24;
    const url = `${BASE_URL}/heating/temporary/${type}/${totalHours}`;
    fetch(url)
        .then(response => {
            if (!response.ok) throw new Error('Network response was not ok');
            snackbar('Temporary setting applied');
            setTimeout(() => location.reload(), 1000);
        })
        .catch(() => snackbar('Failed to set temporary value', true));
}

function refresh() {
    location.reload();
}

// Snackbar notification
function snackbar(message, isError = false) {
    const container = document.getElementById('snackbar');
    if (!container) {
        console.log(isError ? 'Error:' : 'Info:', message);
        return;
    }

    container.textContent = message;
    container.className = `snackbar show ${isError ? 'error' : ''}`;

    setTimeout(() => {
        container.className = 'snackbar';
    }, isError ? 4000 : 2000);
}

// Expose functions to global scope for onclick handlers
window.changeActiveTabTo = changeActiveTabTo;
window.toggleDevice = toggleDevice;
window.runAction = runAction;
window.runActionShelly4PM = runActionShelly4PM;
window.runActionValueChange = runActionValueChange;
window.slideDevice = slideDevice;
window.runReconnect = runReconnect;
window.setTemporary = setTemporary;
window.refresh = refresh;
window.toggleDarkMode = toggleDarkMode;
