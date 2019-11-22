package main

import (
  "github.com/open-nebula/captain/dockercntrl"
  "github.com/open-nebula/spinner/spinresp"
  "github.com/open-nebula/comms"
  "github.com/google/uuid"
  "log"
)

func main() {
  dialurl := "wss://c2a13350.ngrok.io/spin"
  socket, err := comms.EstablishSocket(dialurl)
  if err != nil {return}
  var resp spinresp.Response
  socket.Start(resp)
  reader := socket.Reader()
  writer := socket.Writer()
  for i := 1; i <= 5; i++ {
    container := dockercntrl.Config{
      Image: "busybox",
      Cmd: []string{"echo", "hello"},
      Tty: true,
      Name: uuid.New().String(),
      Env: []string{},
      Port: 0,
      Limits: &dockercntrl.Limits{
        CPUShares: 2,
      },
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
