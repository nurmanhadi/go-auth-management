package handler

import (
	"auth-management/internal/service"
	"auth-management/pkg/dto"
	"auth-management/pkg/response"
	"net/http"

	"github.com/goccy/go-json"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}
func (h *UserHandler) UserRegister(w http.ResponseWriter, r *http.Request) {
	request := new(dto.UserRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		panic(response.Except(http.StatusBadRequest, "failed decode body to json"))
	}
	err := h.userService.UserRegister(request)
	if err != nil {
		panic(err)
	}
	response.Success(w, http.StatusCreated, "OK", r.URL.Path)
}
func (h *UserHandler) UserLogin(w http.ResponseWriter, r *http.Request) {
	request := new(dto.UserRequest)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		panic(response.Except(http.StatusBadRequest, "failed decode body to json"))
	}
	result, err := h.userService.UserLogin(request)
	if err != nil {
		panic(err)
	}
	response.Success(w, http.StatusOK, result, r.URL.Path)
}
func (h *UserHandler) UserGenerateToken(w http.ResponseWriter, r *http.Request) {
	request := new(dto.UserGenerateToken)
	if err := json.NewDecoder(r.Body).Decode(request); err != nil {
		panic(response.Except(http.StatusBadRequest, "failed decode body to json"))
	}
	result, err := h.userService.UserGenerateToken(request)
	if err != nil {
		panic(err)
	}
	response.Success(w, http.StatusOK, result, r.URL.Path)
}
