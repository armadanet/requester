package main

import (
  "github.com/armadanet/captain/dockercntrl"
  "github.com/armadanet/spinner/spinresp"
  "github.com/armadanet/comms"
  "github.com/google/uuid"
  "log"
  "os"
  "fmt"
)

type BeaconResponse struct {
  Valid         bool    `json:"Valid"`  // true if find a spinner
  Token         string  `json:"Token"`
  Ip            string  `json:"Ip"`
  OverlayName   string  `json:"OverlayName"`
  ContainerName string  `json:"ContainerName"`
}

func main() {
  URL := os.Getenv("URL")
  container_name := os.Getenv("CONTAINER_NAME")

  // query the beacon for a spinner to submit the job
  fmt.Println("Query Beacon for a spinner...")
  var res BeaconResponse
  err := comms.SendGetRequest(URL, &res)
  if err != nil {
    log.Println(err)
    return
  }
  // join the selected spinner overlay network
  state, err := dockercntrl.New()
  if err != nil {
    log.Println(err)
    return
  }
  fmt.Println("Found spinner: "+res.ContainerName+". Now join its overlay network...")
  err = state.JoinSwarmAndOverlay(res.Token, res.Ip, container_name, res.OverlayName)
  if err != nil {
    log.Println(err)
    return
  }

  // access spinner through
  fmt.Println("Now sending the task")
  dialurl := "ws://"+res.ContainerName+":5912/spin"
  socket, err := comms.EstablishSocket(dialurl)
  if err != nil {return}
  var resp spinresp.Response
  socket.Start(resp)
  reader := socket.Reader()
  writer := socket.Writer()
  for i := 1; i <= 1; i++ {
    container := dockercntrl.Config{
      Image: "docker.io/codyperakslis/armada-cargo-test",
      Cmd: []string{"./main"},
      Tty: true,
      Name: uuid.New().String(),
      Env: []string{},
      Port: 0,
      Limits: &dockercntrl.Limits{
        CPUShares: 2,
      },
      Storage: true,
    }
    writer <- container
  }
  for {
    select {
    case data, ok := <- reader:
      if !ok {break}
      response, ok := data.(*spinresp.Response)  // convert to type spinresp.Response
      if !ok {break}
      fmt.Println(response)
    }
  }
}
