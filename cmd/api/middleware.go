package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func (a *applicationDependencies)recoverPanic(next http.Handler) http.Handler  {
   return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
       defer func() {
           err := recover();
           if err != nil {
               w.Header().Set("Connection", "close")
               a.serverErrorResponse(w, r, fmt.Errorf("%s", err))
           }
       }()
       next.ServeHTTP(w,r)
   })  
}

func (a *applicationDependencies)rateLimit(next http.Handler) http.Handler {
// Define a rate limiter struct
    type client struct {
        limiter *rate.Limiter
        lastSeen  time.Time  
}
var mu sync.Mutex       
  var clients = make(map[string]*client)  

  go func() {
      for {
          time.Sleep(time.Minute)
          mu.Lock() 

          // delete any entry not seen in three minutes
          for ip, client := range clients {
              if time.Since(client.lastSeen) > 3 * time.Minute {
                  delete(clients, ip)
              }
          }
        mu.Unlock()   
        }
 }()
return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

    if a.config.limiter.enabled {

 // get the IP address
     ip, _, err := net.SplitHostPort(r.RemoteAddr)
     if err != nil {
         a.serverErrorResponse(w, r, err)
         return
     }

     mu.Lock() 
     // check if ip address already in map, if not add it
     _, found := clients[ip]
    if !found {
        clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(a.config.limiter.rps),
                                        a.config.limiter.burst)}
    }
// Update the last seem for the client
 clients[ip].lastSeen = time.Now()

// Check the rate limit status
 if !clients[ip].limiter.Allow() {
     mu.Unlock()    
     a.rateLimitExceededResponse(w, r)
     return
 }

 mu.Unlock()   
 }
 next.ServeHTTP(w, r)
})

}
