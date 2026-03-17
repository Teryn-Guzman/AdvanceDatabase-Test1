package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)
func (a *applicationDependencies)routes() http.Handler  {

   // setup a new router
   router := httprouter.New()
   // handle 404
   router.NotFound = http.HandlerFunc(a.notFoundResponse)
  // handle 405
   router.MethodNotAllowed = http.HandlerFunc(a.methodNotAllowedResponse)
   // setup routes
   router.HandlerFunc(http.MethodPost, "/v1/customers", a.createCustomerHandler)
   router.HandlerFunc(http.MethodGet, "/v1/customers/:id", a.displayCustomerHandler)
   router.HandlerFunc(http.MethodPatch,"/v1/customers/:id", a.updateCustomerHandler)
   router.HandlerFunc(http.MethodDelete,"/v1/customers/:id", a.deleteCustomerHandler)
   router.HandlerFunc(http.MethodGet,"/v1/customers", a.listCustomersHandler)
   router.Handler(http.MethodGet,"/v1/observability/customers/metrics", expvar.Handler())


   // Request sent first to recoverPanic() then sent to rateLimit()
    // finally it is sent to the router.
    return  a.metrics(a.recoverPanic(a.enableCORS(a.rateLimit(router))))
}