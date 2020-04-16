[kubernetes API object style]: https://kubernetes.io/docs/concepts/overview/working-with-objects/kubernetes-objects/#required-fields
[variant]: ../docs/glossary.md#variant

# Demo: configure breakfast


Define a place to work:

<<<<<<< HEAD
<!-- @makeWorkplace @test -->
=======
<!-- @makeWorkplace @testAgainstLatestRelease -->
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
```
DEMO_HOME=$(mktemp -d)
```

Make a place to put the base breakfast configuration:

<<<<<<< HEAD
<!-- @baseDir @test -->
=======
<!-- @baseDir @testAgainstLatestRelease -->
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
```
mkdir -p $DEMO_HOME/breakfast/base
```

Make a `kustomization` to define what goes into
breakfast.  This breakfast has coffee and pancakes:

<<<<<<< HEAD
<!-- @baseKustomization @test -->
=======
<!-- @baseKustomization @testAgainstLatestRelease -->
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
```
cat <<EOF >$DEMO_HOME/breakfast/base/kustomization.yaml
resources:
- coffee.yaml
- pancakes.yaml
EOF
```

Here's a _coffee_ type.  Give it a `kind` and `metdata/name` field
to conform to [kubernetes API object style]; no other
file or definition is needed:

<<<<<<< HEAD
<!-- @coffee @test -->
=======
<!-- @coffee @testAgainstLatestRelease -->
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
```
cat <<EOF >$DEMO_HOME/breakfast/base/coffee.yaml
kind: Coffee
metadata:
  name: morningCup
temperature: lukewarm
data:
  greeting: "Good Morning!"
EOF
```

The `name` field merely distinguishes this instance of
coffee from others (if there were any).

Likewise, define _pancakes_:
<<<<<<< HEAD
<!-- @pancakes @test -->
=======
<!-- @pancakes @testAgainstLatestRelease -->
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
```
cat <<EOF >$DEMO_HOME/breakfast/base/pancakes.yaml
kind: Pancakes
metadata:
  name: comfort
stacksize: 3
topping: none
EOF
```

Make a custom [variant] of breakfast for Alice, who
likes her coffee hot:

<<<<<<< HEAD
<!-- @aliceOverlay @test -->
=======
<!-- @aliceOverlay @testAgainstLatestRelease -->
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
```
mkdir -p $DEMO_HOME/breakfast/overlays/alice

cat <<EOF >$DEMO_HOME/breakfast/overlays/alice/kustomization.yaml
commonLabels:
  who: alice
<<<<<<< HEAD
bases:
=======
resources:
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
- ../../base
patchesStrategicMerge:
- temperature.yaml
EOF

cat <<EOF >$DEMO_HOME/breakfast/overlays/alice/temperature.yaml
kind: Coffee
metadata:
  name: morningCup
temperature: hot!
EOF
```

And likewise a [variant] for Bob, who wants _five_ pancakes, with strawberries:

<<<<<<< HEAD
<!-- @bobOverlay @test -->
=======
<!-- @bobOverlay @testAgainstLatestRelease -->
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
```
mkdir -p $DEMO_HOME/breakfast/overlays/bob

cat <<EOF >$DEMO_HOME/breakfast/overlays/bob/kustomization.yaml
commonLabels:
  who: bob
<<<<<<< HEAD
bases:
=======
resources:
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
- ../../base
patchesStrategicMerge:
- topping.yaml
EOF

cat <<EOF >$DEMO_HOME/breakfast/overlays/bob/topping.yaml
kind: Pancakes
metadata:
  name: comfort
stacksize: 5
topping: strawberries
EOF
```

One can now generate the configs for Alice’s breakfast:

<<<<<<< HEAD
<!-- @generateAlice @test -->
=======
<!-- @generateAlice @testAgainstLatestRelease -->
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
```
kustomize build $DEMO_HOME/breakfast/overlays/alice
```

Likewise for Bob:

<<<<<<< HEAD
<!-- @generateBob @test -->
=======
<!-- @generateBob @testAgainstLatestRelease -->
>>>>>>> Create "add application" odo  pipeline sub-comment (#51)
```
kustomize build $DEMO_HOME/breakfast/overlays/bob
```
