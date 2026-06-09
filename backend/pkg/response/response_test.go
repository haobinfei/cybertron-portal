package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestSuccess(t *testing.T) {
	r := setupRouter()
	r.GET("/test", func(c *gin.Context) {
		Success(c, map[string]string{"key": "value"})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "success", resp.Message)
	assert.Equal(t, map[string]interface{}{"key": "value"}, resp.Data)
}

func TestSuccess_NilData(t *testing.T) {
	r := setupRouter()
	r.GET("/test", func(c *gin.Context) {
		Success(c, nil)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "success", resp.Message)
}

func TestSuccessWithMessage(t *testing.T) {
	r := setupRouter()
	r.GET("/test", func(c *gin.Context) {
		SuccessWithMessage(c, "操作成功", map[string]int{"count": 42})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	var resp Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "操作成功", resp.Message)
}

func TestError(t *testing.T) {
	r := setupRouter()
	r.GET("/test", func(c *gin.Context) {
		Error(c, 400, "参数错误")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	var resp Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.Code)
	assert.Equal(t, "参数错误", resp.Message)
}

func TestPage(t *testing.T) {
	r := setupRouter()
	r.GET("/test", func(c *gin.Context) {
		Page(c, []string{"a", "b", "c"}, int64(100), 1, 10)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	var resp Response
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)

	pageData, ok := resp.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, float64(100), pageData["total"])
	assert.Equal(t, float64(1), pageData["page"])
	assert.Equal(t, float64(10), pageData["page_size"])
}
