**This project is in ALPHA status and is subject to change.**

# EnOcean CCU Server

The EnOcean CCU server can be used to add EnOcean devices to the Homey via the Hoematic App. The server will be 
discovered as another CCU by the App.

The serial used by this CCU is currently hardcoded to abcdefg. This will be configurable in the future to be able to run 
multiple instances.

# Install

The server is implemented as a single binary. However it is required to run an MQTT broker on the same device as 
this server is running.

## Install on Raspberry Pi3

Install Mosquitto MQTT broker

```
sudo apt install mosquitto
```

Save the enoceanccu binay file under /home/pi

Make the fie executabble.

```
chmod +x /home/pi/enoceanccu
```

Create service in systemd
```
sudo bash -c "cat >/etc/systemd/system/enoceanccu.service" <<'EOF'
[Unit]
Description=EnOcean CCU Server
Requires=mosquitto.service
After=mosquitto.service

[Service]
Type=simple
User=pi
WorkingDirectory=/home/pi
ExecStart=/home/pi/enoceanccu --device /dev/ttyUSB0 --devices-config /home/pi/devices.json --serial abcdefg

[Install]
WantedBy=multi-user.target
EOF
```

Relead systemd configuration

```
sudo systemctl daemon-reload
```

Enable the service to be started on boot

```
sudo systemctl enable enoceanccu
```

Create file with enocean devices

The devices are read from /home/pi/devices.json

The format of the file is as follows:

```
{
  "<enocean id>": {
    "rcv-eeps": [
      "eepf60201"
    ],
    "hm-type": "HM-PB-4-WM",
    "hm-address": "<Name as it will appear in Homey>"
  },
  "<another enocean id>": {
    "rcv-eeps": [
      "eepf60201"
    ],
    "hm-type": "HM-PB-4-WM",
    "hm-address": "<Name as this device will appear in Homey>"
  }
}
```

Example:

```
{
  "fecd2345": {
    "rcv-eeps": [
      "eepf60201"
    ],
    "hm-type": "HM-PB-4-WM",
    "hm-address": "Taster 1"
  },
  "abcdef1278": {
    "rcv-eeps": [
      "eepf60201"
    ],
    "hm-type": "HM-PB-4-WM",
    "hm-address": "Taster 2"
  }
}
```

Start the server

```
sudo systemctl start enoceanccu
```
