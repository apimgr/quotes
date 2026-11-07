package programming

import "testing"

func TestJokeService(t *testing.T) {
	if _, err := GetRandomJoke(); err == nil {
		t.Fatal("random joke available before loading")
	}
	if err := LoadJokes([]byte("{")); err == nil {
		t.Fatal("accepted invalid JSON")
	}
	if err := LoadJokes([]byte("[]")); err == nil {
		t.Fatal("accepted empty data")
	}
	if err := LoadJokes([]byte(`[{"id":1,"joke":"Bug","category":"code"}]`)); err != nil {
		t.Fatal(err)
	}
	if GetTotalCount() != 1 || len(GetAllJokes()) != 1 {
		t.Fatal("loaded joke count mismatch")
	}
	if j, err := GetRandomJoke(); err != nil || j == nil {
		t.Fatalf("random joke: %v, %v", j, err)
	}
	if j, err := GetJokeByID(1); err != nil || j.Joke != "Bug" {
		t.Fatalf("lookup: %v, %v", j, err)
	}
	if _, err := GetJokeByID(99); err == nil {
		t.Fatal("accepted unknown ID")
	}
}
