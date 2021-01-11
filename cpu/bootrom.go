package cpu

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func LoadBootrom(f *os.File) {
	b, err := ioutil.ReadFile(f.Name())
	if err != nil{
		panic(err)
	}
	sum := md5.Sum(b)
	actual := fmt.Sprintf("%x", sum)
	expected := "32fbbd84168d3482956eb3c5051637f5"
	if actual != expected {
		log.Panicf("bootrom checksum does not match, actual: %s, expected %s", actual, expected)
	}
	for i := range b {
		c.ram[i] = b[i]
	}
}