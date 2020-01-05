package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"sigs.k8s.io/kustomize/api/kv"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/yaml"
)

type plugin struct {
	h                *resmap.PluginHelpers
	types.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	Data             []AWSParameterStoreSecretItem `json:"data,omitempty" yaml:"data,omitempty"`
}

type AWSParameterStoreSecretItem struct {
	Name          string `json:"name,omitempty" yaml:"name,omitempty"`
	ParameterName string `json:"parameterName,omitempty" yaml:"parameterName,omitempty"`
}

//nolint: golint, deadcode, gochecknoglobals, unused
var KustomizePlugin plugin

func (p *plugin) Config(h *resmap.PluginHelpers, c []byte) error {
	p.h = h
	return yaml.Unmarshal(c, p)
}

func (p *plugin) Generate() (resmap.ResMap, error) {
	args := types.SecretArgs{}
	args.Name = p.Name
	args.Namespace = p.Namespace

	params, err := p.getParameters(p.Data)
	if err != nil {
		return nil, err
	}

	for name, value := range params {
		args.LiteralSources = append(
			args.LiteralSources, fmt.Sprintf("%s=%s", name, value))
	}

	return p.h.ResmapFactory().FromSecretArgs(
		kv.NewLoader(p.h.Loader(), p.h.Validator()), nil, args)
}

func (p *plugin) getParameters(data []AWSParameterStoreSecretItem) (map[string]string, error) {
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
