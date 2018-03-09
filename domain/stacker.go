package domain

import (
	"fmt"
	"log"
)

type Stacker struct {
	pickedContainer *FreightContainer
	Name            string
	isFree          bool
}

func (stacker *Stacker) startWork(containerToStacker <-chan *FreightContainer,
	rejectedContainer chan<- *FreightContainer, containerToContainerArea chan *FreightContainer,
		requestContainerSpace chan<- string, containerAreaClearance <-chan bool) error {

	for {
		select {
			case freightContainer := <- containerToStacker:
				if stacker.isFree {

					log.Println(fmt.Sprintf("picked container %s", freightContainer))
					stacker.pickedContainer = freightContainer
					requestContainerSpace <- stacker.Name
					log.Println("requested space at container area")
				} else {
					log.Println("stacker is already in use. will reject it.")
					rejectedContainer <- freightContainer
				}
			case clearance := <- containerAreaClearance:
				if clearance {
					containerToContainerArea <- stacker.pickedContainer
					stacker.pickedContainer = nil
					stacker.isFree = true

				} else {
					log.Fatal("negative clearance for container area. must. not. happen.")
				}

		}

	}

}
