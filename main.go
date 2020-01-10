package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	yaml "gopkg.in/yaml.v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type ObjectMeta struct {
	Name      string `json:"name,omitempty" yaml:"name,omitempty"`
	Namespace string `json:"namespace,omitempty" yaml:"namespace,omitempty"`
}

type AWSParameterStoreSecret struct {
	Kind       string `json:"kind,omitempty" yaml:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Data       []AWSParameterStoreSecretItem `json:"data,omitempty" yaml:"data,omitempty"`
}

type AWSParameterStoreSecretItem struct {
	Name          string `json:"name,omitempty" yaml:"name,omitempty"`
	ParameterName string `json:"parameterName,omitempty" yaml:"parameterName,omitempty"`
}

type Secret struct {
	Kind       string `json:"kind,omitempty" yaml:"kind,omitempty"`
	APIVersion string `json:"apiVersion,omitempty" yaml:"apiVersion,omitempty"`
	ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Data       map[string]string `json:"data,omitempty" yaml:"data,omitempty"`
	Type       string            `json:"type,omitempty" yaml:"type,omitempty"`
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		log.Fatal("config file name is required as an argument")
	}

	confFile := args[0]

	confData, err := ioutil.ReadFile(confFile)
	if err != nil {
		log.Fatalf("failed to read conf file: %s, %s", confFile, err)
	}

	conf := AWSParameterStoreSecret{}

	err = yaml.Unmarshal(confData, &conf)
	if err != nil {
		log.Fatalf("failed to unmarshal the config: %s, %s", confFile, err)
	}

	params, err := getParameters(conf.Data)
	if err != nil {
		log.Fatalf("failed to get parameters: %s", err)
	}

	secret := Secret{
		Kind:       "Secret",
		APIVersion: "v1",
		ObjectMeta: ObjectMeta{
			Name:      conf.Name,
			Namespace: conf.Namespace,
		},
		Data: map[string]string{},
		Type: "Opaque",
	}

	for k, v := range params {
		secret.Data[k] = base64.StdEncoding.EncodeToString([]byte(v))
	}

	out, err := yaml.Marshal(secret)
	if err != nil {
		log.Fatalf("failed to marshal the secret: %s", err)
	}

	fmt.Print(string(out))
}

func getParameters(data []AWSParameterStoreSecretItem) (map[string]string, error) {
	paramNameMap := map[string]string{}
	for _, item := range data {
		paramNameMap[item.ParameterName] = item.Name
	}

	sess, err := session.NewSession()
	if err != nil {
		return nil, err
	}

	ssmsvc := ssm.New(sess)

	input := ssm.GetParametersInput{
		WithDecryption: aws.Bool(true),
	}

	for _, item := range data {
		input.Names = append(input.Names, aws.String(item.ParameterName))
	}

	output, err := ssmsvc.GetParameters(&input)
	if err != nil {
		return nil, err
	}

	params := map[string]string{}

	for _, p := range output.Parameters {
		name, ok := paramNameMap[aws.StringValue(p.Name)]
		if !ok {
			continue
		}

		params[name] = aws.StringValue(p.Value)
	}

	return params, nil
}
