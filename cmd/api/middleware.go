package main

import (
	"expvar"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// metricsResponseWriter is a custom http.ResponseWriter that allows us to capture the status code of the response and whether the headers have been written or not.
type metricsResponseWriter struct {
    wrapped    http.ResponseWriter   // wrapped is the original http.ResponseWriter that we are wrapping
    statusCode int         // the status code of the response
    headerWritten bool    // headerWritten is a boolean that indicates whether the headers have been written or not
}


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

func (a *applicationDependencies) enableCORS (next http.Handler) http.Handler {                             
   return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Add the Vary header to prevent caching of the CORS response
        w.Header().Add("Vary", "Origin")
        w.Header().Add("Vary", "Access-Control-Request-Method")

        origin := r.Header.Get("Origin")

        // Check if the Origin header is present and if it matches any of the trusted origins
        if origin != "" {
            for i := range a.config.cors.trustedOrigins {
               if origin == a.config.cors.trustedOrigins[i] {
                 w.Header().Set("Access-Control-Allow-Origin", origin)  

                 //check if the request is a preflight request
                 if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
                    w.Header().Set("Access-Control-Allow-Methods",
                             "OPTIONS, PUT, PATCH, DELETE")
                    w.Header().Set("Access-Control-Allow-Headers",
                             "Authorization, Content-Type")
        
                    // For preflight requests, we respond with a 200 OK status and return without calling the next handler
                    w.WriteHeader(http.StatusOK)
                    return
                } 

                                        
                  break
        }
   }
}


        next.ServeHTTP(w, r)
    })
}

func (a *applicationDependencies) metrics (next http.Handler) http.Handler {                             
   // Setup our variable to track the metrics
 var (
      totalRequestsReceived = expvar.NewInt("total_requests_received")
      totalResponsesSent    = expvar.NewInt("total_responses_sent")
      totalProcessingTimeMicroseconds = expvar.NewInt("total_processing_time_μs")
      totalResponsesSentByStatus = expvar.NewMap("total_responses_sent_by_status")
 )
 
 return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    // start is when we receive the request and start processing it
        start := time.Now()
        // update our request received counter
        totalRequestsReceived.Add(1)

        mw := newMetricsResponseWriter(w)

        next.ServeHTTP(mw, r)

        // after the next handler has finished processing the request, we update our response sent counter and our total processing time counter
        totalResponsesSent.Add(1)

        totalResponsesSentByStatus.Add(strconv.Itoa(mw.statusCode), 1)

        // duration is the time it took to process the request in microseconds
        duration := time.Since(start).Microseconds()
        
        // update our total processing time counter
        totalProcessingTimeMicroseconds.Add(duration)
    })
}

func newMetricsResponseWriter(w http.ResponseWriter) *metricsResponseWriter {
    return &metricsResponseWriter {
        wrapped: w,
        statusCode: http.StatusOK,
    }
}

func (mw *metricsResponseWriter) Header() http.Header {
    return mw.wrapped.Header()
}

func (mw *metricsResponseWriter) WriteHeader(statusCode int) {
    mw.wrapped.WriteHeader(statusCode)
if !mw.headerWritten {
        mw.statusCode = statusCode
        mw.headerWritten = true
    }
}

func (mw *metricsResponseWriter) Write(b []byte) (int, error) {
    mw.headerWritten = true
    return mw.wrapped.Write(b)
}

func (mw *metricsResponseWriter) Unwrap() http.ResponseWriter {
    return mw.wrapped
}
