package config

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

func CreateGasConfig(gasMap map[string]uint64) (*GasCost, error) {
	baseOps := &BaseOperationCost{}
	err := mapstructure.Decode(gasMap, baseOps)
	if err != nil {
		return nil, err
	}

	elrondOps := &ElrondAPICost{}
	err = mapstructure.Decode(gasMap, elrondOps)
	if err != nil {
		return nil, err
	}

	err = checkForZeroUint64Fields(elrondOps)
	if err != nil {
		return nil, err
	}

	bigIntOps := &BigIntAPICost{}
	err = mapstructure.Decode(gasMap, bigIntOps)
	if err != nil {
		return nil, err
	}

	err = checkForZeroUint64Fields(bigIntOps)
	if err != nil {
		return nil, err
	}

	ethOps := &EthAPICost{}
	err = mapstructure.Decode(gasMap, ethOps)
	if err != nil {
		return nil, err
	}

	err = checkForZeroUint64Fields(ethOps)
	if err != nil {
		return nil, err
	}

	cryptOps := &CryptoAPICost{}
	err = mapstructure.Decode(gasMap, cryptOps)
	if err != nil {
		return nil, err
	}

	err = checkForZeroUint64Fields(cryptOps)
	if err != nil {
		return nil, err
	}

	gasCost := &GasCost{
		BaseOperationCost: *baseOps,
		BigIntAPICost:     *bigIntOps,
		EthAPICost:        *ethOps,
		ElrondAPICost:     *elrondOps,
		CryptoAPICost:     *cryptOps,
	}

	return gasCost, nil
}

func checkForZeroUint64Fields(arg interface{}) error {
	v := reflect.ValueOf(arg)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() != reflect.Uint64 {
			continue
		}
		if field.Uint() == 0 {
			name := v.Type().Field(i).Name
			return errors.New(fmt.Sprintf("field %s has the value 0", name))
		}
	}

	return nil
}
