package static

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func PerformRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestEmptyDirectory(t *testing.T) {
	// SETUP file
	testRoot, _ := os.Getwd()
	f, err := ioutil.TempFile(testRoot, "")
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f.Name())
	f.WriteString("Gin Web Framework")
	f.Close()

	dir, filename := filepath.Split(f.Name())

	router := gin.New()
	router.Use(ServeRoot("/", dir))
	router.GET("/", func(c *gin.Context) {
		c.String(200, "index")
	})
	router.GET("/a", func(c *gin.Context) {
		c.String(200, "a")
	})
	router.GET("/"+filename, func(c *gin.Context) {
		c.String(200, "this is not printed")
	})
	w := PerformRequest(router, "GET", "/")
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "index")

	w = PerformRequest(router, "GET", "/"+filename)
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "Gin Web Framework")

	w = PerformRequest(router, "GET", "/"+filename+"a")
	assert.Equal(t, w.Code, 404)

	w = PerformRequest(router, "GET", "/a")
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "a")

	router2 := gin.New()
	router2.Use(ServeRoot("/static", dir))
	router2.GET("/"+filename, func(c *gin.Context) {
		c.String(200, "this is printed")
	})

	w = PerformRequest(router2, "GET", "/")
	assert.Equal(t, w.Code, 404)

	w = PerformRequest(router2, "GET", "/static")
	assert.Equal(t, w.Code, 404)
	router2.GET("/static", func(c *gin.Context) {
		c.String(200, "index")
	})

	w = PerformRequest(router2, "GET", "/static")
	assert.Equal(t, w.Code, 200)

	w = PerformRequest(router2, "GET", "/"+filename)
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "this is printed")

	w = PerformRequest(router2, "GET", "/static/"+filename)
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "Gin Web Framework")
}

func TestIndex(t *testing.T) {
	// SETUP file
	testRoot, _ := os.Getwd()
	f, err := os.Create(path.Join(testRoot, "index.html"))
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(f.Name())
	f.WriteString("index")
	f.Close()

	dir, filename := filepath.Split(f.Name())

	router := gin.New()
	router.Use(ServeRoot("/", dir))

	w := PerformRequest(router, "GET", "/"+filename)
	assert.Equal(t, w.Code, 301)

	w = PerformRequest(router, "GET", "/")
	assert.Equal(t, w.Code, 200)
	assert.Equal(t, w.Body.String(), "index")
}
