package main

import(
    "fmt"
    "net"
    "flag"
    "strconv"
    "strings"
    "bufio"
    "regexp"
    "log"
    "os"
    )


func main (){

    port := flag.Int("port",-1,"Port to listen on")
    threadCount := flag.Int("threadCount", 100, "Available thread count")

    flag.Parse()

    if *port<0 {
        log.Fatal("Must pass port via -port x flag")
    }

    fmt.Println("Server Running on port "+strconv.Itoa(*port))

    sharedChan := make(chan *net.TCPConn, *threadCount)
    killChan := make(chan int)
    tcpChan := make(chan *net.TCPConn)

    fmt.Println("Firing up "+strconv.Itoa(*threadCount)+" goroutines")

    tcpAddr, _ := net.ResolveTCPAddr("tcp", os.Getenv("IP_ADDRESS")+":"+strconv.Itoa(*port))
    tcpListener, err := net.ListenTCP("tcp", tcpAddr)

    if err != nil {
        log.Fatal(err)
    }

    go func(){
        for {
            tcpConn, _ := tcpListener.AcceptTCP()
            tcpChan <- tcpConn
        }
    }()
    for i := 0; i < *threadCount; i++ {
        go connectionHandler(sharedChan,killChan)
    }

    for {
        select {
            case tcpConn := <- tcpChan:
                select {
                    case sharedChan <- tcpConn:
                    default:
                }
            case <- killChan:
                return
        }
    }
}

func connectionHandler(sharedChan chan *net.TCPConn, killChan chan int) {
    for {
        tcpConn := <- sharedChan
        message, _ := bufio.NewReader(tcpConn).ReadString('\n')
        message = strings.TrimSpace(message)

        if len(message) > 5 && message[:5] == "HELO " {
            rgx, _ := regexp.Compile("^(.*):(\\d+)$")
            addr := rgx.FindStringSubmatch(tcpConn.LocalAddr().String())
            tcpConn.Write([]byte(message[5:]+"\nIP:"+addr[1]+"\nPort:"+addr[2]+"\nStudentID:08506426\n"))
        }

        if message == "KILL_SERVICE" {
            killChan <- 1
        }
        tcpConn.Close()
    }
}
