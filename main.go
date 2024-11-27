package main

import (
	"encoding/json"
	"fmt"
	"time"
	"strconv"
	"log"
	"net/http"
	"github.com/google/uuid"
	"math"
	"strings"
	"regexp"
)

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"` 
}

type Receipt struct {
	Retailer     string  `json:"retailer"`
	PurchaseDate string  `json:"purchaseDate"`
	PurchaseTime string  `json:"purchaseTime"`
	Items        []Item  `json:"items"`
	Total        string  `json:"total"` 
}


type ReceiptResponse struct {
	ID string `json:"id"`
}

var receiptPointsMap = make(map[string]int)

// Calculates amount of points for length of retailer name
func pointsForRetailerName(retailerName string) int {
	re := regexp.MustCompile(`[^a-zA-Z0-9]`)
	strippedName := re.ReplaceAllString(retailerName, "")
	return len(strippedName)
}

//Calculates if total on receipt is whole dollars
func pointsForReceiptTotal(total string) int{
	floatTotal, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return 0
	}
	if floatTotal == math.Floor(floatTotal) {
		return 50
	} else{
		return 0
	}
}

//Calculates if total is a multiple of .25 (also know as quarter)
func pointsForReceiptTotalMultipleOfQuarter(total string) int{
	fltTotal, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return 0
	}
	if math.Mod(fltTotal,.25) == 0.0 && fltTotal !=0.0 {
		return 25
	} else{
		return 0
	}
}

//Calculates points per two items
func pointsPerTwoItems(items []Item) int{
	length := len(items)
	return (length / 2) * 5
}

//Calculates points for short description on an item
func pointsForShortDescription(item Item) int{
	strippedShortDesc := strings.TrimSpace(item.ShortDescription)
	price, err := strconv.ParseFloat(item.Price, 64)
	if err != nil {
        fmt.Println("Error converting price to float:", err)
        return 0
    }
	if len(strippedShortDesc) % 3 == 0 && len(strippedShortDesc) != 0{
		points := math.Ceil(price * 0.2)
        return int(points)  
	} else{
		return 0
	}
}

//Get total points for all short descriptions in item array
func totalPointsForShortDescription(items []Item) int{
	res := 0 
	for _, item := range items {
		res += pointsForShortDescription(item)
	}
	return res
}

//Gets Points for when item has odd day
func pointsForOddDay(purchaseDate string)int{
	parsedDate, err := time.Parse("2006-01-02", purchaseDate)
	if err != nil {
		return 0
	}
	day := parsedDate.Day()
	if day%2 != 0{
		return 6
	 } else {
		return 0
	 }
}

//Gets points based on purchase time 
func pointsForPurchaseTime(purchaseTime string)int{
	t, err := time.Parse("15:04", purchaseTime)
	if err != nil {
		return 0
	}
	if  t.Hour() > 14 && t.Hour() < 16 || (t.Hour() == 14 && t.Minute() > 0) {
		return 10
	} else {
		return 0
	}
}


// Function that returns total points
func calculateTotalPoints(receipt Receipt) int {
	return pointsForPurchaseTime(receipt.PurchaseTime) +
		pointsForOddDay(receipt.PurchaseDate) +
		totalPointsForShortDescription(receipt.Items) +
		pointsPerTwoItems(receipt.Items) +
		pointsForReceiptTotalMultipleOfQuarter(receipt.Total) +
		pointsForRetailerName(receipt.Retailer) +
		pointsForReceiptTotal(receipt.Total)
}

// Function to generate a unique UUID that is not already in the receiptPointsMap
func uniqueReceiptUUID() string {
	for {
		newUUID := uuid.New().String()
		
		if _, exists := receiptPointsMap[newUUID]; !exists {
			return newUUID
		}
	
	}
}



func main() {

	http.HandleFunc("/receipts/process", func(w http.ResponseWriter, r *http.Request) {
		//check to see if I have a post rqruest
		if r.Method == http.MethodPost {
			
			var receipt Receipt
			decoder := json.NewDecoder(r.Body)
			err := decoder.Decode(&receipt)
			if err != nil {
				http.Error(w, "Failed to parse receipt JSON", http.StatusBadRequest)
				return
			}
			receiptID := uniqueReceiptUUID()
			totalPoints:=calculateTotalPoints(receipt)
			receiptPointsMap[receiptID] = totalPoints
			response := ReceiptResponse{
				ID: receiptID,
			}

			w.Header().Set("Content-Type", "application/json")

			w.WriteHeader(http.StatusOK)
			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				http.Error(w, "Failed to encode response JSON", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/receipts/", func(w http.ResponseWriter, r *http.Request) {
		
		if r.Method == http.MethodGet {
			parts := strings.Split(r.URL.Path, "/")
			if len(parts) < 4 || parts[len(parts)-1] != "points" {
				http.Error(w, "Invalid path or method", http.StatusBadRequest)
				return
			}

			receiptID := parts[len(parts)-2]

			points, exists := receiptPointsMap[receiptID]
			if !exists {
				http.Error(w, "Receipt not found", http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			err := json.NewEncoder(w).Encode(map[string]int{"points": points})
			if err != nil {
				http.Error(w, "Failed to encode response JSON", http.StatusInternalServerError)
			}
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}


