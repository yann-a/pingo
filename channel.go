package main

import "sync"

// InfiniteChannel implements the Channel interface with an infinite buffer between the input and the output.
type channel struct {
	input   chan value
	request chan chan value
}

func createChannel(wg *sync.WaitGroup) channel {
	ch := channel{
		input:   make(chan value),
		request: make(chan chan value),
	}
	go ch.infiniteBuffer(wg)

	return ch
}

func (ch *channel) infiniteBuffer(wg *sync.WaitGroup) {
	var buffer []value
	var requestsBuffer []chan value

	for {
		select {
		case elem := <-ch.input:
			if len(requestsBuffer) > 0 { // There's a pending request
				requestsBuffer[0] <- elem
				close(requestsBuffer[0]) // We close the private channel
				requestsBuffer = requestsBuffer[1:]
				// no wg.Done(), the number of processes stays constant as we move on to the message listener
			} else {
				buffer = append(buffer, elem)
				wg.Done() // We wait for a request
			}

		case request := <-ch.request:
			if len(buffer) > 0 { // There are elements waiting to be sent
				request <- buffer[0]
				buffer = buffer[1:]
			} else {
				requestsBuffer = append(requestsBuffer, request)
				wg.Done() // We wait for elements to send
			}
		}
	}
}
