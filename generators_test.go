package mkpage

import (
	"testing"
)

func TestScanArgs(t *testing.T) {
	expectedGenerator := `one`
	expectedParams := []string{`two`, `three`}
	src := `one two three`
	generator, params := scanArgs(src)

	if generator != expectedGenerator {
		t.Errorf("expected %q, got %q from %q", expectedGenerator, generator, src)
	}
	if len(params) != len(expectedParams) {
		t.Errorf("expected %d, got %d from %+v", len(expectedParams), len(params), params)
		t.FailNow()
	}
	for i, val := range expectedParams {
		if val != params[i] {
			t.Errorf("expected param(%d) %q, got %q from %+v", i, val, params[i], params)
		}
	}

	expectedGenerator = `python3`
	expectedParams = []string{`eprints_html_view.py`, `3`, `et el.`}
	src = `python3 eprints_html_view.py 3 "et el."`
	generator, params = scanArgs(src)
	if generator != expectedGenerator {
		t.Errorf("expected %q, got %q from %q", expectedGenerator, generator, src)
	}
	if len(params) != len(expectedParams) {
		t.Errorf("expected %d, got %d from %+v", len(expectedParams), len(params), params)
		t.FailNow()
	}
	for i, val := range expectedParams {
		if val != params[i] {
			t.Errorf("expected param(%d) %q, got %q from %+v", i, val, params[i], params)
		}
	}
}
