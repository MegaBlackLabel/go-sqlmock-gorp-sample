package main

import (
	"fmt"
	_ "log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gorp.v1"
)

func TestGetUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	ts := GetMainEngine(dbMap)

	req, err := http.NewRequest("GET", "/api/v1/users", nil)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when making API request", err)
	}

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "TestUser1").AddRow(2, "TestUser2")

	mock.ExpectQuery("^SELECT (.+) FROM user$").WillReturnRows(rows)

	resp := httptest.NewRecorder()
	ts.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
	assert.JSONEq(t, resp.Body.String(), `[{"id":1,"name":"TestUser1"},{"id":2,"name":"TestUser2"} ]`)
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}

func TestGetSpecificUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		fmt.Sprintf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}
	ts := GetMainEngine(dbMap)

	req, err := http.NewRequest("GET", "/api/v1/users/1", nil)
	if err != nil {
		t.Fatalf("an error '%s' was not expected when making API request", err)
	}

	rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "TestUser1")

	//mock.ExpectQuery("^SELECT (.+) FROM user where id=?$").WithArgs(1).WillReturnRows(rows)
	//mock.ExpectQuery("^SELECT (.+) FROM user where id=? LIMIT 1$").WithArgs(1).WillReturnRows(rows)
	//mock.ExpectQuery(`SELECT \* FROM user WHERE id=? LIMIT 1`).WithArgs(1).WillReturnRows(rows)
	//mock.ExpectQuery("^SELECT (.+) FROM user where id=? LIMIT 1").WithArgs(1).WillReturnRows(rows)
	// NOTE: Have to escape "?"
	mock.ExpectQuery("^SELECT \\* FROM user WHERE id=\\? LIMIT 1$").WithArgs("1").WillReturnRows(rows)

	resp := httptest.NewRecorder()
	ts.ServeHTTP(resp, req)
	assert.Equal(t, resp.Code, 200)
	fmt.Println(resp.Body.String())
	assert.JSONEq(t, resp.Body.String(), `{"id":1,"name":"TestUser1"}`)

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expections: %s", err)
	}
}
