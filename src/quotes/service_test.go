package quotes

import "testing"

func TestQuoteService(t *testing.T) {
	if _, err := GetRandomQuote(); err == nil {
		t.Fatal("GetRandomQuote succeeded before loading data")
	}
	if err := LoadQuotes([]byte("{")); err == nil {
		t.Fatal("LoadQuotes accepted invalid JSON")
	}
	if err := LoadQuotes([]byte("[]")); err == nil {
		t.Fatal("LoadQuotes accepted an empty dataset")
	}
	if err := LoadQuotes([]byte(`[{"id":1,"quote":"Keep going","author":"A","category":"general"},{"id":2,"quote":"Build","author":"B","category":"work"}]`)); err != nil {
		t.Fatal(err)
	}
	if got := GetTotalCount(); got != 2 {
		t.Fatalf("GetTotalCount() = %d, want 2", got)
	}
	if got, err := GetRandomQuote(); err != nil || got == nil {
		t.Fatalf("GetRandomQuote() = %v, %v", got, err)
	}
	if got := GetAllQuotes(); len(got) != 2 {
		t.Fatalf("GetAllQuotes() returned %d quotes", len(got))
	}
	if got, err := GetQuoteByID(2); err != nil || got.Author != "B" {
		t.Fatalf("GetQuoteByID(2) = %v, %v", got, err)
	}
	if _, err := GetQuoteByID(99); err == nil {
		t.Fatal("GetQuoteByID accepted an unknown ID")
	}
	if got := GetQuotesByCategory("general"); len(got) != 1 || got[0].ID != 1 {
		t.Fatalf("GetQuotesByCategory() = %v", got)
	}
	if got := GetQuotesByAuthor("B"); len(got) != 1 || got[0].ID != 2 {
		t.Fatalf("GetQuotesByAuthor() = %v", got)
	}
	if got, err := GetAllJSON(); err != nil || len(got) == 0 {
		t.Fatalf("GetAllJSON() = %s, %v", got, err)
	}
}
