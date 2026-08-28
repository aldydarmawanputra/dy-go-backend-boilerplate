package pagination

import "testing"

func TestParse(t *testing.T) {
	cases := []struct {
		limit, offset string
		wantL, wantO  int
	}{
		{"", "", 20, 0},
		{"50", "10", 50, 10},
		{"0", "-5", 20, 0},
		{"1000", "0", 20, 0},
		{"abc", "xyz", 20, 0},
	}
	for _, tc := range cases {
		p := Parse(tc.limit, tc.offset)
		if p.Limit != tc.wantL || p.Offset != tc.wantO {
			t.Fatalf("Parse(%q,%q) = {%d,%d}, want {%d,%d}", tc.limit, tc.offset, p.Limit, p.Offset, tc.wantL, tc.wantO)
		}
	}
}

func TestNewMeta(t *testing.T) {
	m := NewMeta(Params{Limit: 20, Offset: 0}, 5, 42)
	if m.Total != 42 || m.Count != 5 || m.Limit != 20 {
		t.Fatalf("meta = %+v", m)
	}
}
