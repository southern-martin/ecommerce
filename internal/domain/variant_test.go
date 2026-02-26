package domain

import (
	"testing"

	"github.com/google/uuid"
)

func TestBuildCombinationKeySortsAttributes(t *testing.T) {
	attrA := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	attrB := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	optA := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	optB := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")

	key, err := BuildCombinationKey([]VariantOptionValue{
		{AttributeID: attrB, OptionID: optB},
		{AttributeID: attrA, OptionID: optA},
	})
	if err != nil {
		t.Fatalf("BuildCombinationKey() error = %v", err)
	}

	want := "11111111-1111-1111-1111-111111111111:aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa|22222222-2222-2222-2222-222222222222:bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"
	if key != want {
		t.Fatalf("BuildCombinationKey() = %s, want %s", key, want)
	}
}

func TestBuildCombinationKeyRejectsDuplicateAttribute(t *testing.T) {
	attr := uuid.New()
	_, err := BuildCombinationKey([]VariantOptionValue{
		{AttributeID: attr, OptionID: uuid.New()},
		{AttributeID: attr, OptionID: uuid.New()},
	})
	if err == nil {
		t.Fatal("BuildCombinationKey() error = nil, want duplicate attribute error")
	}
}
