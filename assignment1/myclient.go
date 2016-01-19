package main

import (
	"net"
	"strings"
	"bufio"
	"io"
	"fmt"
	)

type Command struct {
	Fields []string
	Result chan string
}

type File struct {
	numbytes int64
	version int64
	exptime int64
}


func fileServer(commands chan Command) {
		var data = make(map[string]string)
		for cmd := range commands {
			if len(cmd.Fields)<2 {
				cmd.Result <- "Expected atleast 2 arguments"
				continue
			}

			fmt.Println("GOT command", cmd)

			switch cmd.Fields[0] {
		case "SET":
			if len(cmd.Fields) != 3 {
				cmd.Result <- "Invalid command"
				continue	
			}
			key := cmd.Fields[1]
			value := cmd.Fields[2]
			data[key] = value
			cmd.Result <- ""
		case "GET":
			key := cmd.Fields[1]
			value := data[key]
			//io.WriteString(conn, value+"\n")
			cmd.Result <- value
		case "DEL":
			key := cmd.Fields[1]
			delete(data, key)
			cmd.Result <- ""
		default:
			cmd.Result <- "Invalid command "+ cmd.Fields[0] +"\n"
		}


		}		
}

func handle(commands chan Command, conn net.Conn) {
		defer conn.Close()

		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
		ln := scanner.Text()
		fs := strings.Fields(ln)

		result := make(chan string)
		commands <- Command{
			Fields: fs,
			Result: result,
		}

		io.WriteString(conn, <-result+"\n")
		
	}

}

	
func serverMain() 	{
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	} 
	defer ln.Close()

	commands := make(chan Command)
	go fileServer(commands)

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		} 

		//io.WriteString(conn, fmt.Sprint("Hello World\n", time.Now(), "\n"))
		go handle(commands, conn)

		//conn.Close()

	}
}	

func main() {
	serverMain()
}
