package integration

import (
	"net/http"
	"testing"
)

func TestMediaOperations(t *testing.T) {
	requireServices(t, map[string]string{
		"auth":  AuthBaseURL,
		"media": MediaBaseURL,
	})

	user := createTestUser(t)

	t.Run("Health_check", func(t *testing.T) {
		resp := httpGet(t, MediaBaseURL+"/health", "")
		defer resp.Body.Close()
		assertStatus(t, resp, http.StatusOK)

		t.Log("Media service health check passed")
	})

	t.Run("List_media", func(t *testing.T) {
		resp := httpGetWithUserID(t, MediaBaseURL+"/api/v1/media", user.AccessToken, user.UserID)
		defer resp.Body.Close()
		assertStatus(t, resp, http.StatusOK)

		t.Log("Listed media successfully")
	})
}
