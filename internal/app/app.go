package app

import (
	"context"
	"log"
	"net"
	"user_service/internal/config"
	desc "user_service/pkg/user_v1"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type App struct {
	config *config.Config
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

func NewApp(ctx context.Context) (*App,error){
	a := &App{}

	err := a.InitDeps(ctx)
	if err != nil {
		return nil,err
	}

	return a, nil
}



func (a *App)InitDeps(ctx context.Context)error{
	inits := []func(context.Context)error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}
	for _, f := range inits{
		if err := f(ctx); err != nil{
			return err
		}
	}
	return nil
}

func (a *App)initConfig(_ context.Context)error{
	//для локального запуска
	if err := godotenv.Load("local.env"); err != nil {
      log.Println("Warning: local.env not found, using system env")
  }


	cfg := config.LoadConfig()

	a.config = &cfg
	
	return nil

}


func(a *App)initServiceProvider(_ context.Context)error{

		a.serviceProvider = NewServiceProvider(a.config)

		return nil

}


func(a *App)initGRPCServer(ctx context.Context)error{
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	reflection.Register(a.grpcServer)

	desc.RegisterUserV1Server(a.grpcServer,a.serviceProvider.UserAPI())

	return nil

}

func(a *App)Run()error{
	return a.runGRPCServer()
}

func (a *App)runGRPCServer()error{

	log.Println("GRPC server is running on:",a.serviceProvider.config.GRPC.Addr())

	list,err := net.Listen("tcp",a.serviceProvider.config.GRPC.Addr())
	if err != nil{
		return err}
	
	err = a.grpcServer.Serve(list)
	if err!= nil{
		return err
	}

	return nil

}