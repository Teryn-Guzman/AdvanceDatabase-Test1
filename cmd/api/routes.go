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

    // Health check route
    router.HandlerFunc(http.MethodGet, "/v1/health", a.healthCheckHandler)
    router.HandlerFunc(http.MethodPost, "/v1/health", a.healthCheckHandler)
   
   // Customers routes
   router.HandlerFunc(http.MethodPost, "/v1/customers", a.createCustomerHandler)
   router.HandlerFunc(http.MethodGet, "/v1/customers/:id", a.displayCustomerHandler)
   router.HandlerFunc(http.MethodPatch,"/v1/customers/:id", a.updateCustomerHandler)
   router.HandlerFunc(http.MethodDelete,"/v1/customers/:id", a.deleteCustomerHandler)
   router.HandlerFunc(http.MethodGet,"/v1/customers", a.listCustomersHandler)
   
   // Tables routes
   router.HandlerFunc(http.MethodPost, "/v1/tables", a.createTableHandler)
   router.HandlerFunc(http.MethodGet, "/v1/tables/:id", a.displayTableHandler)
   router.HandlerFunc(http.MethodPatch, "/v1/tables/:id", a.updateTableHandler)
   router.HandlerFunc(http.MethodDelete, "/v1/tables/:id", a.deleteTableHandler)
   router.HandlerFunc(http.MethodGet, "/v1/tables", a.listTablesHandler)
   
   // Shifts routes
   router.HandlerFunc(http.MethodPost, "/v1/shifts", a.createShiftHandler)
   router.HandlerFunc(http.MethodGet, "/v1/shifts/:id", a.displayShiftHandler)
   router.HandlerFunc(http.MethodPatch, "/v1/shifts/:id", a.updateShiftHandler)
   router.HandlerFunc(http.MethodDelete, "/v1/shifts/:id", a.deleteShiftHandler)
   router.HandlerFunc(http.MethodGet, "/v1/shifts", a.listShiftsHandler)
   
   // Time Slots routes
   router.HandlerFunc(http.MethodPost, "/v1/time-slots", a.createTimeSlotHandler)
   router.HandlerFunc(http.MethodGet, "/v1/time-slots/:id", a.displayTimeSlotHandler)
   router.HandlerFunc(http.MethodPatch, "/v1/time-slots/:id", a.updateTimeSlotHandler)
   router.HandlerFunc(http.MethodDelete, "/v1/time-slots/:id", a.deleteTimeSlotHandler)
   router.HandlerFunc(http.MethodGet, "/v1/time-slots", a.listTimeSlotsHandler)
   
   // Reservations routes
   router.HandlerFunc(http.MethodPost, "/v1/reservations", a.createReservationHandler)
   router.HandlerFunc(http.MethodGet, "/v1/reservations/:id", a.displayReservationHandler)
   router.HandlerFunc(http.MethodPatch, "/v1/reservations/:id", a.updateReservationHandler)
   router.HandlerFunc(http.MethodDelete, "/v1/reservations/:id", a.deleteReservationHandler)
   router.HandlerFunc(http.MethodGet, "/v1/reservations", a.listReservationsHandler)
   
   // Reservation Table Assignments routes
   router.HandlerFunc(http.MethodPost, "/v1/reservation-table-assignments", a.createReservationTableAssignmentHandler)
   router.HandlerFunc(http.MethodGet, "/v1/reservation-table-assignments", a.displayReservationTableAssignmentHandler)
   router.HandlerFunc(http.MethodDelete, "/v1/reservation-table-assignments", a.deleteReservationTableAssignmentHandler)
   router.HandlerFunc(http.MethodGet, "/v1/reservation-table-assignments/list", a.listReservationTableAssignmentsHandler)
   
   // Special Requests routes
   router.HandlerFunc(http.MethodPost, "/v1/special-requests", a.createSpecialRequestHandler)
   router.HandlerFunc(http.MethodGet, "/v1/special-requests/:id", a.displaySpecialRequestHandler)
   router.HandlerFunc(http.MethodPatch, "/v1/special-requests/:id", a.updateSpecialRequestHandler)
   router.HandlerFunc(http.MethodDelete, "/v1/special-requests/:id", a.deleteSpecialRequestHandler)
   router.HandlerFunc(http.MethodGet, "/v1/special-requests", a.listSpecialRequestsHandler)
   
   // Waitstaff routes
   router.HandlerFunc(http.MethodPost, "/v1/waitstaff", a.createWaitstaffHandler)
   router.HandlerFunc(http.MethodGet, "/v1/waitstaff/:id", a.displayWaitstaffHandler)
   router.HandlerFunc(http.MethodPatch, "/v1/waitstaff/:id", a.updateWaitstaffHandler)
   router.HandlerFunc(http.MethodDelete, "/v1/waitstaff/:id", a.deleteWaitstaffHandler)
   router.HandlerFunc(http.MethodGet, "/v1/waitstaff", a.listWaitstaffHandler)
   
   // Shift Table Assignments routes
   router.HandlerFunc(http.MethodPost, "/v1/shift-table-assignments", a.createShiftTableAssignmentHandler)
   router.HandlerFunc(http.MethodGet, "/v1/shift-table-assignments", a.displayShiftTableAssignmentHandler)
   router.HandlerFunc(http.MethodPatch, "/v1/shift-table-assignments", a.updateShiftTableAssignmentHandler)
   router.HandlerFunc(http.MethodDelete, "/v1/shift-table-assignments", a.deleteShiftTableAssignmentHandler)
   router.HandlerFunc(http.MethodGet, "/v1/shift-table-assignments/list", a.listShiftTableAssignmentsHandler)
   
    router.Handler(http.MethodGet,"/v1/metrics", expvar.Handler())

     // Request sent first to recoverPanic() then sent to rateLimit()
     // finally it is sent to the router.
     // Compression is applied last (outermost) so it compresses the final response
     return  a.metrics(a.recoverPanic(a.enableCORS(a.rateLimit(a.compress(router)))))
}