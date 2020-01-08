Generates Secret from AWS SSM Parameter Store.

## Usage

```
cat <<'EOF' >./kustomization.yaml
generators:
- awsparameterstoresecret.yaml
```

```
cat <<'EOF' >./awsparameterstoresecret.yaml
apiVersion: kustomize.daisaru11.dev/v1
kind: AWSParameterStoreSecret
metadata:
  name: my-secret
  namespace: default
data:
  - name: FOO
    parameterName: /my-secrets/foo
```

```
aws ssm put-parameter \
  --name "/my-secrets/foo" \
  --value "FOO_SECRET" \
  --type "String"
```

```
kustomize build --enable_alpha_plugins .
apiVersion: v1
data:
  FOO: Rk9PX1NFQ1JFVA==
kind: Secret
metadata:
  name: my-secret-6bbd9759bt
  namespace: default
type: Opaque
```