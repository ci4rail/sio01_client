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

package restart

import (
	api "github.com/ci4rail/io4edge-client-go/core/v1alpha2"
)

// Restart simulates a restart
func Restart() (res *api.CoreResponse, doreset bool) {

	res = &api.CoreResponse{
		Id:            api.CommandId_RESTART,
		Status:        api.Status_OK,
		RestartingNow: true,
	}
	return res, true
}
