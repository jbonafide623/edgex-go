/*******************************************************************************
 * Copyright 2017 Dell Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License
 * is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
 * or implied. See the License for the specific language governing permissions and limitations under
 * the License.
 *******************************************************************************/
package metadata

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/edgexfoundry/edgex-go/internal/core/metadata/errors"
	"github.com/edgexfoundry/edgex-go/internal/core/metadata/interfaces"
	"github.com/edgexfoundry/edgex-go/internal/core/metadata/operators/device"
	"github.com/edgexfoundry/edgex-go/internal/core/metadata/operators/device_profile"
	"github.com/edgexfoundry/edgex-go/internal/pkg/db"
	"github.com/edgexfoundry/edgex-go/internal/pkg/httperror"
	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/gorilla/mux"
)

func restGetAllDeviceProfiles(w http.ResponseWriter, _ *http.Request) {
	op := device_profile.NewGetAllExecutor(Configuration.Service, dbClient, LoggingClient)
	res, err := op.Execute()
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.StatusRequestEntityTooLargeErrorConcept{}}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	json.NewEncoder(w).Encode(&res)
}

func restAddDeviceProfile(w http.ResponseWriter, r *http.Request) {
	var dp models.DeviceProfile

	if err := json.NewDecoder(r.Body).Decode(&dp); err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusBadRequestErrorConcept{})
		return
	}

	if Configuration.Writable.EnableValueDescriptorManagement {
		// Check if the device profile name is unique so that we do not create ValueDescriptors for a DeviceProfile that
		// will fail during the creation process later on.
		nameOp := device_profile.NewGetProfileName(dp.Name, dbClient)
		_, err := nameOp.Execute()
		// The operator will return an ItemNotFound error if the DeviceProfile can not be found.
		if err == nil {
			httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusConflictErrorConcept{})
			return
		}

		op := device_profile.NewAddValueDescriptorExecutor(r.Context(), vdc, LoggingClient, dp.DeviceResources...)
		err = op.Execute()
		if err != nil {
			httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.DeviceProfileServiceClientErrorConcept{Err: err}}, httperror.StatusInternalServerErrorConcept{})
			return
		}
	}

	addDeviceProfile(dp, dbClient, w)
}

func restUpdateDeviceProfile(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var from models.DeviceProfile
	if err := json.NewDecoder(r.Body).Decode(&from); err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusBadRequestErrorConcept{})
		return
	}

	if Configuration.Writable.EnableValueDescriptorManagement {
		vdOp := device_profile.NewUpdateValueDescriptorExecutor(from, dbClient, vdc, LoggingClient, r.Context())
		err := vdOp.Execute()
		if err != nil {
			httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.DeviceProfileServiceClientErrorConcept{Err: err}, httperror.DeviceProfileNotFoundErrorConcept{}, httperror.ValueDescriptorsInUseErrorConcept{}, httperror.DeviceProfileInvalidStateErrorConcept{}}, httperror.StatusInternalServerErrorConcept{})
			return
		}
	}

	op := device_profile.NewUpdateDeviceProfileExecutor(dbClient, from)
	dp, err := op.Execute()
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.DeviceProfileNotFoundErrorConcept{}, httperror.DeviceProfileInvalidStateErrorConcept{}}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	// Notify Associates
	err = notifyProfileAssociates(dp, dbClient, http.MethodPut)
	if err != nil {
		// Log the error but do not change the response to the client. We do not want this to affect the overall status
		// of the operation
		LoggingClient.Warn("Error while notifying profile associates of update: ", err.Error())
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("true"))
}

func restGetProfileByProfileId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var did = vars["id"]

	op := device_profile.NewGetProfileID(did, dbClient)
	res, err := op.Execute()
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.DatabaseNotFoundErrorConcept{}}, httperror.StatusInternalServerErrorConcept{})
		return
	}
	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	json.NewEncoder(w).Encode(res)
}

func restDeleteProfileByProfileId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var did = vars["id"]

	op := device_profile.NewDeleteByIDExecutor(dbClient, did)
	err := op.Execute()
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.DeviceProfileNotFoundErrorConcept{}, httperror.DeviceProfileInvalidStateErrorConcept{}}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("true"))
}

// Delete the device profile based on its name
func restDeleteProfileByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, err := url.QueryUnescape(vars[NAME])
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusBadRequestErrorConcept{})
		return
	}

	op := device_profile.NewDeleteByNameExecutor(dbClient, n)
	err = op.Execute()
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.DeviceProfileNotFoundErrorConcept{}, httperror.DeviceProfileInvalidStateErrorConcept{}}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("true"))
}

func restAddProfileByYaml(w http.ResponseWriter, r *http.Request) {
	f, _, err := r.FormFile("file")
	if err != nil {
		if err == http.ErrMissingFile {
			err = errors.NewErrEmptyFile("YAML")
		}
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.DeviceProfileMissingFileErrorConcept{}}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusInternalServerErrorConcept{})
		return
	}
	if len(data) == 0 {
		httperror.ToHttpError(w, LoggingClient, errors.NewErrEmptyFile("YAML"), []httperror.ErrorConceptType{}, httperror.StatusBadRequestErrorConcept{})
		return
	}

	var dp models.DeviceProfile

	err = yaml.Unmarshal(data, &dp)
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	// Avoid using the 'addDeviceProfile' function because we need to be backwards compatibility for API response codes.
	// The difference is the mapping of 'ErrContractInvalid' to a '409(Conflict)' rather than a '400(Bad request).
	// Disregarding backwards compatibility, the 'addDeviceProfile' function is the correct implementation to use in the
	// 'ErrContractInvalid' since a '400(Bad Request)' is the correct response.
	op := device_profile.NewAddDeviceProfileExecutor(dp, dbClient)
	id, err := op.Execute()

	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.DeviceProfileBadRequestErrorConcept{}, httperror.DeviceProfileContractInvalidErrorConcept{}, httperror.DuplicateIdentifierErrorConcept{}, httperror.DeviceProfileEmptyNameErrorConcept{}}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(id))
}

// Add a device profile with YAML content
// The YAML content is passed as a string in the http request
func restAddProfileByYamlRaw(w http.ResponseWriter, r *http.Request) {
	// Get the YAML string
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	var dp models.DeviceProfile

	err = yaml.Unmarshal(body, &dp)
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusServiceUnavailableErrorConcept{})
		return
	}

	addDeviceProfile(dp, dbClient, w)
}

// This function centralizes the common logic for adding a device profile to the database and dealing with the return
func addDeviceProfile(dp models.DeviceProfile, dbClient interfaces.DBClient, w http.ResponseWriter) {
	op := device_profile.NewAddDeviceProfileExecutor(dp, dbClient)
	id, err := op.Execute()

	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.DeviceProfileContractInvalidErrorConcept{}, httperror.DeviceProfileBadRequestErrorConcept{}, httperror.DuplicateIdentifierErrorConcept{}, httperror.DeviceProfileEmptyNameErrorConcept{}}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(id))
}

func restGetProfileByModel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	an, err := url.QueryUnescape(vars[MODEL])
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusBadRequestErrorConcept{})
		return
	}

	op := device_profile.NewGetModelExecutor(an, dbClient)
	res, err := op.Execute()
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	json.NewEncoder(w).Encode(res)
}

func restGetProfileWithLabel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	label, err := url.QueryUnescape(vars[LABEL])
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusBadRequestErrorConcept{})
		return
	}

	op := device_profile.NewGetLabelExecutor(label, dbClient)
	res, err := op.Execute()
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	json.NewEncoder(w).Encode(res)
}

func restGetProfileByManufacturerModel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	man, err := url.QueryUnescape(vars[MANUFACTURER])
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusBadRequestErrorConcept{})
		return
	}

	mod, err := url.QueryUnescape(vars[MODEL])
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusBadRequestErrorConcept{})
		return
	}

	op := device_profile.NewGetManufacturerModelExecutor(man, mod, dbClient)
	res, err := op.Execute()
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	json.NewEncoder(w).Encode(res)
}

func restGetProfileByManufacturer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	man, err := url.QueryUnescape(vars[MANUFACTURER])
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusBadRequestErrorConcept{})
		return
	}

	op := device_profile.NewGetManufacturerExecutor(man, dbClient)
	res, err := op.Execute()
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	json.NewEncoder(w).Encode(res)
}

func restGetProfileByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dn, err := url.QueryUnescape(vars[NAME])
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusBadRequestErrorConcept{})
		return
	}

	// Get the device
	op := device_profile.NewGetProfileName(dn, dbClient)
	res, err := op.Execute()
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.DatabaseNotFoundErrorConcept{}}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	json.NewEncoder(w).Encode(res)
}

func restGetYamlProfileByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name, err := url.QueryUnescape(vars[NAME])
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusBadRequestErrorConcept{})
		return
	}

	// Check for the device profile
	op := device_profile.NewGetProfileName(name, dbClient)
	dp, err := op.Execute()
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.DatabaseNotFoundErrorConcept{}}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	// Marshal into yaml
	out, err := yaml.Marshal(dp)
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

/*
 * Implementation: https://groups.google.com/forum/#!topic/golang-nuts/EZHtFOXA8UE
 * Response:
 * 	- 200: database generated identifier for the new device profile
 *	- 400: YAML file is empty
 *	- 409: an associated command's name is a duplicate for the profile or if the name is determined to not be uniqe with regard to others
 * 	- 503: Server Error
 */
func restGetYamlProfileById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars[ID]

	// Check if the device profile exists
	op := device_profile.NewGetProfileID(id, dbClient)
	dp, err := op.Execute()
	if err != nil {
		if err == db.ErrNotFound {
			httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusNotFoundErrorConcept{})
			w.Write([]byte(nil))
			return
		} else {
			httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusInternalServerErrorConcept{})
		}
		LoggingClient.Error(err.Error())
		return
	}

	// Marshal the device profile into YAML
	out, err := yaml.Marshal(dp)
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

// Notify the associated device services for changes in the device profile
func notifyProfileAssociates(dp models.DeviceProfile, dl device.DeviceLoader, action string) error {
	// Get the devices
	op := device.NewProfileIdExecutor(Configuration.Service, dl, LoggingClient, dp.Id)
	d, err := op.Execute()
	if err != nil {
		LoggingClient.Error(err.Error())
		return err
	}

	// Get the services for each device
	// Use map as a Set
	dsMap := map[string]models.DeviceService{}
	ds := []models.DeviceService{}
	for _, device := range d {
		// Only add if not there
		if _, ok := dsMap[device.Service.Id]; !ok {
			dsMap[device.Service.Id] = device.Service
			ds = append(ds, device.Service)
		}
	}

	if err := notifyAssociates(ds, dp.Id, action, models.PROFILE); err != nil {
		LoggingClient.Error(err.Error())
		return err
	}

	return nil
}
