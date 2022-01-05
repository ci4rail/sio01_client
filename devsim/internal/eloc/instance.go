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
	"github.com/hashicorp/mdns"
)

// Eloc represents the Easylocate functionality
type Eloc struct {
	deviceID                string
	statusServerMdnsService *mdns.MDNSService
	statusServerMdnsServer  *mdns.Server
	loc                     chan location
	haveServerConnection    bool
	havePosition            bool
}

// NewInstance creates a new Easylocate simulator instance
func NewInstance(deviceID string, statusServerPort int, locationServerAddress string) (*Eloc, error) {
	e := &Eloc{
		deviceID: deviceID,
		loc:      make(chan location),
	}
	err := e.startMdns(statusServerPort)
	if err != nil {
		return nil, err
	}

	err = e.locationClient(locationServerAddress)
	if err != nil {
		return nil, err
	}
	e.locationGenerator()

	return e, nil
}
