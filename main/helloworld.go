package main

import (
	"encoding/json"
	"fmt"
	"gogogo"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Person struct {
	_id   string `bson:"_id"`
	Name  string
	Phone string
}

type DB_Store struct {
	StoreId   string `bson:"_id"`
	StoreUrl  string `bson:"url"`
	StoreName string `bson:"name"`
}

type DB_Menu struct {
	StoreId string                 `bson:"_id"`
	Menu    map[string]interface{} `bson:"menu"`
}

type DB_MainOrder struct {
	StoreId   string `bson:"store_id"`
	StoreName string `bson:store_name`
}

type LatLng struct {
	Lat float64
	Lng float64
}

type Store struct {
	Name        string `json:"name"`
	Id          string `json:"id"`
	Arrive_time string `json:"arrive_time"`
}
type StoreArray []Store

type SearchStoreResults struct {
	Stores StoreArray `json:"stores"`
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/searchStore", searchStore).Methods("POST")
	router.HandleFunc("/searchMenu", searchMenu).Methods("POST")
	router.HandleFunc("/hostOrder", hostOrder).Methods("POST")
	router.HandleFunc("/getHostOrder", getHostOrder).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))

}
func getHostOrder(w http.ResponseWriter, r *http.Request) {
	db_mainOrder := findHostOrderFromDB()
	var resultMap map[string]interface{}
	resultMap = make(map[string]interface{})
	resultMap["result"] = true
	resultMap["order_list"] = db_mainOrder

	output, output_err := json.Marshal(resultMap)

	if output_err != nil {
		fmt.Fprint(w, "[]")
		fmt.Println(output_err.Error())
		return
	}

	fmt.Fprint(w, string(output))

}
func hostOrder(w http.ResponseWriter, r *http.Request) {
	request, request_err := ioutil.ReadAll(r.Body)

	if request_err != nil {
		fmt.Println(request_err)
		fmt.Fprint(w, "[]")
		return
	}

	var request_map map[string]interface{}

	json.Unmarshal(request, &request_map)

	var storeId string = ""

	if request_map != nil && request_map["id"].(string) != "" {
		storeId = request_map["id"].(string)
	} else {
		fmt.Fprint(w, "{}")
		return
	}
	fmt.Print("sasa")
	var store DB_Store = findStoreFromDB(storeId)

	var Main_menu map[string]interface{}

	fmt.Print(store.StoreUrl)

	if err := json.Unmarshal([]byte(mSearchMenu(store.StoreUrl)), &Main_menu); err != nil {
		panic(err)
	}

	// var menus interface{}

	// menus = Main_menu["menus"].(string)

	// var toppings interface{}

	// toppings = Main_menu["topping"].(string)

	db_menu := DB_Menu{}

	db_menu.StoreId = store.StoreId

	db_menu.Menu = Main_menu

	saveOneDataToDB(db_menu, "Menu")

	db_mainorder := DB_MainOrder{}

	db_mainorder.StoreId = store.StoreId

	db_mainorder.StoreName = store.StoreName

	saveOneDataToDB(db_mainorder, "MainOrder")

	fmt.Fprint(w, "{\"result\":true}")

}
func searchStore(w http.ResponseWriter, r *http.Request) {
	request, request_err := ioutil.ReadAll(r.Body)

	if request_err != nil {
		fmt.Println(request_err)
		fmt.Fprint(w, "[]")
		return
	}

	var request_map map[string]interface{}

	json.Unmarshal(request, &request_map)

	var address string = "高雄市,自強三路,5號"

	if request_map != nil && request_map["address"].(string) != "" {
		address = request_map["address"].(string)
	}

	latlng := addressConvertTolatlng(address)

	if latlng.Lat == 0 || latlng.Lng == 0 {
		fmt.Println("latlng null")
		fmt.Fprint(w, "[]")
		return
	}

	fpStoreUrl := "https://www.foodpanda.com.tw/restaurants/lat/" + fmt.Sprintf("%f", latlng.Lat) + "/lng/" + fmt.Sprintf("%f", latlng.Lng)

	fpStoreUrlRespBody, respBodyErr := goquery.NewDocument(fpStoreUrl)

	if respBodyErr != nil {
		fmt.Print(respBodyErr.Error())
	}
	// goquery_resp := fp_storeUrl_respBody.Find(".hreview-aggregate.url")

	var results StoreArray
	var db_storeArray []interface{}
	fpStoreUrlRespBody.Find(".hreview-aggregate.url").Each(func(i int, s *goquery.Selection) {

		store := Store{}
		store.Name = s.Find(".name.fn").Text()
		store.Id = s.Find(".vendor-picture.b-lazy").AttrOr("data-vendor-id", "-1")
		store.Arrive_time = strings.TrimSpace(s.Find(".badge-info").Get(0).FirstChild.Data)

		db_store := DB_Store{}
		db_store.StoreId = s.Find(".vendor-picture.b-lazy").AttrOr("data-vendor-id", "-1")
		db_store.StoreUrl = s.Get(0).Attr[0].Val
		db_store.StoreName = s.Find(".name.fn").Text()

		fmt.Print(db_store.StoreUrl)
		results = append(results, store)
		db_storeArray = append(db_storeArray, db_store)

	})

	saveStoreToDB(db_storeArray)

	output, output_err := json.Marshal(results)

	if output_err != nil {
		fmt.Fprint(w, "[]")
		fmt.Println(output_err.Error())
		return
	}
	// fmt.Println(i, s.Get(0).Attr[0].Val)

	fmt.Fprint(w, string(output))
}

func addressConvertTolatlng(address string) LatLng {

	url := "https://maps.googleapis.com/maps/api/geocode/json?address=" + address + "&key=AIzaSyCMVJRvRF1-B1_ByOfiM-AZERa10pQYcv4"
	resp, _ := http.Get(url)
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	var raw map[string]interface{}

	var final map[string]interface{}

	defer resp.Body.Close()

	json.Unmarshal(bodyBytes, &raw)
	// raw["count"] = 1
	// faf, _ := raw["status"]

	// var detailmap map[string]interface{} = raw["results"]

	// out, _ := json.Marshal(raw)

	var ResultStatus string = raw["status"].(string)

	if ResultStatus != "OK" {
		return LatLng{}
	}

	var addResults []interface{} = raw["results"].([]interface{})

	for _, v := range addResults {

		if reflect.TypeOf(v).String() == "map[string]interface {}" {
			final = v.(map[string]interface{})

			var geometry map[string]interface{} = final["geometry"].(map[string]interface{})

			var location map[string]interface{} = geometry["location"].(map[string]interface{})
			apiresult := LatLng{}
			apiresult.Lat = location["lat"].(float64)
			apiresult.Lng = location["lng"].(float64)

			return apiresult
		}
	}
	return LatLng{}
}
func LinkHttpTest() {

	var name = "Allen"
	fmt.Println(name)
	myFunc := gogogo.GetHelloWorld
	fmt.Println(myFunc())

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Worldf")
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func saveStoreToDB(stores []interface{}) {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("Hello").C("Store")

	for _, store := range stores {
		_, err = c.Upsert(store, store)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	// result := Person{}
	// err = c.Find(bson.M{"name": "Ale"}).One(&result)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// fmt.Println("Phone:", result.Phone)
}

func saveOneDataToDB(data interface{}, tableName string) {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)

	c := session.DB("Hello").C(tableName)

	err = c.Insert(data)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func findStoreFromDB(id string) DB_Store {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("Hello").C("Store")

	store := DB_Store{}
	store.StoreId = id

	fmt.Print("ID:", id)
	var findMap map[string]interface{}

	jsonDecode, _ := json.Marshal(store)

	json.Unmarshal(jsonDecode, findMap)

	resp := c.Find(bson.M{"_id": id}).Select(bson.M{"_id": id, "url": "", "name": ""})
	respStore := DB_Store{}
	_ = resp.One(&respStore)

	return respStore
}

func findHostOrderFromDB() []DB_MainOrder {
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("Hello").C("MainOrder")

	// mainOrder := DB_MainOrder{}

	resp := c.Find(nil)
	respStore := []DB_MainOrder{}
	_ = resp.All(&respStore)

	return respStore
}

func searchMenu(w http.ResponseWriter, r *http.Request) {

	request, request_err := ioutil.ReadAll(r.Body)

	if request_err != nil {
		fmt.Fprint(w, "{}")
		return
	}

	var resqusetJson map[string]interface{}

	json.Unmarshal(request, &resqusetJson)

	store_id := resqusetJson["store_id"].(string)

	var store DB_Store = findStoreFromDB(store_id)

	url := "https://www.foodpanda.com.tw" + store.StoreUrl
	// url := "https://www.foodpanda.com.tw" + resqusetJson["restaurant_url"].(string)
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Print(string(body))

	doc, err := goquery.NewDocument(url)

	if err != nil {
		fmt.Fprint(w, "{}")
		return
	}

	fmt.Fprint(w, doc.Find((".where-wrapper")).Get(0).Attr[1].Val)
}
func mSearchMenu(storeUrl string) string {

	if storeUrl == "" {
		return ""
	}

	url := "https://www.foodpanda.com.tw" + storeUrl
	// body, _ := ioutil.ReadAll(resp.Body)
	// fmt.Print(string(body))

	doc, err := goquery.NewDocument(url)

	if err != nil {
		return ""
	}

	return doc.Find((".where-wrapper")).Get(0).Attr[1].Val
}
