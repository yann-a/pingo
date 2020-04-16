package main

import "sync"

// InfiniteChannel implements the Channel interface with an infinite buffer between the input and the output.
type channel struct {
	input             chan value
  request           chan chan value
}

func createChannel(wg *sync.WaitGroup) channel {
	ch := channel{
		input:  make(chan value),
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
      if len(requestsBuffer) > 0 { // une requête en attente
        requestsBuffer[0] <- elem
        close(requestsBuffer[0]) // on ferme le channel privé
        requestsBuffer = requestsBuffer[1:]
        // pas de wg.Done(), on garde le même nombre de process en passant la main au récepteur du message
      } else {
				buffer = append(buffer, elem)
        wg.Done() // on se met en pause en attendant une requête
      }
		case request := <-ch.request:
			if len(buffer) > 0 { // des éléments sont en attente
        request <- buffer[0]
        buffer = buffer[1:]
      } else {
				requestsBuffer = append(requestsBuffer, request)
        wg.Done() // on se met en pause en attendant un push
      }
		}
	}
}
