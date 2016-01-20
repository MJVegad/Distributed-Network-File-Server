package main

import (
	"net"
	"strings"
	"bufio"
	"io"
	"strconv"
	"time"
	"sync"
	)

type Command struct {
	Fields []string
	Content []byte
	Result chan string
}

type File struct {
	Numbytes uint64
	Version uint64
	Exptime int64
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
var mux sync.RWMutex
var filerepo = make(map[string]File)


func handle(conn net.Conn) {
		for {
		nr := bufio.NewReader(conn)
		ln,_:= nr.ReadString('\n')
		fs := strings.Fields(ln)

		if len(fs)<1 {
			continue
		}
		
		switch fs[0] {
			
		case "write":	
			nb, _ := strconv.Atoi(fs[2])
			data := make([]byte,nb)
			_,_ = io.ReadFull(nr, data)
			
			if len(fs)<=4 {
				key := fs[1]
				var expt int64
				if len(key)<=250 {
				if len(fs)==4 {
					i,_ := strconv.ParseInt(fs[3],10,64)
					expt = time.Now().Unix() + i
				}
				var vson uint64
				mux.Lock()	
				if val, ok := filerepo[key]; ok {
					numbytes1,_ := strconv.ParseUint(fs[2],10,64)
					filerepo[key] = File{Numbytes: numbytes1, Version: val.Version,Exptime: expt, Content: data}
				    vson = filerepo[key].Version
				} else {
					numbytes1,_ := strconv.ParseUint(fs[2],10,64)
					filerepo[key] = File{Numbytes: numbytes1, Version: 0,Exptime: expt, Content: data}
					vson = 0
				}
				mux.Unlock()
				io.WriteString(conn, "OK "+strconv.FormatUint(vson,10)+"\r\n")	
		        } else {
		        	io.WriteString(conn,"ERR_INTERNAL\r\n")
		        }
			} else {
				io.WriteString(conn,"ERR_CMD_ERR\r\n")
			}		
		case "read":
			if len(fs)==2 {
				key := fs[1]
				if len(key)<=250 {
					mux.RLock()
					if filerepo[key].Exptime!=0 && filerepo[key].Exptime < time.Now().Unix() {
						mux.RUnlock()
						io.WriteString(conn,"ERR_FILE_NOT_FOUND\r\n")
					} else {		
						mux.RLock()
					if val, ok := filerepo[key]; ok {
					mux.RUnlock()	
					if val.Exptime!=0 {
						io.WriteString(conn, "CONTENTS "+strconv.FormatUint(val.Version,10)+" "+strconv.FormatUint(val.Numbytes,10)+" "+strconv.FormatInt(filerepo[key].Exptime - time.Now().Unix(),10)+"\r\n"+string(filerepo[key].Content)+"\r\n")
					} else {
						io.WriteString(conn, "CONTENTS "+strconv.FormatUint(val.Version,10)+" "+strconv.FormatUint(val.Numbytes,10)+"\r\n"+string(filerepo[key].Content)+"\r\n")	
					}
					} else {
						io.WriteString(conn, "ERR_FILE_NOT_FOUND\r\n")
					}
					}
					} else {
						io.WriteString(conn,"ERR_INTERNAL\r\n")
					}	
					
			} else {
				io.WriteString(conn,"ERR_CMD_ERR\r\n")
			}
		case "cas":
			if len(fs)<=5 {
				key := fs[1]
				var expt int64
				if len(fs)==5 {
					i,_ := strconv.ParseInt(fs[4],10,64)
					expt = time.Now().Unix() + i
				}	
				
				if len(key)<=250 {
				mux.Lock()	
				if filerepo[key].Exptime!=0 && filerepo[key].Exptime < time.Now().Unix() {
				mux.Unlock()		
						io.WriteString(conn,"ERR_FILE_NOT_FOUND\r\n")
					} else {	
				nb, _ := strconv.ParseUint(fs[3],10,64)
				data := make([]byte,nb)
				_,_ = io.ReadFull(nr, data)
				mux.Lock()
				if val, ok := filerepo[key]; ok {
				mux.Unlock()	
					version := strconv.FormatUint(val.Version,10)
					if version == fs[2] {
						numbytes1,_ := strconv.ParseUint(fs[3],10,64)
						mux.Lock()
						filerepo[key] = File{Numbytes: numbytes1, Version: val.Version+1, Exptime: expt, Content: data}
						mux.Unlock()
						io.WriteString(conn, "OK "+strconv.FormatUint(val.Version+1,10)+"\r\n")
				} else {
					io.WriteString(conn, "ERR_VERSION\r\n")
				}    
				} else {
					io.WriteString(conn,"ERR_FILE_NOT_FOUND\r\n")
				}
				}
			} else {
				io.WriteString(conn, "ERR_INTERNAL\r\n")
			}	
			} else {
				io.WriteString(conn,"ERR_CMD_ERR\r\n")
			}
		case "delete":
			if len(fs)==2 {
				key := fs[1]
				if len(key)<=250 {
					mux.Lock()
					if filerepo[key].Exptime!=0 && filerepo[key].Exptime < time.Now().Unix() {
						io.WriteString(conn,"ERR_FILE_NOT_FOUND\r\n")
					} else {
				if _, ok := filerepo[key]; ok {
					delete(filerepo, key)
					mux.Unlock()
					io.WriteString(conn, "OK\r\n")
				} else {
					io.WriteString(conn, "ERR_FILE_NOT_FOUND\r\n")
				}
				}
				} else {
					io.WriteString(conn,"ERR_INTERNAL\r\n")
				}
				} else {
					io.WriteString(conn,"ERR_CMD_ERR\r\n")
			}	
		default:
			io.WriteString(conn, "ERR_CMD_ERR\r\n")
				
		}
		}	
		
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
		go handle(conn)
	   

		//conn.Close()

	}
}	

func main() {
	serverMain()
}
