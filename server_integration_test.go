//go:build integration

package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestGetAllExpense(t *testing.T) {
	seedExpense(t)
	var ex []Expense

	res := request(http.MethodGet, uri("expenses"), nil)
	err := res.Decode(&ex)

	assert.Nil(t, err)
	assert.EqualValues(t, http.StatusOK, res.StatusCode)
	assert.Greater(t, len(ex), 0)
}

func TestCreateExpense(t *testing.T) {
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie",
		"amount": 79,
		"note": "night market promotion discount 10 bath", 
		"tags": ["food", "beverage"]
	}`)
	var e Expense

	res := request(http.MethodPost, uri("expenses"), body)
	err := res.Decode(&e)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.NotEqual(t, 0, e.ID)
	assert.Equal(t, "strawberry smoothie", e.Title)
	assert.Equal(t, float32(79), e.Amount)
	assert.Equal(t, "night market promotion discount 10 bath", e.Note)
	assert.Contains(t, e.Tags, "food")
}

func TestGetExpenseByID(t *testing.T) {
	c := seedExpense(t)

	var latest Expense
	res := request(http.MethodGet, uri("expenses", strconv.Itoa(c.ID)), nil)
	err := res.Decode(&latest)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, c.ID, latest.ID)
	assert.NotEmpty(t, latest.Title)
	assert.NotEmpty(t, latest.Amount)
	assert.NotEmpty(t, latest.Note)
	assert.NotEmpty(t, latest.Tags)
}

func TestUpdateExpenseByID(t *testing.T) {
	body := bytes.NewBufferString(`{
		"title": "strawberry smoothie [edited]",
		"amount": 799,
		"note": "night market promotion discount 50 bath", 
		"tags": ["food", "beverage"]
	}`)
	c := seedExpense(t)

	var update Expense
	res := request(http.MethodPut, uri("expenses", strconv.Itoa(c.ID)), body)
	err := res.Decode(&update)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusCreated, res.StatusCode)
	assert.Equal(t, c.ID, update.ID)
	assert.Equal(t, "strawberry smoothie [edited]", update.Title)
	assert.Equal(t, float32(799), update.Amount)
	assert.Equal(t, "night market promotion discount 50 bath", update.Note)
	assert.Contains(t, update.Tags, "food")
}

func seedExpense(t *testing.T) Expense {
	var c Expense
	body := bytes.NewBufferString(`{
		"title": "apple smoothie",
    "amount": 89,
    "note": "no discount",
    "tags": ["beverage"]
	}`)
	err := request(http.MethodPost, uri("expenses"), body).Decode(&c)
	if err != nil {
		t.Fatal("can't create expense", err)
	}
	return c
}

func uri(paths ...string) string {
	host := "http://localhost:2565"
	if paths == nil {
		return host
	}

	url := append([]string{host}, paths...)
	return strings.Join(url, "/")
}

func request(method, url string, body io.Reader) *Response {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", os.Getenv("AUTH_TOKEN"))
	req.Header.Add("Content-Type", "application/json")
	client := http.Client{}
	res, err := client.Do(req)
	return &Response{res, err}
}

type Response struct {
	*http.Response
	err error
}

func (r *Response) Decode(v interface{}) error {
	if r.err != nil {
		return r.err
	}

	return json.NewDecoder(r.Body).Decode(v)
}

// ❯ AUTH_TOKEN="Basic YXBpZGVzaWduOjQ1Njc4" go test -v -tags=integration ./...
// === RUN   TestGetAllExpense
// --- PASS: TestGetAllExpense (0.35s)
// === RUN   TestCreateExpense
// --- PASS: TestCreateExpense (0.15s)
// === RUN   TestGetExpenseByID
// --- PASS: TestGetExpenseByID (0.30s)
// === RUN   TestUpdateExpenseByID
// --- PASS: TestUpdateExpenseByID (0.31s)
// PASS
// ok  	github.com/openmymai/assessment	1.667s
