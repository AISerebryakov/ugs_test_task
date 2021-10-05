package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pretcat/ugc_test_task/repositories"

	"github.com/pretcat/ugc_test_task/config"
	"github.com/pretcat/ugc_test_task/http"
	"github.com/pretcat/ugc_test_task/logger"
	buildmng "github.com/pretcat/ugc_test_task/managers/buildings"
	categmng "github.com/pretcat/ugc_test_task/managers/categories"
	companmng "github.com/pretcat/ugc_test_task/managers/companies"
	"github.com/pretcat/ugc_test_task/pg"
	buildrepos "github.com/pretcat/ugc_test_task/repositories/buildings"
	categrepos "github.com/pretcat/ugc_test_task/repositories/categories"
	companrepos "github.com/pretcat/ugc_test_task/repositories/companies"
)

var (
	conf config.Config

	pgClient pg.Client

	categoryRepos categrepos.Repository
	companyRepos  companrepos.Repository
	buildingRepos buildrepos.Repository

	companyMng  companmng.Manager
	buildingMng buildmng.Manager
	categoryMng categmng.Manager

	httpApi *http.Api

	shutdownServiceTimeout = 5 * time.Second
)

func main() {
	var err error
	conf, err = config.New()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error while creating config: %v\n", err)
		os.Exit(1)
	}
	if err := initLogger(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "error while init logger: %v\n", err)
		os.Exit(1)
	}
	if err := initRepositories(); err != nil {
		logger.Msg("error while init repositories").Error(err.Error())
		os.Exit(1)
	}
	if err := initManagers(); err != nil {
		logger.Msg("error while init managers").Error(err.Error())
		os.Exit(1)
	}
	httpApi = http.NewApi(http.Config{
		Host:              conf.HttpServer.Host,
		Port:              conf.HttpServer.Port,
		MetricsPort:       conf.HttpServer.MetricsPort,
		DebugPort:         conf.HttpServer.DebugPort,
		ReadTimeout:       conf.HttpServer.ReadTimeout,
		ReadHeaderTimeout: conf.HttpServer.ReadHeaderTimeout,
		WriteTimeout:      conf.HttpServer.WriteTimeout,
		IdleTimeout:       conf.HttpServer.IdleTimeout,
		MaxHeaderBytes:    conf.HttpServer.MaxHeaderBytes,
		CompanyManager:    companyMng,
		BuildingManager:   buildingMng,
		CategoryManager:   categoryMng,
	})
	httpApi.Start(func(err error) {
		logger.Msg("error while start http api").Error(err.Error())
		shutdownService()
		os.Exit(1)
	})
	handleOsSignals()
}

func initLogger() (err error) {
	return logger.Init(logger.Config{
		Path:   conf.Logger.Path,
		Stdout: conf.Logger.Stdout,
		Stderr: conf.Logger.Stderr,
		Level:  logger.LevelFromString(conf.Logger.Level),
	})
}

func initRepositories() (err error) {
	pgConfig := pg.Config{
		Host:     conf.Pg.Host,
		Port:     conf.Pg.Port,
		Database: conf.Pg.DbName,
		User:     conf.Pg.User,
		Password: conf.Pg.Password,
	}
	pgClient, err = pg.Connect(context.Background(), pgConfig)
	if err != nil {
		return fmt.Errorf("connect to pg database: %v", err)
	}

	repositories.SetClient(pgClient)
	if err = repositories.CreateDatabase(); err != nil {
		return fmt.Errorf("create database: %v", err)
	}

	buildingRepos, err = buildrepos.New(pgClient)
	if err != nil {
		return fmt.Errorf("init 'building' repository: %v", err)
	}

	categoryRepos, err = categrepos.New(pgClient)
	if err != nil {
		return fmt.Errorf("init 'category' repository: %v", err)
	}

	companyRepos, err = companrepos.New(pgClient, categoryRepos)
	if err != nil {
		return fmt.Errorf("init 'company' repository: %v", err)
	}
	return nil
}

func initManagers() (err error) {
	companyMng, err = companmng.New(companmng.Config{
		CompanyRepos: companyRepos,
	})
	if err != nil {
		return fmt.Errorf("error while creating company manager: %v", err)
	}

	buildingMng, err = buildmng.New(buildmng.Config{
		BuildingRepos: buildingRepos,
	})
	if err != nil {
		return fmt.Errorf("error while creating building manager: %v", err)
	}

	categoryMng, err = categmng.New(categmng.Config{
		CategoryRepos: categoryRepos,
	})
	if err != nil {
		return fmt.Errorf("error while creating category manager: %v", err)
	}
	return nil
}

func handleOsSignals() {
	osSignals := make(chan os.Signal)
	defer close(osSignals)
	signal.Notify(osSignals, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	for {
		<-osSignals
		shutdownService()
	}
}

func shutdownService() {
	ctx, cancel := context.WithTimeout(context.Background(), shutdownServiceTimeout)
	defer cancel()
	httpApi.Shutdown(ctx)
	if err := repositories.Stop(ctx); err != nil {
		logger.Msg("shutdown shared repository").Error(err.Error())
	}
	if err := buildingRepos.Stop(ctx); err != nil {
		logger.Msg("shutdown 'building' repository").Error(err.Error())
	}
	logger.Info("shutdown 'building' repository")
	if err := categoryRepos.Stop(ctx); err != nil {
		logger.Msg("shutdown 'category' repository").Error(err.Error())
	}
	logger.Info("shutdown 'category' repository")
	if err := companyRepos.Stop(ctx); err != nil {
		logger.Msg("shutdown 'company' repository").Error(err.Error())
	}
	logger.Info("shutdown 'company' repository")
}
