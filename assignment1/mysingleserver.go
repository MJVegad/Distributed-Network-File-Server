package main

import "net"
	

func main() {
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
			bytes := make([]byte, 1024)
			n, err := conn.Read(bytes)
			if err != nil {
				break
			}
			_, err1 := conn.Write(bytes[:n])
			if err1 != nil {
				break
			}
		}

		conn.Close()

	}

	
}
