package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"github.com/yanzay/tbot"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {

	token := os.Getenv("TELEGRAM_TOKEN")
	bot, err := tbot.NewServer(token)
	if err != nil {
		log.Fatal(err)
	}

	// Use whitelist for Auth middleware, allow to interact only with user1 and user2
	//whitelist := []string{"test0", "test1"}
	//bot.AddMiddleware(tbot.NewAuth(whitelist))

	// start your Bot
	bot.HandleFunc("/start", StartHandler)

	// Get your current location
	bot.HandleFunc("/getLocation", GetLocationHandler)

	// Set default handler if you want to process unmatched input
	bot.HandleDefault(EchoHandler)

	// Start listening for messages
	err = bot.ListenAndServe()
	log.Fatal(err)

}

func EchoHandler(message *tbot.Message) {
	//fmt.Println("Echo Handler")
	fmt.Println(message.Location)

	if message.Location.Latitude == 0 || message.Location.Longitude == 0 {
		message.RequestLocationButton("please first share your location with bot", "Current Location")
	} else {
		response, err := http.Get("https://api.openweathermap.org/data/2.5/weather?lat=" + fmt.Sprintf("%f", message.Location.Latitude) + "&lon=" + fmt.Sprintf("%f", message.Location.Longitude) + "&appid=d9231942f60fb580636c94e95f19f0e8")
		if err != nil {
			fmt.Printf("The HTTP request failed with error %s\n", err)
		} else {
			data, _ := ioutil.ReadAll(response.Body)
			//fmt.Println(string(data))

			fmt.Println(gjson.Get(string(data), "weather.0.main"))

			weatherMainCurrent := gjson.Get(string(data), "weather.0.main").Str
			weatherDescriptionCurrent := gjson.Get(string(data), "weather.0.description").Str
			//weatherIconCurrent := gjson.Get(string(data), "weather.0.icon").Str

			message.Reply("main = " + weatherMainCurrent)
			message.Reply("description = " + weatherDescriptionCurrent)
			//message.ReplyPhoto("http://openweathermap.org/img/w/"+weatherIconCurrent+".png", "weather Icon")
		}
	}

}

func StartHandler(message *tbot.Message) {
	message.RequestLocationButton("please share your location with bot", "Current Location")
}

func GetLocationHandler(message *tbot.Message) {
	//message.Reply("GetLocationHandler")
	message.RequestLocationButton("please share your location with bot", "Current Location")
}
