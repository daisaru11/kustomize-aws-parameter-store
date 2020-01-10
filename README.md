Generates Secret from AWS SSM Parameter Store.

## Usage

```
$ cat <<'EOF' >./kustomization.yaml
generators:
- awsparameterstoresecret.yaml
EOF
```

```
$ cat <<'EOF' >./awsparameterstoresecret.yaml
apiVersion: kustomize.daisaru11.dev/v1
kind: AWSParameterStoreSecret
metadata:
  name: my-secret
  namespace: default
data:
  - name: FOO
    parameterName: /my-secrets/foo
EOF
```

```
$ aws ssm put-parameter \
  --name "/my-secrets/foo" \
  --value "FOO_SECRET" \
  --type "String"
```

```
$ kustomize build --enable_alpha_plugins .
apiVersion: v1
data:
  FOO: Rk9PX1NFQ1JFVA==
kind: Secret
metadata:
  name: my-secret
  namespace: default
type: Opaque
```

## Install

Use `go get`:

```
$ go get -u github.com/daisaru11/kustomize-aws-parameter-store
$ cp $(go env GOPATH)/bin/kustomize-aws-parameter-store \
  ${XDG_CONFIG_HOME:-$HOME/.config}/kustomize/plugin/kustomize.daisaru11.dev/v1/awsparameterstoresecret/AWSParameterStoreSecret
```

Or run the install script:

```
$ curl -sL https://raw.githubusercontent.com/daisaru11/kustomize-aws-parameter-store/master/hack/install.sh | bash
```