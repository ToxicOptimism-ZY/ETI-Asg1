package main

//==================== Imports ====================
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//==================== Structures & Variables ====================
type Trip struct {
	TripID      int
	PickUp      string //char(6)
	DropOff     string //char(6)
	DriverID    int
	PassengerID int
	Status      string
}

var db *sql.DB

//==================== Auxiliary Functions ====================
// Check for valid key within query string
func validKey(r *http.Request) bool {
	v := r.URL.Query()
	if key, ok := v["key"]; ok {
		if key[0] == "2c78afaf-97da-4816-bbee-9ad239abb298" {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

//==================== Database functions ====================

// Creates new trip record
func CreateTrip(db *sql.DB, t Trip) {

	// TripID is auto incremented
	query := fmt.Sprintf("INSERT INTO Trip (PickUp, DropOff, DriverID, PassengerID, `Status`) VALUES ('%s','%s',%d,%d,'%s')",
		t.PickUp, t.DropOff, t.DriverID, t.PassengerID, t.Status)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

// Get list of trips by passenger
func GetByPassengerID(db *sql.DB, passengerID int) ([]Trip, string) {
	query := fmt.Sprintf("SELECT * FROM Trip where PassengerID = '%d'", passengerID)

	// Get all results
	results, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}

	var trips []Trip
	errMsg := "placeholder" //Temporary placeholder till any results existing determined

	// Loop through results
	for results.Next() {
		// Map a row to a Trip
		var trip Trip
		err := results.Scan(&trip.TripID, &trip.PickUp, &trip.DropOff, &trip.DriverID, &trip.PassengerID, &trip.Status)
		if err != nil {
			panic(err.Error())
		}

		errMsg = ""
		// Append mapped trip to trip array
		trips = append(trips, trip)
	}

	// If no result
	if errMsg != "" {
		errMsg = "No trips made by passenger"
	}
	return trips, errMsg
}

// Get passenger's current trip through the passengerID and trip Status (Not "Completed")
func GetPassengerCurrentTrip(db *sql.DB, passengerID int, status string) (Trip, string) {
	query := fmt.Sprintf("SELECT * FROM Trip where PassengerID = %d and `Status` = '%s'", passengerID, status)

	//Get first result, only one should exist
	results := db.QueryRow(query)

	var trip Trip
	var errMsg string

	// Map result to Trip
	switch err := results.Scan(&trip.TripID, &trip.PickUp, &trip.DropOff, &trip.DriverID, &trip.PassengerID, &trip.Status); err {
	case sql.ErrNoRows: //If no result
		errMsg = "Trip does not exist"
	case nil:
	default:
		panic(err.Error())
	}

	return trip, errMsg
}

// Get driver's current trip through the driverID and trip Status (Not "Completed")et passenger's current trip through the passengerID and trip Status (Not "Completed")
func GetDriverCurrentTrip(db *sql.DB, driverID int, status string) (Trip, string) {
	query := fmt.Sprintf("SELECT * FROM Trip where DriverID = %d and `Status` = '%s'", driverID, status)

	// Get first result, only one should exist
	results := db.QueryRow(query)

	var trip Trip
	var errMsg string

	// Map result to a Trip
	switch err := results.Scan(&trip.TripID, &trip.PickUp, &trip.DropOff, &trip.DriverID, &trip.PassengerID, &trip.Status); err {
	case sql.ErrNoRows: // If no result
		errMsg = "Trip does not exist"
	case nil:
	default:
		panic(err.Error())
	}

	return trip, errMsg
}

// Get trip details by trip ID
func GetByTripID(db *sql.DB, tripID int) (Trip, string) {
	query := fmt.Sprintf("SELECT * FROM Trip where TripID = %d", tripID)

	// Get first result, only one exists
	results := db.QueryRow(query)

	var trip Trip
	var errMsg string

	// Map result to a trip
	switch err := results.Scan(&trip.TripID, &trip.PickUp, &trip.DropOff, &trip.DriverID, &trip.PassengerID, &trip.Status); err {
	case sql.ErrNoRows: //If no result
		errMsg = "Trip does not exist"
	case nil:
	default:
		panic(err.Error())
	}

	return trip, errMsg
}

// Update trip details by Trip ID
func UpdateTrip(db *sql.DB, tripID int, t Trip) {
	// Update all details
	query := fmt.Sprintf("UPDATE Trip SET PickUp = '%s', DropOff = '%s', DriverID = %d, `Status` = '%s' WHERE TripID = %d",
		t.PickUp, t.DropOff, t.DriverID, t.Status, tripID)

	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

// Delete trip details by trip ID
func DeleteTrip(db *sql.DB, tripID int) string {
	query := fmt.Sprintf("DELETE FROM Trip WHERE TripID=%d", tripID)

	_, err := db.Query(query)
	var errMsg string

	if err != nil {
		errMsg = "Trip does not exist"
	}
	return errMsg
}

//==================== HTTP Functions ====================

//Post method for a trip record
func CreateTripRecord(w http.ResponseWriter, r *http.Request) {
	// Valid key for API check
	if !validKey(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)

	if err == nil { // If no error

		// Map json to trip
		var trip Trip
		json.Unmarshal([]byte(reqBody), &trip)

		// Check if all non-null information exist
		if trip.PickUp == "" || trip.DropOff == "" || trip.DriverID == 0 || trip.PassengerID == 0 || trip.Status == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply all neccessary driver information "))
		} else { // all not null
			// Run db CreateTrip function
			CreateTrip(db, trip)
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("201 - Trip created: " + trip.PickUp + " to " + trip.DropOff))
		}

	} else { //incorrect format
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply trip information in JSON format"))
	}
}

// Help function that calls appropriate function in accordance to parameters in the query string
func GetTripsQueryStringValidator(w http.ResponseWriter, r *http.Request) {

	// Get query string parameters
	queryString := r.URL.Query()
	_, okPassenger := queryString["passengerid"]
	_, okDriver := queryString["driverid"]
	_, okStatus := queryString["status"]

	// If passengerId and password passed in, get current trip for passengers function
	if okPassenger && okStatus {
		// Run HTTP GetCurrentTripDetailsForPassenger function
		GetCurrentTripDetailsForPassenger(w, r)
		return
	} else if okDriver && okStatus { // If driverId and password passed in, get current trips for driver function
		// Run HTTP GetCurrentTripDetailsForDrivers function
		fmt.Print("Access")
		GetCurrentTripDetailsForDriver(w, r)
		return
	} else if okPassenger { // If passengerID only passed in, get all trips made by a passenger
		// Run HTTP GetTripDetailsByPassengerID function
		GetTripDetailsByPassengerID(w, r)
		return
	} else { //else no appropriate function
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Required parameters not found"))
		return
	}
}

// Get trip details with passenger ID
func GetTripDetailsByPassengerID(w http.ResponseWriter, r *http.Request) {

	// Valid key for API check
	if !validKey(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	// Get query string parameters of passengerID
	queryString := r.URL.Query()
	var passengerid int
	fmt.Sscan(queryString["passengerid"][0], &passengerid)

	var trips []Trip
	var errMsg string

	// Run db GetByPassengerID function
	trips, errMsg = GetByPassengerID(db, passengerid)
	if errMsg != "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - " + errMsg))
	} else {
		// Return trips array
		json.NewEncoder(w).Encode(trips)
	}
}

// Get current trip details for passenger
func GetCurrentTripDetailsForPassenger(w http.ResponseWriter, r *http.Request) {
	// Valid key for API check
	if !validKey(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	// Get query string parameters of passenger ID and status
	queryString := r.URL.Query()
	var passengerid int
	fmt.Sscan(queryString["passengerid"][0], &passengerid)

	var trip Trip
	var errMsg string

	// Run db GetPassengerCurrentTrip function
	trip, errMsg = GetPassengerCurrentTrip(db, passengerid, queryString["status"][0])
	if errMsg == "Trip does not exist" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - " + errMsg))
	} else {
		// Return trip
		json.NewEncoder(w).Encode(trip)
	}
}

// Get current trip details for driver
func GetCurrentTripDetailsForDriver(w http.ResponseWriter, r *http.Request) {

	// Valid key for API check
	if !validKey(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	// Get query string parameters of driver ID and status
	queryString := r.URL.Query()
	var driverid int
	fmt.Sscan(queryString["driverid"][0], &driverid)

	var trip Trip
	var errMsg string

	// Run db GetDriverCurrentTrip function
	trip, errMsg = GetDriverCurrentTrip(db, driverid, queryString["status"][0])
	if errMsg == "Trip does not exist" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - " + errMsg))
	} else {
		// Return trip
		json.NewEncoder(w).Encode(trip)
	}
}

// Get trip details with trip ID
func GetTripDetailsByTripID(w http.ResponseWriter, r *http.Request) {

	// Valid key for API check
	if !validKey(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	// Get param for tripid
	params := mux.Vars(r)
	var tripid int
	fmt.Sscan(params["tripid"], &tripid)

	var trip Trip
	var errMsg string

	// Run db GetByTripId function
	trip, errMsg = GetByTripID(db, tripid)
	if errMsg == "Trip does not exist" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - " + errMsg))
	} else {
		// Return trip
		json.NewEncoder(w).Encode(trip)
	}
}

// (PUT) Update all trip details together
func UpdateTripDetails(w http.ResponseWriter, r *http.Request) {

	// Valid key for API check
	if !validKey(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	// Get param for tripID
	params := mux.Vars(r)
	var tripid int
	fmt.Sscan(params["tripid"], &tripid)

	reqBody, err := ioutil.ReadAll(r.Body)

	if err == nil {
		// Retrieve new object
		var trip Trip
		json.Unmarshal([]byte(reqBody), &trip)

		// Check non-nullable attributes are not null
		if trip.PickUp == "" || trip.DropOff == "" || trip.DriverID == 0 || trip.PassengerID == 0 || trip.Status == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply all trip trip information "))
		} else { // All not null
			// Run db UpdateTrip function
			UpdateTrip(db, tripid, trip)
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("202 - Trip details updated"))
		}

	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply trip information in JSON format"))
	}
}

// Delete trip by TripID
func DeleteTripRecord(w http.ResponseWriter, r *http.Request) {
	if !validKey(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	// Get query string parameters
	queryString := r.URL.Query()
	if _, ok := queryString["useraccess"]; !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Required parameter not found"))
		return
	}

	// If user not authorized, reject
	if queryString["useraccess"][0] != "Admin" && queryString["useraccess"][0] != "Passenger" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Unauthorized user"))
		return
	}

	// Get trip id to delete
	params := mux.Vars(r)
	var tripid int
	fmt.Sscan(params["tripid"], &tripid)

	var trip Trip

	// Run db DeleteDriver function
	errMsg := DeleteTrip(db, tripid)
	if errMsg == "Trip does not exist" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No trip found"))
	} else {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("202 - Trip deleted: " + trip.PickUp + " to " + trip.DropOff))
	}
}

//==================== Main ====================
func main() {

	// Open connection
	var err error
	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/asg1")

	// Handle error
	if err != nil {
		panic(err.Error())
	}

	// Define url functions
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/trips", CreateTripRecord).Methods("POST")
	router.HandleFunc("/api/v1/trips", GetTripsQueryStringValidator).Methods("GET")
	router.HandleFunc("/api/v1/trips/{tripid}", GetTripDetailsByTripID).Methods("GET")
	router.HandleFunc("/api/v1/trips/{tripid}", UpdateTripDetails).Methods("PUT")
	router.HandleFunc("/api/v1/trips/{tripid}", DeleteTripRecord).Methods("DELETE")

	fmt.Println("Trip Service operating on port 5002")
	log.Fatal(http.ListenAndServe(":5002", router))

}
