package utils

import (
	"os"
	"path"

	"github.com/urfave/cli/v2"
)

var (
	// kubernetes flags
	KubernetesFlags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "kubernetes.incluster",
			Usage:   "If this app is in cluster or not",
			Value:   false,
			EnvVars: []string{"KUBERNETES_INCLUSTER"},
		},
		&cli.StringFlag{
			Name:    "kubernetes.kubeconfig.path",
			Usage:   "The path of kubeconfig",
			Value:   path.Join(os.Getenv("HOME"), "/.kube/config"),
			EnvVars: []string{"KUBERNETES_KUBECONFIG_PATH"},
		},
	}
	// databse flags
	DBFlags = []cli.Flag{
		&cli.StringFlag{
			Name:    "db.type",
			Usage:   "The type of database",
			Value:   "postgres",
			EnvVars: []string{"DB_TYPE"},
		},
		&cli.StringFlag{
			Name:    "db.host",
			Usage:   "The host of database",
			Value:   "localhost",
			EnvVars: []string{"DB_HOST"},
		},
		&cli.StringFlag{
			Name:    "db.port",
			Usage:   "The port of database",
			Value:   "5432",
			EnvVars: []string{"DB_PORT"},
		},
		&cli.StringFlag{
			Name:    "db.user",
			Usage:   "The user of database",
			Value:   "postgres",
			EnvVars: []string{"DB_USER"},
		},
		&cli.StringFlag{
			Name:    "db.password",
			Usage:   "The password of database",
			Value:   "postgres",
			EnvVars: []string{"DB_PASSWORD"},
		},
		&cli.StringFlag{
			Name:    "db.dbname",
			Usage:   "The database name of database",
			Value:   "titond",
			EnvVars: []string{"DB_DBNAME"},
		},
	}
	// http flags
	HTTPFlags = []cli.Flag{
		&cli.StringFlag{
			Name:    "http.host",
			Usage:   "The host of server",
			Value:   "0.0.0.0",
			EnvVars: []string{"HTTP_HOST"},
		},
		&cli.StringFlag{
			Name:    "http.port",
			Usage:   "The port of server",
			Value:   "8080",
			EnvVars: []string{"HTTP_PORT"},
		},
	}
)
