package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/jsonapi"
	jwtmiddleware "github.com/jgillich/jwt-middleware"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/faststack/machinestack/model"
)

func TestMachineInfo(t *testing.T) {
	machine := model.Machine{
		Name:   "TestMachineInfo",
		Image:  "ubuntu/trusty",
		Driver: "lxd",
		UserID: int64(testToken.Claims.(jwt.MapClaims)["id"].(float64)),
	}

	if err := testDB.Insert(&machine); err != nil {
		t.Fatal(err)
	}

	payload, _ := jsonapi.MarshalOne(&machine)
	buf, _ := json.Marshal(payload)

	r, err := http.NewRequest("GET", "/machines/TestMachineInfo", bytes.NewBuffer(buf))
	if err != nil {
		t.Fatal(err)
	}
	*r = *r.WithContext(context.WithValue(r.Context(), jwtmiddleware.TokenContextKey{}, testToken))

	rr := httptest.NewRecorder()
	router := httprouter.New()
	router.Handle("GET", "/machines/:name", testHandler.MachineInfo)
	router.ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var info model.Machine
	if err := jsonapi.UnmarshalPayload(rr.Body, &info); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(machine, info) {
		t.Errorf("machine info returned wrong content: got '%v' want '%v", info, machine)
	}
}
