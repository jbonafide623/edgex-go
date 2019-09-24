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
	"errors"
	"net/http"
	"net/url"

	"github.com/edgexfoundry/edgex-go/internal/pkg/httperror"
	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/gorilla/mux"
)

func restGetProvisionWatchers(w http.ResponseWriter, _ *http.Request) {
	res, err := dbClient.GetAllProvisionWatchers()
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusServiceUnavailableErrorConcept{})
		return
	}

	// Check the length
	if len(res) > Configuration.Service.MaxResultCount {
		err := errors.New("Max limit exceeded")
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusRequestEntityTooLargeErrorConcept{})
		return
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	json.NewEncoder(w).Encode(&res)
}

func restDeleteProvisionWatcherById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id string = vars[ID]

	// Check if the provision watcher exists
	pw, err := dbClient.GetProvisionWatcherById(id)
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusNotFoundErrorConcept{})
		return
	}

	err = deleteProvisionWatcher(pw, w)
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusNotFoundErrorConcept{})
		return
	}
	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	w.Write([]byte("true"))
}

func restDeleteProvisionWatcherByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, err := url.QueryUnescape(vars[NAME])
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusBadRequestErrorConcept{})
		return
	}

	// Check if the provision watcher exists
	pw, err := dbClient.GetProvisionWatcherByName(n)
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.DatabaseNotFoundErrorConcept{}}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	if err = deleteProvisionWatcher(pw, w); err != nil {
		LoggingClient.Error("Problem deleting provision watcher: " + err.Error())
		return
	}
	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("true"))
}

// Delete the provision watcher
func deleteProvisionWatcher(pw models.ProvisionWatcher, w http.ResponseWriter) error {
	if err := dbClient.DeleteProvisionWatcherById(pw.Id); err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusServiceUnavailableErrorConcept{})
		return err
	}

	if err := notifyProvisionWatcherAssociates(pw, http.MethodDelete); err != nil {
		LoggingClient.Error("Problem notifying associated device services to provision watcher: " + err.Error())
	}

	return nil
}

func restGetProvisionWatcherById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id string = vars[ID]

	res, err := dbClient.GetProvisionWatcherById(id)
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.DatabaseNotFoundErrorConcept{}}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	json.NewEncoder(w).Encode(res)
}

func restGetProvisionWatcherByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, err := url.QueryUnescape(vars[NAME])
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusBadRequestErrorConcept{})
		return
	}

	res, err := dbClient.GetProvisionWatcherByName(n)
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.DatabaseNotFoundErrorConcept{}}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	json.NewEncoder(w).Encode(res)
}

func restGetProvisionWatchersByProfileId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var pid string = vars[ID]

	// Check if the device profile exists
	if _, err := dbClient.GetDeviceProfileById(pid); err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusNotFoundErrorConcept{})
		return
	}

	res, err := dbClient.GetProvisionWatchersByProfileId(pid)
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	json.NewEncoder(w).Encode(res)
}

func restGetProvisionWatchersByProfileName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pn, err := url.QueryUnescape(vars[NAME])
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusBadRequestErrorConcept{})
		return
	}

	// Check if the device profile exists
	dp, err := dbClient.GetDeviceProfileByName(pn)
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.DatabaseNotFoundErrorConcept{}}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	res, err := dbClient.GetProvisionWatchersByProfileId(dp.Id)
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	json.NewEncoder(w).Encode(res)
}

func restGetProvisionWatchersByServiceId(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var sid string = vars[ID]

	// Check if the device service exists
	if _, err := dbClient.GetDeviceServiceById(sid); err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusNotFoundErrorConcept{})
		return
	}

	res, err := dbClient.GetProvisionWatchersByServiceId(sid)
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	json.NewEncoder(w).Encode(res)
}

func restGetProvisionWatchersByServiceName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sn, err := url.QueryUnescape(vars[NAME])
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusBadRequestErrorConcept{})
		return
	}

	// Check if the device service exists
	ds, err := dbClient.GetDeviceServiceByName(sn)
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.DatabaseNotFoundErrorConcept{}}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	// Get the provision watchers
	res, err := dbClient.GetProvisionWatchersByServiceId(ds.Id)
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusNotFoundErrorConcept{})
		return
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	json.NewEncoder(w).Encode(res)
}

func restGetProvisionWatchersByIdentifier(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	k, err := url.QueryUnescape(vars[KEY])
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusBadRequestErrorConcept{})
		return
	}
	v, err := url.QueryUnescape(vars[VALUE])
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusBadRequestErrorConcept{})
		return
	}

	res, err := dbClient.GetProvisionWatchersByIdentifier(k, v)
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	json.NewEncoder(w).Encode(res)
}

func restAddProvisionWatcher(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var pw models.ProvisionWatcher
	var err error

	if err = json.NewDecoder(r.Body).Decode(&pw); err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusServiceUnavailableErrorConcept{})
		return
	}

	// Check if the name exists
	if pw.Name == "" {
		err = errors.New("No name provided for new provision watcher")
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusConflictErrorConcept{})
		return
	}

	// Check if the device profile exists
	// Try by ID
	var profile models.DeviceProfile
	if pw.Profile.Id != "" {
		profile, err = dbClient.GetDeviceProfileById(pw.Profile.Id)
	}
	if pw.Profile.Id == "" || err != nil {
		// Try by name
		if profile, err = dbClient.GetDeviceProfileByName(pw.Profile.Name); err != nil {
			httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.ProvisionWatcherDeviceServiceNotFoundErrorConcept{}}, httperror.StatusServiceUnavailableErrorConcept{})
			return
		}
	}
	pw.Profile = profile

	// Check if the device service exists
	// Try by ID
	var service models.DeviceService
	if pw.Service.Id != "" {
		service, err = dbClient.GetDeviceServiceById(pw.Service.Id)
	}
	if pw.Service.Id == "" || err != nil {
		// Try by name
		if service, err = dbClient.GetDeviceServiceByName(pw.Service.Name); err != nil {
			httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.ProvisionWatcherDeviceServiceNotFoundErrorConcept{}}, httperror.StatusServiceUnavailableErrorConcept{})
			return
		}
	}
	pw.Service = service

	id, err := dbClient.AddProvisionWatcher(pw)
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.DatabaseNotUniqueErrorConcept{}}, httperror.StatusServiceUnavailableErrorConcept{})
		return
	}

	// Notify Associates
	if err = notifyProvisionWatcherAssociates(pw, http.MethodPost); err != nil {
		LoggingClient.Error("Problem with notifying associating device services for the provision watcher: " + err.Error())
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(id))
}

// Update the provision watcher object
// ID is used first for identification, then name
// The service and profile cannot be updated
func restUpdateProvisionWatcher(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var from models.ProvisionWatcher
	if err := json.NewDecoder(r.Body).Decode(&from); err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusServiceUnavailableErrorConcept{})
		return
	}

	// Check if the provision watcher exists
	// Try by ID
	to, err := dbClient.GetProvisionWatcherById(from.Id)
	if err != nil {
		// Try by name
		if to, err = dbClient.GetProvisionWatcherByName(from.Name); err != nil {
			httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.DatabaseNotFoundErrorConcept{}}, httperror.StatusServiceUnavailableErrorConcept{})
			return
		}
	}

	if err := updateProvisionWatcherFields(from, &to, w); err != nil {
		LoggingClient.Error("Problem updating provision watcher: " + err.Error())
		return
	}

	if err := dbClient.UpdateProvisionWatcher(to); err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusServiceUnavailableErrorConcept{})
		return
	}

	// Notify Associates
	if err := notifyProvisionWatcherAssociates(to, http.MethodPut); err != nil {
		LoggingClient.Error("Problem notifying associated device services for provision watcher: " + err.Error())
	}
	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("true"))
}

// Update the relevant fields of the provision watcher
func updateProvisionWatcherFields(from models.ProvisionWatcher, to *models.ProvisionWatcher, w http.ResponseWriter) error {
	if from.Identifiers != nil {
		to.Identifiers = from.Identifiers
	}
	if from.Origin != 0 {
		to.Origin = from.Origin
	}
	if from.Name != "" {
		// Check that the name is unique
		checkPW, err := dbClient.GetProvisionWatcherByName(from.Name)
		if err != nil {
			// DuplicateProvisionWatcherErrorConcept will evaluate to true if the ID is a duplicate
			httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.ProvisionWatcherDuplicateErrorConcept{CurrentPWId: checkPW.Id, UpdatedPWId: to.Id}}, httperror.StatusServiceUnavailableErrorConcept{})
		}
		to.Name = from.Name
	}

	return nil
}

// Notify the associated device services for the provision watcher
func notifyProvisionWatcherAssociates(pw models.ProvisionWatcher, action string) error {
	// Get the device service for the provision watcher
	ds, err := dbClient.GetDeviceServiceById(pw.Service.Id)
	if err != nil {
		return err
	}

	// Notify the service
	if err = notifyAssociates([]models.DeviceService{ds}, pw.Id, action, models.PROVISIONWATCHER); err != nil {
		return err
	}

	return nil
}
