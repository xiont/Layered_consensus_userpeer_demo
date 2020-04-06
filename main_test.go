package main

import (
	"fmt"
	"testing"
)

type ss struct {
	n string
}


func Test_Chan(t *testing.T) {

	bh := make(chan *ss,2)


	bh <- nil
	fmt.Println(len(bh))
	bh <- &ss{"111"}
	fmt.Println(len(bh))

	c := <- bh
	if c == nil {

	}

	fmt.Println(c)

	d := <- bh
	fmt.Println(*d)

}

//func fun1(){
//	bh = make(chan bool,1)
//
//}
//
//func fun2(){
//	bh <- true
//}
//
//func fun3(){
//	bh <- true
//}