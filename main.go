package main

import (

	"fmt"
	"VuforiaSample/vuforia"
)

func main(){
	client,err := vuforia.Init("access key","secret key")
	if err != nil{
		fmt.Printf("error %s",err.Error())
		return
	}

	idsItem, err := client.TargetIds()
	fmt.Printf("number of items %d \n",len(idsItem))
	for _,item := range idsItem{
		fmt.Printf("id %s \n",item)
	}
}