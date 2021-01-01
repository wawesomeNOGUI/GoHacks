package main

import (
	"fmt"
  "strconv"
	"net"
	"sort"
  "strings"
  "bufio"
  "os"
  "io"
  "time"
)

func worker(ports, results chan int) {
	for p := range ports {
		//address := fmt.Sprintf("scanme.nmap.org:%d", p)
    address := fmt.Sprintf("162.200.58.171:%d", p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
      fmt.Printf("%d closed\n", p)
			continue
		}
		conn.Close()
		results <- p
    fmt.Printf("%d open\n", p)
	}
}

// MustReadStdin blocks until input is received from stdin
func MustReadStdin() string {
  r := bufio.NewReader(os.Stdin)

  var in string
  for {
    var err error
    in, err = r.ReadString('\n')
    if err != io.EOF {
      if err != nil {
        panic(err)
      }
    }
    in = strings.TrimSpace(in)
    if len(in) > 0 {
      break
    }
  }

  fmt.Println("")

  return in
}

func main() {
	ports := make(chan int, 200)
	results := make(chan int)
	var openports []int

  fmt.Println("How many ports do you want to scan? 65535 max, 0-1024 reserved: ")
  var number2Scan, err = strconv.Atoi( MustReadStdin() )
    if err != nil {
      panic(err)
    }

  //For timing the program
  start := time.Now()

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	go func() {
		for i := 1; i <= number2Scan; i++ {
			ports <- i
		}
	}()

	for i := 0; i < number2Scan; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports)
	close(results)
	sort.Ints(openports)
  fmt.Printf("\n\n\n\n")

	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}

  fmt.Printf("Open Ports: " + strconv.Itoa( len(openports) ) + " Out of: " + strconv.Itoa( number2Scan ) + " ports scanned")

  //End timer and print result
  t := time.Now()
  elapsed := t.Sub(start)
  fmt.Printf("\n")
  fmt.Println(elapsed)
}
