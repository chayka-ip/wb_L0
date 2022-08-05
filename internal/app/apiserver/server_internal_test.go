package apiserver

import (
	"bytes"
	"encoding/json"
	"level_zero/internal/app/cache/lru_cache"
	models "level_zero/internal/app/model"
	"level_zero/store/teststore"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func GetTestConfig() *Config {
	return &Config{
		LogLevel:    "debug",
		CacheSize:   1024,
		DatabaseURL: "user=postgres password=postgres host=localhost dbname=orders_db port=5432 sslmode=disable",
	}
}

type DataTest struct {
	name         string
	payload      interface{}
	expectedCode int
}

func NewTestServer() *server {
	config := GetTestConfig()
	store := teststore.New()
	logger := getNewLogger(config)
	stanInfo := config.GetNatsInfo()
	cache, err := lru_cache.NewCache(config.CacheSize, store, logger)
	if err != nil {
		log.Fatal("Failed to create cache")
	}
	return NewServer(store, cache, logger, &stanInfo)
}

func TestServer_HandleOrderCreate(t *testing.T) {
	s := NewTestServer()

	testCases := []DataTest{
		{
			name:         "successfully created",
			payload:      models.TestOrder().Data,
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid payload",
			payload:      models.TestOrderBadData().Data,
			expectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req := httptest.NewRequest(http.MethodPost, routeOrderCreate, b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}

	// test can't add duplicates

	p := models.TestOrder().Data
	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(p)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, routeOrderCreate, b)
	s.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusCreated)

	rec = httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	assert.Equal(t, rec.Code, http.StatusBadRequest)
}

func TestServer_HandleGetOrderByUid(t *testing.T) {
	s := NewTestServer()

	ord := models.TestOrder()
	addRec := httptest.NewRecorder()
	b := &bytes.Buffer{}
	json.NewEncoder(b).Encode(ord)
	addReq := httptest.NewRequest(http.MethodPost, routeOrderCreate, b)
	s.ServeHTTP(addRec, addReq)
	assert.Equal(t, addRec.Code, http.StatusCreated)

	testCases := []DataTest{
		{
			name:         "get existing",
			payload:      ord.OrderUID,
			expectedCode: http.StatusOK,
		},
		{
			name:         "get not existing",
			payload:      models.TestOrder().OrderUID,
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, routeGetOrderByUid, nil)
			q := req.URL.Query()
			q.Add("id", tc.payload.(string))
			req.URL.RawQuery = q.Encode()
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
