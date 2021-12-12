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
type Driver struct {
	DriverID        int
	Username        string //Unique
	Password        string //(Not Retrieved)
	FirstName       string
	LastName        string
	MobileNo        string //varchar(8)
	EmailAddress    string //Unique
	NRIC            string //1 letter, 7 digits, 1 checksum letter, Unique (Not Retrieved)
	CarLicencePlate string //S, 2 letters, 4 numbers, 1 checksum letter, Unique
	Status          string
}

var db *sql.DB

//==================== Auxiliary Functions ====================
// Check for valid key within query string
func validKey(r *http.Request) bool {
	v := r.URL.Query()
	if key, ok := v["key"]; ok {
		if key[0] == "2c78afaf-97da-4816-bbee-9ad239abb297" {
			return true
		} else {
			return false //invalid
		}
	} else {
		return false //invalid
	}
}

//==================== Database functions ====================

// Function to check that attributes that neeed to be unique are indeed unique
func UniquenessValidation(db *sql.DB, username string, emailAddress string, nric string, carLicencePlate string) string {

	// Get any entry with the same attributes as the ones provided
	query := fmt.Sprintf("SELECT * FROM Driver where Username = '%s' or EmailAddress = '%s' or  NRIC = '%s' or carLicencePlate = '%s'", username, emailAddress, nric, carLicencePlate)

	//Retrieve first result
	results := db.QueryRow(query)

	errMsg := ""
	var userName string
	var email string
	var ic string
	var licenceNo string

	//Do not retrieve data outside of whats neccessary
	var throwAway int
	var throwAway2 string

	//Retrieve first result and see what triggered the uniqueness invalid if any
	//Map results retrieved to a driver
	switch err := results.Scan(&throwAway, &userName, &throwAway2, &throwAway2, &throwAway2, &throwAway2, &email, &ic, &licenceNo, &throwAway2); err {
	case sql.ErrNoRows:
	case nil:
		// Get respective error message
		if userName == username {
			errMsg += "Username already in use. "
		}
		if email == emailAddress {
			errMsg += "Email Address already in use. "
		}
		if ic == nric {
			errMsg += "NRIC already in use. "
		}
		if licenceNo == carLicencePlate {
			errMsg += "Licence Plate already in use. "
		}
	default:
		panic(err.Error())
	}

	return errMsg
}

// Creates new driver account
func CreateDriver(db *sql.DB, d Driver) {

	//DriverID is auto incremented
	query := fmt.Sprintf("INSERT INTO Driver (Username, `Password`, FirstName, LastName, MobileNo, EmailAddress, NRIC, CarLicencePlate, `Status`) VALUES ('%s','%s','%s','%s','%s','%s','%s','%s','%s')",
		d.Username, d.Password, d.FirstName, d.LastName, d.MobileNo, d.EmailAddress, d.NRIC, d.CarLicencePlate, d.Status)

	_, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}
}

//Retrieve driver details by username and password
func Login(db *sql.DB, username string, password string) (Driver, string) {
	query := fmt.Sprintf("SELECT * FROM Driver where Username = '%s' and `Password` = '%s'", username, password)

	//Get first entry, only one exists
	results := db.QueryRow(query)

	var driver Driver
	var errMsg string
	var throwAway string //do not retrieve NRIC as that is personal data

	//Retrieve neccessary data
	//Map results retrieved to a driver
	switch err := results.Scan(&driver.DriverID, &driver.Username, &driver.Password, &driver.FirstName, &driver.LastName, &driver.MobileNo, &driver.EmailAddress, &throwAway, &driver.CarLicencePlate, &driver.Status); err {
	case sql.ErrNoRows: //No rows found
		errMsg = "Account does not exist"
	case nil:
	default:
		panic(err.Error())
	}

	return driver, errMsg
}

// Get list of drivers by status - "Available" or "Unavailable"
func GetByStatus(db *sql.DB, status string) ([]Driver, string) {
	query := fmt.Sprintf("SELECT * FROM Driver where `Status` = '%s'", status)

	// Get all results
	results, err := db.Query(query)

	if err != nil {
		panic(err.Error())
	}

	var drivers []Driver
	errMsg := "placeholder" //Temporary placeholder till any results existing determined
	var throwAway string

	// Loop through results
	for results.Next() {
		// Map a row to a driver
		var driver Driver
		err := results.Scan(&driver.DriverID, &driver.Username, &driver.Password, &driver.FirstName, &driver.LastName, &driver.MobileNo, &driver.EmailAddress, &throwAway, &driver.CarLicencePlate, &driver.Status)
		if err != nil {
			panic(err.Error())
		}

		errMsg = ""

		// Append mapped driver to driver array
		drivers = append(drivers, driver)
	}

	// If no result
	if errMsg != "" {
		errMsg = "No account of status: " + status
	}

	return drivers, errMsg
}

// Get driver details by passenger ID
func GetByID(db *sql.DB, driverID int) (Driver, string) {
	query := fmt.Sprintf("SELECT * FROM Driver where driverID = %d", driverID)

	// Get first result, only one exists
	results := db.QueryRow(query)

	var driver Driver
	var errMsg string
	var throwAway string

	// Map result to a driver
	switch err := results.Scan(&driver.DriverID, &driver.Username, &driver.Password, &driver.FirstName, &driver.LastName, &driver.MobileNo, &driver.EmailAddress, &throwAway, &driver.CarLicencePlate, &driver.Status); err {
	case sql.ErrNoRows: //If no result
		errMsg = "Account does not exist"
	case nil:
	default:
		panic(err.Error())
	}

	return driver, errMsg
}

// (PUT) Update driver details by driver ID
func UpdateDriver(db *sql.DB, driverID int, d Driver) {
	// Update all details including *STATUS, retrieve most recent status before updating
	query := fmt.Sprintf("UPDATE Driver SET Username = '%s', `Password` = '%s', FirstName = '%s', LastName = '%s', MobileNo = '%s', EmailAddress = '%s', CarLicencePlate = '%s', Status = '%s' WHERE DriverID = %d",
		d.Username, d.Password, d.FirstName, d.LastName, d.MobileNo, d.EmailAddress, d.CarLicencePlate, d.Status, driverID)

	_, err := db.Query(query)
	if err != nil {
		panic(err.Error())
	}
}

// Delete driver details by driver ID
func DeleteDriver(db *sql.DB, driverID int) string {
	query := fmt.Sprintf("DELETE FROM Passenger WHERE PassengerID=%d", driverID)

	_, err := db.Query(query)
	var errMsg string

	if err != nil {
		errMsg = "Account does not exist"
	}
	return errMsg
}

//==================== HTTP Functions ====================

// Post method for a driver account
func CreateDriverAccount(w http.ResponseWriter, r *http.Request) {

	// Valid key for API check
	if !validKey(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)

	if err == nil { // If no error

		// Map json to driver
		var driver Driver
		json.Unmarshal([]byte(reqBody), &driver)

		// Check if all non-null information exist
		if driver.Username == "" || driver.Password == "" || driver.FirstName == "" || driver.LastName == "" || driver.MobileNo == "" || driver.EmailAddress == "" || driver.NRIC == "" || driver.CarLicencePlate == "" || driver.Status == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply all neccessary driver information "))
		} else { //all not null

			// Run db UniquenessValidation function
			errMsg := UniquenessValidation(db, driver.Username, driver.EmailAddress, driver.NRIC, driver.CarLicencePlate)
			if errMsg != "" { //not unique
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte(
					"422 - " + errMsg))
			} else { //unique
				// Run db CreateDriver function
				CreateDriver(db, driver)
				w.WriteHeader(http.StatusCreated)
				w.Write([]byte("201 - Driver account created: " + driver.Username))
			}
		}

	} else { //incorrect format
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply driver information in JSON format"))
	}
}

// Helper function that calls appropriate function in accordance to their parameters
func GetDriversQueryStringValidator(w http.ResponseWriter, r *http.Request) {

	// Get query string parameters
	queryString := r.URL.Query()
	_, okUser := queryString["username"]
	_, okPassword := queryString["password"]
	_, okStatus := queryString["status"]

	// If user and password passed in, log in function
	if okUser && okPassword {
		// Run HTTP GetDriverDetails function
		GetDriverDetails(w, r)
		return
	} else if okStatus { //If status passed in, do find by status
		// Run HTTP GetDriversByStatus function
		GetDriversByStatus(w, r)
		return
	} else { //else no appropriate function
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Required parameters not found"))
		return
	}
}

//Get driver details with username and password
func GetDriverDetails(w http.ResponseWriter, r *http.Request) {

	// Valid key for API check
	if !validKey(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	// Get query string parameters
	queryString := r.URL.Query()

	var driver Driver
	var errMsg string

	// Run db Login function
	driver, errMsg = Login(db, queryString["username"][0], queryString["password"][0])
	if errMsg == "Account does not exist" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No account found"))
	} else {
		// Return driver
		json.NewEncoder(w).Encode(driver)
	}

}

// Get driver details by status
func GetDriversByStatus(w http.ResponseWriter, r *http.Request) {

	// Valid key for API check
	if !validKey(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	// Get query string parameters of status
	queryString := r.URL.Query()

	var drivers []Driver
	var errMsg string

	// Run db GetByStatus function
	drivers, errMsg = GetByStatus(db, queryString["status"][0])
	if errMsg != "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - " + errMsg))
	} else {
		// Return driver array
		json.NewEncoder(w).Encode(drivers)
	}
}

// Get driver details with driverID
func GetDriverDetailsByID(w http.ResponseWriter, r *http.Request) {

	// Valid key for API check
	if !validKey(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	// Get param for driverid
	params := mux.Vars(r)
	var driverid int
	fmt.Sscan(params["driverid"], &driverid)

	var driver Driver
	var errMsg string

	// Run db GetByID function
	driver, errMsg = GetByID(db, driverid)
	if errMsg == "Account does not exist" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No account found"))
	} else {
		// Return driver
		json.NewEncoder(w).Encode(driver)
	}
}

// Update all driver details together, for *STATUS, retrieve latest before updating
func UpdateDriverDetails(w http.ResponseWriter, r *http.Request) {

	// Valid key for API check
	if !validKey(r) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 - Invalid key"))
		return
	}

	// Get param for driverID
	params := mux.Vars(r)
	var driverid int
	fmt.Sscan(params["driverid"], &driverid)

	reqBody, err := ioutil.ReadAll(r.Body)

	if err == nil {
		// Retrieve new object
		var driver Driver
		json.Unmarshal([]byte(reqBody), &driver)

		// Check non-nullable attributes are not null
		if driver.Username == "" || driver.Password == "" || driver.FirstName == "" || driver.LastName == "" || driver.MobileNo == "" || driver.EmailAddress == "" || driver.CarLicencePlate == "" || driver.Status == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("422 - Please supply all neccessary driver information "))
		} else { // All not null
			// Run db UpdateDriver function
			UpdateDriver(db, driverid, driver)
			w.WriteHeader(http.StatusAccepted)
			w.Write([]byte("202 - Account details updated"))
		}
	} else {
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("422 - Please supply driver information in JSON format"))
	}
}

// Delete driver by driverID
func DeleteDriverAccount(w http.ResponseWriter, r *http.Request) {
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

	// Get driver id to delete
	params := mux.Vars(r)
	var driverid int
	fmt.Sscan(params["driverid"], &driverid)

	var driver Driver

	// Run db DeleteDriver function
	errMsg := DeleteDriver(db, driverid)
	if errMsg == "Account does not exist" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - No account found"))
	} else {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("202 - Driver account deleted: " + driver.Username))
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

	router.HandleFunc("/api/v1/drivers", CreateDriverAccount).Methods("POST")
	router.HandleFunc("/api/v1/drivers", GetDriversQueryStringValidator).Methods("GET")
	router.HandleFunc("/api/v1/drivers/{driverid}", GetDriverDetailsByID).Methods("GET")
	router.HandleFunc("/api/v1/drivers/{driverid}", UpdateDriverDetails).Methods("PUT")
	router.HandleFunc("/api/v1/drivers/{driverid}", DeleteDriverAccount).Methods("DELETE")

	fmt.Println("Driver Service operating on port 5001")
	log.Fatal(http.ListenAndServe(":5001", router))

}
