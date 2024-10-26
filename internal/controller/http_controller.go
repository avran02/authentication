package controller

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/avran02/authentication/internal/dto"
	"github.com/avran02/authentication/internal/service"
)

type HTTPController interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	RefreshTokens(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type httpController struct {
	service service.Service
}

func (c *httpController) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiError(w, http.StatusBadRequest, err)
		return
	}

	id, accessToken, refreshToken, expTime, err := c.service.Register(r.Context(), req.Username, req.Password, req.Email)
	if err != nil {
		apiError(w, http.StatusInternalServerError, err)
		return
	}

	resp := dto.RegisterResponse{
		ID:          id,
		AccessToken: accessToken,
	}

	c.setRefreshTokenCookie(w, refreshToken, expTime)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		apiError(w, http.StatusInternalServerError, err)
		return
	}
}

func (c *httpController) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiError(w, http.StatusBadRequest, err)
		return
	}

	id, accessToken, refreshToken, expTime, err := c.service.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		apiError(w, http.StatusInternalServerError, err)
		return
	}

	resp := dto.LoginResponse{
		ID:          id,
		AccessToken: accessToken,
	}

	c.setRefreshTokenCookie(w, refreshToken, expTime)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		apiError(w, http.StatusInternalServerError, err)
		return
	}
}

func (c *httpController) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	cookies := r.Cookies()
	var refreshToken string
	for _, cookie := range cookies {
		if cookie.Name == "refreshToken" {
			refreshToken = cookie.Value
		}
	}
	newAccessToken, newRefreshToken, expTime, err := c.service.RefreshTokens(r.Context(), refreshToken)
	if err != nil {
		apiError(w, http.StatusInternalServerError, err)
		return
	}

	resp := dto.RefreshTokenResponse{
		AccessToken: newAccessToken,
	}

	c.setRefreshTokenCookie(w, newRefreshToken, expTime)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		apiError(w, http.StatusInternalServerError, err)
		return
	}
}

func (c *httpController) Logout(w http.ResponseWriter, r *http.Request) {
	var req dto.LogoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apiError(w, http.StatusBadRequest, err)
		return
	}

	ok, err := c.service.Logout(r.Context(), req.AccessToken)
	if err != nil {
		apiError(w, http.StatusInternalServerError, err)
		return
	}

	resp := dto.LogoutResponse{
		OK: ok,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		apiError(w, http.StatusInternalServerError, err)
		return
	}
}

func (c *httpController) setRefreshTokenCookie(w http.ResponseWriter, refreshToken string, expTime time.Time) {
	cookie := http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		Path:     "/",
		Expires:  expTime,
		HttpOnly: true,
		// Secure:   true,
		SameSite: http.SameSiteDefaultMode,
	}
	http.SetCookie(w, &cookie)
}

func newHTTPController(service service.Service) HTTPController {
	return &httpController{
		service: service,
	}
}
