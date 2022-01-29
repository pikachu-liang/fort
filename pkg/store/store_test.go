package store

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

type testCase = struct {
	Name  string
	Store Store
	path  string
}

func prepareTestCases(t *testing.T) []testCase {
	badgerStore, badgerPath := createBadgerStore(t)
	bboltStore, bblotPath := createBBoltStore(t)

	return []testCase{
		{
			"Badger",
			badgerStore,
			badgerPath,
		},
		{
			"BBolt",
			bboltStore,
			bblotPath,
		},
	}
}

// TestStore tests basic operations for each storage driver
func TestStore(t *testing.T) {
	testCases := prepareTestCases(t)
	for _, testCase := range testCases {
		defer cleanUp(testCase.Store, testCase.path)
		TestBasicInteractions(testCase.Store, t)
	}
}

// TestStoreConcurrent tests concurrent access operations for each storage driver
func TestStoreConcurrent(t *testing.T) {
	testCases := prepareTestCases(t)
	for _, testCase := range testCases {
		defer cleanUp(testCase.Store, testCase.path)
		goroutineCount := 1000
		TestConcurrentInteractions(t, goroutineCount, testCase.Store)
	}
}

// TestErrors tests some error cases.
func TestErrors(t *testing.T) {
	testCases := prepareTestCases(t)
	for _, testCase := range testCases {
		store := testCase.Store
		path := testCase.path
		defer cleanUp(store, path)
		err := store.Set("", nil)
		if err == nil {
			t.Error("Expect ErrKeyInvalid")
		}
		_, err = store.Get("")
		if err == nil {
			t.Error("Expect ErrKeyInvalid")
		}
		err = store.Delete("")
		if err == nil {
			t.Error("Expect ErrKeyInvalid")
		}
		_, err = store.Get("foo")
		if err == nil {
			t.Error("Expect ErrKeyNotFound")
		}
	}
}

// TestClose tests if the close method returns any errors.
func TestClose(t *testing.T) {
	testCases := prepareTestCases(t)
	for _, testCase := range testCases {
		store := testCase.Store
		path := testCase.path
		defer os.RemoveAll(path)
		err := store.Close()
		if err != nil {
			t.Error(err)
		}
	}
}

// TestBadgerStoreNonExistingDir tests whether the BadgerDB store can create the given directory on its own.
// When using BadgerDB directly, it requires the given path to exist and to be writeable.
func TestBadgerStoreNonExistingDir(t *testing.T) {
	tmpDir := os.TempDir() + "/BadgerDB"
	err := os.RemoveAll(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	options := BadgerStoreOptions{
		Dir: tmpDir,
	}
	store, err := NewBadgerStore(options)
	defer cleanUp(store, tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	err = store.Set("foo", []byte("bar"))
	if err != nil {
		t.Error(err)
	}
}

// TestBBoltStoreNonExistingDir tests whether the BBoltStore can create the given directory on its own.
func TestBBoltStoreNonExistingDir(t *testing.T) {
	tmpDir := os.TempDir() + "/bbolt"
	err := os.RemoveAll(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	options := BBoltStoreOptions{
		Path: tmpDir,
	}
	store, err := NewBBoltStore(options)
	defer cleanUp(store, tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	err = store.Set("foo", []byte("bar"))
	if err != nil {
		t.Error(err)
	}
}

func createBadgerStore(t *testing.T) (Store, string) {
	dir := generateRandomTempDBpath(t, "BadgerDB")
	options := BadgerStoreOptions{
		Dir: dir,
	}
	store, err := NewBadgerStore(options)
	if err != nil {
		t.Fatal(err)
	}
	return store, dir
}

func createBBoltStore(t *testing.T) (Store, string) {
	path := generateRandomTempDBpath(t, "BBoltDB") + "bbolt_store.db"
	options := BBoltStoreOptions{
		Path: path,
	}
	store, err := NewBBoltStore(options)
	if err != nil {
		t.Fatal(err)
	}
	return store, path
}

func generateRandomTempDBpath(t *testing.T, folderName string) string {
	path, err := ioutil.TempDir(os.TempDir(), folderName)
	if err != nil {
		t.Fatalf("Generating random DB path failed: %v", err)
	}
	return path
}

// cleanUp cleans up (deletes) the database files that have been created during a test.
// If an error occurs the test is NOT marked as failed.
func cleanUp(store Store, path string) {
	err := store.Close()
	if err != nil {
		log.Printf("Error during cleaning up after a test (during closing the store): %v\n", err)
	}
	err = os.RemoveAll(path)
	if err != nil {
		log.Printf("Error during cleaning up after a test (during removing the data directory): %v\n", err)
	}
}
