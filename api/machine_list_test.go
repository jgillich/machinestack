package api

import (
	"context"
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

func TestMachineList(t *testing.T) {
	machine := model.Machine{
		Name:   "TestMachineList",
		Image:  "ubuntu/trusty",
		Driver: "lxd",
		UserID: testToken.Claims.(jwt.MapClaims)["id"].(int64),
	}
	if err := testDB.Insert(&machine); err != nil {
		t.Fatal(err)
	}

	r, err := http.NewRequest("GET", "/machines", nil)
	if err != nil {
		t.Fatal(err)
	}
	*r = *r.WithContext(context.WithValue(r.Context(), jwtmiddleware.TokenContextKey{}, testToken))

	rr := httptest.NewRecorder()
	router := httprouter.New()
	router.Handle("GET", "/machines", testHandler.MachineList)
	router.ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	machines, err := jsonapi.UnmarshalManyPayload(rr.Body, reflect.TypeOf(new(model.Machine)))
	if err != nil {
		t.Fatal(err)
	}

	found := false
	for _, mm := range machines {
		m, ok := mm.(*model.Machine)
		if !ok {
			t.Errorf("returned model is not machine: '%v'", mm)
		}
		if m.Name == machine.Name {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("machine was not included in list")
	}
}
