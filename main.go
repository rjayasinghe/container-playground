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
	containerArea := domain.NewContainerArea(10, "yolo-area")

	go containerArea.StartWork(requestContainerSpace, containerAreaClearance, containerToContainerArea)

	go stacker.StartWork(containerToStacker, rejectedContainer, containerToContainerArea, requestContainerSpace, containerAreaClearance)

	<- forever
}
