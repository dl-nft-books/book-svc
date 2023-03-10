package config

import (
	"crypto/ecdsa"
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/crypto"
	"gitlab.com/distributed_lab/figure"
	"gitlab.com/distributed_lab/kit/comfig"
	"gitlab.com/distributed_lab/kit/kv"
	"gitlab.com/distributed_lab/logan/v3/errors"
)

const deploySignatureYamlKey = "deploy_signature"

type DeploySignatureConfigurator interface {
	DeploySignatureConfig() *DeploySignatureConfig
}

type DeploySignatureConfig struct {
	PrivateKey    *ecdsa.PrivateKey `fig:"eth_signer,required"`
	InitialOffset int64             `fig:"initial_offset,required"`
}

type deploySignatureConfigurator struct {
	getter kv.Getter
	once   comfig.Once
}

func NewDeploySignatureConfigurator(getter kv.Getter) DeploySignatureConfigurator {
	return &deploySignatureConfigurator{
		getter: getter,
	}
}

func (c *deploySignatureConfigurator) DeploySignatureConfig() *DeploySignatureConfig {
	return c.once.Do(func() interface{} {
		conf := DeploySignatureConfig{}

		if err := figure.
			Out(&conf).
			With(figure.BaseHooks, ecdsaHook).
			From(kv.MustGetStringMap(c.getter, deploySignatureYamlKey)).
			Please(); err != nil {
			panic(errors.Wrap(err, "failed to figure out deploy signature config"))
		}

		return &conf
	}).(*DeploySignatureConfig)
}

var ecdsaHook = figure.Hooks{
	"*ecdsa.PrivateKey": func(value interface{}) (reflect.Value, error) {
		switch v := value.(type) {
		case string:
			privKey, err := crypto.HexToECDSA(v)
			if err != nil {
				return reflect.Value{}, errors.Wrap(err, "invalid hex private key")
			}
			return reflect.ValueOf(privKey), nil
		default:
			return reflect.Value{}, fmt.Errorf("unsupported conversion from %T", value)
		}
	},
}
