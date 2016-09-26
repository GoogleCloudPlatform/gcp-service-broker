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
//
////////////////////////////////////////////////////////////////////////////////
//

package brokerapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"code.cloudfoundry.org/lager"
	"github.com/gorilla/mux"
	"gcp-service-broker/auth"
	"gcp-service-broker/brokerapi/brokers/models"
)

const provisionLogKey = "provision"
const deprovisionLogKey = "deprovision"
const bindLogKey = "bind"
const unbindLogKey = "unbind"
const lastOperationLogKey = "lastOperation"

const instanceIDLogKey = "instance-id"
const instanceDetailsLogKey = "instance-details"
const bindingIDLogKey = "binding-id"

const invalidServiceDetailsErrorKey = "invalid-service-details"
const invalidBindDetailsErrorKey = "invalid-bind-details"
const invalidUnbindDetailsErrorKey = "invalid-unbind-details"
const invalidDeprovisionDetailsErrorKey = "invalid-deprovision-details"
const instanceLimitReachedErrorKey = "instance-limit-reached"
const instanceAlreadyExistsErrorKey = "instance-already-exists"
const bindingAlreadyExistsErrorKey = "binding-already-exists"
const instanceMissingErrorKey = "instance-missing"
const bindingMissingErrorKey = "binding-missing"
const asyncRequiredKey = "async-required"
const planChangeNotSupportedKey = "plan-change-not-supported"
const unknownErrorKey = "unknown-error"
const invalidRawParamsKey = "invalid-raw-params"
const appGuidNotProvidedErrorKey = "app-guid-not-provided"

const statusUnprocessableEntity = 422

type BrokerCredentials struct {
	Username string
	Password string
}

func New(serviceBroker models.ServiceBroker, logger lager.Logger, brokerCredentials BrokerCredentials) http.Handler {
	router := mux.NewRouter()
	AttachRoutes(router, serviceBroker, logger)
	return auth.NewWrapper(brokerCredentials.Username, brokerCredentials.Password).Wrap(router)
}

func AttachRoutes(router *mux.Router, serviceBroker models.ServiceBroker, logger lager.Logger) {
	handler := serviceBrokerHandler{serviceBroker: serviceBroker, logger: logger}
	router.HandleFunc("/v2/catalog", handler.catalog).Methods("GET")

	router.HandleFunc("/v2/service_instances/{instance_id}", handler.provision).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{instance_id}", handler.deprovision).Methods("DELETE")
	router.HandleFunc("/v2/service_instances/{instance_id}/last_operation", handler.lastOperation).Methods("GET")
	router.HandleFunc("/v2/service_instances/{instance_id}", handler.update).Methods("PATCH")

	router.HandleFunc("/v2/service_instances/{instance_id}/service_bindings/{binding_id}", handler.bind).Methods("PUT")
	router.HandleFunc("/v2/service_instances/{instance_id}/service_bindings/{binding_id}", handler.unbind).Methods("DELETE")
}

type serviceBrokerHandler struct {
	serviceBroker models.ServiceBroker
	logger        lager.Logger
}

func (h serviceBrokerHandler) catalog(w http.ResponseWriter, req *http.Request) {
	catalog := CatalogResponse{
		Services: h.serviceBroker.Services(),
	}

	h.respond(w, http.StatusOK, catalog)
}

func (h serviceBrokerHandler) provision(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["instance_id"]

	logger := h.logger.Session(provisionLogKey, lager.Data{
		instanceIDLogKey: instanceID,
	})

	var details models.ProvisionDetails
	if err := json.NewDecoder(req.Body).Decode(&details); err != nil {
		logger.Error(invalidServiceDetailsErrorKey, err)
		h.respond(w, statusUnprocessableEntity, ErrorResponse{
			Description: err.Error(),
		})
		return
	}

	acceptsIncompleteFlag, _ := strconv.ParseBool(req.URL.Query().Get("accepts_incomplete"))

	logger = logger.WithData(lager.Data{
		instanceDetailsLogKey: details,
	})

	provisionResponse, err := h.serviceBroker.Provision(instanceID, details, acceptsIncompleteFlag)

	if err != nil {
		switch err {
		case models.ErrRawParamsInvalid:
			logger.Error(invalidRawParamsKey, err)
			h.respond(w, 422, ErrorResponse{
				Description: err.Error(),
			})
		case models.ErrInstanceAlreadyExists:
			logger.Error(instanceAlreadyExistsErrorKey, err)
			h.respond(w, http.StatusConflict, EmptyResponse{})
		case models.ErrInstanceLimitMet:
			logger.Error(instanceLimitReachedErrorKey, err)
			h.respond(w, http.StatusInternalServerError, ErrorResponse{
				Description: err.Error(),
			})
		case models.ErrAsyncRequired:
			logger.Error(asyncRequiredKey, err)
			h.respond(w, 422, ErrorResponse{
				Error:       "AsyncRequired",
				Description: err.Error(),
			})
		default:
			logger.Error(unknownErrorKey, err)
			h.respond(w, http.StatusInternalServerError, ErrorResponse{
				Description: err.Error(),
			})
		}
		return
	}

	if provisionResponse.IsAsync {
		h.respond(w, http.StatusAccepted, ProvisioningResponse{
			DashboardURL: provisionResponse.DashboardURL,
		})
	} else {
		h.respond(w, http.StatusCreated, ProvisioningResponse{
			DashboardURL: provisionResponse.DashboardURL,
		})
	}
}

func (h serviceBrokerHandler) update(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["instance_id"]

	var details models.UpdateDetails
	if err := json.NewDecoder(req.Body).Decode(&details); err != nil {
		h.logger.Error(invalidServiceDetailsErrorKey, err)
		h.respond(w, statusUnprocessableEntity, ErrorResponse{
			Description: err.Error(),
		})
		return
	}

	acceptsIncompleteFlag, _ := strconv.ParseBool(req.URL.Query().Get("accepts_incomplete"))

	isAsync, err := h.serviceBroker.Update(instanceID, details, acceptsIncompleteFlag)
	if err != nil {
		switch err {
		case models.ErrAsyncRequired:
			h.logger.Error(asyncRequiredKey, err)
			h.respond(w, 422, ErrorResponse{
				Error:       "AsyncRequired",
				Description: err.Error(),
			})
			return

		case models.ErrPlanChangeNotSupported:
			h.logger.Error(planChangeNotSupportedKey, err)
			h.respond(w, 422, ErrorResponse{
				Error:       "PlanChangeNotSupported",
				Description: err.Error(),
			})
			return

		default:
			h.logger.Error(unknownErrorKey, err)
			h.respond(w, http.StatusInternalServerError, ErrorResponse{
				Description: err.Error(),
			})
			return
		}
	}

	statusCode := http.StatusOK
	if isAsync {
		statusCode = http.StatusAccepted
	}
	h.respond(w, statusCode, struct{}{})
}

func (h serviceBrokerHandler) deprovision(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["instance_id"]
	logger := h.logger.Session(deprovisionLogKey, lager.Data{
		instanceIDLogKey: instanceID,
	})

	details := models.DeprovisionDetails{
		PlanID:    req.FormValue("plan_id"),
		ServiceID: req.FormValue("service_id"),
	}
	asyncAllowed := req.FormValue("accepts_incomplete") == "true"

	isAsync, err := h.serviceBroker.Deprovision(instanceID, details, asyncAllowed)
	if err != nil {
		switch err {
		case models.ErrInstanceDoesNotExist:
			logger.Error(instanceMissingErrorKey, err)
			h.respond(w, http.StatusGone, EmptyResponse{})
		case models.ErrAsyncRequired:
			logger.Error(asyncRequiredKey, err)
			h.respond(w, 422, ErrorResponse{
				Error:       "AsyncRequired",
				Description: err.Error(),
			})
		default:
			logger.Error(unknownErrorKey, err)
			h.respond(w, http.StatusInternalServerError, ErrorResponse{
				Description: err.Error(),
			})
		}
		return
	}

	if isAsync {
		h.respond(w, http.StatusAccepted, EmptyResponse{})
	} else {
		h.respond(w, http.StatusOK, EmptyResponse{})
	}
}

func (h serviceBrokerHandler) bind(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["instance_id"]
	bindingID := vars["binding_id"]

	logger := h.logger.Session(bindLogKey, lager.Data{
		instanceIDLogKey: instanceID,
		bindingIDLogKey:  bindingID,
	})

	var details models.BindDetails
	if err := json.NewDecoder(req.Body).Decode(&details); err != nil {
		logger.Error(invalidBindDetailsErrorKey, err)
		h.respond(w, statusUnprocessableEntity, ErrorResponse{
			Description: err.Error(),
		})
		return
	}

	binding, err := h.serviceBroker.Bind(instanceID, bindingID, details)
	if err != nil {
		switch err {
		case models.ErrInstanceDoesNotExist:
			logger.Error(instanceMissingErrorKey, err)
			h.respond(w, http.StatusNotFound, ErrorResponse{
				Description: err.Error(),
			})
		case models.ErrBindingAlreadyExists:
			logger.Error(bindingAlreadyExistsErrorKey, err)
			h.respond(w, http.StatusConflict, ErrorResponse{
				Description: err.Error(),
			})
		case models.ErrAppGuidNotProvided:
			logger.Error(appGuidNotProvidedErrorKey, err)
			h.respond(w, statusUnprocessableEntity, ErrorResponse{
				Description: err.Error(),
			})
		default:
			logger.Error(unknownErrorKey, err)
			h.respond(w, http.StatusInternalServerError, ErrorResponse{
				Description: err.Error(),
			})
		}
		return
	}

	h.respond(w, http.StatusCreated, binding)
}

func (h serviceBrokerHandler) unbind(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["instance_id"]
	bindingID := vars["binding_id"]

	logger := h.logger.Session(unbindLogKey, lager.Data{
		instanceIDLogKey: instanceID,
		bindingIDLogKey:  bindingID,
	})

	details := models.UnbindDetails{
		PlanID:    req.FormValue("plan_id"),
		ServiceID: req.FormValue("service_id"),
	}

	if err := h.serviceBroker.Unbind(instanceID, bindingID, details); err != nil {
		switch err {
		case models.ErrInstanceDoesNotExist:
			logger.Error(instanceMissingErrorKey, err)
			h.respond(w, http.StatusGone, EmptyResponse{})
		case models.ErrBindingDoesNotExist:
			logger.Error(bindingMissingErrorKey, err)
			h.respond(w, http.StatusGone, EmptyResponse{})
		default:
			logger.Error(unknownErrorKey, err)
			h.respond(w, http.StatusInternalServerError, ErrorResponse{
				Description: err.Error(),
			})
		}
		return
	}

	h.respond(w, http.StatusOK, EmptyResponse{})
}

func (h serviceBrokerHandler) lastOperation(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	instanceID := vars["instance_id"]

	logger := h.logger.Session(lastOperationLogKey, lager.Data{
		instanceIDLogKey: instanceID,
	})

	logger.Info("starting-check-for-operation")

	lastOperation, err := h.serviceBroker.LastOperation(instanceID)

	if err != nil {
		switch err {
		case models.ErrInstanceDoesNotExist:
			logger.Error(instanceMissingErrorKey, err)
			h.respond(w, http.StatusNotFound, ErrorResponse{
				Description: err.Error(),
			})
		default:
			logger.Error(unknownErrorKey, err)
			h.respond(w, http.StatusInternalServerError, ErrorResponse{
				Description: err.Error(),
			})
		}

		return
	}

	logger.WithData(lager.Data{"state": lastOperation.State}).Info("done-check-for-operation")

	lastOperationResponse := LastOperationResponse{
		State:       string(lastOperation.State),
		Description: lastOperation.Description,
	}

	h.respond(w, http.StatusOK, lastOperationResponse)
}

func (h serviceBrokerHandler) respond(w http.ResponseWriter, status int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	err := encoder.Encode(response)
	if err != nil {
		h.logger.Error("encoding response", err, lager.Data{"status": status, "response": response})
	}
}
