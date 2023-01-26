package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Type of the stored cars

type car struct {
	ID            string `json:"id"`
	Licence_plate string `json:"licence_plate"`
	Owner         string `json:"owner"`
	Power         int    `json:"power"`
}

// Very basic authentication with one user

const Password = "12345"

var Login = false

// If the string given as an argument matches the Password constant,
// the Login variable is set to true thus enabling the user
// to use the deleting and modifying features.

func login(context *gin.Context) {
	password := context.Param("password")
	if password == Password {
		context.IndentedJSON(http.StatusOK, gin.H{"message": "You logged in."})
		Login = true
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{"message": "Incorrect password."})
}

func logout(context *gin.Context) {
	if Login {
		Login = false
		context.IndentedJSON(http.StatusOK, gin.H{"message": "You logged out."})
		return
	}
	context.IndentedJSON(http.StatusOK, gin.H{"message": "Log in to log out."})
}

// Datatsructure storing the cars

var cars = []car{
	{ID: "0", Licence_plate: "ABC123", Owner: "Gábor", Power: 140},
	{ID: "1", Licence_plate: "ABC321", Owner: "Zoltán", Power: 200},
	{ID: "2", Licence_plate: "CBA123", Owner: "Béla", Power: 300},
	{ID: "3", Licence_plate: "CBA321", Owner: "Jani", Power: 100},
	{ID: "4", Licence_plate: "OBJ140", Owner: "Vlad", Power: 580},
}

// Finds the first unused ID to then assign it to a newly added car

func nextID() int {
	ID := 0
	carindex, err := getCarByID(strconv.Itoa(ID))
	for err == nil {
		ID++
		carindex, err = getCarByID(strconv.Itoa(ID))
	}
	carindex++
	return ID
}

// Sends all the cars currently stored in a JSON format

func getcars(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, cars)
}

// Adds a new car to the database
// As an input expects a JSON with one element
// Only available after logging in

func addcar(context *gin.Context) {

	if !Login {
		context.IndentedJSON(http.StatusOK, gin.H{"message": "You are not logged in."})
		return
	}

	var newcar car
	if err := context.BindJSON(&newcar); err != nil {
		context.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Incorrect JSON format."})
		return
	}

	newcar.ID = strconv.Itoa(nextID())

	cars = append(cars, newcar)
	context.IndentedJSON(http.StatusCreated, cars)
}

// Returns the index of a car in the cars datastructur by an ID given as a parameter

func getCarByID(id string) (int, error) {
	for i, b := range cars {
		if b.ID == id {
			return i, nil
		}
	}

	return 0, errors.New("Car not found")
}

// Sends a single element JSON of the car with the matching ID

func carByID(context *gin.Context) {

	id := context.Param("id")
	carindex, err := getCarByID(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Car not found."})
		return
	}
	context.IndentedJSON(http.StatusOK, cars[carindex])
}

// Deletes the car with the matching ID from the datastructure
// Only available after logging in

func deleteCar(context *gin.Context) {

	if !Login {
		context.IndentedJSON(http.StatusOK, gin.H{"message": "You are not logged in."})
		return
	}

	id := context.Param("id")
	carindex, err := getCarByID(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Car not found."})
		return
	}

	cars[carindex] = cars[len(cars)-1]
	cars = cars[:len(cars)-1]

	context.IndentedJSON(http.StatusOK, cars)
}

// Modifies a car's data in the datastructure
// As an input expects a JSON with one element
// Modifies the car that has the same ID as in the parameter JSON
// Hence anything is modifiable except the ID
// Only available after logging in

func modifyCar(context *gin.Context) {
	if !Login {
		context.IndentedJSON(http.StatusOK, gin.H{"message": "You are not logged in."})
		return
	}

	var newcar car

	if err := context.BindJSON(&newcar); err != nil {
		return
	}

	id := newcar.ID
	carindex, err := getCarByID(id)

	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"message": "Car not found."})
		return
	}

	cars[carindex] = newcar
	context.IndentedJSON(http.StatusOK, cars)
}

func main() {
	router := gin.Default()

	// GET request to return all the cars in JSON
	router.GET("/cars", getcars)

	// GET request to return one car in JSON by an ID parameter
	router.GET("/car/:id", carByID)

	// GET request to log in
	router.GET("/login/:password", login)

	// GET request to log out
	router.GET("/logout", logout)

	// POST request to add a new element
	router.POST("/cars", addcar)

	// PATCH request to delete an element
	router.PATCH("/delete/:id", deleteCar)

	// PATCH request to modify an element
	router.PATCH("/modify", modifyCar)

	router.Run("localhost:1000")
}
