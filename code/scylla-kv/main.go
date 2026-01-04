package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gocql/gocql"
)

// Item represents the data structure in our table
type Item struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
}

/*
cluster.Consistency = gocql.One:
This tells the driver: "I only need one server to confirm the write."
In a 3-node cluster, this is the fastest setting. If you set it to Quorum,
it would wait for 2 nodes.

keyExists function:
Notice I used SELECT key FROM items. I did not use SELECT *.
Since ScyllaDB stores data on disk, fetching the whole row (the JSON blob)
is expensive. Fetching just the key allows Scylla to often answer solely
from the in-memory Bloom Filter or Key Cache without reading the heavy
JSON data from the disk.
*/

func main() {
	// 1. CONNECT TO THE CLUSTER
	// We connect to "127.0.0.1" because we are running this from the host machine,
	// pointing to the Docker container mapped on port 9042.
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "key_value_store"
	cluster.Consistency = gocql.One // "One" is fastest (good for non-critical logging)
	cluster.ProtoVersion = 4

	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Could not connect to ScyllaDB: %v", err)
	}
	defer session.Close()

	fmt.Println("✅ Connected to ScyllaDB")

	// Test Data
	myKey := "user_123"
	myValue := `{"name": "Alice", "role": "admin"}` // storing JSON as a string

	// 2. CREATE (INSERT)
	err = createItem(session, myKey, myValue)
	if err != nil {
		log.Fatalf("Failed to insert: %v", err)
	}
	fmt.Printf("📝 Created key: %s\n", myKey)

	// 3. CHECK IF EXISTS
	exists := keyExists(session, myKey)
	fmt.Printf("🔍 Does key '%s' exist? %v\n", myKey, exists)

	// 4. READ
	item, err := readItem(session, myKey)
	if err != nil {
		log.Printf("Failed to read: %v", err)
	} else {
		fmt.Printf("📖 Read Value: %s (Created at: %v)\n", item.Value, item.CreatedAt)
	}

	// 5. DELETE
	err = deleteItem(session, myKey)
	if err != nil {
		log.Fatalf("Failed to delete: %v", err)
	}
	fmt.Printf("❌ Deleted key: %s\n", myKey)

	// 6. VERIFY DELETION
	exists = keyExists(session, myKey)
	fmt.Printf("🔍 Does key '%s' exist? %v\n", myKey, exists)
}

// createItem inserts a new row
func createItem(s *gocql.Session, key, value string) error {
	// CQL: INSERT INTO items ...
	query := "INSERT INTO items (key, value, created_at) VALUES (?, ?, ?)"
	return s.Query(query, key, value, time.Now()).Exec()
}

// readItem retrieves the value and timestamp
func readItem(s *gocql.Session, key string) (Item, error) {
	var i Item
	i.Key = key

	// CQL: SELECT ...
	query := "SELECT value, created_at FROM items WHERE key = ? LIMIT 1"

	// Scan matches the order of columns in the SELECT statement
	err := s.Query(query, key).Scan(&i.Value, &i.CreatedAt)
	return i, err
}

// keyExists checks efficiently using the Bloom Filter (no full row read)
func keyExists(s *gocql.Session, key string) bool {
	var tempKey string
	// We only select the KEY itself. This is highly optimized in Scylla
	// because it hits the Bloom Filter first.
	err := s.Query("SELECT key FROM items WHERE key = ? LIMIT 1", key).Scan(&tempKey)

	if err == gocql.ErrNotFound {
		return false
	} else if err != nil {
		log.Printf("Error checking existence: %v", err)
		return false
	}
	return true
}

// deleteItem removes the row
func deleteItem(s *gocql.Session, key string) error {
	return s.Query("DELETE FROM items WHERE key = ?", key).Exec()
}
