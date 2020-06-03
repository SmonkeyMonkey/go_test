package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUsersHandler(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest(http.MethodGet, "/users", nil)
	usersHandler(w, r)
	if w.Code != http.StatusOK {
		t.Errorf("got %v status code", w.Code)
	}
}
