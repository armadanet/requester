package requester

import (
  "google.golang.org/grpc"
  "github.com/armadanet/spinner/spincomm"
  "context"
  "log"
)

func  Run(dialurl string) error {
  var opts []grpc.DialOption
  opts = append(opts, grpc.WithInsecure())
  conn, err := grpc.Dial(dialurl, opts...)
  if err != nil {return err}
  defer conn.Close()
  log.Println("Connected")
  client := spincomm.NewSpinnerClient(conn)

  ctx := context.Background()
  

  request := &spincomm.TaskRequest{
    TaskId: &spincomm.UUID{Value: "test_request"},
    Image: "docker.io/codyperakslis/armada-cargo-test",
    Command: []string{"echo", "hello"},
    Tty: true,
    Limits: &spincomm.TaskLimits{CpuShares: 2},
  }
  stream, err := client.Request(ctx, request)
  if err != nil {return err}

  for {
    taskLog, err := stream.Recv()
    if err != nil {return err}
    log.Println(taskLog)
  }
  log.Println("Loop ended unexpectantly")
  return nil
}
