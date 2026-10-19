package ai

import (
	"reflect"
	"testing"
)

func TestFilterTickersDropsInventedCodes(t *testing.T) {
	valid := []string{"AWAB", "WGBX", "TELE"}
	got := filterTickers([]string{"CBE", "awab", "AWB", " WGBX ", "AWAB"}, valid)
	want := []string{"AWAB", "WGBX"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v want %v", got, want)
	}
	if got := filterTickers([]string{"CBE"}, nil); len(got) != 0 {
		t.Fatalf("with no known tickers nothing should pass, got %v", got)
	}
}
