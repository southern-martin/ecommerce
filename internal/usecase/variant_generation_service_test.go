package usecase

import (
	"testing"

	"github.com/google/uuid"
)

func TestGenerateMatrix(t *testing.T) {
	svc := NewVariantGenerationService(nil, nil)
	productID := uuid.New()

	axes := []VariantAxisInput{
		{
			AttributeID: uuid.New(),
			OptionIDs:   []uuid.UUID{uuid.New(), uuid.New()},
		},
		{
			AttributeID: uuid.New(),
			OptionIDs:   []uuid.UUID{uuid.New(), uuid.New(), uuid.New()},
		},
	}

	variants, err := svc.GenerateMatrix(GenerateVariantsInput{
		ProductID:       productID,
		Axes:            axes,
		BasePriceMinor:  9999,
		InitialStockQty: 10,
	})
	if err != nil {
		t.Fatalf("GenerateMatrix() error = %v", err)
	}
	if len(variants) != 6 {
		t.Fatalf("GenerateMatrix() variants = %d, want 6", len(variants))
	}

	seen := map[string]struct{}{}
	for _, v := range variants {
		if v.CombinationKey == "" {
			t.Fatal("variant has empty combination key")
		}
		if _, ok := seen[v.CombinationKey]; ok {
			t.Fatalf("duplicate combination key detected: %s", v.CombinationKey)
		}
		seen[v.CombinationKey] = struct{}{}
	}
}

func TestGenerateMatrixRejectsDuplicateAxis(t *testing.T) {
	svc := NewVariantGenerationService(nil, nil)
	axisID := uuid.New()

	_, err := svc.GenerateMatrix(GenerateVariantsInput{
		ProductID: uuid.New(),
		Axes: []VariantAxisInput{
			{AttributeID: axisID, OptionIDs: []uuid.UUID{uuid.New()}},
			{AttributeID: axisID, OptionIDs: []uuid.UUID{uuid.New()}},
		},
	})
	if err == nil {
		t.Fatal("GenerateMatrix() error = nil, want duplicate axis error")
	}
}

func TestGenerateMatrixRejectsLargeMatrix(t *testing.T) {
	svc := NewVariantGenerationService(nil, nil)
	svc.SetMaxVariants(4)

	_, err := svc.GenerateMatrix(GenerateVariantsInput{
		ProductID: uuid.New(),
		Axes: []VariantAxisInput{
			{AttributeID: uuid.New(), OptionIDs: []uuid.UUID{uuid.New(), uuid.New(), uuid.New()}},
			{AttributeID: uuid.New(), OptionIDs: []uuid.UUID{uuid.New(), uuid.New()}},
		},
	})
	if err == nil {
		t.Fatal("GenerateMatrix() error = nil, want matrix too large error")
	}
}
