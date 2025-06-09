package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"user_service/src/internal/adaptors/persistance"
	"user_service/src/internal/core/dto"
	pkgmiddleware "user_service/src/internal/interfaces/input/api/rest/middleware"
	"user_service/src/internal/usecase"
	errorhandling "user_service/src/pkg/error_handling"
	pkgresponse "user_service/src/pkg/response"

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
		errorhandling.HandlerError(w, "failed to register user", http.StatusBadRequest, err)
		return
	}

	// register user using request data
	err = u.userService.RegisterUser(ctx, user)
	if err != nil {
		errorhandling.HandlerError(w, "failed to register user", http.StatusBadRequest, err)
		return
	}

	response := pkgresponse.StandardResponse{
		Status:  "SUCCESS",
		Message: "User Registered Successfully ",
	}
	pkgresponse.WriteResponse(w, http.StatusOK, response)
}

// LOGIN
func (u *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	// get request body
	var req dto.UserDetails
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorhandling.HandlerError(w, "failed to login user", http.StatusBadRequest, err)
		return
	}

	// verify user credentials
	loginResp, err := u.userService.LoginUser(ctx, req)
	if err != nil {
		errorhandling.HandlerError(w, "failed to login user", http.StatusBadRequest, err)
		return
	}

	// create a new session_id
	sessionID := uuid.New().String()

	// store session in Redis: key=session:<sessionID>, value=userID, expires in 24h
	err = persistance.RedisClient.Set(ctx, "session:"+sessionID, loginResp.Id.String(), 24*time.Hour).Err()
	if err != nil {
		errorhandling.HandlerError(w, "failed to create session", http.StatusInternalServerError, err)
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

	response := pkgresponse.StandardResponse{
		Status:  "SUCCESS",
		Data:    loginResp,
		Message: "User Logged-in Successfully ",
	}
	pkgresponse.WriteResponse(w, http.StatusOK, response)
}

func (u *UserHandler) ProfileHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	userID, ok := pkgmiddleware.GetUserIDFromContext(r.Context())
	if !ok {
		errorhandling.HandlerError(w, "failed to get user profile", http.StatusBadRequest, fmt.Errorf("couldn't get user_id from context"))
		return
	}

	userDetails, err := u.userService.GetUserProfile(ctx, userID)
	if err != nil {
		errorhandling.HandlerError(w, "failed to get user profile", http.StatusBadRequest, err)
		return
	}

	pkgresponse.WriteResponse(w, http.StatusOK, pkgresponse.StandardResponse{
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
		errorhandling.HandlerError(w, "failed to update user profile", http.StatusBadRequest, nil)
		return
	}

	var reqData dto.UserDetails
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		errorhandling.HandlerError(w, "failed to update user profile", http.StatusBadRequest, nil)
		return
	}

	userInfo, err := u.userService.UpdateUserProfile(ctx, userId, reqData)
	if err != nil {
		errorhandling.HandlerError(w, "failed to update user data", http.StatusBadRequest, nil)
		return
	}

	response := pkgresponse.StandardResponse{
		Status:  "SUCCESS",
		Data:    userInfo,
		Message: "user profile updated successfully ",
	}
	pkgresponse.WriteResponse(w, http.StatusOK, response)
}
