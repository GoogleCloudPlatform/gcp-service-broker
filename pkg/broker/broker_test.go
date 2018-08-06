// Copyright the Service Broker Project Authors.
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

package broker

import (
	"fmt"

	"github.com/spf13/viper"
)

func ExampleBrokerService_EnabledProperty() {
	service := BrokerService{
		Name: "left-handed-smoke-sifter",
	}

	fmt.Println(service.EnabledProperty())

	// Output: service.left-handed-smoke-sifter.enabled
}

func ExampleBrokerService_DefinitionProperty() {
	service := BrokerService{
		Name: "left-handed-smoke-sifter",
	}

	fmt.Println(service.DefinitionProperty())

	// Output: service.left-handed-smoke-sifter.definition
}

func ExampleBrokerService_UserDefinedPlansProperty() {
	service := BrokerService{
		Name: "left-handed-smoke-sifter",
	}

	fmt.Println(service.UserDefinedPlansProperty())

	// Output: service.left-handed-smoke-sifter.plans
}

func ExampleBrokerService_IsEnabled() {
	service := BrokerService{
		Name: "left-handed-smoke-sifter",
	}

	viper.Set(service.EnabledProperty(), true)
	fmt.Println(service.IsEnabled())

	viper.Set(service.EnabledProperty(), false)
	fmt.Println(service.IsEnabled())

	// Output: true
	// false
}
