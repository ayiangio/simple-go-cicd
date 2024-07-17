package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestCreateItem tests the CreateItem function to ensure it works correctly.
func TestCreateItem(t *testing.T) {
	clearItems()
	item := Item{
		Name:   "Test Item",
		Rating: 4.5,
		Stock:  10,
		Price:  25.5,
	}

	itemJSON, _ := json.Marshal(item)
	req, err := http.NewRequest("POST", "/items", bytes.NewBuffer(itemJSON))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(createItem)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status %d; got %d", http.StatusCreated, rr.Code)
	}

	var responseJSON JSONResponse
	if err := json.NewDecoder(rr.Body).Decode(&responseJSON); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	responseItem, ok := responseJSON.Data.(map[string]interface{})
	if !ok {
		t.Errorf("Expected response data to be an item; got %v", responseJSON.Data)
	}

	id, ok := responseItem["id"].(float64) // JSON numbers are float64
	if !ok {
		t.Errorf("Expected item ID to be a number; got %v", responseItem["id"])
	}

	if id != 1 {
		t.Errorf("Expected item ID to be 1; got %d", int(id))
	}
}
// TestDeleteItem tests the DeleteItem function to ensure it works correctly.
func TestDeleteItem(t *testing.T) {
	clearItems()
	addSampleItem()

	req, err := http.NewRequest("DELETE", "/items/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteItem)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status %d; got %d", http.StatusNoContent, rr.Code)
	}

	// Check if item is deleted
	req, err = http.NewRequest("GET", "/items", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler = http.HandlerFunc(getItems)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status %d; got %d", http.StatusOK, rr.Code)
	}

	var responseJSON JSONResponse
	if err := json.NewDecoder(rr.Body).Decode(&responseJSON); err != nil {
		t.Errorf("Error decoding response body: %v", err)
	}

	itemsList, ok := responseJSON.Data.([]interface{})
	if !ok {
		t.Errorf("Expected items list to be an array; got %v", responseJSON.Data)
	}

	if len(itemsList) != 0 {
		t.Errorf("Expected items list to be empty after deletion; got %v", itemsList)
	}
}

func clearItems() {
	items = []Item{}
	nextID = 1
}

func addSampleItem() {
	item := Item{
		ID:     1,
		Name:   "Sample Item",
		Rating: 4.0,
		Stock:  5,
		Price:  20.0,
	}
	items = append(items, item)
}
