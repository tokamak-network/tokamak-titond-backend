package main

import (
	"log"
	"os"

	"github.com/tokamak-network/tokamak-titond-backend/cmd/utils"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/api"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/db"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/http"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/kubernetes"
	"github.com/urfave/cli/v2"
)

var app = &cli.App{
	Name:  "titond",
	Usage: "The titond command line interface",
}

func init() {
	app.Action = titond

	app.Flags = append(app.Flags, utils.KubernetesFlags...)
	app.Flags = append(app.Flags, utils.DBFlags...)
	app.Flags = append(app.Flags, utils.HTTPFlags...)
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func titond(ctx *cli.Context) error {
	k8sClient, err := kubernetes.NewKubernetes(&kubernetes.Config{
		InCluster:      ctx.Bool("kubernetes.incluster"),
		KubeconfigPath: ctx.String("kubernetes.kubeconfig.path"),
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

	apis := api.NewTitondAPI(k8sClient, dbClient)
	http := http.NewHTTPServer(&http.Config{
		Host: ctx.String("http.host"),
		Port: ctx.String("http.port"),
	}, apis)
	http.Run()

	return nil
}
