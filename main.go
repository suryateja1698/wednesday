package main

import (
	"fmt"

	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

type user struct {
	UserID       int    `json:"userid,omitempty"`
	Name         string `json:"name,omitempty"`
	Email        string `json:"email,omitempty"`
	Password     string `json:"password,omitempty"`
	MobileNumber string `json:"mobile_number,omitempty"`
}
type cabs struct {
	CabID               int    `json:"cabid,omitempty"`
	CabNumber           string `json:"cab_no,omitempty"`
	DriverName          string `json:"driver_name,omitempty"`
	DriverNumber        string `json:"driver_number,omitempty"`
	Status              string `json:"status,omitempty"`
	CurrentLocation     string `json:"current,omitempty"`
	WillingToGoLocation string `json:"willing,omitempty"`
}
type bookings struct {
	BookingID int    `json:"bookingid,omitempty"`
	UserID    int    `json:"user_id,omitempty"`
	CabID     int    `json:"cab_id,omitempty"`
	FromLoc   string `json:"from_loc,omitempty"`
	ToLoc     string `json:"to_loc,omitempty"`
	BookedAt  string `json:"booked_at,omitempty"`
	Status    string `json:"status,omitempty"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to HomePage!")
	fmt.Println("Endpoint Hit: HomePage")
}

func handleRequests() {
	log.Println("Server is starting")

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/allusers", allusers).Methods("GET")
	myRouter.HandleFunc("/get/{current}/{willing}", getcarsavailable).Methods("GET")
	myRouter.HandleFunc("/book/{cabid}", bookcar).Methods("GET")
	myRouter.HandleFunc("/cabs/{current}", getcabs).Methods("GET")
	myRouter.HandleFunc("/getpast/{user_id}", pastbookings).Methods("GET")

	log.Fatal(http.ListenAndServe(":1079", myRouter))
}

func main() {

	db, err = gorm.Open("mysql", "root:*****@tcp(127.0.0.1:3306)/My?charset=utf8&parseTime=True")

	if err != nil {
		log.Println("Unable to connect to database")
	} else {
		log.Println("Connection Successful")
	}

	db.AutoMigrate(&user{}, &cabs{}, &bookings{})
	handleRequests()
}

func allusers(w http.ResponseWriter, r *http.Request) {
	usr := []user{}
	db.Find(&usr)
	fmt.Println("All the users that are in the database: ", usr)
	json.NewEncoder(w).Encode(usr)

}

func getcarsavailable(w http.ResponseWriter, r *http.Request) {
	var cab []cabs
	vars := mux.Vars(r)
	key := vars["current"]
	key1 := vars["willing"]

	db.Find(&cab)
	for _, cabs := range cab {
		if cabs.CurrentLocation == key && cabs.WillingToGoLocation == key1 {
			fmt.Println("Avaialble cabs for the destination are: ", cabs)
			fmt.Println("Enter cab id to book that particular cab...")

			json.NewEncoder(w).Encode(cabs)
		}
	}
}

func bookcar(w http.ResponseWriter, r *http.Request) {

	var cab []cabs
	vars := mux.Vars(r)
	key := vars["cabid"]

	db.Find(&cab)
	for _, cabs := range cab {
		k, err := strconv.Atoi(key)
		if err == nil {
			if cabs.CabID == k {
				fmt.Println("The cab you selected is: ")

				json.NewEncoder(w).Encode(cabs)
				fmt.Println(cabs)
				fmt.Println("Have a safe journey")
			}
		}
	}
}

func getcabs(w http.ResponseWriter, r *http.Request) {
	var cab []cabs
	vars := mux.Vars(r)
	key := vars["current"]

	db.Find(&cab)
	for _, cabs := range cab {
		if cabs.CurrentLocation == key {
			fmt.Println(cabs)
			fmt.Println("Current Location:", key)
			json.NewEncoder(w).Encode(cabs)
		}
	}

}

func pastbookings(w http.ResponseWriter, r *http.Request) {
	var booking []bookings
	vars := mux.Vars(r)
	key := vars["user_id"]

	db.Find(&booking)
	for _, bookings := range booking {
		s, err := strconv.Atoi(key)
		if err == nil {
			if bookings.UserID == s {
				fmt.Println(bookings)
				fmt.Println("The bookings made by the user are:", bookings)
				json.NewEncoder(w).Encode(bookings)
			}
		}
	}
	fmt.Println("Press Ctrl+c to close")
}
