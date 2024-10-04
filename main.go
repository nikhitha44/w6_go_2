package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

// Recipe Management
//Defining datatypes - Recipe

type Recipe struct {
	ID           int      `json:"id"`
	Title        string   `json:"title"`
	Ingredients  []string `json:"ingredients"`
	Instructions string   `json:"instructions"`
	Status       string   `json:"status"` // "pending" or "completed"
}

var (
	recipes = make(map[int]Recipe)
	nextID  = 1
)

// Creating a new recipe - CREATE

func createRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe Recipe
	if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	recipe.ID = nextID
	nextID++
	recipes[recipe.ID] = recipe

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(recipe)
}

// Review the recipes available by using specific ID - READ

func getRecipes(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr != "" {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
			return
		}
		recipe, ok := recipes[id]
		if !ok {
			http.Error(w, "Recipe not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(recipe)
		return
	}

	var recipeList []Recipe
	for _, recipe := range recipes {
		recipeList = append(recipeList, recipe)
	}
	json.NewEncoder(w).Encode(recipeList)
}

// Updating the existing recipes by using ID - UPDATE

func updateRecipe(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing recipe ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	recipe, ok := recipes[id]
	if !ok {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	var updatedRecipe Recipe
	if err := json.NewDecoder(r.Body).Decode(&updatedRecipe); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Update only relevant fields

	recipe.Title = updatedRecipe.Title
	recipe.Ingredients = updatedRecipe.Ingredients
	recipe.Instructions = updatedRecipe.Instructions
	recipe.Status = updatedRecipe.Status
	recipes[id] = recipe

	json.NewEncoder(w).Encode(recipe)
}

//  Deleting a recipe by ID - DELETE

func deleteRecipe(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Missing recipe ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid recipe ID", http.StatusBadRequest)
		return
	}

	if _, ok := recipes[id]; !ok {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	}

	delete(recipes, id)
	w.WriteHeader(http.StatusNoContent)

}

func main() {
	http.HandleFunc("/recipes", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			createRecipe(w, r)
		case "GET":
			getRecipes(w, r)
		case "PUT":
			updateRecipe(w, r)
		case "DELETE":
			deleteRecipe(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// start the HTTP server

	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
