## Supporting Custom Resources (defined by a CRD)

This tutorial shows how to add transformer configurations to support a custom resource.

Create a workspace by
<<<<<<< HEAD
<!-- @createws @test -->
=======
<!-- @createws @testAgainstLatestRelease -->
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
```
DEMO_HOME=$(mktemp -d)
```

### Adding a custom resource

Consider a CRD of kind `MyKind` with fields
- `.spec.secretRef.name` reference a Secret
- `.spec.beeRef.name` reference an instance of CRD `Bee`
- `.spec.containers.command` as the list of container commands
- `.spec.selectors` as the label selectors

Add the following file to configure the transformers for the above fields
<<<<<<< HEAD
<!-- @addConfig @test -->
=======
<!-- @addConfig @testAgainstLatestRelease -->
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
```
mkdir $DEMO_HOME/kustomizeconfig
cat > $DEMO_HOME/kustomizeconfig/mykind.yaml << EOF

commonLabels:
- path: spec/selectors
  create: true
  kind: MyKind

nameReference:
- kind: Bee
  fieldSpecs:
  - path: spec/beeRef/name
    kind: MyKind
- kind: Secret
  fieldSpecs:
  - path: spec/secretRef/name
    kind: MyKind

varReference:
- path: spec/containers/command
  kind: MyKind
- path: spec/beeRef/action
  kind: MyKind

EOF
```

### Apply config

Create a file with some resources that
includes an instance of `MyKind`:

<<<<<<< HEAD
<!-- @createResource @test -->
=======
<!-- @createResource @testAgainstLatestRelease -->
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
```
cat > $DEMO_HOME/resources.yaml << EOF
apiVersion: v1
kind: Secret
metadata:
  name: crdsecret
data:
  PATH: YmJiYmJiYmIK
---
apiVersion: v1beta1
kind: Bee
metadata:
  name: bee
spec:
  action: fly
---
apiVersion: jingfang.example.com/v1beta1
kind: MyKind
metadata:
  name: mykind
spec:
  secretRef:
    name: crdsecret
  beeRef:
    name: bee
    action: \$(BEE_ACTION)
  containers:
  - command:
    - "echo"
    - "\$(BEE_ACTION)"
    image: myapp
EOF
```

Create a kustomization referring to it:

<<<<<<< HEAD
<!-- @createKustomization @test -->
=======
<!-- @createKustomization @testAgainstLatestRelease -->
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
```
cat > $DEMO_HOME/kustomization.yaml << EOF
resources:
- resources.yaml

namePrefix: test-

commonLabels:
  foo: bar

vars:
- name: BEE_ACTION
  objref:
    kind: Bee
    name: bee
    apiVersion: v1beta1
  fieldref:
    fieldpath: spec.action
EOF
```

Use the customized transformer configurations by specifying them
in the kustomization file:
<<<<<<< HEAD
<!-- @addTransformerConfigs @test -->
=======
<!-- @addTransformerConfigs @testAgainstLatestRelease -->
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
```
cat >> $DEMO_HOME/kustomization.yaml << EOF
configurations:
- kustomizeconfig/mykind.yaml
EOF
```

Run `kustomize build` and verify that the namereference is correctly resolved.

<<<<<<< HEAD
<!-- @build @test -->
=======
<!-- @build @testAgainstLatestRelease -->
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
```
test 2 == \
$(kustomize build $DEMO_HOME | grep -A 2 ".*Ref" | grep "test-" | wc -l); \
echo $?    
```

Run `kustomize build` and verify that the vars correctly resolved.

<<<<<<< HEAD
<!-- @verify @test -->
=======
<!-- @verify @testAgainstLatestRelease -->
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
```
test 0 == \
$(kustomize build $DEMO_HOME | grep "BEE_ACTION" | wc -l); \
echo $?    
```
