package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/zmeditation/hackernews/graph"
	"github.com/zmeditation/hackernews/graph/generated"

	"database/sql"
    
    _ "github.com/go-sql-driver/mysql"
    "github.com/golang-migrate/migrate/v4"
    "github.com/golang-migrate/migrate/v4/database/mysql"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db, _ := sql.Open("mysql", "root:@tcp(host:port)/dbname?multiStatements=true")
    driver, _ := mysql.WithInstance(db, &mysql.Config{})
    m, _ := migrate.NewWithDatabaseInstance(
        "file:///migrations",
        "mysql", 
        driver,
    )
    
    m.Steps(2)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
