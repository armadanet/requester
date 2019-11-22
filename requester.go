package main

import (
  "github.com/open-nebula/captain/dockercntrl"
  "github.com/open-nebula/spinner/spinresp"
  "github.com/open-nebula/comms"
  "log"
)

func main() {
  dialurl := "wss://4fbf1747.ngrok.io/spin"
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
