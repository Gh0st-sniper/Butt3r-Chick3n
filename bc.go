package main

import (

	"fmt"
	"net"
	//"golang.org/x/crypto/ssh"
	"bufio"
	//"sync"
	"os"
)


func main(){


var ip string
var wordlist string

ip = os.Args[1]
wordlist = os.Args[2]
//fmt.Println(len(os.Args))

if len(os.Args) > 2 {


	knock(ip)
	readW(wordlist)

}else{

	fmt.Println("Usage: ./bc <Victim IP> <Wordlist>")
	os.Exit(1)
}


 
	
}


func readW(file string)(data []string, err error){

	buffer,err := os.Open(file)
	if err != nil {
		return 
	} 

	defer buffer.Close()

	scanner := bufio.NewScanner(buffer)
	for scanner.Scan(){
		data = append(data,scanner.Text())
		fmt.Println(data)
	}

	return 
}

func knock(ip string){

	full_ip := ip + ":22"
	conn, err := net.Dial("tcp",full_ip)
	if err != nil {
		fmt.Println("Is the port open ? or are you connected to the internet")
				
	}else{
		fmt.Println(" PORT 22 ::: IS UP")
		
	}
	conn.Close()	
}








