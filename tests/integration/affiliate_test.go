package integration

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAffiliateOperations(t *testing.T) {
	requireServices(t, map[string]string{
		"auth":      AuthBaseURL,
		"affiliate": AffiliateBaseURL,
	})

	user := createTestUser(t)
	var affiliateID string

	t.Run("Create_affiliate", func(t *testing.T) {
		affiliateBody := map[string]interface{}{
			"name":            "Test Affiliate",
			"email":           user.Email,
			"commission_rate": 10.0,
		}

		resp := httpPostWithUserID(t, AffiliateBaseURL+"/api/v1/affiliates", affiliateBody, user.AccessToken, user.UserID)
		defer resp.Body.Close()
		require.Contains(t, []int{http.StatusOK, http.StatusCreated}, resp.StatusCode, "failed to create affiliate")

		data := readBodyMap(t, resp)
		id := assertJSONField(t, data, "id")
		idStr, ok := id.(string)
		require.True(t, ok, "affiliate id should be a string")
		assert.NotEmpty(t, idStr)

		affiliateID = idStr
		t.Logf("Created affiliate: %s", affiliateID)
	})

	t.Run("Get_affiliate", func(t *testing.T) {
		require.NotEmpty(t, affiliateID, "affiliate ID must be set from Create_affiliate")

		resp := httpGetWithUserID(t, fmt.Sprintf("%s/api/v1/affiliates/%s", AffiliateBaseURL, affiliateID), user.AccessToken, user.UserID)
		defer resp.Body.Close()
		assertStatus(t, resp, http.StatusOK)

		data := readBodyMap(t, resp)
		assertJSONField(t, data, "id")

		t.Logf("Retrieved affiliate: %s", affiliateID)
	})

	t.Run("List_affiliates", func(t *testing.T) {
		resp := httpGetWithUserID(t, AffiliateBaseURL+"/api/v1/affiliates", user.AccessToken, user.UserID)
		defer resp.Body.Close()
		assertStatus(t, resp, http.StatusOK)

		t.Log("Listed affiliates successfully")
	})

	t.Run("Create_referral_link", func(t *testing.T) {
		require.NotEmpty(t, affiliateID, "affiliate ID must be set from Create_affiliate")

		linkBody := map[string]interface{}{
			"product_id": "00000000-0000-0000-0000-000000000001",
			"campaign":   "summer_sale",
		}

		resp := httpPostWithUserID(t,
			fmt.Sprintf("%s/api/v1/affiliates/%s/links", AffiliateBaseURL, affiliateID),
			linkBody, user.AccessToken, user.UserID,
		)
		defer resp.Body.Close()
		require.Contains(t, []int{http.StatusOK, http.StatusCreated}, resp.StatusCode, "failed to create referral link")

		t.Log("Created referral link successfully")
	})
}
