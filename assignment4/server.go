package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/MJVegad/cs733/assignment4/fs"
	"net"
	"os"
	"strconv"
	//"strings"
)

var crlf = []byte{'\r', '\n'}

type RaftNodeMsg struct {
	Data Rnmsg
}

func encode(data Rnmsg) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	le := RaftNodeMsg{Data: data}
	err := enc.Encode(le)
	return buf.Bytes(), err
}

func decode(dbytes []byte) (Rnmsg, error) {
	//fmt.Printf("dbytes->%v\n",dbytes)
	buf := bytes.NewBuffer(dbytes)
	enc := gob.NewDecoder(buf)
	var le RaftNodeMsg
	err := enc.Decode(&le)
	//fmt.Printf("le.data->%v\n",le.Data)
	return le.Data, err
}

/*func decode(dbytes []byte) (interface{}, error) {
	buf := bytes.NewBuffer(dbytes)
	enc := gob.NewDecoder(buf)
	var le LogEntry
	err := enc.Decode(&le)
	return le.Data, err
}*/

func check(obj interface{}) {
	if obj != nil {
		fmt.Println(obj)
		os.Exit(1)
	}
}

type Rnmsg struct {
	ClientId int64
	Msg      *fs.Msg
}

func reply(conn *net.TCPConn, msg *fs.Msg, paddress string) bool {
	var err error
	write := func(data []byte) {
		if err != nil {
			return
		}
		_, err = conn.Write(data)
	}
	var resp string
	switch msg.Kind {
	case 'C': // read response
		resp = fmt.Sprintf("CONTENTS %d %d %d", msg.Version, msg.Numbytes, msg.Exptime)
	case 'O':
		resp = "OK "
		if msg.Version > 0 {
			resp += strconv.Itoa(msg.Version)
		}
	case 'F':
		resp = "ERR_FILE_NOT_FOUND"
	case 'V':
		resp = "ERR_VERSION " + strconv.Itoa(msg.Version)
	case 'M':
		resp = "ERR_CMD_ERR"
	case 'I':
		resp = "ERR_INTERNAL"
	case 'P':
		resp = "ERR_REDIRECT " + paddress
	default:
		fmt.Printf("Unknown response kind '%c'", msg.Kind)
		return false
	}
	resp += "\r\n"
	write([]byte(resp))
	if msg.Kind == 'C' {
		write(msg.Contents)
		write(crlf)
	}
	return err == nil
}

func serve(conn *net.TCPConn, rn RaftNode, clientId int64, connch map[int64]chan *fs.Msg, ff *fs.FS, gversion *int) {
	reader := bufio.NewReader(conn)
	for {
		msg, msgerr, fatalerr := fs.GetMsg(reader)
		if fatalerr != nil || msgerr != nil {
			reply(conn, &fs.Msg{Kind: 'M'}, "")
			conn.Close()
			break
		}

		if msgerr != nil {
			if (!reply(conn, &fs.Msg{Kind: 'M'}, "")) {
				conn.Close()
				break
			}
		}

		//fmt.Printf("Inside serve, msg:%c\n", msg.Kind)
		if msg.Kind != 'r' {
			tempmsg := Rnmsg{clientId, msg}
			rnmsgbytes, err1 := encode(tempmsg)
			if err1 != nil {
				//fmt.Printf("encoding error:%v", err1)
				return
			}
			rn.Append(rnmsgbytes)
			//fmt.Printf("waiting for response for client %v..!!\n", clientId)
			dmsg := <-connch[clientId]
			//response := fs.ProcessMsg(dmsg,ff,gversion)
			//				fmt.Printf("Received on client channel -> %v, msg kind -> %v\n",clientId, dmsg.Kind)
			if dmsg.Kind == 'P' {
				//fmt.Printf("Voted For in serve->%v\n", rn.sm.votedFor)
				for j := 0; j < len(temp.Peers); j++ {
					if temp.Peers[j].Id == rn.sm.votedFor {
						reply(conn, dmsg, temp.Peers[j].ServerAddress)
					}
				}

			} else {
				//fmt.Printf("msg kind not P on client %v\n", clientId)
			//	response := fs.ProcessMsg(dmsg, ff, gversion)
				if !reply(conn, dmsg, "") {
					conn.Close()
					break
				}
			}

		} else {
			response := fs.ProcessMsg(msg, ff, gversion)
			if !reply(conn, response, "") {
				conn.Close()
				break
			}
		}
	}
}

func (rn *RaftNode) listenToCommitChannels(connch map[int64]chan *fs.Msg, ff *fs.FS, gversion *int) {
	for {
		//fmt.Printf("for begins...!!!\n")
		ci := <-rn.CommitChannel()
		//fmt.Printf("In serve, ci->%v\n", ci)
		rnmsgcc, err := decode(ci.Data)
		if ci.Err != nil {
			connch[rnmsgcc.ClientId] <- &fs.Msg{'P', rnmsgcc.Msg.Filename, rnmsgcc.Msg.Contents, rnmsgcc.Msg.Numbytes, rnmsgcc.Msg.Exptime, rnmsgcc.Msg.Version}
			continue
		}
		if err != nil {
		//	fmt.Printf("In decode, rnmsgcc.err->%v\n", err)
			return
		}
		
		response := fs.ProcessMsg(rnmsgcc.Msg, ff, gversion)
		
		_,ok := connch[rnmsgcc.ClientId]
		
		//fmt.Printf("Received on channel for client -> %v, msg kind -> %v\n", rnmsgcc.ClientId, rnmsgcc.Msg.Kind)
		if ok==true {
			connch[rnmsgcc.ClientId] <- response
		} 	
		//fmt.Printf("for finishes..!!\n")
	}
}

func initRaftNode(i int) (rn RaftNode) {
	initRaftStateFile("PersistentData_" + strconv.Itoa((i+1)*100))
	prepareRaftNodeConfigObj()
	if i == 1 {
		return New(Config{peers, int64((i + 1) * 100), "PersistentData_" + strconv.Itoa((i+1)*100), 4000, 800}, "config.json")
	} else {
		return New(Config{peers, int64((i + 1) * 100), "PersistentData_" + strconv.Itoa((i+1)*100), 4000, 800}, "config.json")
	}
}

func serverMain(i int) {
	rn := initRaftNode(i)
	var clientId int64 = int64((i + 1) * 100)
	var ff = &fs.FS{Dir: make(map[string]*fs.FileInfo, 1000)}
	var gversion = 0 // global version
	connch := make(map[int64]chan *fs.Msg)
	gob.Register(RaftNodeMsg{})
	go rn.processEvents()
	go rn.listenToCommitChannels(connch, ff, &gversion)
	tcpaddr, err := net.ResolveTCPAddr("tcp", temp.Peers[i].ServerAddress)
	check(err)
	tcp_acceptor, err := net.ListenTCP("tcp", tcpaddr)
	check(err)
	for {
		tcp_conn, err := tcp_acceptor.AcceptTCP()
		connendch := make(chan *fs.Msg)
		connch[clientId] = connendch
		check(err)
		go serve(tcp_conn, rn, clientId, connch, ff, &gversion)
		clientId++
	}
}

func main() {

	//fserver, err := exec.Command("./fserver")
	//fserver.Stdout = os.Stdout
	//fserver.Stdin = os.Stdin
	//fserver.Start()
	//serverindex := os.Args
	/*for i := 0 ; i < 5; i++ {
		fs[i], err := exec.Command("./fs")
		fs[i].Stdout = os.Stdout
		fs[i].Stdin = os.Stdin
		fs[i].Start()
	}*/
	serverindex := os.Args[1]
	si, _ := strconv.Atoi(serverindex)
	//fmt.Printf("Server %v started.\n", si)
	serverMain(si)
}
