package main

import (
	"github.com/rjayasinghe/container-playgrond/domain"
)

func main() {
	forever := make(chan bool)

	containerToStacker := make(chan *domain.FreightContainer)
	rejectedContainer := make(chan *domain.FreightContainer)
	containerToContainerArea := make(chan *domain.FreightContainer)
	requestContainerSpace := make(chan string)
	containerAreaClearance := make(chan bool)

	stacker := domain.Stacker{Name : "stacker 01" }
	go stacker.StartWork(containerToStacker, rejectedContainer, containerToContainerArea, requestContainerSpace, containerAreaClearance)

	<- forever
}
