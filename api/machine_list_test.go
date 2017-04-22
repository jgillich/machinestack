package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/google/jsonapi"
	"github.com/julienschmidt/httprouter"
	"gitlab.com/faststack/machinestack/model"
)

func TestMachineList(t *testing.T) {
	machine := model.Machine{
		Name:   "TestMachineList",
		Image:  "ubuntu/trusty",
		Driver: "lxd",
		UserID: testToken.Claims.(jwt.MapClaims)["id"].(int),
	}
	if err := testDB.Insert(&machine); err != nil {
		t.Fatal(err)
	}

	r, err := http.NewRequest("GET", "/machines", nil)
	if err != nil {
		t.Fatal(err)
	}
	*r = *r.WithContext(context.WithValue(r.Context(), UserContextKey, testToken))

	rr := httptest.NewRecorder()
	router := httprouter.New()
	router.Handle("GET", "/machines", testHandler.MachineList)
	router.ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var machines []model.Machine
	if err := jsonapi.UnmarshalPayload(rr.Body, machines); err != nil {
		t.Fatal(err)
	}

	found := false
	for _, m := range machines {
		if m.Name == machine.Name {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("machine was not included in list")
	}
}
