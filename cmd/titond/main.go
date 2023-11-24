package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/tokamak-network/tokamak-titond-backend/cmd/utils"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/api"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/db"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/http"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/kubernetes"
	"github.com/tokamak-network/tokamak-titond-backend/pkg/services"
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
			Usage:   "Greet a person",
			Action:  checkSwagger,
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

func checkSwagger(ctx *cli.Context) error {
	fmt.Println("Check swagger")
	http := http.NewHTTPServer(&http.Config{
		Host: ctx.String("http.host"),
		Port: ctx.String("http.port"),
	}, nil)

	numAPIs := len(http.R.Routes())

	for _, route := range http.R.Routes() {
		fmt.Printf("%-6s %-25s %s\n", route.Method, route.Path, route.Handler)
	}
	fmt.Println("Total APIs: ", numAPIs)

	jsonData, err := ioutil.ReadFile("./docs/swagger.json")
	if err != nil {
		fmt.Println("Read docs file:", err)
		return err
	}

	var data map[string]interface{}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println("Docs data:", err)
		return err
	}
	paths, exist := data["paths"]
	if !exist {
		fmt.Println("Failed to fetch api description in swagger file")
		return errors.New("")
	}
	switch v := paths.(type) {
	case map[string]interface{}:
		fmt.Println("num api in swagger", len(v))
		if len(v) != numAPIs-1 {
			return errors.New("the docs file was not updated")
		}
	default:
		return errors.New("")
	}
	return nil
}
