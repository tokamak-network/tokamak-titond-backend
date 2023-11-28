package main

import (
	"log"
	"os"

	"github.com/tokamak-network/tokamak-titond-backend/cmd/utils"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/api"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/db"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/http"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/kubernetes"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/services"
	apptest "github.com/tokamak-network/tokamak-titond-backend/test"
	"github.com/urfave/cli/v2"
)

var app = &cli.App{
	Name:  "titond",
	Usage: "The titond command line interface",
}

func init() {
	app.Commands = []*cli.Command{
		{
			Name:    "check-swagger",
			Aliases: []string{"g"},
			Usage:   "check swagger",
			Action:  apptest.CheckSwagger,
		},
	}
	app.Action = titond
	app.Flags = append(app.Flags, utils.TitondFlags...)
	app.Flags = append(app.Flags, utils.KubernetesFlags...)
	app.Flags = append(app.Flags, utils.DBFlags...)
	app.Flags = append(app.Flags, utils.ServicesFlags...)
	app.Flags = append(app.Flags, utils.HTTPFlags...)
}

// @title Titond
// @version 1.0
// @description Titond-backend application
// @host localhost:8080
func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func titond(ctx *cli.Context) error {
	k8sClient, err := kubernetes.NewKubernetes(&kubernetes.Config{
		InCluster:      ctx.Bool("kubernetes.incluster"),
		KubeconfigPath: ctx.String("kubernetes.kubeconfig.path"),
		ManifestPath:   ctx.String("kubernetes.manifest.path"),
	})
	if err != nil {
		return err
	}

	var dbClient db.Client
	if ctx.String("db.type") == "postgres" {
		dbClient, err = db.NewPostgresql(&db.Config{
			Host:     ctx.String("db.host"),
			Port:     ctx.String("db.port"),
			User:     ctx.String("db.user"),
			Password: ctx.String("db.password"),
			DBName:   ctx.String("db.dbname"),
		})
		if err != nil {
			return err
		}
	}

	fileManager := services.NewS3(&services.S3Config{
		BucketName: ctx.String("services.s3.bucket"),
		AWSRegion:  ctx.String("services.s3.region"),
	})

	apis := api.NewTitondAPI(k8sClient, dbClient, fileManager, &api.Config{
		Namespace:              ctx.String("titond.namespace"),
		ContractsRpcUrl:        ctx.String("titond.contracts.rpc.url"),
		ContractsTargetNetwork: ctx.String("titond.contracts.target.network"),
		ContractsDeployerKey:   ctx.String("titond.contracts.deployer.key"),
	})

	http := http.NewHTTPServer(&http.Config{
		Host: ctx.String("http.host"),
		Port: ctx.String("http.port"),
	}, apis)
	http.Run()

	return nil
}
