package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"project/helper"
	"project/internals/service"
	"project/internals/storage"
)

type HttpService struct {
	service service.Storer
}

func NewHttpService(serive service.Storer) *HttpService {
	return &HttpService{service: serive}
}

func (h *HttpService) HttpCreate(w http.ResponseWriter, r *http.Request) {

	var DTO storage.UserDTO

	if err := json.NewDecoder(r.Body).Decode(&DTO); err != nil {
		helper.Error(w, err, http.StatusBadRequest)
		return
	}

	user, err := h.service.Create(
		DTO.FirstName, 
		DTO.LastName, 
		DTO.Phone, 
		DTO.CarDTO)
	if err != nil {
		helper.Error(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(user); err != nil {
		helper.Error(w, err, http.StatusInternalServerError)
		return
	}

}

func (h *HttpService) HttpGet(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	user, err := h.service.Get(id)
	if err != nil {
		helper.Error(w, err, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(user); err != nil {
		helper.Error(w, err, http.StatusInternalServerError)
		return
	}

}

func (h *HttpService) HttpGetAll(w http.ResponseWriter, r *http.Request) {
	users := h.service.GetAll()

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(users); err != nil {
		helper.Error(w, err, http.StatusInternalServerError)
		return
	}
}

func (h *HttpService) HttpUpdate(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	

	var DTO storage.UserDTO

	if err := json.NewDecoder(r.Body).Decode(&DTO); err != nil {
		helper.Error(w, err, http.StatusBadRequest)
		return
	}

	user, err := h.service.Update(
		id,
		DTO.FirstName,
		DTO.LastName,
		DTO.Phone,
		DTO.CarDTO)

	if err != nil {
		if errors.Is(err, helper.UserNotFound) ||
			errors.Is(err, helper.CarNotFound) {
			helper.Error(w, err, http.StatusNotFound)
			return

		}

		helper.Error(w, err, http.StatusInternalServerError)
		return

	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(user)
}

func (h *HttpService) HttpDelete(w http.ResponseWriter, r *http.Request) {

	id := r.PathValue("id")

	err := h.service.Delete(id)
	if err != nil {
		if errors.Is(err, helper.UserNotFound) {
			helper.Error(w, err, http.StatusNotFound)
			return
		}

		helper.Error(w, err, http.StatusInternalServerError)
		return

	}

	w.WriteHeader(http.StatusNoContent)

}
