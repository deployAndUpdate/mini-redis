package server

import (
	"bufio"
	"log"
	"mini-redis/internal/store"
	"net"
	"strings"
)

type Server struct {
	store *store.Store
}

func New(s *store.Store) *Server {
	return &Server{store: s}
}

func (srv *Server) Start(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Println("Server started at", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept err:", err)
			continue
		}
		go srv.handleConn(conn)
	}
}

func (srv *Server) handleConn(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		args := strings.Split(line, " ")

		switch strings.ToUpper(args[0]) {
		case "SET":
			if len(args) != 3 {
				conn.Write([]byte("ERR wrong number of arguments\n"))
				continue
			}
			srv.store.Set(args[1], args[2])
			conn.Write([]byte("OK\n"))

		case "GET":
			if len(args) != 2 {
				conn.Write([]byte("ERR wrong number of arguments\n"))
				continue
			}
			if val, ok := srv.store.Get(args[1]); ok {
				conn.Write([]byte(val + "\n"))
			} else {
				conn.Write([]byte("(nil)\n"))
			}

		case "DEL":
			if len(args) != 2 {
				conn.Write([]byte("ERR wrong number of arguments\n"))
				continue
			}
			srv.store.Del(args[1])
			conn.Write([]byte("OK\n"))

		case "SETTOCHAN":
			if len(args) != 4 {
				conn.Write([]byte("ERR wrong number of arguments\n"))
				continue
			}
			srv.store.SetToChan(args[1], args[2], args[3])
			conn.Write([]byte("OK\n"))

		case "READFROMCHAN":
			if len(args) != 2 {
				conn.Write([]byte("ERR wrong number of arguments\n"))
				continue
			}
			srv.store.ReadFromChan(args[1])
			conn.Write([]byte("OK\n"))

		case "CLOSECHAN":
			if len(args) != 2 {
				conn.Write([]byte("ERR wrong number of arguments\n"))
				continue
			}
			srv.store.CloseChan(args[1])
			conn.Write([]byte("OK\n"))

		case "CREATECHAN":
			if len(args) != 2 {
				conn.Write([]byte("ERR wrong number of arguments\n"))
				continue
			}
			srv.store.CreateChan(args[1])
			conn.Write([]byte("OK\n"))

		default:
			conn.Write([]byte("ERR unknown command\n"))
		}
	}
}
