package main

import "errors"
import "fmt"
import "io/ioutil"
import "net/http"
import "strconv"
import "encoding/json"

//main service class
type TopJsonService struct {
	requestSelector RequestSelector
}

//start and init service
func (serviceState *TopJsonService) Start(listenPort int) error {
	serviceState.requestSelector = RequestSelector{}
	serviceState.requestSelector.Init()
	fmt.Println(":" + strconv.Itoa(listenPort))
	http.HandleFunc("/webtop", serviceState.ServeHTTP)
	http.HandleFunc("/index.html", serviceState.ServePage)
	http.HandleFunc("/samplejson", serviceState.ReturnDummyReq)
	retVal := http.ListenAndServe(":"+strconv.Itoa(listenPort), nil)
	return retVal
}

//parse json message
func parseJsonRequest(request *http.Request) (Request, error) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		return nil, errors.New("Can't read request body")
	}
	var reqData Request
	err = json.Unmarshal(body, reqData)
	if err != nil {
		return nil, errors.New("Can't parse json")
	}
	return reqData, nil
}

//serve http responce in different thread
func (service *TopJsonService) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	if reqData, err := parseJsonRequest(request); err == nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
	} else {
		service.requestSelector.Dispatch(reqData, responseWriter, request)
	}

}

//serve main page request
func (service *TopJsonService) ServePage(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.Header().Set("Content-Type: text/html", "*")
	content, err := ioutil.ReadFile("index.html")
	if err != nil {
		responseWriter.Write([]byte("Can't find start page"))
		return
	}
	responseWriter.Write(content)
}

func (service *TopJsonService) ReturnDummyReq(responseWriter http.ResponseWriter, request *http.Request) {
	//this is service is not need lock
	serviceState := ServiceStateRequest{}
	js, err := json.Marshal(serviceState)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(js)
	return

}
