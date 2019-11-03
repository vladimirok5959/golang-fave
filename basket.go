package main

import (
	"fmt"
	"time"

	"golang-fave/engine/basket"
)

func basket_clean_do(sb *basket.Basket, stop chan bool) {
	sb.Cleanup()
}

func basket_clean_start(sb *basket.Basket) (chan bool, chan bool) {
	ch := make(chan bool)
	stop := make(chan bool)

	go func() {
		for {
			select {
			case <-time.After(30 * time.Minute):
				// Cleanup every 30 minutes
				basket_clean_do(sb, stop)
			case <-ch:
				ch <- true
				return
			}
		}
	}()
	return ch, stop
}

func basket_clean_stop(ch, stop chan bool) {
	for {
		select {
		case stop <- true:
		case ch <- true:
			<-ch
			return
		case <-time.After(3 * time.Second):
			fmt.Println("Basket error: force exit by timeout after 3 seconds")
			return
		}
	}
}
