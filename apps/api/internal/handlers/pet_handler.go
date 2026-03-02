package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/my-pets/api/internal/models"
)

// In-memory store (replace with DB later)
var pets = []models.Pet{
	{ID: 1, Name: "Firulais", Species: "dog", Breed: "Labrador", Age: 3, Owner: "Juan", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	{ID: 2, Name: "Michi", Species: "cat", Breed: "Siamese", Age: 2, Owner: "Maria", CreatedAt: time.Now(), UpdatedAt: time.Now()},
}

var nextID uint = 3

func GetPets(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"data":  pets,
		"total": len(pets),
	})
}

func GetPet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	for _, pet := range pets {
		if pet.ID == uint(id) {
			c.JSON(http.StatusOK, gin.H{"data": pet})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "pet not found"})
}

func CreatePet(c *gin.Context) {
	var pet models.Pet
	if err := c.ShouldBindJSON(&pet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pet.ID = nextID
	pet.CreatedAt = time.Now()
	pet.UpdatedAt = time.Now()
	nextID++

	pets = append(pets, pet)
	c.JSON(http.StatusCreated, gin.H{"data": pet})
}

func UpdatePet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var input models.Pet
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i, pet := range pets {
		if pet.ID == uint(id) {
			input.ID = pet.ID
			input.CreatedAt = pet.CreatedAt
			input.UpdatedAt = time.Now()
			pets[i] = input
			c.JSON(http.StatusOK, gin.H{"data": pets[i]})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "pet not found"})
}

func DeletePet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	for i, pet := range pets {
		if pet.ID == uint(id) {
			pets = append(pets[:i], pets[i+1:]...)
			c.JSON(http.StatusOK, gin.H{"message": "pet deleted"})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "pet not found"})
}
