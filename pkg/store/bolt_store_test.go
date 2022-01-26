package store

import (
	"io/ioutil"
	"testing"
)

func TestPutGet(t *testing.T) {
	dir, err := ioutil.TempDir("", "storetest_")
	if err != nil {
		t.Fatal(err)
	}
	store := NewBoltStore()

	store.Init(dir)

	k := []byte("mykey")
	if err := store.Put(k, []byte("neo")); err != nil {
		t.Fatal(err)
	}

	if val, err := store.Get(k); err != nil {
		t.Fatal(err)
	} else if string(val) != "neo" {
		t.Errorf("Expected 'neo'. Found: %s", string(val))
	}

	if err := store.Put(k, []byte("the one")); err != nil {
		t.Fatal(err)
	}

	if val, err := store.Get(k); err != nil {
		t.Fatal(err)
	} else if string(val) != "the one" {
		t.Errorf("Expected 'the one'. Found: %s", string(val))
	}
}
