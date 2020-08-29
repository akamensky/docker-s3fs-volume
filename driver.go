package main

import (
	"fmt"
	"github.com/docker/go-plugins-helpers/volume"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type S3Driver struct {
	volumes    map[string]string
	m          *sync.Mutex
	mountPoint string
}

func NewS3Driver() *S3Driver {
	return &S3Driver{
		volumes: make(map[string]string),
		m:       &sync.Mutex{},
		// TODO: Create it?
		mountPoint: "/mnt/mounts",
	}
}

func (d *S3Driver) Create(r *volume.CreateRequest) error {
	log.Print("Create volume ", r)

	d.m.Lock()
	defer d.m.Unlock()

	if _, ok := d.volumes[r.Name]; ok {
		return nil
	}

	volumePath := filepath.Join(d.mountPoint, r.Name)

	_, err := os.Lstat(volumePath)
	if err != nil {
		log.Printf("Error %s %v", volumePath, err.Error())
		return fmt.Errorf("error: %s: %s", volumePath, err.Error())
	}

	d.volumes[r.Name] = volumePath

	return nil
}

func (d *S3Driver) List() (*volume.ListResponse, error) {
	log.Print("List")

	volumes := make([]*volume.Volume, 0)

	for name, path := range d.volumes {
		volumes = append(volumes, &volume.Volume{
			Name:       name,
			Mountpoint: path,
		})
	}

	return &volume.ListResponse{Volumes: volumes}, nil
}

func (d *S3Driver) Get(r *volume.GetRequest) (*volume.GetResponse, error) {
	log.Print("Get ", r)

	if path, ok := d.volumes[r.Name]; ok {
		return &volume.GetResponse{
			Volume: &volume.Volume{
				Name:       r.Name,
				Mountpoint: path,
			},
		}, nil
	}
	return nil, fmt.Errorf("volume named %s not found", r.Name)
}

func (d *S3Driver) Remove(r *volume.RemoveRequest) error {
	log.Print("Remove volume ", r)

	d.m.Lock()
	defer d.m.Unlock()

	if _, ok := d.volumes[r.Name]; ok {
		delete(d.volumes, r.Name)
	}

	return nil
}

func (d *S3Driver) Path(r *volume.PathRequest) (*volume.PathResponse, error) {
	log.Print("Get volume path", r)

	if path, ok := d.volumes[r.Name]; ok {
		return &volume.PathResponse{
			Mountpoint: path,
		}, nil
	}
	return nil, fmt.Errorf("volume named %s not found", r.Name)
}

func (d *S3Driver) Mount(r *volume.MountRequest) (*volume.MountResponse, error) {
	log.Print("Mount volume ", r)

	if path, ok := d.volumes[r.Name]; ok {
		return &volume.MountResponse{
			Mountpoint: path,
		}, nil
	}

	return nil, fmt.Errorf("volume named %s not found", r.Name)
}

func (d *S3Driver) Unmount(r *volume.UnmountRequest) error {
	log.Print("Unmount ", r)
	return nil
}

func (d *S3Driver) Capabilities() *volume.CapabilitiesResponse {
	log.Print("Capabilities")

	return &volume.CapabilitiesResponse{
		Capabilities: volume.Capability{
			Scope: "local",
		},
	}
}
