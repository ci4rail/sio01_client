# sio01_client
Client examples and protobuf files for SIO01 device

## Protobuf definitions
The `proto` folder contains
* [tracelet_location.proto](proto/tracelet_location.proto): The message sent by the SIO01 device to the location server
* [tracelet_status.proto](proto/tracelet_status.proto): The messages exchanged between the monitoring system and the SIO01 device

## Examples

The `examples` folder contains some simple examples in python:
* [location_server](examples/location_server.py): Receive location messages from SIO01 devices and print them to stdout.
* [status_client](examples/status_client.py): Send a status request to SIO01 device and print result.

### Usage

Prerequisites:
* Linux machine
* Python >= 3.9 installed

```bash
cd examples
pip3 install -r requirements.txt
```

#### Run location server:
```bash
./location_server.py
```
On the SIO01, configure parameter `loc-srv` to the IP address of the machine executing `location_server.py` and port 11002, e.g. `192.168.0.200:11002`.

#### Run status client:
```bash
./status_client.py <device-id>
```
Where `device-id` is the id set in the SIO01 via the persistent parameter `device-id`.

Note: The machine executing `status_client.py` must be in the same network as the SIO01.
