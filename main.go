package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/shouta0715/simple-bank/gapi"
	"github.com/shouta0715/simple-bank/pb"
	"github.com/shouta0715/simple-bank/util"

	"github.com/golang-migrate/migrate/v4"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	db "github.com/shouta0715/simple-bank/db/sqlc"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	_ "github.com/shouta0715/simple-bank/doc/statik"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	runDBMigrations(config.MigrationURL, config.DBSource)

	store := db.NewStore(conn)
	go runGrpcServer(config, store)
	runGatewayServer(config, store)

}

func runDBMigrations(migrationURL, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)

	if err != nil {
		log.Fatal("cannot create migration:", err)
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("cannot migrate db:", err)
	}

	log.Println("migration completed")
}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)

	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	log.Println("starting gRPC server on", listener.Addr().String())

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal("cannot start server:", err)
	}

}

func runGatewayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)

	if err != nil {
		log.Fatal("cannot create server:", err)
	}

	grpcMux := runtime.NewServeMux()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)

	if err != nil {
		log.Fatal("cannot register gateway server:", err)
	}

	mux := http.NewServeMux()

	// gRPCを受け取る
	mux.Handle("/", grpcMux)

	statikFs, err := fs.New()
	if err != nil {
		log.Fatal("cannot create statik file system:", err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFs))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	log.Println("starting HTTP gateway server on", listener.Addr().String())

	err = http.Serve(listener, mux)

	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}

// func runGinServer(config util.Config, store db.Store) {
// 	server, err := api.NewServer(config, store)

// 	if err != nil {
// 		log.Fatal("cannot create server:", err)
// 	}

// 	err = server.Start(config.HTTPServerAddress)

// 	if err != nil {
// 		log.Fatal("cannot start server:", err)
// 	}
// }
