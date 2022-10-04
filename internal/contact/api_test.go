package contact

import (
	"github.com/qiangxue/go-rest-api/internal/auth"
	"github.com/qiangxue/go-rest-api/internal/entity"
	"github.com/qiangxue/go-rest-api/internal/test"
	"github.com/qiangxue/go-rest-api/pkg/log"
	"net/http"
	"testing"
	"time"
)

func TestAPI(t *testing.T) {
	logger, _ := log.NewForTest()
	router := test.MockRouter(logger)
	repo := &mockRepository{items: []entity.contact{
		{"123", "contact123", time.Now(), time.Now()},
	}}
	RegisterHandlers(router.Group(""), NewService(repo, logger), auth.MockAuthHandler, logger)
	header := auth.MockAuthHeader()

	tests := []test.APITestCase{
		{"get all", "GET", "/contacts", "", nil, http.StatusOK, `*"total_count":1*`},
		{"get 123", "GET", "/contacts/123", "", nil, http.StatusOK, `*contact123*`},
		{"get unknown", "GET", "/contacts/1234", "", nil, http.StatusNotFound, ""},
		{"create ok", "POST", "/contacts", `{"name":"test"}`, header, http.StatusCreated, "*test*"},
		{"create ok count", "GET", "/contacts", "", nil, http.StatusOK, `*"total_count":2*`},
		{"create auth error", "POST", "/contacts", `{"name":"test"}`, nil, http.StatusUnauthorized, ""},
		{"create input error", "POST", "/contacts", `"name":"test"}`, header, http.StatusBadRequest, ""},
		{"update ok", "PUT", "/contacts/123", `{"name":"contactxyz"}`, header, http.StatusOK, "*contactxyz*"},
		{"update verify", "GET", "/contacts/123", "", nil, http.StatusOK, `*contactxyz*`},
		{"update auth error", "PUT", "/contacts/123", `{"name":"contactxyz"}`, nil, http.StatusUnauthorized, ""},
		{"update input error", "PUT", "/contacts/123", `"name":"contactxyz"}`, header, http.StatusBadRequest, ""},
		{"delete ok", "DELETE", "/contacts/123", ``, header, http.StatusOK, "*contactxyz*"},
		{"delete verify", "DELETE", "/contacts/123", ``, header, http.StatusNotFound, ""},
		{"delete auth error", "DELETE", "/contacts/123", ``, nil, http.StatusUnauthorized, ""},
	}
	for _, tc := range tests {
		test.Endpoint(t, router, tc)
	}
}
