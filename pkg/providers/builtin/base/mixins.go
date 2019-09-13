// Copyright 2019 the Service Broker Project Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package base

import (
	"context"
	"encoding/json"

	"github.com/GoogleCloudPlatform/gcp-service-broker/db_service/models"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/varcontext"
	"github.com/pivotal-cf/brokerapi"
)

type synchronousBase struct{}

// PollInstance does nothing but return an error because Base services are
// provisioned synchronously so this method should not be called.
func (b *synchronousBase) PollInstance(ctx context.Context, instance models.ServiceInstanceDetails) (bool, error) {
	return true, brokerapi.ErrAsyncRequired
}

// ProvisionsAsync indicates if provisioning must be done asynchronously.
func (b *synchronousBase) ProvisionsAsync() bool {
	return false
}

// DeprovisionsAsync indicates if deprovisioning must be done asynchronously.
func (b *synchronousBase) DeprovisionsAsync() bool {
	return false
}

// AsynchronousInstanceMixin sets ProvisionAsync and DeprovisionsAsync functions
// to be true.
type AsynchronousInstanceMixin struct{}

// ProvisionsAsync indicates if provisioning must be done asynchronously.
func (b *AsynchronousInstanceMixin) ProvisionsAsync() bool {
	return true
}

// DeprovisionsAsync indicates if deprovisioning must be done asynchronously.
func (b *AsynchronousInstanceMixin) DeprovisionsAsync() bool {
	return true
}

// NoOpBindMixin does a no-op binding. This can be used when you still want a
// service to be bindable but nothing is required server-side to support it.
// For example, when the service requires no authentication.
type NoOpBindMixin struct{}

// Bind does a no-op bind.
func (m *NoOpBindMixin) Bind(ctx context.Context, vc *varcontext.VarContext) (map[string]interface{}, error) {
	return make(map[string]interface{}), nil
}

// Unbind does a no-op unbind.
func (m *NoOpBindMixin) Unbind(ctx context.Context, instance models.ServiceInstanceDetails, creds models.ServiceBindingCredentials) error {
	return nil
}

// MergedInstanceCredsMixin adds the BuildInstanceCredentials function that
// merges the OtherDetails of the bind and instance records.
type MergedInstanceCredsMixin struct{}

// BuildInstanceCredentials combines the bind credentials with the connection
// information in the instance details to get a full set of connection details.
func (b *MergedInstanceCredsMixin) BuildInstanceCredentials(ctx context.Context, bindRecord models.ServiceBindingCredentials, instanceRecord models.ServiceInstanceDetails) (map[string]interface{}, error) {
	vc, err := varcontext.Builder().
		MergeJsonObject(json.RawMessage(bindRecord.OtherDetails)).
		MergeJsonObject(json.RawMessage(instanceRecord.OtherDetails)).
		Build()
	if err != nil {
		return nil, err
	}

	return vc.ToMap(), nil
}
