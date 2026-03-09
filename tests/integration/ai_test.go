package integration

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAIOperations(t *testing.T) {
	requireServices(t, map[string]string{
		"auth": AuthBaseURL,
		"ai":   AiBaseURL,
	})

	user := createTestUser(t)

	t.Run("Health_check", func(t *testing.T) {
		resp := httpGet(t, AiBaseURL+"/health", "")
		defer resp.Body.Close()
		assertStatus(t, resp, http.StatusOK)

		t.Log("AI service health check passed")
	})

	t.Run("Get_recommendations", func(t *testing.T) {
		resp := httpGetWithUserID(t, AiBaseURL+"/api/v1/ai/recommendations?user_id="+user.UserID, user.AccessToken, user.UserID)
		defer resp.Body.Close()
		require.Contains(t, []int{http.StatusOK, http.StatusNotFound, http.StatusServiceUnavailable}, resp.StatusCode,
			"recommendations endpoint should return OK, not found, or service unavailable")

		t.Logf("AI recommendations returned status: %d", resp.StatusCode)
	})
}
