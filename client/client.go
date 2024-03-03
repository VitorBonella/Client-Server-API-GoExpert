package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main(){
	
	jsonCotacao, cotacao, err := callServer()
	if err != nil{
		fmt.Println("Error calling server")
	}

	fmt.Println(jsonCotacao)

	if err := os.WriteFile("cotacao.txt",[]byte(fmt.Sprintf("Dólar: {%s}",cotacao)),0644); err != nil{
		fmt.Println("Error writing cotacao to a file")
	}

}


func callServer() (string, string, error){

	ctx, cancel := context.WithTimeout(context.Background(), 700*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		fmt.Println("Error creating request", err)
		return "","",err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error getting response from server", err)
		return "","",err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil{
		fmt.Println("Error reading body",err)
		return "","",err
	}

	respString := string(respBytes)

	finalResp := strings.Split(respString,"\n\n BID-> ")
	if len(finalResp) != 2{
		fmt.Println("Error processing response")
		return "","",err
	}

	return finalResp[0],finalResp[1],nil

}