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
	"net/http"
	"net/url"
	"strconv"

	"github.com/edgexfoundry/edgex-go/internal/core/metadata/operators/addressable"
	"github.com/edgexfoundry/edgex-go/internal/pkg/httperror"
	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/models"
	"github.com/gorilla/mux"
)

func restGetAllAddressables(w http.ResponseWriter, _ *http.Request) {
	op := addressable.NewAddressableLoadAll(Configuration.Service, dbClient, LoggingClient)
	addressables, err := op.Execute()
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.StatusRequestEntityTooLargeErrorConcept{}}, httperror.StatusInternalServerErrorConcept{})
		return
	}
	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	err = json.NewEncoder(w).Encode(&addressables)
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.StatusInternalServerErrorConcept{})
		return
	}
}

// Add a new addressable
// The name must be unique
func restAddAddressable(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var a models.Addressable
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.EmptyAddressableNameErrorConcept{})
		return
	}

	op := addressable.NewAddExecutor(dbClient, a)
	id, err := op.Execute()
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.DuplicateIdentifierErrorConcept{}, httperror.EmptyAddressableNameErrorConcept{}}, httperror.StatusInternalServerErrorConcept{})
		return
	}
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(id))
	if err != nil {
		LoggingClient.Error(err.Error())
		return
	}
}

// Update addressable by ID or name (ID used first)
func restUpdateAddressable(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var a models.Addressable
	err := json.NewDecoder(r.Body).Decode(&a)
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{}, httperror.EmptyAddressableNameErrorConcept{})
		return
	}

	op := addressable.NewUpdateExecutor(dbClient, a)
	err = op.Execute()
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.AddressableNotFoundErrorConcept{}, httperror.AddressableInUseErrorConcept{}}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("true"))
	if err != nil {
		LoggingClient.Error(err.Error())
		return
	}
}

func restGetAddressableById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id string = vars["id"]
	op := addressable.NewIdExecutor(dbClient, id)
	result, err := op.Execute()
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.DatabaseNotFoundErrorConcept{}}, httperror.StatusInternalServerErrorConcept{})
		return
	}
	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	json.NewEncoder(w).Encode(result)
}

func restDeleteAddressableById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id string = vars[ID]

	op := addressable.NewDeleteByIdExecutor(dbClient, id)
	err := op.Execute()
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.AddressableNotFoundErrorConcept{}, httperror.AddressableInUseErrorConcept{}}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	w.Write([]byte("true"))
}

func restDeleteAddressableByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name, err := url.QueryUnescape(vars[NAME])
	// Problems unescaping
	if err != nil {
		LoggingClient.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	op := addressable.NewDeleteByNameExecutor(dbClient, name)
	err = op.Execute()
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.AddressableNotFoundErrorConcept{}, httperror.AddressableInUseErrorConcept{}}, httperror.StatusInternalServerErrorConcept{})
		return
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("true"))
}

func restGetAddressableByName(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	dn, err := url.QueryUnescape(vars[NAME])
	if err != nil {
		LoggingClient.Error(err.Error())
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	op := addressable.NewNameExecutor(dbClient, dn)
	result, err := op.Execute()
	if err != nil {
		httperror.ToHttpError(w, LoggingClient, err, []httperror.ErrorConceptType{httperror.DatabaseNotFoundErrorConcept{}}, httperror.AddressableNotFoundByNameErrorConcept{})
		return
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	json.NewEncoder(w).Encode(result)
}

func restGetAddressableByTopic(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	t, err := url.QueryUnescape(vars[TOPIC])
	if err != nil {
		LoggingClient.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	op := addressable.NewTopicExecutor(dbClient, t)
	res, err := op.Execute()
	if err != nil {
		LoggingClient.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	json.NewEncoder(w).Encode(res)
}
func restGetAddressableByPort(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var strp string = vars[PORT]
	p, err := strconv.Atoi(strp)
	if err != nil {
		LoggingClient.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	op := addressable.NewPortExecutor(dbClient, p)
	res, err := op.Execute()
	if err != nil {
		LoggingClient.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	json.NewEncoder(w).Encode(res)
}
func restGetAddressableByPublisher(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	p, err := url.QueryUnescape(vars[PUBLISHER])
	if err != nil {
		LoggingClient.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	op := addressable.NewPublisherExecutor(dbClient, p)
	res, err := op.Execute()
	if err != nil {
		LoggingClient.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	json.NewEncoder(w).Encode(res)
}
func restGetAddressableByAddress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	a, err := url.QueryUnescape(vars[ADDRESS])
	if err != nil {
		LoggingClient.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	op := addressable.NewAddressExecutor(dbClient, a)
	res, err := op.Execute()
	if err != nil {
		LoggingClient.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set(clients.ContentType, clients.ContentTypeJSON)
	json.NewEncoder(w).Encode(res)
}
