package init

import (
	"strings"
	"testing"
)

func TestReplaceManagedBlock(t *testing.T) {
	block := managedBegin + "\nnew\n" + managedEnd

	tests := []struct {
		name      string
		existing  string
		want      string
		wantFound bool
	}{
		{
			name:      "both markers replaced in place",
			existing:  "before\n" + managedBegin + "\nold\n" + managedEnd + "\nafter\n",
			want:      "before\n" + block + "\nafter\n",
			wantFound: true,
		},
		{
			name:      "no markers appends and reports not found",
			existing:  "# My Guide\n",
			want:      "# My Guide\n\n" + block + "\n",
			wantFound: false,
		},
		{
			name:      "begin marker only is not a block",
			existing:  "# My Guide\n" + managedBegin + "\n",
			want:      "# My Guide\n" + managedBegin + "\n\n" + block + "\n",
			wantFound: false,
		},
		{
			name:      "end marker only is not a block",
			existing:  "# My Guide\n" + managedEnd + "\n",
			want:      "# My Guide\n" + managedEnd + "\n\n" + block + "\n",
			wantFound: false,
		},
		{
			name:      "end before begin is not a block",
			existing:  managedEnd + "\nx\n" + managedBegin + "\n",
			want:      managedEnd + "\nx\n" + managedBegin + "\n\n" + block + "\n",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, found := replaceManagedBlock(tt.existing, block)
			if found != tt.wantFound {
				t.Errorf("found = %v, want %v", found, tt.wantFound)
			}
			if got != tt.want {
				t.Errorf("merged = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestReplaceManagedBlock_idempotent(t *testing.T) {
	block := managedBegin + "\nnew\n" + managedEnd

	once, _ := replaceManagedBlock("# My Guide\n", block)
	twice, found := replaceManagedBlock(once, block)
	if !found {
		t.Error("second pass should find the block it just appended")
	}
	if twice != once {
		t.Errorf("second pass changed the file: %q, want %q", twice, once)
	}
	if n := strings.Count(twice, managedBegin); n != 1 {
		t.Errorf("got %d begin markers, want 1: %q", n, twice)
	}
}
