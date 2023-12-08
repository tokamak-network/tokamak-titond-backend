package api

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	apptypes "github.com/tokamak-network/tokamak-titond-backend/pkg/types"
)

func MakeDeployerName(id uint) string {
	return fmt.Sprintf("deployer-%d", id)
}

func generateVolumePath(name string, networkID, componentID uint) string {
	return fmt.Sprintf("%s-%d-%d", name, networkID, componentID)
}

func generateNamespace(networkID uint) string {
	return fmt.Sprintf("namespace-%d", networkID)
}

func ConvertStructToMap(obj interface{}) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	val := reflect.ValueOf(obj)

	if val.Kind() != reflect.Struct {
		return nil, errors.New("the type of object is not struct")
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		fieldName := typ.Field(i).Name
		fieldValue := val.Field(i).Interface()
		result[fieldName] = fieldValue
	}
	return result, nil
}

func generatePublcURL(name string, networkID, componentID uint) string {
	return fmt.Sprintf("%s-%d-%d.titond-holesky.tokamak.network", name, networkID, componentID)
}

func checkDependency(status bool) error {
	if !status {
		return apptypes.ErrComponentDependency
	}
	return nil
}

func generateKey() (sk, address string) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		fmt.Printf("generateKey err : %s\n", err)
		return "", ""
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	sk = hexutil.Encode(privateKeyBytes)[2:]

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("error casting public key to ECDSA")
		return "", ""
	}

	address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return
}

func exportAddressFromPrivateKey(privateKey string) string {
	sk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		fmt.Printf("export address err : %s : %s\n", privateKey, err)
		return ""
	}

	publicKey := sk.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		fmt.Println("error casting public key to ECDSA")
		return ""
	}

	return crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
}
