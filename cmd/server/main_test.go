package main

import "testing"

func TestRequiredDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")

	_, err := requiredDatabaseURL()
	if err == nil {
		t.Fatal("se esperaba un error cuando DATABASE_URL está vacía")
	}
}
