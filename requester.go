package requester

import (
	"context"
	"encoding/json"
	"github.com/armadanet/spinner/spincomm"
	"google.golang.org/grpc"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

func Run(dialurl string) error {
	submitTenTask(dialurl)
	return nil
}

// Submit the task request to spinner
func submitTask(client spincomm.SpinnerClient, request *spincomm.TaskRequest) error {
	ctx := context.Background()
	log.Println("sending task")
	log.Println(request)
	stream, err := client.Request(ctx, request)
	if err != nil {
		log.Println(err)
		return err
	}

	for {
		taskLog, err := stream.Recv()
		if err != nil {
            log.Println(err)
			return err
		}
		log.Println(taskLog)
	}
	log.Println("Loop ended unexpectantly")
	return nil
}

// Submit 10 times
func submitTenTask(dialurl string) error {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(dialurl, opts...)
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Println("Connected")
	client := spincomm.NewSpinnerClient(conn)
	var wg sync.WaitGroup

	taskRequest := GenTaskRequest()
	for i := 0; i < 10; i++ {
		taskID := "task_" + strconv.Itoa(i)
		lat := rand.Float64() * 100
		lon := rand.Float64() * 100
		request := CopyRequest(taskRequest, taskID, lat, lon)
		wg.Add(1)
		go worker(client, request, &wg)
		time.Sleep(3 * time.Second)
	}
	wg.Wait()
	return nil
}

func worker(client spincomm.SpinnerClient, request *spincomm.TaskRequest, wg *sync.WaitGroup) {
	submitTask(client, request)
	wg.Done()
}

// Copy from original task request
func CopyRequest(tq *spincomm.TaskRequest, taskID string, lat float64, lon float64) *spincomm.TaskRequest {
	req := new(spincomm.TaskRequest)
	buffer, _ := json.Marshal(tq)
	json.Unmarshal([]byte(buffer), req)
	req.TaskId = &spincomm.UUID{Value: taskID}
	req.GetTaskspec().DataSources = &spincomm.Location{Lat: lat, Lon: lon}
	return req
}

// A new task request
func GenTaskRequest() *spincomm.TaskRequest {
	taskSpec := spincomm.TaskSpec{
		Filters:     []string{"Resource", "Affinity"},
		Sort:        "Geolocation",
		ResourceMap: map[string]*spincomm.ResourceRequirement{},
		Ports:       map[string]string{},
		IsPublic:    false,
		NumReplicas: 1,
		CargoSpec: &spincomm.CargoReq{
			Size:     1,
			NReplica: 3,
		},
		DataSources: &spincomm.Location{Lat: 40.0196, Lon: -90.2402},
	}
	taskSpec.ResourceMap["CPU"] = &spincomm.ResourceRequirement{
		Weight:    0.5,
		Requested: 0,
		Required:  true,
	}
	taskSpec.ResourceMap["Memory"] = &spincomm.ResourceRequirement{
		Weight:    0.5,
		Requested: 1,
		Required:  true,
	}
	taskSpec.Ports["8080"] = ""

	request := &spincomm.TaskRequest{
		AppId:    &spincomm.UUID{Value: "App_1"},
		Image:    "zhiying12/goface-new",
		Command:  []string{},
		Tty:      true,
		Limits:   &spincomm.TaskLimits{CpuShares: 2},
		Taskspec: &taskSpec,
		Port:     8080,
		TaskId:   &spincomm.UUID{Value: "task_0"},
	}
	return request
}
