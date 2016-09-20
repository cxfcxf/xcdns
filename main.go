package main

import (
//    "fmt"
    "log"
    "os"
    "syscall"
    "os/signal"
    "strconv"
    "github.com/miekg/dns"
)

const SOA string = "google.com. 1000 SOA ns.stackpath.net. admin.stackpath.net. 1 4294967294 4294967293 4294967295 100"

const A string = "www.google.com IN A 192.168.1.1"


func main() {

    rr, _ := dns.NewRR(SOA)
    rrx := rr.(*dns.SOA)

    ra, _ := dns.NewRR(A)
    rxa := ra.(*dns.A)


    dns.HandleFunc("www.google.com", func(w dns.ResponseWriter, r *dns.Msg) {
            m := new(dns.Msg)
            m.SetReply(r)
            m.Answer = []dns.RR{rxa}
            m.Authoritative = true
            m.Ns = []dns.RR{rrx}
            w.WriteMsg(m)
        })

    port := 8053
    protocols := []string{"tcp", "udp"}
    for _, proto := range protocols {
        go func(p string) {
            srv := &dns.Server{Addr: ":" + strconv.Itoa(port), Net: p}
            if err := srv.ListenAndServe(); err != nil {
                log.Fatalf("Failed to set tcp listener %s\n", err.Error())
            }
        }(proto)
    }

    sig := make(chan os.Signal)
    signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
    s := <- sig
    log.Fatalf("Signal (%v) received, stopping\n", s)
}


