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
	"fmt"

	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/broker"
	"github.com/GoogleCloudPlatform/gcp-service-broker/pkg/validation"
)

const (
	// InstanceIDKey is the key used by Instance identifier BrokerVariables.
	InstanceIDKey = "instance_id"

	// ZoneKey is the key used by Zone BrokerVariables.
	ZoneKey = "zone"

	// RegionKey is the key used by Region BrokerVariables.
	RegionKey = "region"

	// AuthorizedNetworkKey is the key used to define authorized networks.
	AuthorizedNetworkKey = "authorized_network"
)

// UniqueArea defines an umbrella under which identifiers must be unique.
type UniqueArea string

const (
	// ZoneArea indicates uniqueness per zone
	ZoneArea UniqueArea = "per zone"
	// RegionArea indicates uniqueness per region
	RegionArea UniqueArea = "per region"
	// ProjectArea indicates uniqueness per project
	ProjectArea UniqueArea = "per project"
	// GlobalArea indicates global uniqueness
	GlobalArea UniqueArea = "globally"
)

// InstanceID creates an InstanceID broker variable with key InstanceIDKey.
// It accepts lower-case InstanceIDs with hyphens.
func InstanceID(minLength, maxLength int, uniqueArea UniqueArea) broker.BrokerVariable {
	return broker.BrokerVariable{
		FieldName: InstanceIDKey,
		Type:      broker.JsonTypeString,
		Details:   fmt.Sprintf("The name of the instance. The name must be unique %s.", uniqueArea),
		Default:   "gsb-${counter.next()}-${time.nano()}",
		Constraints: validation.NewConstraintBuilder().
			MinLength(minLength).
			MaxLength(maxLength).
			Pattern("^[a-z]([-0-9a-z]*[a-z0-9]$)*").
			Build(),
	}
}

// Zone creates a variable that accepts GCP zones.
func Zone(defaultLocation, supportedLocationsURL string) broker.BrokerVariable {
	return broker.BrokerVariable{
		FieldName: ZoneKey,
		Type:      broker.JsonTypeString,
		Details:   fmt.Sprintf("The zone to create the instance in. Supported zones can be found here: %s.", supportedLocationsURL),
		Default:   defaultLocation,
		Constraints: validation.NewConstraintBuilder().
			Pattern("^[A-Za-z][-a-z0-9A-Z]+$").
			Build(),
	}
}

// Region creates a variable that accepts GCP regions.
func Region(defaultLocation, supportedLocationsURL string) broker.BrokerVariable {
	return broker.BrokerVariable{
		FieldName: RegionKey,
		Type:      broker.JsonTypeString,
		Details:   fmt.Sprintf("The region to create the instance in. Supported regions can be found here: %s.", supportedLocationsURL),
		Default:   defaultLocation,
		Constraints: validation.NewConstraintBuilder().
			Pattern("^[A-Za-z][-a-z0-9A-Z]+$").
			Build(),
	}
}

// AuthorizedNetwork returns a variable used to attach resources to a
// user-defined network.
func AuthorizedNetwork() broker.BrokerVariable {
	return broker.BrokerVariable{
		FieldName: AuthorizedNetworkKey,
		Type:      broker.JsonTypeString,
		Details:   "The name of the VPC network to attach the instance to.",
		Default:   "default",
		Constraints: validation.NewConstraintBuilder().
			Examples("default", "projects/MYPROJECT/global/networks/MYNETWORK").
			Build(),
	}
}
