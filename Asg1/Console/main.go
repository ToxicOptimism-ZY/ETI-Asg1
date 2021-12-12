package main

//==================== Imports ====================
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
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

type Trip struct {
	TripID      int
	PickUp      string //char(6)
	DropOff     string //char(6)
	DriverID    int
	PassengerID int
	Status      string
}

// struct to track user's current and past decisions
type Session struct {
	Usertype         string
	BreadCrumbOption []string
	BreadCrumbMenu   []string
}

const passengerURL = "http://localhost:5000/api/v1/passengers"
const driverURL = "http://localhost:5001/api/v1/drivers"
const tripURL = "http://localhost:5002/api/v1/trips"
const keyPass = "2c78afaf-97da-4816-bbee-9ad239abb296"
const keyDriver = "2c78afaf-97da-4816-bbee-9ad239abb297"
const keyTrip = "2c78afaf-97da-4816-bbee-9ad239abb298"

//==================== Auxiliary Functions ====================
// Add option selected and new current menu to session
func breadCrumbAppend(option string, menu string, session *Session) {
	session.BreadCrumbOption = append(session.BreadCrumbOption, option)
	session.BreadCrumbMenu = append(session.BreadCrumbMenu, menu)
}

// Remove last item from session
func breadCrumbPop(session *Session) {
	session.BreadCrumbOption = session.BreadCrumbOption[:len(session.BreadCrumbOption)-1]
	session.BreadCrumbMenu = session.BreadCrumbMenu[:len(session.BreadCrumbMenu)-1]
}

// Add form function to retrieve option, and do respective back and exit actions if neccessary
func form(msg string, attribute *string, session *Session) bool {
	fmt.Println(msg)
	fmt.Scanln(attribute)
	if *attribute == "b" {
		breadCrumbPop(session)
		return true
	} else if *attribute == "0" {
		breadCrumbAppend("0", "Exit", session)
		return true
	}
	return false
}

// Recursive function to reverse the order of a list
// Not done in database query in the event to allow for
// microservice users wanting a chronological order version
func reverseTrip(input []Trip) []Trip {
	if len(input) == 0 {
		return input
	}
	return append(reverseTrip(input[1:]), input[0]) //Remove and append to the end
}

//==================== API Callers ====================

//==================== Passengers API Callers ====================

// Create passenger account
func CreatePassengerAccount(passenger Passenger) string {
	// Set up url
	url := passengerURL + "?key=" + keyPass

	// Convert to Json
	jsonValue, _ := json.Marshal(passenger)

	// Post with object
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))

	var errMsg string

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// Get fail or success msg
		if response.StatusCode == 401 {
			errMsg = string(data)
		} else if response.StatusCode == 422 {
			errMsg = string(data)
		} else {
			errMsg = "Success"
		}
	}

	response.Body.Close()

	return errMsg
}

// Log in to passenger account
func GetPassengerDetails(username string, password string) (string, Passenger) {

	// Set up url
	url := passengerURL + "?key=" + keyPass + "&username=" + username + "&password=" + password

	// Get method
	response, err := http.Get(url)

	var passenger Passenger
	var errMsg string

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// Get fail or success msg
		if response.StatusCode == 401 {
			errMsg = string(data)
		} else if response.StatusCode == 404 {
			errMsg = string(data)
		} else {
			errMsg = "Success"
			json.Unmarshal([]byte(data), &passenger) // Convert json to passenger details
		}
	}

	response.Body.Close()

	return errMsg, passenger
}

// Get passenger info by passenger ID
func GetPassengerDetailsByID(passengerID int) (string, Passenger) {

	// Set up URL
	url := passengerURL + "/" + strconv.Itoa(passengerID) + "?key=" + keyPass

	// Get method
	response, err := http.Get(url)

	var passenger Passenger
	var errMsg string

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// Get fail or success msg
		if response.StatusCode == 401 {
			errMsg = string(data)
		} else if response.StatusCode == 404 {
			errMsg = string(data)
		} else {
			errMsg = "Success"
			json.Unmarshal([]byte(data), &passenger) // Convert json to passenger details
		}
	}

	response.Body.Close()

	return errMsg, passenger
}

// Update passenger info by passenger id
func UpdatePassengerDetails(passengerID int, p Passenger) string {

	// Set up url
	url := passengerURL + "/" + strconv.Itoa(passengerID) + "?key=" + keyPass

	// Convert passenger object to json format
	jsonValue, _ := json.Marshal(p)

	// Put method
	request, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonValue))

	request.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)

	var errMsg string

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// Get fail or success msg
		if response.StatusCode == 401 {
			errMsg = string(data)
		} else if response.StatusCode == 422 {
			errMsg = string(data)
		} else {
			errMsg = "Success"
		}
	}

	response.Body.Close()

	return errMsg
}

// Delete passenger by passenger id
func DeletePassengerAccount(passengerID int, userType string) string {

	// Set up url
	url := passengerURL + "/" + strconv.Itoa(passengerID) + "?key=" + keyPass + "&useraccess=" + userType

	// Delete passenger
	request, _ := http.NewRequest(http.MethodDelete, url, nil)

	client := &http.Client{}
	response, err := client.Do(request)

	var errMsg string

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// Get fail or success msg
		if response.StatusCode == 401 {
			errMsg = string(data)
		} else if response.StatusCode == 404 {
			errMsg = string(data)
		} else {
			errMsg = "Success"
		}
	}

	response.Body.Close()

	return errMsg
}

//==================== Drivers API Callers ====================
// Create driver account
func CreateDriverAccount(driver Driver) string {

	// Set up url
	url := driverURL + "?key=" + keyDriver

	// Convert driver into json
	jsonValue, _ := json.Marshal(driver)

	// Post method
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))

	var errMsg string

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// Get fail or success msg
		if response.StatusCode == 401 {
			errMsg = string(data)
		} else if response.StatusCode == 422 {
			errMsg = string(data)
		} else {
			errMsg = "Success"
		}
	}
	response.Body.Close()

	return errMsg
}

// Log in to driver account
func GetDriverDetails(username string, password string) (string, Driver) {
	// Set up URL
	url := driverURL + "?key=" + keyDriver + "&username=" + username + "&password=" + password

	// Get method
	response, err := http.Get(url)

	var driver Driver
	var errMsg string

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// Get fail or success msg
		if response.StatusCode == 401 {
			errMsg = string(data)
		} else if response.StatusCode == 404 {
			errMsg = string(data)
		} else {
			errMsg = "Success"
			json.Unmarshal([]byte(data), &driver) // Convert json to driver details
		}
	}
	response.Body.Close()

	return errMsg, driver
}

// Get list of drivers with a certain status
func GetDriverByStatus(status string) (string, []Driver) {

	// Set up url
	url := driverURL + "?key=" + keyDriver + "&status=" + status
	// Get method
	response, err := http.Get(url)

	var drivers []Driver
	var errMsg string

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// Get success or fail msg
		if response.StatusCode == 401 {
			errMsg = string(data)
		} else if response.StatusCode == 404 {
			errMsg = string(data)
		} else {
			errMsg = "Success"
			json.Unmarshal([]byte(data), &drivers) // Convert json to list of drivers
		}
	}

	response.Body.Close()

	return errMsg, drivers
}

// Get driver details by driver ID
func GetDriverDetailsByID(driverID int) (string, Driver) {

	// Set up URL
	url := driverURL + "/" + strconv.Itoa(driverID) + "?key=" + keyDriver

	// Get method
	response, err := http.Get(url)

	var driver Driver
	var errMsg string

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// Get fail or success msg
		if response.StatusCode == 401 {
			errMsg = string(data)
		} else if response.StatusCode == 404 {
			errMsg = string(data)
		} else {
			errMsg = "Success"
			json.Unmarshal([]byte(data), &driver) // Convert json to driver detail
		}
	}

	response.Body.Close()

	return errMsg, driver

}

// Update driver details by driver ID
func UpdateDriverDetails(driverID int, d Driver) string {

	// Set up url
	url := driverURL + "/" + strconv.Itoa(driverID) + "?key=" + keyDriver

	// Convert driver object to json
	jsonValue, _ := json.Marshal(d)

	// Put method
	request, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)

	var errMsg string

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// Get success or fail msg
		if response.StatusCode == 401 {
			errMsg = string(data)
		} else if response.StatusCode == 422 {
			errMsg = string(data)
		} else {
			errMsg = "Success"
		}
	}

	response.Body.Close()

	return errMsg
}

// Delete driver by driver ID
func DeleteDriverAccount(driverID int, userType string) string {

	// Set up url
	url := driverURL + "/" + strconv.Itoa(driverID) + "?key=" + keyDriver + "&useraccess=" + userType

	// Delete method
	request, _ := http.NewRequest(http.MethodDelete, url, nil)
	client := &http.Client{}
	response, err := client.Do(request)

	var errMsg string

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// Get sucess or fail msg
		if response.StatusCode == 401 {
			errMsg = string(data)
		} else if response.StatusCode == 404 {
			errMsg = string(data)
		} else {
			errMsg = "Success"
		}
	}

	response.Body.Close()

	return errMsg
}

//==================== Trips API Callers ====================
// Create trip record
func CreateTripRecord(trip Trip) string {
	//Set up url
	url := tripURL + "?key=" + keyTrip

	// Convert trip to json
	jsonValue, _ := json.Marshal(trip)

	// Post method
	response, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))

	var errMsg string

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// Get fail or success msg
		if response.StatusCode == 401 {
			errMsg = string(data)
		} else if response.StatusCode == 422 {
			errMsg = string(data)
		} else {
			errMsg = "Success"
		}
	}

	response.Body.Close()

	return errMsg
}

// Get all trips made by passenger id
func GetTripDetailsByPassengerID(passengerID int) (string, []Trip) {

	// Set up url
	url := tripURL + "?key=" + keyTrip + "&passengerid=" + strconv.Itoa(passengerID)

	// Get method
	response, err := http.Get(url)

	var trips []Trip
	var errMsg string

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// Get success or fail msg
		if response.StatusCode == 401 {
			errMsg = string(data)
		} else if response.StatusCode == 404 {
			errMsg = string(data)
		} else {
			errMsg = "Success"
			json.Unmarshal([]byte(data), &trips) //Convert json to array of trips
		}
	}

	response.Body.Close()

	return errMsg, trips
}

// Get current trip in progress by passenger id
func GetCurrentTripDetailsForPassenger(passengerID int, status string) (string, Trip) {

	// Set up url
	url := tripURL + "?key=" + keyTrip + "&passengerid=" + strconv.Itoa(passengerID) + "&status=" + status

	// Get method
	response, err := http.Get(url)

	var trip Trip
	var errMsg string

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// Get success or fail msg
		if response.StatusCode == 401 {
			errMsg = string(data)
		} else if response.StatusCode == 404 {
			errMsg = string(data)
		} else {
			errMsg = "Success"
			json.Unmarshal([]byte(data), &trip) // Convert json to trip details
		}
	}

	response.Body.Close()

	return errMsg, trip
}

// Get current trip in progress by driver id
func GetCurrentTripDetailsForDriver(driverID int, status string) (string, Trip) {

	// Set up url
	url := tripURL + "?key=" + keyTrip + "&driverid=" + strconv.Itoa(driverID) + "&status=" + status

	// Get method
	response, err := http.Get(url)

	var trip Trip
	var errMsg string

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// Get success or fail msg
		if response.StatusCode == 401 {
			errMsg = string(data)
		} else if response.StatusCode == 404 {
			errMsg = string(data)
		} else {
			errMsg = "Success"
			json.Unmarshal([]byte(data), &trip) // Convert json to trip details
		}
	}

	response.Body.Close()

	return errMsg, trip
}

// Get trip by trip ID
func GetTripDetailsByID(tripID int) (string, Trip) {

	// Set up url
	url := tripURL + "/" + strconv.Itoa(tripID) + "?key=" + keyTrip

	// Get method
	response, err := http.Get(url)

	var trip Trip
	var errMsg string

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// Get success or fail msg
		if response.StatusCode == 401 {
			errMsg = string(data)
		} else if response.StatusCode == 404 {
			errMsg = string(data)
		} else {
			errMsg = "Success"
			json.Unmarshal([]byte(data), &trip) // Convert json to trip object
		}
	}

	response.Body.Close()

	return errMsg, trip
}

// Update trip details by trip ID
func UpdateTripDetails(tripID int, t Trip) string {

	// Set up url
	url := tripURL + "/" + strconv.Itoa(tripID) + "?key=" + keyTrip

	// Convert trip to json
	jsonValue, _ := json.Marshal(t)

	// Put method
	request, _ := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(jsonValue))
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(request)

	var errMsg string

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// Get success or fail msg
		if response.StatusCode == 401 {
			errMsg = string(data)
		} else if response.StatusCode == 422 {
			errMsg = string(data)
		} else {
			errMsg = "Success"
		}
	}

	response.Body.Close()

	return errMsg
}

// Delete trip by TripID
func DeleteTripRecord(tripID int, userType string) string {
	// Set up url
	url := tripURL + "/" + strconv.Itoa(tripID) + "?key=" + keyTrip + "&useraccess=" + userType

	// Delete method
	request, _ := http.NewRequest(http.MethodDelete, url, nil)
	client := &http.Client{}
	response, err := client.Do(request)

	var errMsg string

	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// Get success or fail msg
		if response.StatusCode == 401 {
			errMsg = string(data)
		} else if response.StatusCode == 404 {
			errMsg = string(data)
		} else {
			errMsg = "Success"
		}
	}

	response.Body.Close()

	return errMsg
}

//==================== Authentication Menus ====================

// 1st menu to select user
func selectUserMenu(session *Session) {

	var option string

	// Menu
	fmt.Println("------ Pick a user type ------")
	fmt.Println("[1] Passenger")
	fmt.Println("[2] Driver")
	fmt.Println("[0] Exit application")
	fmt.Println("------------------------------")
	fmt.Println("Enter your option: ")
	fmt.Scanln(&option)

	// Options
	switch option {
	case "1":
		breadCrumbAppend(option, "EntryMenu", session)
		session.Usertype = "Passenger"
	case "2":
		breadCrumbAppend(option, "EntryMenu", session)
		session.Usertype = "Driver"
	case "0":
		breadCrumbAppend(option, "Exit", session)
	default:
	}
}

// 2nd menu to login/ sign up as passenger/ driver
func entryMenu(session *Session) {

	var option string

	// Menu
	fmt.Println("------ Welcome " + session.Usertype + "! ------")
	fmt.Println("[1] Login")
	fmt.Println("[2] Sign Up")
	fmt.Println("[0] Exit application")
	fmt.Println("----------------------")
	if exitEarly := form("Enter your option (b to back): ", &option, session); exitEarly {
		return
	}

	// Options - Head to respective menu
	if option == "1" {
		if session.Usertype == "Passenger" {
			breadCrumbAppend(option, "LoginPassengerMenu", session)
		} else if session.Usertype == "Driver" {
			breadCrumbAppend(option, "LoginDriverMenu", session)
		}
	} else if option == "2" {
		if session.Usertype == "Passenger" {
			breadCrumbAppend(option, "SignUpPassengerMenu", session)
		} else if session.Usertype == "Driver" {
			breadCrumbAppend(option, "SignUpDriverMenu", session)
		}
	}

}

// Login menu to get user info
func loginMenu(session *Session) (string, string) {
	var username string
	var password string

	// Menu
	fmt.Println("------ Log in ------")
	fmt.Println("Please fill in the following details (b to back, 0 to exit).")
	if exitEarly := form("Username: ", &username, session); exitEarly {
		return "b", "b"
	}
	if exitEarly := form("Password: ", &password, session); exitEarly {
		return "b", "b"
	}
	fmt.Println("--------------------")

	return username, password //return filled up options for appropriate passenger or driver log in
}

// Login for passengers
func loginPassenger(session *Session) Passenger {

	var p Passenger
	var throwAway Passenger

	username, password := loginMenu(session) // Call login menu to get user info

	// If exit early
	if username == "b" {
		return throwAway
	}

	// Call Api caller function to get passenger object
	errMsg, p := GetPassengerDetails(username, password)
	if errMsg != "Success" {
		fmt.Println(errMsg)
	} else {
		breadCrumbAppend(errMsg, "PassengerMenu", session) //If logged in, head to passenger
		return p
	}

	return throwAway
}

// Login for drivers
func loginDriver(session *Session) Driver {

	var d Driver
	var throwAway Driver

	username, password := loginMenu(session) // Call login menu to get user info

	// If exit early
	if username == "b" {
		return throwAway
	}

	// Call api caller to get driver object
	errMsg, d := GetDriverDetails(username, password)
	if errMsg != "Success" {
		fmt.Println(errMsg)
	} else {
		breadCrumbAppend(errMsg, "DriverMenu", session) // If logged in, head to driver
		return d
	}

	return throwAway
}

// Sign up for passengers
func signUpPassengerMenu(session *Session) Passenger {

	var p Passenger
	var throwAway Passenger // Store changed attribute

	// Form
	fmt.Println("------ Sign up ------")
	fmt.Println("Please fill in the following details (b to back, 0 to exit).")
	if exitEarly := form("Username: ", &p.Username, session); exitEarly {
		return throwAway
	}
	if exitEarly := form("Password: ", &p.Password, session); exitEarly {
		return throwAway
	}
	if exitEarly := form("First Name: ", &p.FirstName, session); exitEarly {
		return throwAway
	}
	if exitEarly := form("Last Name: ", &p.LastName, session); exitEarly {
		return throwAway
	}
	if exitEarly := form("Mobile No: ", &p.MobileNo, session); exitEarly {
		return throwAway
	}
	if exitEarly := form("Email Address: ", &p.EmailAddress, session); exitEarly {
		return throwAway
	}
	fmt.Println("--------------------")

	// Call api caller to create a new passenger object
	errMsg := CreatePassengerAccount(p)
	if errMsg != "Success" {
		fmt.Println(errMsg)
	} else {
		//Retrieve auto generated passenger id
		_, p := GetPassengerDetails(p.Username, p.Password)

		// Remove sign up breadcrumb, replace with login
		breadCrumbPop(session)
		breadCrumbAppend("2", "LoginPassengerMenu", session)
		breadCrumbAppend(errMsg, "PassengerMenu", session)

		return p
	}

	return throwAway
}

// Sign up for drivers
func signUpDriverMenu(session *Session) Driver {
	var d Driver
	var throwAway Driver

	// Form
	fmt.Println("------ Sign up ------")
	fmt.Println("Please fill in the following details (b to back, 0 to exit).")
	if exitEarly := form("Username: ", &d.Username, session); exitEarly {
		return throwAway
	}
	if exitEarly := form("Password: ", &d.Password, session); exitEarly {
		return throwAway
	}
	if exitEarly := form("First Name: ", &d.FirstName, session); exitEarly {
		return throwAway
	}
	if exitEarly := form("Last Name: ", &d.LastName, session); exitEarly {
		return throwAway
	}
	if exitEarly := form("Mobile No: ", &d.MobileNo, session); exitEarly {
		return throwAway
	}
	if exitEarly := form("Email Address: ", &d.EmailAddress, session); exitEarly {
		return throwAway
	}
	if exitEarly := form("NRIC: ", &d.NRIC, session); exitEarly {
		return throwAway
	}
	if exitEarly := form("Car Licence No: ", &d.CarLicencePlate, session); exitEarly {
		return throwAway
	}
	fmt.Println("--------------------")

	d.Status = "Available"

	// Call api caller to create new driver object
	errMsg := CreateDriverAccount(d)
	if errMsg != "Success" {
		fmt.Println(errMsg)
	} else {
		//Retrieve auto generated driver id
		_, d := GetDriverDetails(d.Username, d.Password)

		// Remove sign up breadcrumb, replace with login
		breadCrumbPop(session)
		breadCrumbAppend("2", "LoginDriverMenu", session)
		breadCrumbAppend(errMsg, "DriverMenu", session)

		return d
	}

	return throwAway
}

//==================== Passenger Menus ====================

// Passenger Menu
func passengerMenu(session *Session, p Passenger) {

	var option string
	var errMsg string
	var trip Trip
	currentTrip := false

	// Get current trip for passengers if any
	if errMsg, trip = GetCurrentTripDetailsForPassenger(p.PassengerID, "Waiting"); errMsg == "Success" {
		currentTrip = true
	} else if errMsg, trip = GetCurrentTripDetailsForPassenger(p.PassengerID, "Travelling"); errMsg == "Success" {
		currentTrip = true
	}

	// Menu
	fmt.Println("------ Welcome " + p.Username + " ------")
	fmt.Println("[1] Update My Details")
	fmt.Println("[2] Get Trip History")
	// Print appropriate action based on current trip status
	if !currentTrip {
		fmt.Println("[3] Start New Trip")
	} else {
		fmt.Println("")
		fmt.Println("Current Trip: ")
		tripLister(trip)
	}
	fmt.Println("[0] Exit application")
	fmt.Println("----------------------")
	if exitEarly := form("Enter your option (b to back): ", &option, session); exitEarly {
		return
	}

	// Options - Head to appropriate action menu
	if option == "1" {
		breadCrumbAppend(option, "UpdatePassengerDetailsMenu", session)
	} else if option == "2" {
		breadCrumbAppend(option, "ListTripsMenu", session)
	} else if option == "3" && !currentTrip {
		breadCrumbAppend(option, "CreateTripMenu", session)
	}
}

// Update passenger details menu
func updatePassengerDetailsMenu(session *Session, p Passenger) Passenger {

	var option string
	throwAway := p

	// Menu
	fmt.Println("------ Update details for " + p.Username + " ------")
	fmt.Println("[1] Username: " + p.Username)
	fmt.Println("[2] Password: ")
	fmt.Println("[3] First Name: " + p.FirstName)
	fmt.Println("[4] Last Name: " + p.LastName)
	fmt.Println("[5] Mobile No: " + p.MobileNo)
	fmt.Println("[6] Email Address: " + p.EmailAddress)
	fmt.Println("----------------------")
	if exitEarly := form("Enter your option (b to back): ", &option, session); exitEarly {
		return p
	}

	// Option to update + new data
	if option == "1" {
		if exitEarly := form("New Username (b to back): ", &throwAway.Username, session); exitEarly {
			return p
		}
	} else if option == "2" {
		if exitEarly := form("New Password (b to back): ", &throwAway.Password, session); exitEarly {
			return p
		}
	} else if option == "3" {
		if exitEarly := form("New First Name (b to back): ", &throwAway.FirstName, session); exitEarly {
			return p
		}
	} else if option == "4" {
		if exitEarly := form("New Last Name (b to back): ", &throwAway.LastName, session); exitEarly {
			return p
		}
	} else if option == "5" {
		if exitEarly := form("New Mobile No (b to back): ", &throwAway.MobileNo, session); exitEarly {
			return p
		}
	} else if option == "6" {
		if exitEarly := form("New Email Address (b to back): ", &throwAway.EmailAddress, session); exitEarly {
			return p
		}
	}

	// Call api caller for updating passenger details
	errMsg := UpdatePassengerDetails(p.PassengerID, p)
	if errMsg != "Success" {
		fmt.Println(errMsg)
	} else {
		p = throwAway // Update actual
	}

	return p
}

// List trips in a reversed chronological order
func listTripsMenu(session *Session, p Passenger) {

	// Menu
	var option string
	fmt.Println("------ Trip History ------")
	var trips []Trip

	// Call api caller to get history of trips made by passengerID
	errMsg, trips := GetTripDetailsByPassengerID(p.PassengerID)

	if errMsg != "Success" && errMsg[:3] != "404" {
		fmt.Println(errMsg)
	} else if errMsg[:3] == "404" {
		fmt.Println("No trips were created yet.")
	} else {
		tripsLister(reverseTrip(trips)) // List reverse trips
	}
	fmt.Println("--------------------------")

	form("Enter your option (b to back): ", &option, session)
}

// Create trip menu
func createTripMenu(session *Session, p Passenger) {

	var trip Trip

	// Form
	fmt.Println("------ Create Trip ------")
	fmt.Println("Please fill in the following details (b to back, 0 to exit).")
	if exitEarly := form("Pick Up Postal Code: ", &trip.PickUp, session); exitEarly {
		return
	}
	if exitEarly := form("Drop Off Postal Code: ", &trip.DropOff, session); exitEarly {
		return
	}

	// Default variables
	trip.PassengerID = p.PassengerID
	trip.Status = "Waiting"

	// Call Api caller to get list of available drivers
	errMsg, drivers := GetDriverByStatus("Available")
	if errMsg != "Success" {
		fmt.Println(errMsg)
	} else {
		// Assign first available driver
		driver := drivers[0]

		fmt.Print(trip)
		fmt.Print(driver)

		// Update driver status by using the API caller
		driver.Status = "Unavailable"
		errMsg := UpdateDriverDetails(driver.DriverID, driver)

		if errMsg != "Success" {
			fmt.Println(errMsg)
		} else {
			trip.DriverID = driver.DriverID
			// Call api caller to create trip
			errMsg := CreateTripRecord(trip)
			if errMsg != "Success" {
				fmt.Println(errMsg)
			} else {
				breadCrumbPop(session) // Complete action
			}
		}

	}
	fmt.Println("--------------------")

}

//==================== Driver Menus ====================

// Driver menu
func driverMenu(session *Session, d Driver) {

	var option string
	var errMsg string
	var trip Trip

	// Get current trip assigned for the driver
	currentTrip := false
	if errMsg, trip = GetCurrentTripDetailsForDriver(d.DriverID, "Waiting"); errMsg == "Success" {
		currentTrip = true
	} else if errMsg, trip = GetCurrentTripDetailsForDriver(d.DriverID, "Travelling"); errMsg == "Success" {
		currentTrip = true
	}

	// Menu
	fmt.Println("------ Welcome " + d.Username + " ------")
	fmt.Println("[1] Update My Details")

	// Based on current trip's status, show appropriate actions
	if !currentTrip {
		fmt.Println("[2] No Trips Currently - Refresh Page")
	} else {
		fmt.Println("")
		fmt.Println("Current Trip: ")
		tripLister(trip)

		if trip.Status == "Waiting" { // If yet to start
			fmt.Println("[2] Start Trip")
		} else if trip.Status == "Travelling" { // If yet to end
			fmt.Println("[2] End Trip")
		}
	}

	fmt.Println("[0] Exit application")
	fmt.Println("----------------------")
	if exitEarly := form("Enter your option (b to back): ", &option, session); exitEarly {
		return
	}

	// Options
	if option == "1" {
		breadCrumbAppend(option, "UpdateDriverDetailsMenu", session)
	} else if option == "2" && !currentTrip {
		return
	} else if option == "2" && trip.Status == "Waiting" {
		// Update trip with new status based on api caller
		trip.Status = "Travelling"
		UpdateTripDetails(trip.TripID, trip)

	} else if option == "2" && trip.Status == "Travelling" {
		// Update trip with new status with api caller
		trip.Status = "Completed"
		UpdateTripDetails(trip.TripID, trip)

		// Update driver with new status with api caller
		d.Status = "Available"
		UpdateDriverDetails(d.DriverID, d)
	}
}

// Update driver details
func updateDriverDetailsMenu(session *Session, d Driver) Driver {

	var option string
	throwAway := d // Store changed attribute

	// Menu
	fmt.Println("------ Update details for " + d.Username + " ------")
	fmt.Println("[1] Username: " + d.Username)
	fmt.Println("[2] Password: ")
	fmt.Println("[3] First Name: " + d.FirstName)
	fmt.Println("[4] Last Name: " + d.LastName)
	fmt.Println("[5] Mobile No: " + d.MobileNo)
	fmt.Println("[6] Email Address: " + d.EmailAddress)
	fmt.Println("[7] Car Licence Plate: " + d.CarLicencePlate)
	fmt.Println("----------------------")
	if exitEarly := form("Enter your option (b to back): ", &option, session); exitEarly {
		return d
	}

	// Options
	if option == "1" {
		if exitEarly := form("New Username (b to back): ", &throwAway.Username, session); exitEarly {
			return d
		}
	} else if option == "2" {
		if exitEarly := form("New Password (b to back): ", &throwAway.Password, session); exitEarly {
			return d
		}
	} else if option == "3" {
		if exitEarly := form("New First Name (b to back): ", &throwAway.FirstName, session); exitEarly {
			return d
		}
	} else if option == "4" {
		if exitEarly := form("New Last Name (b to back): ", &throwAway.LastName, session); exitEarly {
			return d
		}
	} else if option == "5" {
		if exitEarly := form("New Mobile No (b to back): ", &throwAway.MobileNo, session); exitEarly {
			return d
		}
	} else if option == "6" {
		if exitEarly := form("New Email Address (b to back): ", &throwAway.EmailAddress, session); exitEarly {
			return d
		}
	} else if option == "7" {
		if exitEarly := form("New Licence Plate Number (b to back): ", &throwAway.CarLicencePlate, session); exitEarly {
			return d
		}
	}

	// Use api caller to get updated status so it remains unaffected
	var throwAway1 Driver

	_, throwAway1 = GetDriverDetailsByID(d.DriverID)
	d.Status = throwAway1.Status
	throwAway.Status = throwAway1.Status

	// Use api caller to update driver
	errMsg := UpdateDriverDetails(d.DriverID, d)
	if errMsg != "Success" {
		fmt.Println(errMsg)
	} else {
		d = throwAway // Update actual
	}

	return d
}

//==================== Trip Display Menus ====================

// List all trips in array
func tripsLister(trips []Trip) {
	for _, trip := range trips {
		tripLister(trip)
	}
}

// Display all relevant details in trip
func tripLister(trip Trip) {
	fmt.Println(trip.TripID)
	fmt.Println("From " + trip.PickUp + " to " + trip.DropOff)

	// Use api caller to get driver details
	var driver Driver
	_, driver = GetDriverDetailsByID(trip.DriverID)
	fmt.Println("Driver: " + driver.Username + "#" + strconv.Itoa(driver.DriverID))

	// Use api caller to get passenger details
	var passenger Passenger
	_, passenger = GetPassengerDetailsByID(trip.PassengerID)
	fmt.Println("Passenger: " + passenger.Username + "#" + strconv.Itoa(passenger.PassengerID))

	fmt.Println("Current Status: " + trip.Status)
	fmt.Println("")
}

//==================== Main ====================

func main() {
	var session Session
	var passenger Passenger
	var driver Driver

	var throwAwayP Passenger
	var throwAwayD Driver

	breadCrumbAppend("Start", "SelectUserMenu", &session)

	// While not exit
	for {
		currentCrumb := session.BreadCrumbMenu[len(session.BreadCrumbMenu)-1]
		fmt.Println("")
		if currentCrumb == "Exit" { // If exit
			break

		} else if currentCrumb == "SelectUserMenu" {
			selectUserMenu(&session)

		} else if currentCrumb == "EntryMenu" {
			entryMenu(&session)

		} else if currentCrumb == "LoginPassengerMenu" {
			// Ensure new passenger details isnt a failed login attempt
			throwAwayP = loginPassenger(&session)
			if throwAwayP.PassengerID != 0 {
				passenger = throwAwayP
			}

		} else if currentCrumb == "LoginDriverMenu" {
			// Ensure new driver details isnt a failed login attempt
			throwAwayD = loginDriver(&session)
			if throwAwayD.DriverID != 0 {
				driver = throwAwayD
			}

		} else if currentCrumb == "SignUpPassengerMenu" {
			// Ensure new passenger details isnt a failed sign up attempt
			throwAwayP = signUpPassengerMenu(&session)
			if throwAwayP.PassengerID != 0 {
				passenger = throwAwayP
			}

		} else if currentCrumb == "SignUpDriverMenu" {
			// Ensure new driver details isnt a failed sign up attempt
			throwAwayD = signUpDriverMenu(&session)
			if throwAwayD.DriverID != 0 {
				driver = throwAwayD
			}

		} else if currentCrumb == "PassengerMenu" {
			passengerMenu(&session, passenger)
		} else if currentCrumb == "DriverMenu" {
			driverMenu(&session, driver)
		} else if currentCrumb == "UpdatePassengerDetailsMenu" {
			passenger = updatePassengerDetailsMenu(&session, passenger)
		} else if currentCrumb == "ListTripsMenu" {
			listTripsMenu(&session, passenger)
		} else if currentCrumb == "CreateTripMenu" {
			createTripMenu(&session, passenger)
		} else if currentCrumb == "UpdateDriverDetailsMenu" {
			driver = updateDriverDetailsMenu(&session, driver)
		}

	}
}
