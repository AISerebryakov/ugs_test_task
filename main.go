package main

import (
	"fmt"
	"os"
	"ugc_test_task/companymng"
	"ugc_test_task/companyrepos"
	"ugc_test_task/config"
	"ugc_test_task/http"
	"ugc_test_task/logger"
	buildmng "ugc_test_task/managers/buildings"
	categmng "ugc_test_task/managers/categories"
	"ugc_test_task/pg"
	buildrepos "ugc_test_task/repositories/buildings"
	categrepos "ugc_test_task/repositories/categories"
)

var (
	conf config.Config

	categoryRepos categrepos.Repository
	companyRepos  companyrepos.Repository
	buildingRepos buildrepos.Repository

	companyMng  companymng.Manager
	buildingMng buildmng.Manager
	categoryMng categmng.Manager
)

func main() {
	var err error
	conf, err = config.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error while creating config: %v\n", err)
		os.Exit(1)
	}
	if err := initLogger(); err != nil {
		fmt.Fprintf(os.Stderr, "error while init logger: %v\n", err)
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
	httpApi, err := http.NewApi(http.Config{
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
	if err != nil {
		logger.Msg("error while creating http api").Error(err.Error())
		os.Exit(1)
	}
	//todo: handle error
	httpApi.Start(func(err error) {
		logger.Msg("error while start http api").Error(err.Error())
	})
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

	categoryRepos, err = categrepos.New(categrepos.NewConfig(pgConfig))
	if err != nil {
		return fmt.Errorf("init category repository: %v", err)
	}

	companyConf := companyrepos.NewConfig(pgConfig)
	companyConf.CategoryRepos = categoryRepos
	companyRepos, err = companyrepos.New(companyConf)
	if err != nil {
		return fmt.Errorf("init company repository: %v", err)
	}

	buildingRepos, err = buildrepos.New(buildrepos.NewConfig(pgConfig))
	if err != nil {
		return fmt.Errorf("init building repository: %v", err)
	}
	return nil
}

func initManagers() (err error) {
	companyMng, err = companymng.New(companymng.Config{
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
