package app

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/joho/godotenv"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
	"user_service/internal/closer"
	"user_service/internal/config"
	"user_service/internal/interceptor"
	desc "user_service/pkg/user_v1"
	_ "user_service/statik"
)

type App struct {
	config          *config.Config
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
	swaggerServer   *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.InitDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) InitDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
		a.initHttpServer,
		a.initSwaggerServer,
	}
	for _, f := range inits {
		if err := f(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initConfig(_ context.Context) error {
	//для локального запуска
	if err := godotenv.Load("local.env"); err != nil {
		log.Println("Warning: local.env not found, using system env")
	}

	cfg := config.LoadConfig()

	a.config = &cfg

	return nil

}

func (a *App) initServiceProvider(_ context.Context) error {

	a.serviceProvider = NewServiceProvider(a.config)

	return nil

}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()
	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			log.Printf("failed to run grpc server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runHTTPserver()
		if err != nil {
			log.Printf("failed to run http server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		err := a.runSwaggerServer()
		if err != nil {
			log.Printf("failed to run swagger server: %v", err)
		}
	}()
	wg.Wait()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(interceptor.ValidateInterceptor),
	)

	reflection.Register(a.grpcServer)

	desc.RegisterUserV1Server(a.grpcServer, a.serviceProvider.UserAPI())

	return nil

}

func (a *App) initHttpServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:    a.serviceProvider.config.Http.Address(),
		Handler: corsMiddleware.Handler(mux),
	}

	go func() {
		opts := []grpc.DialOption{
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}

		grpcAddr := a.serviceProvider.config.GRPC.Addr()
		for {
			log.Println("Trying to register HTTP Gateway...")
			err := desc.RegisterUserV1HandlerFromEndpoint(ctx, mux, grpcAddr, opts)
			if err != nil {
				log.Println("Failed to connect to gRPC, retrying in 1s:", err)
				time.Sleep(time.Second)
				continue
			}
			log.Println("HTTP Gateway registered successfully")
			break
		}
	}()

	return nil
}

func (a *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/api.swagger.json", serveSwaggerFile("/api.swagger.json"))

	a.swaggerServer = &http.Server{
		Addr:    a.serviceProvider.config.Swagger.Address(),
		Handler: mux,
	}

	return nil
}

func (a *App) runGRPCServer() error {

	log.Println("GRPC server is running on:", a.serviceProvider.config.GRPC.Addr())

	list, err := net.Listen("tcp", a.serviceProvider.config.GRPC.Addr())
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil

}

func (a *App) runHTTPserver() error {
	log.Printf("HTTP server is running on %s", a.serviceProvider.config.Http.Address())

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (a *App) runSwaggerServer() error {
	log.Printf("Swagger server is running on %s", a.serviceProvider.config.Swagger.Address())

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving Swagger file %s", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Open swagger file: %s", path)
		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()
		log.Printf("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Printf("Write swagger file: %s", path)
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Served swagger file: %s", path)

	}
}
