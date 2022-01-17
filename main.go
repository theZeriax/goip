package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	dotenv "github.com/joho/godotenv"
	gomail "gopkg.in/gomail.v2"
)

// Styling
var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

var Bold = "\033[1m"
var Italic = "\033[3m"
var Underline = "\033[4m"
var LineThrough = "\033[9m"

// Message types
func info(message string) {
	fmt.Println(fmt.Sprintf("%s%s[INFO] %s", Bold, Green, Reset) + message)
}

func warn(message string) {
	fmt.Println(fmt.Sprintf("%s%s[WARN] %s", Bold, Yellow, Reset) + message)
}

func error(message string) {
	fmt.Println(fmt.Sprintf("%s%s[ERROR] %s", Bold, Red, Reset) + message)
}

// You can set any port you want.
var port = "3000"

type GeoIP struct {
	Ip          string  `json:"query"`
	CountryCode string  `json:"countryCode"`
	CountryName string  `json:"country"`
	RegionCode  string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zipcode     string  `json:"zip"`
	Lat         float32 `json:"lat"`
	Lon         float32 `json:"lon"`
	Isp         string  `json:"isp"`
	Timezone    string  `json:"timezone"`
}

// Function to get the visitor details, such as IP or country.
func homePage(w http.ResponseWriter, r *http.Request) {
	ip := r.RemoteAddr[:strings.LastIndex(r.RemoteAddr, ":")]
	url := "http://ip-api.com/json/" + ip

	// Fetch API
	req, _ := http.NewRequest("GET", url, nil)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var geoIP GeoIP
	json.Unmarshal(body, &geoIP)

	// Shorter variables
	country_code := geoIP.CountryCode
	country_name := geoIP.CountryName
	region_code := geoIP.RegionCode
	region_name := geoIP.RegionName
	city := geoIP.City
	zipcode := geoIP.Zipcode
	lat := geoIP.Lat
	lon := geoIP.Lon
	isp := geoIP.Isp
	timezone := geoIP.Timezone

	dotenv.Load()

	// Print to the http response
	fmt.Fprintf(w, "Your Data:\n"+
		"\n\n> System Data:"+
		"\nUser-Agent: "+r.UserAgent()+
		"\nIPv4: "+
		r.RemoteAddr[:strings.LastIndex(r.RemoteAddr, ":")]+
		"\nTime: "+time.Now().Format(time.RFC1123)+
		"\nProto: "+r.Proto+
		"\nMethod: "+r.Method+" "+r.URL.String()+

		"\n\n\n> Location Data:"+
		"\nCountry: "+country_name+
		"\nCountry Code: "+country_code+
		"\nRegion: "+region_name+
		"\nRegion Code: "+region_code+
		"\nCity: "+city+
		"\nZipcode: "+zipcode+
		"\nLatitude: "+fmt.Sprintf("%f", lat)+
		"\nLongtitude: "+fmt.Sprintf("%f", lon)+
		"\nISP: "+isp+
		"\nTimezone: "+timezone+"\n\n")

	// Print to the console
	fmt.Println(Bold + Blue + "\nEndpoint Hit (" + ip + "):" + Reset)
	fmt.Println(Red + "User-Agent: " + Reset + r.UserAgent())
	fmt.Println(Red + "IPv4: " + Reset + r.RemoteAddr[:strings.LastIndex(r.RemoteAddr, ":")])
	fmt.Println(Red + "Proto: " + Reset + r.Proto)
	fmt.Println(Red + "Method: " + Reset + r.Method + " " + r.URL.String())
	fmt.Println(Red + "Time: " + Reset + time.Now().Format(time.RFC1123))
	fmt.Println(Red + "Country: " + Reset + country_name)
	fmt.Println(Red + "Region: " + Reset + region_name)
	fmt.Println(Red + "City: " + Reset + city)
	fmt.Println(Red + "Latitude: " + Reset + fmt.Sprintf("%f", lat))
	fmt.Println(Red + "Longtitude: " + Reset + fmt.Sprintf("%f", lon))
	fmt.Println(Red + "ISP: " + Reset + isp)
	fmt.Println(Red + "Timezone: " + Reset + timezone)

	// .env variable are required! Only works with gmail.
	emailfrom := os.Getenv("EMAIL_FROM")
	emailto := os.Getenv("EMAIL_TO")
	pass := os.Getenv("PASSWORD")

	// Compose a new email
	m := gomail.NewMessage()
	m.SetHeader("From", emailfrom)
	m.SetHeader("To", emailto)
	m.SetHeader("Subject", "Found a user: "+
		(r.RemoteAddr[:strings.LastIndex(r.RemoteAddr, ":")])+
		" at "+time.Now().Format(time.RFC1123))
	m.SetBody("text/html", "<h1>Found a new visitor!</h1>"+"<br>> System Data:<br><h3>"+
		"<b>User-Agent</b>: "+r.UserAgent()+"<br>"+
		"<b>IPv4</b>: http://"+
		r.RemoteAddr[:strings.LastIndex(r.RemoteAddr, ":")]+
		"<br>"+
		"<b>Time</b>: "+time.Now().Format(time.RFC1123)+"<br>"+
		"<b>Proto</b>: "+r.Proto+"<br>"+
		"<b>Method</b>: "+r.Method+" "+r.URL.String()+
		"</h3>"+
		"<br>> Location Data:<br><h3>"+
		"<b>Country</b>: "+country_name+"<br>"+
		"<b>Country Code</b>: "+country_code+"<br>"+
		"<b>Region</b>: "+region_name+"<br>"+
		"<b>Region Code</b>: "+region_code+"<br>"+
		"<b>City</b>: "+city+"<br>"+
		"<b>Zipcode</b>: "+zipcode+"<br>"+
		"<b>Latitude</b>: "+fmt.Sprintf("%f", lat)+"<br>"+
		"<b>Longtitude</b>: "+fmt.Sprintf("%f", lon)+"<br>"+
		"<b>ISP</b>: "+isp+"<br>"+
		"<b>Timezone</b>: "+timezone+"<br>"+
		"</h3>"+
		"<img src='https://countryflagsapi.com/png/"+
		country_code+"' alt='flag-"+strings.ToLower(country_name)+"'/>")

	d := gomail.NewDialer("smtp.gmail.com", 587, emailfrom, pass)

	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

// http handler for the / endpoint
func handleRequests() {
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}

func main() {
	info("Server online on http://localhost:" + port)
	handleRequests()
}
