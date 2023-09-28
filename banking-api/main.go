package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/fclebinho/codebank/src/data/grpc/pb"
	"github.com/fclebinho/codebank/src/data/services"
	"github.com/fclebinho/codebank/src/domain/usecases"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
}

func main() {
	db := setupDatabase()
	defer db.Close()

	println("Running gRPC server")

	listener, err := net.Listen(os.Getenv("API_NETWORK"), os.Getenv("API_ADDRESS"))
	if err != nil {
		panic(err)
	}

	usecase := setupTransactionUseCase(db)
	server := setupGrpcServer(usecase)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed ro serve: %v", err)
	}

}

func setupTransactionUseCase(db *sql.DB) usecases.ProcessTranscationUsecase {
	service := services.NewPostgresService(db)
	usecase := usecases.NewProcessTranscationUsecase(service)
	// useCase.KafkaProducer = producer
	return usecase
}

func setupDatabase() *sql.DB {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("error connection database:" + err.Error())
	}

	return db
}

func setupGrpcServer(usecase usecases.ProcessTranscationUsecase) *grpc.Server {
	grpcServer := grpc.NewServer()
	pb.RegisterPaymentServiceServer(grpcServer, &services.TransactionServer{ProcessTranscationUsecase: usecase})
	reflection.Register(grpcServer)

	return grpcServer
}

// creditCard := entities.NewCreditCard()
// creditCard.Name = in.GetCreditCard().GetName()
// creditCard.Number = in.GetCreditCard().GetNumber()
// creditCard.ExpirationMonth = in.GetCreditCard().GetExpirationMonth()
// creditCard.ExpirationYear = in.GetCreditCard().GetExpirationYear()
// creditCard.CVV = in.GetCreditCard().GetCvv()

// err := s.ProcessTranscationUsecase.TransactionService.CreateCreditCard(*creditCard)
// if err != nil {
// 	return nil, err
// }
