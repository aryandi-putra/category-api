package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

var categories = []Category{
	{ID: 1, Name: "Smartphone", Description: "A mobile phone with advanced computing capabilities"},
	{ID: 2, Name: "Laptop", Description: "A portable, all-in-one personal computer with a built-in screen, keyboard, and battery"},
	{ID: 3, Name: "Game Console", Description: "A specialized electronic device, essentially a dedicated computer, designed primarily for playing video games"},
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func parseIDFromPath(r *http.Request) (int, error) {
	idStr := r.PathValue("id")
	return strconv.Atoi(idStr)
}

func getCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func addCategory(w http.ResponseWriter, r *http.Request) {
	// read data from request
	var newCategory Category
	err := json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// add data to categories
	newCategory.ID = len(categories) + 1
	categories = append(categories, newCategory)

	respondWithJSON(w, http.StatusCreated, newCategory)
}

func getCategoryById(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Category Id")
		return
	}

	for _, category := range categories {
		if category.ID == id {
			respondWithJSON(w, http.StatusOK, category)
			return
		}
	}

	respondWithError(w, http.StatusNotFound, "Category not found")
}

func updateCategory(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Category Id")
		return
	}

	// get data from request
	var newCategory Category
	err = json.NewDecoder(r.Body).Decode(&newCategory)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	// loop categories, find match id and change by new category
	for i := range categories {
		if categories[i].ID == id {
			newCategory.ID = id
			categories[i] = newCategory
			respondWithJSON(w, http.StatusOK, newCategory)
			return
		}
	}

	respondWithError(w, http.StatusNotFound, "Category not found")
}

func deleteCategory(w http.ResponseWriter, r *http.Request) {
	id, err := parseIDFromPath(r)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Category ID")
		return
	}

	// loop categories, and delete the category if id found
	for i, category := range categories {
		if category.ID == id {
			// Create a new slice with the previous and next index data.
			categories = append(categories[:i], categories[i+1:]...)
			respondWithJSON(w, http.StatusOK, map[string]string{
				"message": "Category deleted successfully",
			})
			return
		}
	}

	respondWithError(w, http.StatusNotFound, "Category not found")
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /categories", getCategories)
	mux.HandleFunc("POST /categories", addCategory)
	mux.HandleFunc("GET /categories/{id}", getCategoryById)
	mux.HandleFunc("PUT /categories/{id}", updateCategory)
	mux.HandleFunc("DELETE /categories/{id}", deleteCategory)

	// localhost:8080/health
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, http.StatusOK, map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	fmt.Println("Server running on localhost:8080")

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Failed to running server")
	}
}
