package main

import (
  "github.com/open-nebula/captain/dockercntrl"
  "github.com/open-nebula/spinner/spinresp"
  "github.com/open-nebula/comms"
  "log"
  "net/http"
  "io/ioutil"
  "fmt"
  "encoding/json"  
)

type GeoIP struct {
	Ip         string  `json:"ip"`
	Lat        float32 `json:"latitude"`
	Lon        float32 `json:"longitude"`
}
type QueryResp struct {
  Port      int     `json:"port"`
  Ip        string  `json:"ip"`
}

func main() {
  // query beacon first
  var (
  	err      error
  	geo      GeoIP
  	response *http.Response
  	body     []byte
  )
  // get self IP and location info
  response, err = http.Get("http://api.ipstack.com/check?access_key=0bbaa9ccd131225ec08fa2c02c0a3260")
	if err != nil {
		fmt.Println(err)
	}
  body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
  err = json.Unmarshal(body, &geo)
	if err != nil {
		fmt.Println(err)
	}
  response.Body.Close()
  // Send location info to beacon to query the closest spinner
  lat := fmt.Sprintf("%f", geo.Lat)
  lon := fmt.Sprintf("%f", geo.Lon)
  response, err = http.Get("http://c062a166.ngrok.io"+"/query/"+lat+"/"+lon)
  if err != nil {
		log.Println(err)
    return
	}
  body, err = ioutil.ReadAll(response.Body)
  if err != nil {
    log.Println(err)
    return
  }
  var queryResp QueryResp
  err = json.Unmarshal(body, &queryResp)
  if err != nil {
    log.Println(err)
    return
  }
  log.Println("Closest spinner address:\t", queryResp.Ip)
  // "http://" + queryResp.Ip + "/join"

  dialurl := "wss://"+queryResp.Ip+"/spin"
  socket, err := comms.EstablishSocket(dialurl)
  if err != nil {return}
  var resp spinresp.Response
  socket.Start(resp)
  reader := socket.Reader()
  writer := socket.Writer()
  container := dockercntrl.Config{
    Image: "busybox",
    Cmd: []string{"echo", "hello"},
    Tty: true,
    Name: "james",
    Env: []string{},
    Port: 0,
    Limits: &dockercntrl.Limits{
      CPUShares: 2,
    },
  }
  writer <- container
  for {
    select {
    case data, ok := <- reader:
      if !ok {break}
      response, ok := data.(*spinresp.Response)
      if !ok {break}
      log.Println(response)
    }
  }
}
