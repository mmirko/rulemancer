package rulemancer

import (
	"fmt"
	"sync"
	"testing"
)

// TestCreateRoomConcurrent tests concurrent room creation
func TestCreateRoomConcurrent(t *testing.T) {
	engine := NewEngine()

	// Number of concurrent goroutines
	numGoroutines := 100
	roomsPerGoroutine := 10

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Channel to collect all created rooms
	roomsChan := make(chan *Room, numGoroutines*roomsPerGoroutine)

	// Launch multiple goroutines creating rooms concurrently
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < roomsPerGoroutine; j++ {
				room := engine.createRoom(
					fmt.Sprintf("Room-%d-%d", goroutineID, j),
					fmt.Sprintf("Description for room %d-%d", goroutineID, j),
				)
				roomsChan <- room
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(roomsChan)

	// Collect all rooms and verify uniqueness
	seenIDs := make(map[string]bool)
	roomCount := 0

	for room := range roomsChan {
		roomCount++
		if room == nil {
			t.Fatalf("CreateRoom returned nil")
		}

		if room.id == "" {
			t.Errorf("Room %d has empty ID", roomCount)
		}

		if seenIDs[room.id] {
			t.Errorf("Duplicate room ID found: %s", room.id)
		}
		seenIDs[room.id] = true

		if room.name == "" {
			t.Errorf("Room %s has empty name", room.id)
		}
	}

	expectedRooms := numGoroutines * roomsPerGoroutine
	if roomCount != expectedRooms {
		t.Errorf("Expected %d rooms, got %d", expectedRooms, roomCount)
	}

	t.Logf("Successfully created %d unique rooms concurrently", roomCount)
}

// TestCreateRoomConcurrentStress is a more intensive stress test
func TestCreateRoomConcurrentStress(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	engine := NewEngine()

	// Higher numbers for stress testing
	numGoroutines := 500
	roomsPerGoroutine := 20

	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// Track errors in concurrent operations
	type result struct {
		roomID string
		err    error
	}
	resultsChan := make(chan result, numGoroutines*roomsPerGoroutine)

	// Launch goroutines
	for i := 0; i < numGoroutines; i++ {
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < roomsPerGoroutine; j++ {
				room := engine.createRoom(
					fmt.Sprintf("StressRoom-%d-%d", goroutineID, j),
					fmt.Sprintf("Stress test room %d-%d", goroutineID, j),
				)
				if room == nil {
					resultsChan <- result{err: fmt.Errorf("nil room from goroutine %d, iteration %d", goroutineID, j)}
				} else {
					resultsChan <- result{roomID: room.id}
				}
			}
		}(i)
	}

	wg.Wait()
	close(resultsChan)

	// Verify all results
	seenIDs := make(map[string]bool)
	errorCount := 0

	for res := range resultsChan {
		if res.err != nil {
			t.Error(res.err)
			errorCount++
			continue
		}

		if seenIDs[res.roomID] {
			t.Errorf("Duplicate room ID in stress test: %s", res.roomID)
			errorCount++
		}
		seenIDs[res.roomID] = true
	}

	expectedRooms := numGoroutines * roomsPerGoroutine
	actualRooms := len(seenIDs)

	if errorCount > 0 {
		t.Errorf("Encountered %d errors during stress test", errorCount)
	}

	if actualRooms != expectedRooms {
		t.Errorf("Expected %d unique rooms, got %d", expectedRooms, actualRooms)
	}

	t.Logf("Stress test: Successfully created %d unique rooms across %d goroutines", actualRooms, numGoroutines)
}

// TestCreateRoomRaceConditions uses -race detector friendly patterns
func TestCreateRoomRaceConditions(t *testing.T) {
	engine := NewEngine()

	numGoroutines := 50
	var wg sync.WaitGroup
	wg.Add(numGoroutines)

	// All goroutines create rooms with same name/description to test ID generation
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			room := engine.createRoom("SameName", "SameDescription")
			if room == nil {
				t.Errorf("Goroutine %d: CreateRoom returned nil", id)
			}
		}(i)
	}

	wg.Wait()
	t.Log("Race condition test completed successfully")
}
