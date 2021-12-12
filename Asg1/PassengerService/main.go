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
type Passenger struct {
	PassengerID  int
	Username     string //Unique
	Password     string //(Not Retrieved)
	FirstName    string
	LastName     string
	MobileNo     string //varchar(8)
	EmailAddress string //Unique
}

var db *sql.DB

//==================== Auxiliary Functions ====================
func validKey(r *http.Request) bool {
	v := r.URL.Query()
	if key, ok := v["key"]; ok {
		if key[0] == "2c78afaf-97da-4816-bbee-9ad239abb296" {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

//==================== Database functions ====================

// Function to check that attributes that neeed to be unique are indeed unique
func UniquenessValidation(db *sql.DB, username string, emailAddress string) string {

	// Get any entry with the same attributes as the ones provided
	query := fmt.Sprintf("SELECT * FROM Passenger where Username = '%s' or EmailAddress = '%s'", username, emailAddress)

	results := db.QueryRow(query)

	errMsg := ""
	var userName string
	var email string

	//Do not retrieve data outside of whats neccessary
	var throwAway int
	var throwAway2 string

	//Retrieve results and see what triggered the uniqueness invalid if any
	//Map results retrieved to a passenger
	switch err := results.Scan(&throwAway, &userName, &throwAway2, &throwAway2, &throwAway2, &throwAway2, &email); err {
	case sql.ErrNoRows:
	case nil:
		// Get respective error message
		if userName == username {
			errMsg += "Username already in use. "
		}
		if email == emailAddress {
			errMsg += "Email Address already in use."
		}
	default:
		fmt.Printf("The HTTP request failed with error %s\n", err)
		panic(err.Error())
	}

	return errMsg
}

// Creates new passenger account
func CreatePassenger(db *sql.DB, p Passenger) {

	//PassengerID is auto incremented
	query := fmt.Sprintf("INSERT INTO Passenger (Username, `Password`, FirstName, LastName, MobileNo, EmailAddress) VALUES ('%s','%s','%s','%s','%s','%s')",
		p.Username, p.Password, p.FirstName, p.LastName, p.MobileNo, p.EmailAddress)

	_, err := db.Query(query)

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
		panic(err.Error())
	}
}

//Retrieve passenger details by username and password
func Login(db *sql.DB, username string, password string) (Passenger, string) {
	query := fmt.Sprintf("SELECT * FROM Passenger where Username = '%s' and `Password` = '%s'", username, password)

	//Get first entry, only one exists
	results := db.QueryRow(query)

	var passenger Passenger
	var errMsg string

	//Retrieve neccessary data
	//Map results retrieved to a passenger
	switch err := results.Scan(&passenger.PassengerID, &passenger.Username, &passenger.Password, &passenger.FirstName, &passenger.LastName, &passenger.MobileNo, &passenger.EmailAddress); err {
	case sql.ErrNoRows: //No rows found
		errMsg = "Account does not exist"
	case nil:
	default:
		panic(err.Error())
	}

	return passenger, errMsg
}

// Get passenger details by passenger ID
func GetByID(db *sql.DB, passengerID int) (Passenger, string) {
	query := fmt.Sprintf("SELECT * FROM Passenger where PassengerID = %d", passengerID)

	// Get first result, only one exists
	results := db.QueryRow(query)

	var passenger Passenger
	var errMsg string

	// Map result to a passenger
	switch err := results.Scan(&passenger.PassengerID, &passenger.Username, &passenger.Password, &passenger.FirstName, &passenger.LastName, &passenger.MobileNo, &passenger.EmailAddress); err {
	case sql.ErrNoRows: //If no result
		errMsg = "Account does not exist"
	case nil:
	default:
		panic(err.Error())
	}

	return passenger, errMsg
}

// Update passenger details by passenger ID
func UpdatePassenger(db *sql.DB, passengerID int, p Passenger) {
	// Update all details
	query := fmt.Sprintf("UPDATE Passenger SET Username = '%s', `Password` = '%s', FirstName = '%s', LastName = '%s', MobileNo = '%s', EmailAddress = '%s' WHERE PassengerID = %d",
		p.Username, p.Password, p.FirstName, p.LastName, p.MobileNo, p.EmailAddress, passengerID)

	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

// Delete passenger details by passenger ID
func DeletePassenger(db *sql.DB, passengerID int) string {
	query := fmt.Sprintf("DELETE FROM Passenger WHERE PassengerID=%d", passengerID)

	_, err := db.Query(query)
	var errMsg string

	if err != nil {
		errMsg = "Account does not exist"
	}
	return errMsg
}

//==================== HTTP Functions ====================

// Post method for a passenger account
func CreatePassengerAccount(w http.ResponseWriter, r *http.Request) {
	// Valid key for API check
	if !validKey(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)

	if err == nil { // If no error

		// Map json to passenger
		var passenger Passenger
		json.Unmarshal([]byte(reqBody), &passenger)

		// Check if all non-null information exist
		if passenger.Username == "" || passenger.Password == "" || passenger.FirstName == "" || passenger.LastName == "" || passenger.MobileNo == "" || passenger.EmailAddress == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply all neccessary passenger information "))
		} else { //all not null

			// Run db UniquenessValidation function
			errMsg := UniquenessValidation(db, passenger.Username, passenger.EmailAddress)
			if errMsg != "" { //not unique
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte(
					"422 - " + errMsg))
			} else { //unique
				// Run db CreatePassenger function
				CreatePassenger(db, passenger)
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte("201 - Passenger account created: " + passenger.Username))
			}
		}

	} else { //Incorrect format
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply passenger information in JSON format"))
	}
}

//Get passenger details with username and password
func GetPassengerDetails(w http.ResponseWriter, r *http.Request) {

	// Valid key for API check
	if !validKey(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	// Get query string parameters of username and password
	queryString := r.URL.Query()
	if _, ok := queryString["username"]; !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Required parameters not found"))
		return
	}
	if _, ok := queryString["password"]; !ok {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Required parameters not found"))
		return
	}

	var passenger Passenger
	var errMsg string

	// Run db Login function
	passenger, errMsg = Login(db, queryString["username"][0], queryString["password"][0])
	if errMsg == "Account does not exist" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No account found"))
	} else {
		// Return passenger
		json.NewEncoder(w).Encode(passenger)
	}
}

// Get passenger details with passengerID
func GetPassengerDetailsByID(w http.ResponseWriter, r *http.Request) {

	// Valid key for API check
	if !validKey(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	// Get param for passengerid
	params := mux.Vars(r)
	var passengerid int
	fmt.Sscan(params["passengerid"], &passengerid)

	var passenger Passenger
	var errMsg string

	// Run db GetByID function
	passenger, errMsg = GetByID(db, passengerid)
	if errMsg == "Account does not exist" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No account found"))
	} else {
		// Return passenger
		json.NewEncoder(w).Encode(passenger)
	}
}

// (PUT) Update all passenger details together
func UpdatePassengerDetails(w http.ResponseWriter, r *http.Request) {

	// Valid key for API check
	if !validKey(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	// Get param for passengerID
	params := mux.Vars(r)
	var passengerid int
	fmt.Sscan(params["passengerid"], &passengerid)

	reqBody, err := ioutil.ReadAll(r.Body)

	if err == nil {
		// Retrieve new object
		var passenger Passenger
		json.Unmarshal([]byte(reqBody), &passenger)

		// Check non-nullable attributes are not null
		if passenger.Username == "" || passenger.Password == "" || passenger.FirstName == "" || passenger.LastName == "" || passenger.MobileNo == "" || passenger.EmailAddress == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply all neccessary passenger information "))
		} else { // All not null
			// Run db UpdatePassenger function
			UpdatePassenger(db, passengerid, passenger)
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("202 - Account details updated"))
		}
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply passenger information in JSON format"))
	}
}

// Delete passenger by passengerID
func DeletePassengerAccount(w http.ResponseWriter, r *http.Request) {
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
	if queryString["useraccess"][0] != "Admin" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("401 - Unauthorized user"))
		return
	}

	// Get passenger id to delete
	params := mux.Vars(r)
	var passengerid int
	fmt.Sscan(params["passengerid"], &passengerid)

	var passenger Passenger

	// Run db DeletePassenger function
	errMsg := DeletePassenger(db, passengerid)
	if errMsg == "Account does not exist" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No account found"))
	} else {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("202 - Passenger Account Deleted: " + passenger.Username))
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

	router.HandleFunc("/api/v1/passengers", CreatePassengerAccount).Methods("POST")
	router.HandleFunc("/api/v1/passengers", GetPassengerDetails).Methods("GET")
	router.HandleFunc("/api/v1/passengers/{passengerid}", GetPassengerDetailsByID).Methods("GET")
	router.HandleFunc("/api/v1/passengers/{passengerid}", UpdatePassengerDetails).Methods("PUT")
	router.HandleFunc("/api/v1/passengers/{passengerid}", DeletePassengerAccount).Methods("DELETE")

	fmt.Println("Passenger Service operating on port 5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
