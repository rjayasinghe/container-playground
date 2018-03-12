package domain

import (
	"log"
	"fmt"
	"sync"
	"time"
)

type ContainerArea struct {
	Name         	string
	StorageSpace 	[]*FreightContainer
	reservationLock *sync.Mutex
}

func NewContainerArea(capacity int, name string) *ContainerArea {
	containerArea := new (ContainerArea)
	containerArea.Name = name
	containerArea.StorageSpace = make([]*FreightContainer, capacity)
	containerArea.reservationLock = &sync.Mutex{}

	return containerArea
}

func retryLater(requestContainerSpace chan string, stackerName string) {
	log.Println("will wait a bit and try request again")
	time.Sleep(time.Second/2)
	requestContainerSpace <- stackerName
}

func (containerArea *ContainerArea) startWork(
	requestContainerSpace chan string, containerAreaClearance chan<- bool,
	containerToContainerArea <-chan *FreightContainer) {
	for {
		select {
			case stackerName := <- requestContainerSpace:
				if len(containerArea.StorageSpace) < cap(containerArea.StorageSpace) {
					containerArea.reservationLock.Lock()
					log.Println(fmt.Sprintf("stacker %s requested a container space", stackerName))

					//somebody could have taken the last slot between checking and locking
					if len(containerArea.StorageSpace) < cap(containerArea.StorageSpace) {
						containerAreaClearance <- true //TODO this currently works for only one stacker
					} else {
						go retryLater(requestContainerSpace, stackerName)
					}
					log.Println(fmt.Sprintf("sent clearance for stacker  %s to dump one container", stackerName))
				} else { // try again in half a second
					go retryLater(requestContainerSpace, stackerName)
				}

			case freightContainer := <- containerToContainerArea:
				if len(containerArea.StorageSpace) < cap(containerArea.StorageSpace) {
					containerArea.StorageSpace = append(containerArea.StorageSpace, freightContainer)
					containerArea.reservationLock.Unlock()
				} else {
					panic("tried to take a container although capacity is full. this must not happen.")
				}
		}
	}
}