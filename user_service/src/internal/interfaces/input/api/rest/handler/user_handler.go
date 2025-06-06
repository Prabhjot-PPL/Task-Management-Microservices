package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"user_service/src/internal/adaptors/persistance"
	"user_service/src/internal/core/dto"
	pkgmiddleware "user_service/src/internal/interfaces/input/api/rest/middleware"
	"user_service/src/internal/usecase"
	"user_service/src/pkg"

	"github.com/google/uuid"
)

type UserHandler struct {
	userService usecase.Service
}

func NewUserHandler(userService usecase.Service) *UserHandler {
	return &UserHandler{userService: userService}
}

func (u *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	var user dto.UserDetails

	// decode request body
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response := pkg.StandardResponse{
			Status:  "FAILURE",
			Message: "failed to register user",
		}
		pkg.WriteResponse(w, http.StatusBadRequest, response)
		return
	}

	// register user using request data
	err = u.userService.RegisterUser(ctx, user)
	if err != nil {
		response := pkg.StandardResponse{
			Status:  "FAILURE",
			Message: "failed to register user",
		}
		pkg.WriteResponse(w, http.StatusBadRequest, response)
		return
	}

	response := pkg.StandardResponse{
		Status:  "SUCCESS",
		Message: "User Registered Successfully ",
	}
	pkg.WriteResponse(w, http.StatusOK, response)
}

// LOGIN
func (u *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// get request body
	var req dto.UserDetails
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response := pkg.StandardResponse{
			Status:  "FAILURE",
			Message: "failed to login user",
		}
		pkg.WriteResponse(w, http.StatusBadRequest, response)
		return
	}

	// verify user credentials
	loginResp, err := u.userService.LoginUser(ctx, req)
	if err != nil {
		response := pkg.StandardResponse{
			Status:  "FAILURE",
			Message: "failed to login user",
		}
		pkg.WriteResponse(w, http.StatusBadRequest, response)
		return
	}

	// create a new session_id
	sessionID := uuid.New().String()

	// store session in Redis: key=session:<sessionID>, value=userID, expires in 24h
	err = persistance.RedisClient.Set(ctx, "session:"+sessionID, loginResp.Id.String(), 24*time.Hour).Err()
	if err != nil {
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	// Set session_id cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	response := pkg.StandardResponse{
		Status:  "SUCCESS",
		Data:    loginResp,
		Message: "User Logged-in Successfully ",
	}
	pkg.WriteResponse(w, http.StatusOK, response)
}

func (u *UserHandler) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	userID, ok := pkgmiddleware.GetUserIDFromContext(r.Context())
	if !ok {
		response := pkg.StandardResponse{
			Status:  "FAILURE",
			Message: "failed to get user profile",
		}
		pkg.WriteResponse(w, http.StatusBadRequest, response)
		return
	}

	userDetails, err := u.userService.GetUserProfile(ctx, userID)
	if err != nil {
		response := pkg.StandardResponse{
			Status:  "FAILURE",
			Message: "failed to get user profile",
		}
		pkg.WriteResponse(w, http.StatusBadRequest, response)
		return
	}

	pkg.WriteResponse(w, http.StatusOK, pkg.StandardResponse{
		Status:  "success",
		Data:    userDetails,
		Message: "User profile fetched successfully",
	})
}

func (u *UserHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	userId, ok := pkgmiddleware.GetUserIDFromContext(r.Context())
	if !ok {
		response := pkg.StandardResponse{
			Status:  "FAILURE",
			Message: "failed to get update user",
		}
		pkg.WriteResponse(w, http.StatusBadRequest, response)
		return
	}

	var reqData dto.UserDetails
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		response := pkg.StandardResponse{
			Status:  "FAILURE",
			Message: "failed to update user data",
		}
		pkg.WriteResponse(w, http.StatusBadRequest, response)
		return
	}

	userInfo, err := u.userService.UpdateUserProfile(ctx, userId, reqData)
	if err != nil {
		response := pkg.StandardResponse{
			Status:  "FAILURE",
			Message: "failed to update user data",
		}
		pkg.WriteResponse(w, http.StatusBadRequest, response)
		return
	}

	response := pkg.StandardResponse{
		Status:  "SUCCESS",
		Data:    userInfo,
		Message: "user profile updated successfully ",
	}
	pkg.WriteResponse(w, http.StatusOK, response)
}
