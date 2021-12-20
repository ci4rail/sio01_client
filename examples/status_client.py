#!/usr/bin/env python3

""" Example that shows how to get the status of a SIO01 """

import logging
import sys
import socket
import struct
import tracelet_status_pb2

from zeroconf import Zeroconf

TYPE = '_io4edge-eloc._tcp.local.'


def find_device(mdns_name):
    zeroconf = Zeroconf()

    info = zeroconf.get_service_info(TYPE, mdns_name + '-eloc' + '.' + TYPE)
    if info:
        rv = info.parsed_addresses()[0], info.port
    else:
        rv = None, 0

    zeroconf.close()
    return rv


def send_request(s, id):
    req = tracelet_status_pb2.StatusRequest()
    req.id = id
    data = req.SerializeToString()
    hdr = struct.pack('<HL', 0xEDFE, len(data))
    s.sendall(hdr+data)


def recv_response(s):
    hdr = s.recv(6)
    if hdr[0:2] == b'\xfe\xed':
        len = struct.unpack('<L', hdr[2:6])[0]
        proto_data = s.recv(len)
        stat = tracelet_status_pb2.StatusResponse()
        stat.ParseFromString(proto_data)
        return stat
    else:
        raise RuntimeError('bad magic')


def status_request(s, id):
    send_request(s, id)
    stat = recv_response(s)
    print(f'{stat.power_up_count} power ups')
    print(f"has {'' if stat.has_server_connection else 'no '}server connection")
    print(f"has {'valid' if stat.has_time else 'no'} time")
    print(f"has {'valid' if stat.has_position else 'no'} position")
    print(
        f"eloc module status is {'' if stat.eloc_module_status_ok else 'not '}ok")


if __name__ == '__main__':
    if len(sys.argv) < 2:
        print('Please supply mdns name')
        sys.exit(1)
    mdns_name = sys.argv[1]

    ip, port = find_device(mdns_name)
    if ip is None:
        print("Can't find device")
        sys.exit(1)

    print(ip, port)

    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.connect((ip, port))
        status_request(s, 1)
