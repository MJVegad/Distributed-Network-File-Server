package main

import (
	"net"
	"strings"
	"bufio"
	"io"
	//"fmt"
	"strconv"
	)

type Command struct {
	Fields []string
	Content []byte
	Result chan string
}

type File struct {
	Numbytes uint64
	Version uint64
	//Exptime uint64
	Content []byte
}

func Extend(slice []byte, slice1 []byte) []byte {
    n := len(slice)
    n1 := len(slice1)
    slice = slice[0 : n+n1]
    
    for i := 0; i<n1; i++ {
    	slice[n+i] = slice1[i]
    } 
    
    return slice
}

/*func fileServer(commands chan Command) {
	
		var filerepo = make(map[string]File)
		var data = make(map[string]string)
		for cmd := range commands {
			if len(cmd.Fields)<2 {
				cmd.Result <- "Expected atleast 2 arguments"
				continue
			}

			fmt.Println("GOT command::", cmd)

			switch cmd.Fields[0] {
		case "write":
			if len(cmd.Fields) < 3 {
				cmd.Result <- "ERR_CMD_ERR\r\n"
				continue	
			}
			key := cmd.Fields[1]
			if val, ok := filerepo[key]; ok {
				val.Content = cmd.Content
				val.Numbytes,_ = strconv.ParseUint(cmd.Fields[2],10,64)
				
			} else {
				numbytes1,_ := strconv.ParseUint(cmd.Fields[2],10,64)
				filerepo[key] = File{Numbytes: numbytes1, Version: 0, Content: cmd.Content}
			}
			
			//data[key] = value
			
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
}*/

var filerepo = make(map[string]File)

func handle(conn net.Conn) {
		//defer conn.Close()

		//scanner := bufio.NewScanner(conn)
		//for scanner.Scan() {
		nr := bufio.NewReader(conn)
		//ln := scanner.Text()
		//fs := strings.Fields(ln)
		ln,_:= nr.ReadString('\n')
		//ln := string(ln1)
		fs := strings.Fields(ln)

		//result := make(chan string)
		
		switch fs[0] {
			
		case "write":	
			nb, _ := strconv.Atoi(fs[2])
			data := make([]byte,nb)
			_,_ = io.ReadFull(nr, data)
			key := fs[1]
			if val, ok := filerepo[key]; ok {
				numbytes1,_ := strconv.ParseUint(fs[2],10,64)
				filerepo[key] = File{Numbytes: numbytes1, Version: val.Version, Content: data}
				
			} else {
				numbytes1,_ := strconv.ParseUint(fs[2],10,64)
				filerepo[key] = File{Numbytes: numbytes1, Version: 0, Content: data}
			}
		
			io.WriteString(conn, "OK "+strconv.FormatUint(filerepo[key].Version,10)+"\n")		
		case "read":
			key := fs[1]
			if val, ok := filerepo[key]; ok {
			io.WriteString(conn, "CONTENTS "+strconv.FormatUint(val.Version,10)+" "+strconv.FormatUint(filerepo[key].Numbytes,10)+"\n"+string(filerepo[key].Content)+"\n")
			} else {
				io.WriteString(conn, "ERR_FILE_NOT_FOUND\n")
			}	
		case "cas":
			key := fs[1]
			nb, _ := strconv.ParseUint(fs[3],10,64)
			data := make([]byte,nb)
			_,_ = io.ReadFull(nr, data)
			if val, ok := filerepo[key]; ok {
				version := strconv.FormatUint(val.Version,10)
				if strings.Compare(version, fs[2])==0 {
					numbytes1,_ := strconv.ParseUint(fs[3],10,64)
					filerepo[key] = File{Numbytes: numbytes1, Version: val.Version+1, Content: data}
					
					io.WriteString(conn, "OK "+strconv.FormatUint(filerepo[key].Version,10)+"\n")
				} else {
					io.WriteString(conn, "ERR_VERSION\n")
				}    
			} else {
				io.WriteString(conn, "ERR_FILE_NOT_FOUND\n")
			}	
		case "delete":
			key := fs[1]
			if _, ok := filerepo[key]; ok {
				delete(filerepo, key)
				io.WriteString(conn, "OK\n")
			} else {
				io.WriteString(conn, "ERR_FILE_NOT_FOUND\n")
			}
				
		default:
			io.WriteString(conn, "ERR_CMD_ERR\n")
				
		}	
			
		//io.WriteString(conn, ln)
		
	}



	
func serverMain() 	{
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	} 
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		} 

		//io.WriteString(conn, fmt.Sprint("Hello World\n", time.Now(), "\n"))
		for {
		handle(conn)
	    }

		conn.Close()

	}
}	

func main() {
	serverMain()
}
