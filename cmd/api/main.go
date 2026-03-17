package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/Teryn-Guzman/Lab-3/internal/data"
	_ "github.com/lib/pq"
)
const appVersion = "1.0.0"

type serverConfig struct {
    port int 
    environment string
    db struct {
        dsn string
    }

     limiter struct {
        rps float64                      // requests per second
        burst int                        // initial requests possible
        enabled bool                     // enable or disable rate limiter
    }

    cors struct {
        trustedOrigins []string
    }


}

type applicationDependencies struct {
    config serverConfig
    logger *slog.Logger
    customerModel data.CustomerModel
    tableModel data.TableModel
    shiftModel data.ShiftModel
    timeSlotModel data.TimeSlotModel
    reservationModel data.ReservationModel
    reservationTableAssignmentModel data.ReservationTableAssignmentModel
    specialRequestModel data.SpecialRequestModel
    waitstaffModel data.WaitstaffModel
    shiftTableAssignmentModel data.ShiftTableAssignmentModel
}
func main() {
    var settings serverConfig

    flag.IntVar(&settings.port, "port", 4000, "Server port")
    flag.StringVar(&settings.environment, "env", "development",
                  "Environment(development|staging|production)")
    flag.StringVar(&settings.db.dsn, "db-dsn", "postgres://restaurant:restaurant@localhost/restaurant_management_db","PostgreSQL DSN")

    flag.Float64Var(&settings.limiter.rps, "limiter-rps", 2,
                  "Rate Limiter maximum requests per second")

    flag.IntVar(&settings.limiter.burst, "limiter-burst", 5,
                  "Rate Limiter maximum burst")

    flag.BoolVar(&settings.limiter.enabled, "limiter-enabled", true,
                  "Enable rate limiter")

    flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)",
              func(val string) error {
                   settings.cors.trustedOrigins = strings.Fields(val)
                   return nil
              })


    flag.Parse()
    
    logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

    db, err := openDB(settings)
    if err != nil {
    logger.Error(err.Error())
    os.Exit(1)
    }
    defer db.Close()

    logger.Info("database connection pool established")

    expvar.NewString("version").Set(appVersion)

    // the number of active goroutines
    expvar.Publish("goroutines", expvar.Func(func() any {
        return runtime.NumGoroutine()
    }))

    // the database connection pool metrics
    expvar.Publish("database", expvar.Func(func() any {
        return runtime.NumGoroutine()
    }))

   // the current Unix timestamp
   expvar.Publish("timestamp", expvar.Func(func() any {
        return time.Now().Unix()
   }))

    appInstance := &applicationDependencies {
        config: settings,
        logger: logger,
        customerModel: data.CustomerModel{DB: db},
        tableModel: data.TableModel{DB: db},
        shiftModel: data.ShiftModel{DB: db},
        timeSlotModel: data.TimeSlotModel{DB: db},
        reservationModel: data.ReservationModel{DB: db},
        reservationTableAssignmentModel: data.ReservationTableAssignmentModel{DB: db},
        specialRequestModel: data.SpecialRequestModel{DB: db},
        waitstaffModel: data.WaitstaffModel{DB: db},
        shiftTableAssignmentModel: data.ShiftTableAssignmentModel{DB: db},
    }


    err = appInstance.serve()
    if err != nil {
        logger.Error(err.Error())
        os.Exit(1)
    }
}
 
    func openDB(settings serverConfig) (*sql.DB, error) {

        db, err := sql.Open("postgres", settings.db.dsn)
    if err != nil {
        return nil, err
    }
    
    ctx, cancel := context.WithTimeout(context.Background(),
                                       5 * time.Second)
    defer cancel()
    err = db.PingContext(ctx)
    if err != nil {
        db.Close()
        return nil, err
    }

    return db, nil

} 
