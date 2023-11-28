package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/tokamak-network/tokamak-titond-backend/pkg/http"
	"github.com/urfave/cli/v2"
)

func CheckSwagger(ctx *cli.Context) error {
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
	fmt.Println("Total APIs excluding swagger api: ", (numAPIs - 1))

	jsonData, err := ioutil.ReadFile("./api/swagger.json")
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
		numOfApisInDocs := countNumOfAPIsInDocs(v)
		if numOfApisInDocs != numAPIs-1 {
			return errors.New("missing swagger description for api")
		}
	default:
		return errors.New("")
	}
	return nil
}

func countNumOfAPIsInDocs(paths map[string]interface{}) int {
	counter := 0
	for path, metadata := range paths {
		fmt.Println("  [Path]:", path)
		switch v := metadata.(type) {
		case map[string]interface{}:
			counter += len(v)
		default:
		}
	}
	fmt.Println("Api in docs ", counter)
	return counter
}
