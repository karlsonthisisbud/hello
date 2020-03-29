package collector

import (
	"reflect"
	"testing"

	"github.com/i25959341/sku-aggregator/internal/types"
)

func TestCollectVegaware(t *testing.T) {
	tests := []struct {
		name string
		want []types.SKU
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CollectVegaware(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CollectVegaware() = %v, want %v", got, tt.want)
			}
		})
	}
}
