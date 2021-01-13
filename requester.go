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

  taskSpec := spincomm.TaskSpec{
    Filters:     []string{"Resource", "Affinity"},
    Sort:        "Geolocation",
    ResourceMap: map[string]*spincomm.ResourceRequirement{},
    Ports:       map[string]string{},
    IsPublic:    false,
    NumReplicas: 1,
    CargoSpec: &spincomm.CargoReq{
      Size: 1,
      NReplica: 3,
    },
  }
  taskSpec.ResourceMap["CPU"] = &spincomm.ResourceRequirement{
    Weight: 0.5,
    Requested: 0,
    Required: true,
  }
  taskSpec.ResourceMap["Memory"] = &spincomm.ResourceRequirement{
    Weight: 0.5,
    Requested: 1,
    Required: true,
  }
  taskSpec.Ports["8080"] = ""

  request := &spincomm.TaskRequest{
    AppId: &spincomm.UUID{Value: "test_request"},
    Image: "zhiying12/goface-new",
    Command: []string{},
    Tty: true,
    Limits: &spincomm.TaskLimits{CpuShares: 2},
    Taskspec: &taskSpec,
    Port: 8080,
    TaskId: &spincomm.UUID{Value: "test_request_1"},
  }
  log.Println("sending task")
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
