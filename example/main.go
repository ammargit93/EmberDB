package main

import (
	"emberdb/sdk"
	"fmt"
	"log"
)

func main() {
	// Initialize the sdk client
	client := sdk.NewClient("http://localhost:9182")
	defer client.Close()

	// Example 1: Set key-value pairs
	fmt.Println("1. Setting key-value pairs...")
	err := client.SetKey("users", "user:1", map[string]interface{}{
		"name":  "Alice",
		"email": "alice@example.com",
		"age":   28,
	})
	if err != nil {
		log.Printf("Error setting user:1: %v\n", err)
	} else {
		fmt.Println("✓ Set user:1 successfully")
	}

	err = client.SetKey("users", "user:2", map[string]interface{}{
		"name":  "Bob",
		"email": "bob@example.com",
		"age":   35,
	})
	if err != nil {
		log.Printf("Error setting user:2: %v\n", err)
	} else {
		fmt.Println("✓ Set user:2 successfully")
	}

	// Set some other data types
	err = client.SetKey("settings", "theme", "dark")
	if err != nil {
		log.Printf("Error setting theme: %v\n", err)
	} else {
		fmt.Println("✓ Set theme to 'dark'")
	}

	err = client.SetKey("counters", "page_views", 1234)
	if err != nil {
		log.Printf("Error setting page_views: %v\n", err)
	} else {
		fmt.Println("✓ Set page_views to 1234")
	}

	fmt.Println()

	// Example 2: Get key values
	fmt.Println("2. Retrieving key values...")
	value, err := client.GetKey("users", "user:1")
	if err != nil {
		log.Printf("Error getting user:1: %v\n", err)
	} else {
		fmt.Printf("✓ Retrieved user:1: %v\n", value)
	}

	theme, err := client.GetKey("settings", "theme")
	if err != nil {
		log.Printf("Error getting theme: %v\n", err)
	} else {
		fmt.Printf("✓ Retrieved theme: %v\n", theme)
	}

	fmt.Println()

	// Example 3: Get key with metadata
	fmt.Println("3. Retrieving key with metadata...")
	metadata, err := client.GetKeyWithMetadata("users", "user:2")
	if err != nil {
		log.Printf("Error getting metadata for user:2: %v\n", err)
	} else {
		fmt.Printf("✓ Retrieved metadata for user:2:\n")
		fmt.Printf("  Namespace: %s\n", metadata.Namespace)
		fmt.Printf("  Key: %s\n", metadata.Key)
		fmt.Printf("  Value: %v\n", metadata.Value)
	}

	fmt.Println()

	// Example 4: Check if key exists
	fmt.Println("4. Checking if keys exist...")
	exists, err := client.Exists("users", "user:1")
	if err != nil {
		log.Printf("Error checking user:1: %v\n", err)
	} else {
		fmt.Printf("✓ user:1 exists: %v\n", exists)
	}

	exists, err = client.Exists("users", "user:999")
	if err != nil {
		log.Printf("Error checking user:999: %v\n", err)
	} else {
		fmt.Printf("✓ user:999 exists: %v\n", exists)
	}

	fmt.Println()

	// Example 5: Update a key
	fmt.Println("5. Updating a key...")
	err = client.UpdateKey("users", "user:1", map[string]interface{}{
		"name":  "Alice Johnson",
		"email": "alice.johnson@example.com",
		"age":   29,
	})
	if err != nil {
		log.Printf("Error updating user:1: %v\n", err)
	} else {
		fmt.Println("✓ Updated user:1 successfully")
	}

	// Verify the update
	updated, err := client.GetKey("users", "user:1")
	if err != nil {
		log.Printf("Error retrieving updated user:1: %v\n", err)
	} else {
		fmt.Printf("✓ Verified update: %v\n", updated)
	}

	fmt.Println()

	// Example 6: Get all data
	fmt.Println("6. Retrieving all data...")
	allData, err := client.GetAll()
	if err != nil {
		log.Printf("Error getting all data: %v\n", err)
	} else {
		fmt.Println("✓ Retrieved all data:")
		for key, value := range allData {
			fmt.Printf("  %s: %v\n", key, value)
		}
	}

	fmt.Println()

	// Example 7: Delete a key
	fmt.Println("7. Deleting a key...")
	err = client.DeleteKey("counters", "page_views")
	if err != nil {
		log.Printf("Error deleting page_views: %v\n", err)
	} else {
		fmt.Println("✓ Deleted page_views successfully")
	}

	// Verify deletion
	exists, err = client.Exists("counters", "page_views")
	if err != nil {
		log.Printf("Error checking page_views: %v\n", err)
	} else {
		fmt.Printf("✓ page_views exists after deletion: %v\n", exists)
	}

	fmt.Println()

	// Example 8: Error handling - try to get a non-existent key
	fmt.Println("8. Error handling example...")
	_, err = client.GetKey("users", "user:nonexistent")
	if err != nil {
		fmt.Printf("✓ Expected error caught: %v\n", err)
	}

	fmt.Println("\n=== Example Complete ===")
}
