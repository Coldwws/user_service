package app

import (
	"context"
	"log"
	apiUser "user_service/internal/api/user"
	"user_service/internal/client/db"
	"user_service/internal/client/db/pg"
	"user_service/internal/client/db/transaction"
	"user_service/internal/closer"
	"user_service/internal/config"
	"user_service/internal/repository"
	userRepo "user_service/internal/repository/user"
	"user_service/internal/service"
	userService "user_service/internal/service/user"

	"github.com/jackc/pgx/v4/pgxpool"
)

type serviceProvider struct {
	config *config.Config
	pgPool *pgxpool.Pool

	dbClient  db.Client
	txManager db.TxManager

	userRepository repository.UserRepository
	userService    service.UserService

	userAPI *apiUser.Server
}

func NewServiceProvider(cfg *config.Config) *serviceProvider {
	return &serviceProvider{
		config: cfg,
	}
}

var ctx = context.Background()

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.config == nil {
		log.Fatal("config is nil")
	}
	return s.config.PG
}

func (s *serviceProvider) PGPool() *pgxpool.Pool {
	if s.pgPool == nil {
		pool, err := pgxpool.Connect(ctx, s.PGConfig().DSN())

		if err != nil {
			log.Fatalf("Failed to connect database: %v", err)
		}

		err = pool.Ping(ctx)
		if err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}
		closer.Add(func() error {

			pool.Close()
			return nil
		})

		s.pgPool = pool

	}
	return s.pgPool
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("Failed to init pg client %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %v", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}
	return s.dbClient

}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager == nil {
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	}
	return s.txManager
}

func (s *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepo.NewRepository(s.DBClient(ctx))
	}
	return s.userRepository
}

func (s *serviceProvider) UserService(ctx context.Context) service.UserService {
	if s.userService == nil {
		s.userService = userService.NewUserService(s.UserRepository(ctx), s.TxManager(ctx))
	}
	return s.userService
}

func (s *serviceProvider) UserAPI() *apiUser.Server {

	if s.userAPI == nil {
		s.userAPI = apiUser.NewServer(s.UserService(ctx))
	}

	return s.userAPI
}
