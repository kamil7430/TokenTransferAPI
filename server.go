package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kamil7430/TokenTransferAPI/graph"
	"github.com/kamil7430/TokenTransferAPI/graph/model"
	"github.com/kamil7430/TokenTransferAPI/repository"
	"github.com/kamil7430/TokenTransferAPI/service"
	"github.com/vektah/gqlparser/v2/ast"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const port = "8080"

func fatalIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	dbUser := os.Getenv("POSTGRES_USER")
	passwordFile, err := os.ReadFile(os.Getenv("POSTGRES_PASSWORD_FILE"))
	fatalIfError(err)
	dbPassword := strings.TrimSpace(string(passwordFile))
	dbDb := os.Getenv("POSTGRES_DB")
	dbPort := os.Getenv("POSTGRES_DB_PORT")

	dsn := fmt.Sprintf("host=db user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbUser, dbPassword, dbDb, dbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	fatalIfError(err)

	err = db.AutoMigrate(&model.Wallet{})
	fatalIfError(err)

	repo := &repository.DatabaseWalletRepository{}

	// Add initial wallet with 1 000 000 tokens (this will fail if wallet exists)
	err = db.Where(model.Wallet{Address: "0x0000000000000000000000000000000000000000"}).
		FirstOrCreate(&model.Wallet{
			Address: "0x0000000000000000000000000000000000000000",
			Tokens:  1_000_000,
		}).Error
	fatalIfError(err)

	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{
			WalletService: &service.WalletService{
				WalletRepository: repo,
				Database:         db,
			},
		},
	}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
