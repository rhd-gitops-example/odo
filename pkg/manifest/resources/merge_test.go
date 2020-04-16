package resources

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMerge(t *testing.T) {
	mergeTests := []struct {
		src  resources
		dest resources
		want resources
	}{
		{
			src:  resources{"test1": "val1"},
			dest: resources{},
			want: resources{"test1": "val1"},
		},
		{
			src:  resources{"test1": "val1"},
			dest: resources{"test2": "val2"},
			want: resources{"test1": "val1", "test2": "val2"},
		},
		{
			src:  resources{"test1": "val1"},
			dest: resources{"test1": "val2"},
			want: resources{"test1": "val1"},
		},
	}

	for _, tt := range mergeTests {
		result := Merge(tt.src, tt.dest)

		if diff := cmp.Diff(tt.want, result); diff != "" {
			t.Fatalf("failed merge: %s\n", diff)
		}
	}

}
