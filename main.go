package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/iotaledger/iota.go/address"
	. "github.com/iotaledger/iota.go/consts"
)

func generate_adresses(start uint64, end uint64, seed string, result chan []string) {
	addresses := make([]string, end-start)
	for j := start; j < end; j++ {
		address, err := address.GenerateAddress(seed, j, SecurityLevelMedium)
		must(err)
		addresses[j-start] = fmt.Sprintf("%v: %s\n", j, address)
	}
	result <- addresses
}

func main() {
	t1 := time.Now()
	seed := strings.Repeat("N", 81)
	const threads = 8
	total_addresses := 1000
	amount := total_addresses / threads
	var chans [threads]chan []string

	for i := 0; i < threads; i++ {
		chans[i] = make(chan []string)
		if i == threads-1 {
			go generate_adresses(uint64(i*amount), uint64(((i+1)*amount)+total_addresses%threads), seed, chans[i])
		} else {
			go generate_adresses(uint64(i*amount), uint64((i+1)*amount), seed, chans[i])
		}
	}

	var final_addresses []string = nil
	for i := 0; i < threads; i++ {
		addresses := <-chans[i]
		final_addresses = append(final_addresses, addresses...)
	}
	fmt.Println(final_addresses)
	t2 := time.Now()
	diff := t2.Sub(t1)
	fmt.Println(diff)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
