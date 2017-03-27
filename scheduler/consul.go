package scheduler

import (
	"fmt"
	"io"
	"strconv"

	"github.com/faststackco/machinestack/config"
	"github.com/faststackco/machinestack/driver"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/consul/api"
	"github.com/jmcvetta/randutil"
)

type ConsulScheduler struct {
	driverOptions *config.DriverOptions
	health        *api.Health
	catalog       *api.Catalog
	kv            *api.KV
}

func NewConsulScheduler(options *config.DriverOptions) (Scheduler, error) {
	consul, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		return nil, err
	}

	return &ConsulScheduler{
		driverOptions: options,
		catalog:       consul.Catalog(),
		kv:            consul.KV(),
		health:        consul.Health(),
	}, nil
}

func (c *ConsulScheduler) Create(name, image, driverName string) (string, error) {
	hosts, _, err := c.health.Service(driverName, "", true, nil)
	if err != nil {
		return "", err
	}

	if len(hosts) == 0 {
		return "", fmt.Errorf("no hosts found for driver '%v'", driverName)
	}

	var choices []randutil.Choice
	for _, h := range hosts {
		weight, err := strconv.Atoi(h.Node.Meta["weight"])
		if err != nil {
			weight = 1
		}

		choices = append(choices, randutil.Choice{Item: h, Weight: weight})
	}

	choice, err := randutil.WeightedChoice(choices)
	if err != nil {
		return "", err
	}

	entry := choice.Item.(*api.ServiceEntry)

	driver, err := c.newDriver(driverName, entry.Node)
	if err != nil {
		return "", err
	}

	if err := driver.Create(name, image); err != nil {
		return "", err
	}

	return entry.Node.ID, nil
}

func (c *ConsulScheduler) Delete(name, driverName, nodeID string) error {

	node, _, err := c.catalog.Node(nodeID, nil)
	if err != nil {
		return err
	}

	driver, err := c.newDriver(driverName, node.Node)
	if err != nil {
		return err
	}

	if err := driver.Delete(name); err != nil {
		return err
	}

	return nil
}

func (c *ConsulScheduler) Exec(name, driverName, nodeID string, stdin io.ReadCloser, stdout io.WriteCloser, stderr io.WriteCloser, controlHandler func(*websocket.Conn)) error {

	node, _, err := c.catalog.Node(nodeID, nil)
	if err != nil {
		return err
	}

	driver, err := c.newDriver(driverName, node.Node)
	if err != nil {
		return err
	}

	return driver.Exec(name, stdin, stdout, stderr, controlHandler)
}

func (c *ConsulScheduler) newDriver(name string, node *api.Node) (driver.Driver, error) {
	driverOptions := make(map[string]string)
	for key, value := range *c.driverOptions {
		driverOptions[key] = value
	}

	// TODO protocol, port
	remote := fmt.Sprintf("%s:%v", node.Node, 1000)

	driverOptions[fmt.Sprintf("%s.remote", name)] = remote

	return driver.NewDriver(name, driverOptions)
}
