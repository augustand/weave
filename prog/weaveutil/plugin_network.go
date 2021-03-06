/* various weave docker network plugin operations */
package main

import (
	"fmt"

	"github.com/fsouza/go-dockerclient"
)

func createPluginNetwork(args []string) error {
	if len(args) != 2 {
		cmdUsage("create-plugin-network", "<network-name> <driver-name>")
	}
	networkName := args[0]
	driverName := args[1]
	d, err := newDockerClient()
	if err != nil {
		return err
	}
	_, err = d.CreateNetwork(
		docker.CreateNetworkOptions{
			Name:           networkName,
			CheckDuplicate: true,
			Driver:         driverName,
			IPAM:           docker.IPAMOptions{Driver: driverName},
		})
	if err != docker.ErrNetworkAlreadyExists && err != nil {
		// Despite appearances to the contrary, CreateNetwork does
		// sometimes(always?) *not* return ErrNetworkAlreadyExists
		// when the network already exists. Hence we need to check for
		// this explicitly.
		if _, err2 := d.NetworkInfo(networkName); err2 != nil {
			return fmt.Errorf("unable to create network: %s", err)
		}
	}
	return nil
}

func removePluginNetwork(args []string) error {
	if len(args) != 1 {
		cmdUsage("remove-plugin-network", "<network-name>")
	}
	networkName := args[0]
	d, err := newDockerClient()
	if err != nil {
		return err
	}
	err = d.RemoveNetwork(networkName)
	if _, ok := err.(*docker.NoSuchNetwork); !ok && err != nil {
		return fmt.Errorf("unable to remove network: %s", err)
	}
	return nil
}

func newDockerClient() (*docker.Client, error) {
	// API 1.21 is the first version that supports docker network
	// commands
	c, err := docker.NewVersionedClientFromEnv("1.21")
	if err != nil {
		return nil, fmt.Errorf("unable to connect to docker: %s", err)
	}
	_, err = c.Version()
	if err != nil {
		return nil, fmt.Errorf("unable to connect to docker: %s", err)
	}
	return c, nil
}
