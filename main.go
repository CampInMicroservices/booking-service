package main

import (
	"booking-service/api"
	"booking-service/config"
	"booking-service/db"
	"booking-service/proto"
	"context"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//	@title			CampIn Booking Service API
//	@version		1.0
//	@description	This is a booking service server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Šimen Ravnik
//	@contact.email	sr8905@student.uni-lj.si

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		20.13.80.52
// @BasePath	booking-service/v1
func main() {

	// Load configuration settings
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Failed to load config: ", err)
	}

	// Connect to the database
	store, err := db.Connect(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Run DB migration
	runDBMigration(config.MigrationURL, config.DBSource)

	// gRPC timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	// Open connection to gRPC
	conn, err := grpc.DialContext(ctx, config.PaymentServiceGRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Failed to dial gRPC server: ", err)
	}
	defer conn.Close()

	grpcClient := proto.NewPaymentServiceClient(conn)

	// Create a server and setup routes
	server, err := api.NewServer(config, store, grpcClient)
	if err != nil {
		log.Fatal("Failed to create a server: ", err)
	}

	// Start a server
	if err := server.Start(config.ServerAddress); err != nil {
		log.Fatal("Failed to start a server: ", err)
	}
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("Cannot create new migrate instance", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to run migrate up", err)
	}

	log.Println("Db migrated successfully")
}
