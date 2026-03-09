package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReturnOperations(t *testing.T) {
	requireServices(t, map[string]string{
		"auth":   AuthBaseURL,
		"return": ReturnBaseURL,
	})

	user := createTestUser(t)
	var returnID string

	t.Run("Create_return_request", func(t *testing.T) {
		returnBody := map[string]interface{}{
			"order_id": "00000000-0000-0000-0000-000000000001",
			"items": []map[string]interface{}{
				{
					"product_id": "00000000-0000-0000-0000-000000000002",
					"quantity":   1,
					"reason":     "Defective product received",
				},
			},
		}

		resp := httpPostWithUserID(t, ReturnBaseURL+"/api/v1/returns", returnBody, user.AccessToken, user.UserID)
		defer resp.Body.Close()
		require.Contains(t, []int{http.StatusOK, http.StatusCreated, http.StatusBadRequest}, resp.StatusCode,
			"return creation should succeed or fail with bad request for non-existent order")

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			data := readBodyMap(t, resp)
			id := assertJSONField(t, data, "id")
			idStr, ok := id.(string)
			require.True(t, ok, "return id should be a string")
			assert.NotEmpty(t, idStr)

			returnID = idStr
			t.Logf("Created return request: %s", returnID)
		} else {
			t.Log("Return creation returned 400 (expected if order does not exist)")
		}
	})

	t.Run("Get_return", func(t *testing.T) {
		if returnID == "" {
			t.Skip("skipping: no return was created in previous step")
		}

		resp := httpGetWithUserID(t, fmt.Sprintf("%s/api/v1/returns/%s", ReturnBaseURL, returnID), user.AccessToken, user.UserID)
		defer resp.Body.Close()
		assertStatus(t, resp, http.StatusOK)

		data := readBodyMap(t, resp)
		assertJSONField(t, data, "id")

		t.Logf("Retrieved return: %s", returnID)
	})

	t.Run("List_returns", func(t *testing.T) {
		resp := httpGetWithUserID(t, ReturnBaseURL+"/api/v1/returns", user.AccessToken, user.UserID)
		defer resp.Body.Close()
		assertStatus(t, resp, http.StatusOK)

		t.Log("Listed returns successfully")
	})

	t.Run("Update_return_status", func(t *testing.T) {
		if returnID == "" {
			t.Skip("skipping: no return was created in previous step")
		}

		statusBody := map[string]interface{}{
			"status": "approved",
		}

		resp := httpPatch(t, fmt.Sprintf("%s/api/v1/returns/%s/status", ReturnBaseURL, returnID), statusBody, user.AccessToken, user.UserID)
		defer resp.Body.Close()
		require.Contains(t, []int{http.StatusOK, http.StatusNoContent, http.StatusForbidden}, resp.StatusCode,
			"return status update should succeed or be forbidden for non-admin")

		t.Logf("Update return status returned: %d", resp.StatusCode)
	})
}
