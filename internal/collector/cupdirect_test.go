package collector

import (
	"reflect"
	"testing"

	"github.com/i25959341/sku-aggregator/internal/types"
)

func TestCollectCupdirect(t *testing.T) {
	tests := []struct {
		name string
		want []types.SKU
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CollectCupdirect(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CollectCupdirect() = %v, want %v", got, tt.want)
			}
		})
	}
}
