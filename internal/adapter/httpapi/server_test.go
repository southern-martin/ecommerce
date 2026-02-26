package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"ecommerce/catalog-service/internal/infra/memory"
	"ecommerce/catalog-service/internal/usecase"
)

func TestServerEndToEndVariantFlow(t *testing.T) {
	store := memory.NewStore()
	categoryRepo := memory.NewCategoryRepo(store)
	attributeRepo := memory.NewAttributeRepo(store)
	productRepo := memory.NewProductRepo(store)

	server := NewServer(
		usecase.NewCategoryService(categoryRepo),
		usecase.NewAttributeService(attributeRepo, categoryRepo),
		usecase.NewProductService(productRepo, categoryRepo, attributeRepo),
		usecase.NewVariantGenerationService(productRepo, attributeRepo),
		productRepo,
	)
	handler := server.Routes()

	categoryResp := doJSONRequest(t, handler, http.MethodPost, "/categories", map[string]any{
		"name": "Shoes",
		"slug": "shoes",
	})
	categoryID, _ := categoryResp["id"].(string)
	if categoryID == "" {
		t.Fatal("category id is empty")
	}

	treeResp := doJSONRequest(t, handler, http.MethodGet, "/categories/tree", nil)
	treeNodes, ok := treeResp["categories"].([]any)
	if !ok || len(treeNodes) == 0 {
		t.Fatal("category tree should contain at least one node")
	}

	childrenResp := doJSONRequest(t, handler, http.MethodGet, "/categories/"+categoryID+"/children", nil)
	children, ok := childrenResp["children"].([]any)
	if !ok {
		t.Fatal("children response should contain children array")
	}
	if len(children) != 0 {
		t.Fatalf("children count = %d, want 0", len(children))
	}

	attributeResp := doJSONRequest(t, handler, http.MethodPost, "/categories/"+categoryID+"/attributes", map[string]any{
		"name":            "Size",
		"code":            "size",
		"type":            "select",
		"is_variant_axis": true,
	})
	attributeID, _ := attributeResp["id"].(string)
	if attributeID == "" {
		t.Fatal("attribute id is empty")
	}

	optionResp := doJSONRequest(t, handler, http.MethodPost, "/attributes/"+attributeID+"/options", map[string]any{
		"value": "42",
		"label": "EU 42",
	})
	optionID, _ := optionResp["id"].(string)
	if optionID == "" {
		t.Fatal("option id is empty")
	}

	productResp := doJSONRequest(t, handler, http.MethodPost, "/products", map[string]any{
		"name":                "Running Shoe",
		"slug":                "running-shoe",
		"primary_category_id": categoryID,
		"base_price_minor":    9900,
	})
	productID, _ := productResp["id"].(string)
	if productID == "" {
		t.Fatal("product id is empty")
	}

	productGetResp := doJSONRequest(t, handler, http.MethodGet, "/products/"+productID, nil)
	if gotID, _ := productGetResp["id"].(string); gotID != productID {
		t.Fatalf("product get id = %s, want %s", gotID, productID)
	}

	categoryProductsResp := doJSONRequest(t, handler, http.MethodGet, "/categories/"+categoryID+"/products", nil)
	productCount, _ := categoryProductsResp["count"].(float64)
	if int(productCount) != 1 {
		t.Fatalf("category products count = %v, want 1", productCount)
	}

	doJSONRequestExpectStatus(t, handler, http.MethodPut, "/products/"+productID+"/attributes", map[string]any{
		"values": []map[string]any{
			{
				"attribute_id": attributeID,
				"option_id":    optionID,
			},
		},
	}, http.StatusBadRequest)

	generateResp := doJSONRequest(t, handler, http.MethodPost, "/products/"+productID+"/variants/generate", map[string]any{
		"axes": []map[string]any{
			{
				"attribute_id": attributeID,
				"option_ids":   []string{optionID},
			},
		},
		"base_price_minor":  9900,
		"initial_stock_qty": 5,
	})
	count, _ := generateResp["count"].(float64)
	if int(count) != 1 {
		t.Fatalf("generated variants count = %v, want 1", count)
	}

	listResp := doJSONRequest(t, handler, http.MethodGet, "/products/"+productID+"/variants", nil)
	listCount, _ := listResp["count"].(float64)
	if int(listCount) != 1 {
		t.Fatalf("list variants count = %v, want 1", listCount)
	}
}

func doJSONRequestExpectStatus(
	t *testing.T,
	handler http.Handler,
	method, path string,
	body any,
	wantStatus int,
) {
	t.Helper()
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code != wantStatus {
		t.Fatalf("%s %s status=%d, want %d body=%s", method, path, rec.Code, wantStatus, rec.Body.String())
	}
}

func doJSONRequest(t *testing.T, handler http.Handler, method, path string, body any) map[string]any {
	t.Helper()
	var payload []byte
	var err error
	if body != nil {
		payload, err = json.Marshal(body)
		if err != nil {
			t.Fatalf("json.Marshal() error = %v", err)
		}
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(payload))
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	if rec.Code >= 300 {
		t.Fatalf("%s %s failed: status=%d body=%s", method, path, rec.Code, rec.Body.String())
	}
	if rec.Code == http.StatusNoContent {
		return map[string]any{}
	}
	var out map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("json.Unmarshal() error = %v body=%s", err, rec.Body.String())
	}
	return out
}
