package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/jsonapi"

	"gitlab.com/faststack/machinestack/model"
)

func TestMachineCreate(t *testing.T) {
	machine := model.Machine{
		Name: "TestMachineCreate",
	}

	payload, _ := jsonapi.MarshalOne(&machine)
	buf, _ := json.Marshal(payload)

	r, err := http.NewRequest("POST", "/machines", bytes.NewBuffer(buf))
	if err != nil {
		t.Fatal(err)
	}
	*r = *r.WithContext(context.WithValue(r.Context(), UserContextKey, testToken))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(testHandler.MachineCreate)
	handler.ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	if _, ok := mockScheduler.Machines[machine.Name]; !ok {
		t.Error("machine was not created")
	}
}