package main

import ("log"; "net")

type Client struct {
	conn net.Conn
	ch chan<- string
}

func main () {
	ln, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatal(err)
	}
	msgchan := make(chan string)
	addchan := make(chan Client)
	rmchan := make(chan Client)
	go printMessages(msgchan)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleMessages(msgchan, addchan, rmchan)
		go handleConnection(conn, msgchan, addchan, rmchan)
	}
}

func handleConnection(c net.Conn, msgchan chan<- string, addchan chan<- Client, rmchan chan<- Client) {
	buf := make([]byte, 4096)
	ch := make(chan string)
	addchan <- Client{c, ch}

	defer func() {
		rmchan <- Client{c, ch}
	}()

	for {
		n, err := c.Read(buf)
		if err != nil || n == 0 {
			c.Close()
			break;
		}
		//n, err = c.Write(buf[0:n])
		msgchan <- string(buf[0:n])
		if err != nil {
			c.Close()
			break;
		}
	}
	log.Printf("Connection from %v closed.", c.RemoteAddr())
}

func handleMessages(msgchan <-chan string, addchan <-chan Client, rmchan <-chan net.Conn) {
	clients := make(map[net.Conn]chan<- string)
	for {
		select{
			case msg := <-msgchan:
				for _, ch := range clients {
					go func(mch chan<- string) {mch <- "\033[1;33;40m" + msg + "\033[m\r\n" }(ch)
				}
			case client := <-addchan:
				clients[client.conn] = client.ch
			case conn := <-rmchan:
				delete(clients, conn.ch)
		}
	}
}

func printMessages(msgchan <-chan string) {
	for msg := range msgchan {
		log.Printf("new message: %s", msg)
	}
}
