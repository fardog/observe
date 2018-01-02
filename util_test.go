package observe

import "testing"

func TestAnonymizeIP(t *testing.T) {
	cases := [][]string{
		[]string{"175.72.100.10:3000", "175.72.96.0"},
		[]string{"175.72.100.10", "175.72.96.0"},
		[]string{"[2001:0db8:0123:4567:89ab:cdef:1234:5678]:9999", "2001:db8::"},
		[]string{"2001:0db8:0123:4567:89ab:cdef:1234:5678", "2001:db8::"},
	}

	for _, c := range cases {
		if a, err := AnonymizeIP(c[0]); err != nil {
			t.Errorf("got unexptected error: %v", err)
		} else if a != c[1] {
			t.Errorf("expected value %v, got %v", c[1], a)
		}
	}
}
