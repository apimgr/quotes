package anime

import "testing"

func TestAnimeQuoteService(t *testing.T) {
	if _, err := GetRandomQuote(); err == nil {
		t.Fatal("random quote available before loading")
	}
	if err := LoadQuotes([]byte("{")); err == nil {
		t.Fatal("accepted invalid JSON")
	}
	if err := LoadQuotes([]byte("[]")); err == nil {
		t.Fatal("accepted empty data")
	}
	if err := LoadQuotes([]byte(`[{"id":1,"quote":"Go","character":"A","anime":"Show","category":"general"},{"id":2,"quote":"Run","character":"B","anime":"Other","category":"action"}]`)); err != nil {
		t.Fatal(err)
	}
	if GetTotalCount() != 2 || len(GetAllQuotes()) != 2 {
		t.Fatal("loaded quote count mismatch")
	}
	if q, err := GetRandomQuote(); err != nil || q == nil {
		t.Fatalf("random quote: %v, %v", q, err)
	}
	if q, err := GetQuoteByID(1); err != nil || q.Character != "A" {
		t.Fatalf("lookup: %v, %v", q, err)
	}
	if _, err := GetQuoteByID(99); err == nil {
		t.Fatal("accepted unknown ID")
	}
	if got := GetQuotesByCategory("general"); len(got) != 1 {
		t.Fatalf("category results: %v", got)
	}
	if got := GetQuotesByAnime("Other"); len(got) != 1 {
		t.Fatalf("anime results: %v", got)
	}
	if got := GetQuotesByCharacter("A"); len(got) != 1 {
		t.Fatalf("character results: %v", got)
	}
}
