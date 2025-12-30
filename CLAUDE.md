# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build and Run Commands

```bash
# Run the application
go run main.go

# Build the application
go build -o domotique

# Run built binary
./domotique
```

The application runs an HTTP server on port 9998 (configurable in `configuration.yaml`).

**Configuration file lookup order:**
1. `LOGGER_CONFIGURATION_FILE` environment variable
2. `/home/pi/go/src/go-domotique/configuration.yaml` (Raspberry Pi production path)
3. `./configuration.yaml` (local development)

## Architecture Overview

**go-domotique** is a home automation system that manages smart devices (primarily Shelly IoT devices) through MQTT, provides heating/radiator management, and event notifications.

### Core Flow

```
main.go → controller.Controller() → Spawns daemon goroutine + HTTP server
```

### Key Packages

| Package | Responsibility |
|---------|----------------|
| `controller/` | HTTP route registration and system orchestration |
| `daemon/` | MQTT event loop, device state management, message processing |
| `configuration/` | YAML config loading, MariaDB queries for devices/heating/commands |
| `models/` | Domain data structures (Configuration, Devices, Heating, Channels) |
| `heating/` | Temperature/radiator management, heating schedules, web UI |
| `devices/` | Device lookup helpers |
| `events/` | Event logging to MariaDB, Prowl push notifications |
| `utils/` | HTTP client with retry, database helpers, string utilities |
| `logger/` | Colored console output with log levels |

### Channel-Based Communication

The system uses Go channels for cross-goroutine communication:
- `MqttSend`: Send MQTT messages to devices
- `MqttReconnect`: Reconnect MQTT broker
- `MqttDomotiqueDevicePost`: Queue device state updates (batch processed)
- `UpdateConfig`: Trigger configuration reload

### MQTT Device Integration

- Broker: `192.168.222.55:1883`
- Topic format: `shellies/device_{id}/{topic}`
- Topics include: `/relay/0/power`, `/relay/0/ison`, `/temperature`, `/humidity`, `/battery`, etc.
- Message handlers in `daemon/mqttUpdateDevice.go`

### Database Structure

MariaDB stores:
- `domotique` database: devices, heating programs
- `domotiqueStats` database: event logs, statistics
- Key tables: `devices`, `heating`, `heatingProgram`, `heatingLevels`, `deviceActions`

## HTTP Endpoints

- `/heating/status` - Heating system web UI
- `/heating/update`, `/heating/updateValues` - Heating AJAX updates
- `/runAction` - Execute device actions
- `/wifi/` - WiFi device control
- `/configuration/update` - Hot reload configuration
- `/event/new` - Event capture

## Key Files to Understand First

1. `controller/controller.go` - System orchestration, route registration
2. `daemon/mqttDaemon.go` - Core MQTT event loop (`Mqtt_Deamon` function)
3. `daemon/mqttUpdateDevice.go` - MQTT message parsing and device state updates
4. `configuration/configuration.go` - Config loading and initialization
5. `models/configuration.go` - Core data structures

## Development Notes

- All feature configuration is stored in MariaDB, not code
- Devices are identified by `domotiqueId`; `boxId=100` indicates MQTT/Shelly devices
- Heating system requests temperature from sensors every minute
- Use `config.Logger.Info/Debug/Warn/Error()` for logging
- Templates are in `./heating/templates/` and use Go's `html/template`

## Git Commit Convention

Commits follow the pattern: `[PREFIX][AREA] descriptive message`
- Example: `[HF][Heating] request temp from devices every minute`
- Example: `[UP][Radiators] moved to percentage based management`
