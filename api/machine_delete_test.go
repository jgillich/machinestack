package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.com/faststack/machinestack/model"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/julienschmidt/httprouter"
)

func TestMachineDelete(t *testing.T) {
	machine := model.Machine{
		Name:   "TestMachineDelete",
		Image:  "ubuntu/trusty",
		Driver: "lxd",
		UserID: testToken.Claims.(jwt.MapClaims)["id"].(int),
		Node:   "node",
	}

	if err := testDB.Insert(&machine); err != nil {
		t.Fatal(err)
	}

	mockScheduler.Machines["TestMachineDelete"] = "node"

	r, err := http.NewRequest("DELETE", "/machines/TestMachineDelete", nil)
	if err != nil {
		t.Fatal(err)
	}
	*r = *r.WithContext(context.WithValue(r.Context(), UserContextKey, testToken))

	rr := httptest.NewRecorder()
	router := httprouter.New()
	router.Handle("DELETE", "/machines/:name", testHandler.MachineDelete)
	router.ServeHTTP(rr, r)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if _, ok := mockScheduler.Machines["TestMachineDelete"]; ok {
		t.Error("machine was not deleted")
	}
}
