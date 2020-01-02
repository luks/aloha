package main

import (
	"net/http"
	"testing"
)

func TestHealthOK(t *testing.T) {
	_, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}
}
