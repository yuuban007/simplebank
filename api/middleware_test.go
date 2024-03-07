package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/yuuban007/simplebank/token"
)

func TestAuthMiddleware(t *testing.T) {
	testCases := []struct {
		name          string
		setUpAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{}
}

/*
	// Create a new Gin router
	router := gin.New()

	// Create a mock token maker
	mockTokenMaker := &token.MockTokenMaker{}

	// Set up the auth middleware
	router.Use(authMiddleware(mockTokenMaker))

	// Define a test route that requires authentication
	router.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Authenticated")
	})

	// Create a new HTTP request to the test route
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	// Set the authorization header with a valid token
	req.Header.Set("Authorization", "Bearer valid_token")

	// Create a new HTTP response recorder
	res := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(res, req)

	// Assert that the response status code is 200 OK
	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "Authenticated", res.Body.String())
*/
