package main

import(
  "net"
  "io"
)

func handle(src net.Conn, whereTo string){
  dst, err := net.Dial("tcp", whereTo)
  if err != nil {
    panic("Unable to connect to specified host")
  }

  defer dst.Close()

  //run in seperate thread to prevent blocking
  go func(){
    //Copy src output to DeSTination
    if _, err := io.Copy(dst, src); err != nil {
      panic(err)
    }
  }()

  //Copy dst output to source connection
  if _, err := io.Copy(src, dst); err != nil {
    panic(err)
  }
}

func main(){
  //Listen on local port 80
  listener, err := net.Listen("tcp", ":80")
  if err != nil {
    panic("unable to bind to port")
  }

  for{
    //Accept blocks until it receives a connection
    conn, err := listener.Accept()
    if err != nil {
      panic("unable to Accept connection")
    }

    go handle(conn, "192.168.1.254:80")
  }
}
