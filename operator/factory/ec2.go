package factory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/masahiro331/funnel/pod"
	"io"
	"log"
	"net/http"
)

const (
	ImageID = "ami-0bba69335379e17f8"
)

var (
	_ Factory = &EC2Factory{}
	_ Pod     = &EC2{}
)

type EC2Factory struct{}

type EC2 struct {
	c          *Client
	ctx        context.Context
	InstanceId string

	PublicIp  string
	PrivateIp string
}

func (e *EC2) Target() string {
	return e.PublicIp
}

func (e *EC2) Exec(taskId, name string, args []string) error {
	url := fmt.Sprintf("http://%s:6332/exec", e.Target())
	request := pod.ExecRequest{
		TaskId: taskId,
		Name:   name,
		Args:   args,
	}

	b, err := json.Marshal(&request)
	if err != nil {
		return fmt.Errorf("exec request marshal error: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("post exec request error: %v", err)
	}
	defer resp.Body.Close()

	return nil
}

func (e *EC2) Statuses() ([]pod.Status, error) {
	url := fmt.Sprintf("http://%s:6332/statuses", e.Target())
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("get statuses request error: %v", err)
	}
	defer resp.Body.Close()

	var statuses []pod.Status

	if err := json.NewDecoder(resp.Body).Decode(&statuses); err != nil {
		return nil, fmt.Errorf("decode statuses response error: %v", err)
	}

	return statuses, nil
}

func (e *EC2) Pull(taskId string) (pod.Result, error) {
	url := fmt.Sprintf("http://%s:6332/pull", e.Target())
	r := pod.PullRequest{
		TaskId: taskId,
	}
	b, err := json.Marshal(&r)
	if err != nil {
		return pod.Result{}, fmt.Errorf("marshal pull request error: %v", err)
	}

	resp, err := http.Post(url, "", bytes.NewReader(b))
	if err != nil {
		return pod.Result{}, fmt.Errorf("post pull request error: %v", err)
	}
	defer resp.Body.Close()

	buf := bytes.NewBuffer(nil)
	io.Copy(buf, resp.Body)

	return pod.Result{
		Status: pod.Status{
			TaskId: taskId,
			State:  pod.Done,
		},
		Content: buf,
	}, nil
}

func (e *EC2) Delete() error {
	_, err := e.c.TerminateInstances(e.ctx, &ec2.TerminateInstancesInput{
		InstanceIds: []string{e.InstanceId},
	})
	if err != nil {
		return fmt.Errorf("terminate instance error: %v", err)
	}
	log.Printf("terminating... %s", e.InstanceId)
	return nil
}

func (e *EC2) Name() string {
	return e.InstanceId
}

func (e *EC2) Ready() (bool, error) {
	o, err := e.c.DescribeInstanceStatus(e.ctx, &ec2.DescribeInstanceStatusInput{
		InstanceIds: []string{e.InstanceId},
	})
	if err != nil {
		return false, err
	}
	for _, s := range o.InstanceStatuses {
		if s.InstanceStatus.Status == "ok" &&
			s.SystemStatus.Status == "ok" &&
			s.InstanceState.Name == "running" {
			d, err := e.c.DescribeInstances(e.ctx, &ec2.DescribeInstancesInput{
				InstanceIds: []string{e.InstanceId},
			})
			if err != nil {
				return false, err
			}
			for _, r := range d.Reservations {
				for _, i := range r.Instances {
					e.PublicIp = *i.PublicIpAddress
					e.PrivateIp = *i.PrivateIpAddress
				}
			}

			return true, nil
		}
	}
	return false, nil
}

func (e *EC2Factory) Create(number int) ([]Pod, error) {
	ctx := context.Background()
	c, _ := New(ctx, Option{})
	userData := "IyEvYmluL2Jhc2gKc3VkbyB5dW0gaW5zdGFsbCAteSBubWFwCmNkIC9ob21lL2VjMi11c2VyCmN1cmwgLUxPIGh0dHBzOi8vZ2l0aHViLmNvbS9tYXNhaGlybzMzMS9mdW5uZWwvcmVsZWFzZXMvZG93bmxvYWQvMC4wLjEvZnVubmVsXzAuMC4xX0xpbnV4X3g4Nl82NC50YXIuZ3oKdGFyIHh2ZnogZnVubmVsXzAuMC4xX0xpbnV4X3g4Nl82NC50YXIuZ3oKc3VkbyAuL2Z1bm5lbCBwb2Q="
	input := &ec2.RunInstancesInput{
		MaxCount:     toPtr(int32(number)),
		MinCount:     toPtr(int32(1)),
		ImageId:      toPtr(ImageID),
		InstanceType: types.InstanceTypeT2Micro,
		UserData:     &userData,
	}
	o, err := c.RunInstances(ctx, input)
	if err != nil {
		return nil, err
	}

	var pods []Pod
	for _, instance := range o.Instances {
		pods = append(pods, &EC2{
			c:          c,
			ctx:        ctx,
			InstanceId: *instance.InstanceId,
		})
		log.Printf("creating... %s", *instance.InstanceId)
	}
	return pods, nil
}

func toPtr[T any](t T) *T {
	return &t
}

type Option struct {
	AwsSecretKey string
	AwsAccessKey string
	AwsRegion    string
}

type Client struct {
	*ec2.Client
}

func New(ctx context.Context, option Option) (*Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}
	if option.AwsRegion != "" {
		cfg.Region = option.AwsRegion
	}
	c := ec2.NewFromConfig(cfg)

	return &Client{c}, nil
}
