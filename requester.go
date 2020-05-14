package main

import (
  "github.com/armadanet/captain/dockercntrl" // @cargo-connect
  "github.com/armadanet/spinner/spinresp"
  "github.com/armadanet/comms"
  "github.com/google/uuid"
  "log"
)

func main() {
  dialurl := "wss://69de8bbe.ngrok.io/spin"
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
      response, ok := data.(*spinresp.Response)
      if !ok {break}
      log.Println(response)
    }
  }
}
