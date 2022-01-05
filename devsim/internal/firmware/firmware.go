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

package firmware

import (
	"encoding/json"
	"errors"
	"log"
	"strings"

	api "github.com/ci4rail/io4edge-client-go/core/v1alpha2"
)

type firmwareID struct {
	Name    string
	Version string
}

type firmwareHeader struct {
	Name          string
	Version       string
	FirmwareWorks bool `json:"firmware_works"`
}

var (
	fwID            = firmwareID{"default", "1.0.0.alpha"}
	nextChunkNumber = uint32(0)
	nextFlashOffset = uint32(0)
	flash           = make([]byte, 200000)
)

// IdentifyFirmware reports the currently active firmware name and version
func IdentifyFirmware() *api.CoreResponse {

	res := &api.CoreResponse{
		Id:            api.CommandId_IDENTIFY_FIRMWARE,
		Status:        api.Status_OK,
		RestartingNow: false,
		Data: &api.CoreResponse_IdentifyFirmware{
			IdentifyFirmware: &api.IdentifyFirmwareResponse{
				Name:    fwID.Name,
				Version: fwID.Version,
			},
		},
	}
	return res
}

// LoadFirmwareChunk loads the chunk in c to the virtual flash
func LoadFirmwareChunk(c *api.LoadFirmwareChunkCommand) (res *api.CoreResponse, doreset bool) {

	var status = api.Status_OK
	doreset = false
	if nextChunkNumber != c.ChunkNumber {
		status = api.Status_BAD_CHUNK_SEQ
	} else {
		log.Printf("Loading chunk %d @%08x\n", nextChunkNumber, nextFlashOffset)

		// simulate flash programming
		copy(flash[nextFlashOffset:], c.Data)
		nextFlashOffset += uint32(len(c.Data))
		nextChunkNumber++

		if c.IsLastChunk {
			nextFlashOffset = 0
			nextChunkNumber = 0
			header, err := fwHeaderFromFlash(flash)
			if err != nil {
				log.Printf("firmware header not ok %v\n", err)
			} else if !header.FirmwareWorks {
				log.Printf("firmware not working\n")
				doreset = true
			} else {
				log.Printf("activating new firmware %v\n", header)
				fwID.Name = header.Name
				fwID.Version = header.Version
				doreset = true
			}
		}
	}

	res = &api.CoreResponse{
		Id:            api.CommandId_LOAD_FIRMWARE_CHUNK,
		RestartingNow: doreset,
		Status:        status,
	}
	return res, doreset
}

func fwHeaderFromFlash(flash []byte) (*firmwareHeader, error) {
	// find end of json. This works only if json has no nested {}
	idx := strings.Index(string(flash), "}")
	if idx == -1 {
		return nil, errors.New("bad json")
	}
	flashJSON := flash[:idx+1]

	var header firmwareHeader
	err := json.Unmarshal(flashJSON, &header)
	if err != nil {
		return nil, err
	}
	return &header, nil
}
