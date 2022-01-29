package store

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/go-test/deep"
)

// TestBasicInteractions tests if reading from, writing to and deleting from the store works properly.
func TestBasicInteractions(store Store, t *testing.T) {
	key := strconv.FormatInt(rand.Int63(), 10)

	// Initially the key shouldn't exist
	_, err := store.Get(key)
	if err != ErrKeyNotFound {
		t.Errorf("Should return ErrKeyNotFound, but: %v", err)
	}

	// Deleting a non-existing key-value pair should NOT lead to an error
	err = store.Delete(key)
	if err != nil {
		t.Error(err)
	}

	val := []byte("test value")

	// Store an KV pair
	err = store.Set(key, val)
	if err != nil {
		t.Error(err)
	}
	// Storing it again should not lead to an error but just overwrite it
	err = store.Set(key, val)
	if err != nil {
		t.Error(err)
	}

	// Retrieve the object
	expectedVal := val
	actualVal, err := store.Get(key)
	if err != nil {
		t.Error(err)
	}

	if diff := deep.Equal(actualVal, expectedVal); diff != nil {
		t.Error(diff)
	}

	// Delete
	err = store.Delete(key)
	if err != nil {
		t.Error(err)
	}

	// Key-value pair shouldn't exist anymore
	val, err = store.Get(key)
	if err == nil {
		t.Error("should return ErrKeyNotFound")
	}
	if val != nil {
		t.Error("A value was found, but no value was expected")
	}
}

// TestConcurrentInteractions launches a bunch of goroutines that concurrently work with the store.
func TestConcurrentInteractions(t *testing.T, goroutineCount int, store Store) {
	value := []byte("test value")

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(goroutineCount)
	for i := 0; i < goroutineCount; i++ {
		go interactWithStore(store, strconv.Itoa(i), value, t, &waitGroup)
	}
	waitGroup.Wait()

	// Now make sure that all values are in the store
	for i := 0; i < goroutineCount; i++ {
		actual, err := store.Get(strconv.Itoa(i))
		if err != nil {
			t.Errorf("An error occurred during the test: %v", err)
		}
		if diff := deep.Equal(actual, value); diff != nil {
			t.Error(diff)
		}
	}
}

// interactWithStore reads from and writes to the DB. Meant to be executed in a goroutine.
// Does NOT check if the DB works correctly (that's done elsewhere),
// only checks for errors that might occur due to concurrent access.
func interactWithStore(store Store, key string, value []byte, t *testing.T, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()
	_, err := store.Get(key)
	if err == nil {
		t.Errorf("Should return ErrKeyNotFound, but: %v", err)
	}
	err = store.Set(key, value)
	if err != nil {
		t.Error(err)
	}
	_, err = store.Get(key)
	if err != nil {
		t.Error(err)
	}
}
