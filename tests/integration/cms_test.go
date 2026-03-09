package integration

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCMSOperations(t *testing.T) {
	requireServices(t, map[string]string{
		"auth": AuthBaseURL,
		"cms":  CmsBaseURL,
	})

	user := createTestUser(t)
	var pageID string

	t.Run("Create_page", func(t *testing.T) {
		pageBody := map[string]interface{}{
			"title":   "Integration Test Page",
			"slug":    fmt.Sprintf("test-page-%d", randomSuffix()),
			"content": "<h1>Hello</h1><p>This is a test page created by integration tests.</p>",
			"status":  "draft",
		}

		resp := httpPostWithUserID(t, CmsBaseURL+"/api/v1/cms/pages", pageBody, user.AccessToken, user.UserID)
		defer resp.Body.Close()
		require.Contains(t, []int{http.StatusOK, http.StatusCreated}, resp.StatusCode, "failed to create CMS page")

		data := readBodyMap(t, resp)
		id := assertJSONField(t, data, "id")
		idStr, ok := id.(string)
		require.True(t, ok, "page id should be a string")
		assert.NotEmpty(t, idStr)

		pageID = idStr
		t.Logf("Created CMS page: %s", pageID)
	})

	t.Run("Get_page", func(t *testing.T) {
		require.NotEmpty(t, pageID, "page ID must be set from Create_page")

		resp := httpGetWithUserID(t, fmt.Sprintf("%s/api/v1/cms/pages/%s", CmsBaseURL, pageID), user.AccessToken, user.UserID)
		defer resp.Body.Close()
		assertStatus(t, resp, http.StatusOK)

		data := readBodyMap(t, resp)
		assertJSONField(t, data, "id")
		assertJSONField(t, data, "title")

		t.Logf("Retrieved CMS page: %s", pageID)
	})

	t.Run("List_pages", func(t *testing.T) {
		resp := httpGetWithUserID(t, CmsBaseURL+"/api/v1/cms/pages", user.AccessToken, user.UserID)
		defer resp.Body.Close()
		assertStatus(t, resp, http.StatusOK)

		t.Log("Listed CMS pages successfully")
	})

	t.Run("Update_page", func(t *testing.T) {
		require.NotEmpty(t, pageID, "page ID must be set from Create_page")

		updateBody := map[string]interface{}{
			"title":   "Updated Integration Test Page",
			"content": "<h1>Updated</h1><p>This page has been updated by integration tests.</p>",
		}

		resp := httpPatch(t, fmt.Sprintf("%s/api/v1/cms/pages/%s", CmsBaseURL, pageID), updateBody, user.AccessToken, user.UserID)
		defer resp.Body.Close()
		require.Contains(t, []int{http.StatusOK, http.StatusNoContent}, resp.StatusCode, "failed to update CMS page")

		t.Logf("Updated CMS page: %s", pageID)
	})

	t.Run("Delete_page", func(t *testing.T) {
		require.NotEmpty(t, pageID, "page ID must be set from Create_page")

		resp := httpDelete(t, fmt.Sprintf("%s/api/v1/cms/pages/%s", CmsBaseURL, pageID), user.AccessToken, user.UserID)
		defer resp.Body.Close()
		require.Contains(t, []int{http.StatusOK, http.StatusNoContent}, resp.StatusCode, "failed to delete CMS page")

		t.Logf("Deleted CMS page: %s", pageID)
	})
}

// randomSuffix returns a simple numeric suffix for unique slugs.
func randomSuffix() int64 {
	return time.Now().UnixNano()
}
