package main

import (
  "github.com/open-nebula/captain/dockercntrl"
  "github.com/open-nebula/spinner/spinresp"
  "github.com/open-nebula/comms"
  "log"
)

func main() {
  dialurl := "wss://412db34f.ngrok.io/spin"
  socket, err := comms.EstablishSocket(dialurl)
  if err != nil {return}
  comms.Reader(func(data interface{}, ok bool) {
    log.Println("something")
    if !ok {return}
    log.Println(data)
    config, ok := data.(spinresp.Response)
    if !ok {return}
    log.Println(config)
  }, socket.Reader())
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
  for{select{}}
}
