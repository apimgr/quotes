package catalog

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestLoadEmbeddedCatalog(t *testing.T) {
	store, err := LoadEmbedded()
	if err != nil {
		t.Fatal(err)
	}
	counts := store.Categories()
	if len(counts) != len(categoryNames) {
		t.Fatalf("category count = %d, want %d", len(counts), len(categoryNames))
	}
	total := 0
	for _, name := range Names() {
		if counts[name] == 0 {
			t.Errorf("category %q has no records", name)
		}
		total += counts[name]
	}
	all, count, err := store.List("", "", 1, total)
	if err != nil || count != total || len(all) != total {
		t.Fatalf("list all: count=%d, len=%d, err=%v", count, len(all), err)
	}
	for index, quote := range all {
		if quote.ID != index+1 {
			t.Fatalf("ID at index %d = %d, want %d", index, quote.ID, index+1)
		}
		if quote.Quote == "" || quote.Category == "" {
			t.Fatalf("incomplete quote: %#v", quote)
		}
	}
}

func TestNewRejectsInvalidCollections(t *testing.T) {
	valid := fixtureCollections()
	for _, tc := range []struct {
		name string
		mutate func(map[string][]byte)
		want string
	}{
		{"missing", func(collections map[string][]byte) { delete(collections, CategoryAnime) }, "missing anime"},
		{"empty", func(collections map[string][]byte) { collections[CategoryAnime] = []byte("[]") }, "empty"},
		{"malformed", func(collections map[string][]byte) { collections[CategoryAnime] = []byte("{") }, "parse anime"},
		{"invalid ID", func(collections map[string][]byte) { collections[CategoryAnime] = []byte(`[{"id":0,"quote":"x"}]`) }, "invalid record"},
		{"missing text", func(collections map[string][]byte) { collections[CategoryAnime] = []byte(`[{"id":1}]`) }, "invalid record"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			collections := cloneCollections(valid)
			tc.mutate(collections)
			_, err := New(collections)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("New error = %v, want containing %q", err, tc.want)
			}
		})
	}
}

func TestListFilterPaginationAndValidation(t *testing.T) {
	store, err := New(fixtureCollections())
	if err != nil {
		t.Fatal(err)
	}
	quotes, total, err := store.List(CategoryAnime, "CHARACTER 2", 1, 1)
	if err != nil || total != 1 || len(quotes) != 1 || quotes[0].Author != "Character 2" {
		t.Fatalf("filtered list = %#v, total=%d, err=%v", quotes, total, err)
	}
	quotes, total, err = store.List("", "", 2, 1)
	if err != nil || total != 5 || len(quotes) != 1 {
		t.Fatalf("paginated list = %#v, total=%d, err=%v", quotes, total, err)
	}
	quotes, total, err = store.List("", "", int(^uint(0)>>1), 1000)
	if err != nil || total != 5 || len(quotes) != 0 {
		t.Fatalf("out-of-range page = %#v, total=%d, err=%v", quotes, total, err)
	}
	if _, _, err := store.List("invalid", "", 1, 5); err == nil {
		t.Fatal("unknown category was accepted")
	}
	if _, _, err := store.List("", "", 0, 1); err == nil {
		t.Fatal("invalid page was accepted")
	}
	if _, _, err := store.List("", "", 1, 0); err == nil {
		t.Fatal("invalid page size was accepted")
	}
}

func TestGetAndRandom(t *testing.T) {
	store, err := New(fixtureCollections())
	if err != nil {
		t.Fatal(err)
	}
	quote, err := store.Get("", 2)
	if err != nil || quote.Category != CategoryAnime {
		t.Fatalf("Get = %#v, %v", quote, err)
	}
	if _, err := store.Get(CategoryQuotes, 2); err == nil {
		t.Fatal("category mismatch was accepted")
	}
	if _, err := store.Get("", 0); err == nil {
		t.Fatal("unknown ID was accepted")
	}
	for _, category := range []string{"", CategoryProgramming} {
		quote, err := store.Random(category)
		if err != nil {
			t.Fatal(err)
		}
		if category != "" && quote.Category != category {
			t.Fatalf("random category = %q, want %q", quote.Category, category)
		}
	}
	if _, err := store.Random("unknown"); err == nil {
		t.Fatal("unknown category was accepted")
	}
}

func fixtureCollections() map[string][]byte {
	collections := make(map[string][]byte, len(categoryNames))
	for index, category := range categoryNames {
		record := map[string]interface{}{"id": 1, "quote": "Example quote", "category": "group"}
		if category == CategoryAnime {
			record["character"] = "Character 2"
		}
		if category == CategoryChuckNorris || category == CategoryDadJokes || category == CategoryProgramming {
			delete(record, "quote")
			record["joke"] = "Example joke"
		}
		if category == CategoryQuotes {
			record["author"] = "Author"
		}
		if category == CategoryAnime {
			record["id"] = 2
		}
		if index == 4 {
			record["id"] = 3
		}
		content, _ := json.Marshal([]interface{}{record})
		collections[category] = content
	}
	return collections
}

func cloneCollections(input map[string][]byte) map[string][]byte {
	clone := make(map[string][]byte, len(input))
	for key, value := range input {
		clone[key] = append([]byte(nil), value...)
	}
	return clone
}
