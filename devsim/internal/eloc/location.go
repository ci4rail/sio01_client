/*
Copyright © 2022 Ci4Rail GmbH
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

package eloc

import (
	"errors"
	"fmt"
	"time"

	"github.com/ci4rail/io4edge-client-go/client"
	"github.com/ci4rail/io4edge-client-go/transport"
	"github.com/ci4rail/io4edge-client-go/transport/socket"
	"github.com/ci4rail/sio01_host/devsim/internal/eloc/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (e *Eloc) locationClient(locationServerAddress string) error {
	go func() {
		for {
			ch, err := channelFromSocketAddress(locationServerAddress)

			if err == nil {
				for {
					m := &pb.LocationReport{
						ReceiveTs:         timestamppb.Now(),
						TraceletId:        e.deviceID,
						X:                 0.01,
						Y:                 0.02,
						SiteId:            12345,
						LocationSignature: 0x12345678ABCDEF,
					}
					err := ch.WriteMessage(m)
					if err != nil {
						fmt.Errorf("WriteMessage failed", err)
						break
					}
					time.Sleep(1000 * time.Millisecond)
				}
			}

			time.Sleep(500 * time.Millisecond)
		}
	}()
	return nil
}

func (e *Eloc) locationGenerator() {

}

func channelFromSocketAddress(address string) (*client.Channel, error) {
	t, err := socket.NewSocketConnection(address)
	if err != nil {
		return nil, errors.New("can't create connection: " + err.Error())
	}
	ms := transport.NewFramedStreamFromTransport(t)
	ch := client.NewChannel(ms)

	return ch, nil
}
