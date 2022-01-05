/*
Copyright © 2021 Ci4Rail GmbH
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"io"
	"log"
	"time"

	"github.com/ci4rail/sio01_host/devsim/internal/firmware"
	"github.com/ci4rail/sio01_host/devsim/internal/hardware"
	"github.com/ci4rail/sio01_host/devsim/internal/restart"
	"github.com/ci4rail/sio01_host/devsim/pkg/version"
	"github.com/ci4rail/io4edge-client-go/client"
	api "github.com/ci4rail/io4edge-client-go/core/v1alpha2"
	"github.com/ci4rail/io4edge-client-go/transport"
	"github.com/ci4rail/io4edge-client-go/transport/socket"
)

var (
	port = ":9999"
)

func main() {
	log.Printf("io4edge-devsim version: %s listen at port %s\n", version.Version, port)

	listener, err := socket.NewSocketListener(port)
	if err != nil {
		log.Fatalf("Failed to create listener: %s", err)
	}

	if err != nil {
		log.Fatalf("Failed to create devproto: %s", err)
	}
	for {
		conn, err := socket.WaitForSocketConnect(listener)
		if err != nil {
			log.Fatalf("Failed to wait for connection: %s", err)
		}
		log.Printf("new connection!\n")

		ms := transport.NewFramedStreamFromTransport(conn)
		ch := client.NewChannel(ms)

		serveConnection(ch)
		time.Sleep(4 * time.Second) // simulate reboot
	}
}

func serveConnection(ch *client.Channel) {
	defer ch.Close()

	for {
		c := &api.CoreCommand{}
		err := ch.ReadMessage(c, 0)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Fatalf("Failed to read: %s", err)
		}

		var res *api.CoreResponse
		dorestart := false
		switch c.Id {
		case api.CommandId_IDENTIFY_FIRMWARE:
			res = firmware.IdentifyFirmware()
		case api.CommandId_IDENTIFY_HARDWARE:
			res = hardware.IdentifyHardware()
		case api.CommandId_PROGRAM_HARDWARE_IDENTIFICATION:
			res = hardware.ProgramHardwareIdentification(c.GetProgramHardwareIdentification())
		case api.CommandId_LOAD_FIRMWARE_CHUNK:
			res, dorestart = firmware.LoadFirmwareChunk(c.GetLoadFirmwareChunk())
		case api.CommandId_RESTART:
			res, dorestart = restart.Restart()
		default:
			res = &api.CoreResponse{
				Id:     c.Id,
				Status: api.Status_UNKNOWN_COMMAND,
			}
		}

		err = ch.WriteMessage(res)
		if err != nil {
			log.Printf("Failed to write: %s", err)
			return
		}
		if dorestart {
			return
		}
	}
}
