package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type ExchangeRate struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
	gorm.Model
}

type ExchangeRateResponse struct {
	USDBRL ExchangeRate `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/cotacao", handle)
	http.ListenAndServe(":8080", nil)
}

func handle(w http.ResponseWriter, r *http.Request) {

	cotacao, err := getCotacao()
	if err != nil {
		w.Write([]byte("Error to get cotacao"))
		return
	}
	jsonCotacao, err := json.Marshal(cotacao.USDBRL)
	if err != nil {
		w.Write([]byte("Error to get cotacao json"))
		return
	}
	w.Write([]byte(string(jsonCotacao)))

	//open db, if error do nothing
	DB, _ := openDb()


	//save to db
	ctxDB, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	tx := DB.WithContext(ctxDB)
	tx.Create(&cotacao.USDBRL)
	
	//write cotacao to client
	w.Write([]byte("\n\n BID-> " + cotacao.USDBRL.Bid))

}

func getCotacao() (*ExchangeRateResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond) //only can get with 600 ms on my pc
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		fmt.Println("Error creating request", err)
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error getting dolar API", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Check the response status
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Unexpected status code:", resp.Status)
		return nil, err
	}

	exchangeRate := ExchangeRateResponse{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&exchangeRate); err != nil {
		fmt.Println("Error to decode json", err)
		return nil, err
	}

	return &exchangeRate, nil

}

func openDb() (*gorm.DB,error){

	db, err := gorm.Open(sqlite.Open("dolar.db"), &gorm.Config{})
	if err != nil {
	  fmt.Println("failed to connect database",err)
	  return nil,err
	}
  
	// Migrate the schema
	db.AutoMigrate(&ExchangeRate{})

	return db,nil

}