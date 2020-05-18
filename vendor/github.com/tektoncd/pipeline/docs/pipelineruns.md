<!--
---
linkTitle: "PipelineRuns"
weight: 4
---
-->
# PipelineRuns

- [Overview](#overview)
- [Configuring a `PipelineRun`](#configuring-a-pipelinerun)
  - [Specifying the target `Pipeline`](#specifying-the-target-pipeline)
  - [Specifying `Resources`](#specifying-resources)
  - [Speciying `Parameters`](#specifying-parameters)
  - [Specifying custom `ServiceAccount` credentials](#specifying-custom-serviceaccount-credentials)
  - [Mapping `ServiceAccount` credentials to `Tasks`](#mapping-serviceaccount-credentials-to-tasks)
  - [Specifying `TaskRunSpecs`](#specifying-task-run-specs)
  - [Specifying a `Pod` template](#specifying-a-pod-template)
  - [Specifying `Workspaces`](#specifying-workspaces)
  - [Specifying `LimitRange` values](#specifying-limitrange-values)
  - [Configuring a failure timeout](#configuring-a-failure-timeout)
- [Cancelling a `PipelineRun`](#cancelling-a-pipelinerun)
- [Events](events.md#pipelineruns)



## Overview

A `PipelineRun` allows you to instantiate and execute a [`Pipeline`](pipelines.md) on-cluster.
A `Pipeline` specifies one or more `Tasks` in the desired order of execution. A `PipelineRun`
executes the `Tasks` in the `Pipeline` in the order they are specified until all `Tasks` have
executed successfully or a failure occurs.

**Note:** A `PipelineRun` automatically creates corresponding `TaskRuns` for every
`Task` in your `Pipeline`.

The `Status` field tracks the current state of a `PipelineRun`, and can be used to monitor
progress.
This field contains the status of every `TaskRun`, as well as the full `PipelineSpec` used
to instantiate this `PipelineRun`, for full auditability.

## Configuring a `PipelineRun`

A `PipelineRun` definition supports the following fields:

- Required:
  - [`apiVersion`][kubernetes-overview] - Specifies the API version. For example
    `tekton.dev/v1beta1`.
  - [`kind`][kubernetes-overview] - Indicates that this resource object is a `PipelineRun` object.
  - [`metadata`][kubernetes-overview] - Specifies the metadata that uniquely identifies the
    `PipelineRun` object. For example, a `name`.
  - [`spec`][kubernetes-overview] - Specifies the configuration information for
    this `PipelineRun` object.
    - [`pipelineRef` or `pipelineSpec`](#specifying-the-target-pipeline) - Specifies the target [`Pipeline`](pipelines.md).
- Optional:
  - [`resources`](#specifying-resources) - Specifies the [`PipelineResources`](resources.md) to provision
    for executing the target `Pipeline`.
  - [`params`](#specifying-parameters) - Specifies the desired execution parameters for the `Pipeline`.
  - [`serviceAccountName`](#specifying-serviceaccount-credentials) - Specifies a `ServiceAccount`
    object that supplies specific execution credentials for the `Pipeline`.
  - [`serviceAccountNames`](#mapping-serviceaccount-credentials-to-tasks) - Maps specific `serviceAccountName` values
    to `Tasks` in the `Pipeline`. This overrides the credentials set for the entire `Pipeline`.
  - [`taskRunSpec`](#specifying-task-run-specs) - Specifies a list of `PipelineRunTaskSpec` which allows for setting `ServiceAccountName` and [`Pod` template](./podtemplates.md) for each task. This overrides the `Pod` template set for the entire `Pipeline`. 
  - [`timeout`](#configuring-a-failure-timeout) - Specifies the timeout before the `PipelineRun` fails.
  - [`podTemplate`](#pod-template) - Specifies a [`Pod` template](./podtemplates.md) to use as the basis
    for the configuration of the `Pod` that executes each `Task`.

[kubernetes-overview]:
  https://kubernetes.io/docs/concepts/overview/working-with-objects/kubernetes-objects/#required-fields

### Specifying the target `Pipeline`

You must specify the target `Pipeline` that you want the `PipelineRun` to execute, either by referencing 
an existing `Pipeline` definition, or embedding a `Pipeline` definition directly in the `PipelineRun`.

To specify the target `Pipeline` by reference, use the `pipelineRef` field:

```yaml
spec:
  pipelineRef:
    name: mypipeline

```

To embed a `Pipeline` definition in the `PiepelineRun`, use the `pipelineSpec` field:

```yaml
spec:
  pipelineSpec:
    tasks:
    - name: task1
      taskRef:
        name: mytask
```

The `Pipeline` in the [`pipelineSpec` example](../examples/v1beta1/pipelineruns/pipelinerun-with-pipelinespec.yaml)
example displays morning and evening greetings. Once you create and execute it, you can check the logs for its `Pods`:

```bash
kubectl logs $(kubectl get pods -o name | grep pipelinerun-echo-greetings-echo-good-morning)
Good Morning, Bob!

kubectl logs $(kubectl get pods -o name | grep pipelinerun-echo-greetings-echo-good-night)
Good Night, Bob!
```

You can also embed a `Task` definition the embedded `Pipeline` definition:

```yaml
spec:
  pipelineSpec:
    tasks:
    - name: task1
      taskSpec:
        steps:
          ...
```

In the [`taskSpec` in `pipelineSpec` example](../examples/v1beta1/pipelineruns/pipelinerun-with-pipelinespec-and-taskspec.yaml)
it's `Tasks` all the way down!

## Specifying `Resources`

A `Pipeline` requires [`PipelineResources`](resources.md) to provide inputs and store outputs
for the `Tasks` that comprise it. You must provision those resources in the `resources` field
in the `spec` section of the `PipelineRun` definition.

A `Pipeline` may require you to provision a number of different resources. For example:

- When executing a `Pipeline` against a pull request, the triggering
  system must specify the commit-ish of a `git` resource.
- When executing a `Pipeline` manually against your own environment, you
  must provision your GitHub fork using the `git` resource; your image
  registry using the `image` resource; and your Kubernetes cluster using the
  `cluster` resource.

You can reference a `PipelineResources` using the `resourceRef` field:

```yaml
spec:
  resources:
    - name: source-repo
      resourceRef:
        name: skaffold-git
    - name: web-image
      resourceRef:
        name: skaffold-image-leeroy-web
    - name: app-image
      resourceRef:
        name: skaffold-image-leeroy-app
```

You can also embed a `PipelineResource` definition in the `PipelineRun` using the `resourceSpec` field:

```yaml
spec:
  resources:
    - name: source-repo
      resourceSpec:
        type: git
        params:
          - name: revision
            value: v0.32.0
          - name: url
            value: https://github.com/GoogleContainerTools/skaffold
    - name: web-image
      resourceSpec:
        type: image
        params:
          - name: url
            value: gcr.io/christiewilson-catfactory/leeroy-web
    - name: app-image
      resourceSpec:
        type: image
        params:
          - name: url
            value: gcr.io/christiewilson-catfactory/leeroy-app
```

**Note:** All `persistentVolumeClaims` specified within a `PipelineRun` are bound
until their respective `Pods` or the entire `PipelineRun` are deleted. This also applies
to all `persistentVolumeClaims` generated internally.

### Specifying `Parameters`

You can specify `Parameters` that you want to pass to the `Pipeline` during execution,
including different values of the same parameter for different `Tasks` in the `Pipeline`.

For example:

```yaml
spec:
  params:
  - name: pl-param-x
    value: "100"
  - name: pl-param-y
    value: "500"
```

### Specifying custom `ServiceAccount` credentials

You can execute the `Pipeline` in your `PipelineRun` with a specific set of credentials by 
specifying a `ServiceAccount` object name in the `serviceAccountName` field in your `PipelineRun`
definition. If you do not explicitly specify this, the `TaskRuns` created by your `PipelineRun`
will execute with the credentials specified in the `configmap-defaults` `ConfigMap`. If this
default is not specified, the `TaskRuns` will execute with the [`default` service account](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/#use-the-default-service-account-to-access-the-api-server)
set for the target [`namespace`](https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/).

For more information, see [`ServiceAccount`](auth.md).

### Mapping `ServiceAccount` credentials to `Tasks`

If you require more granularity in specifying execution credentials, use the `serviceAccounNames` field to
map a specific `serviceAccountName` value to a specific `Task` in the `Pipeline`. This overrides the global
`serviceAccountName` you may have set for the `Pipeline` as described in the previous section. 

For example, if you specify these mappings:

```yaml
spec:
  serviceAccountName: sa-1
  serviceAccountNames:
    - taskName: build-task
      serviceAccountName: sa-for-build
```

for this `Pipeline`:

```yaml
kind: Pipeline
spec:
  tasks:
    - name: build-task
      taskRef:
        name: build-push
    - name: test-task
      taskRef:
        name: test
```

then `test-task` will execute using the `sa-1` account while `build-task` will execute with `sa-for-build`.

### Specifying a `Pod` template

You can specify a [`Pod` template](podtemplates.md) configuration that will serve as the configuration starting
point for the `Pod` in which the container images specified in your `Tasks` will execute. This allows you to
customize the `Pod` confguration specifically for each `TaskRun`.

In the following example, the `Task` defines a `volumeMount` object named `my-cache`. The `PipelineRun`
provisions this object for the `Task` using a `persistentVolumeClaim` and executes it as user 1001.

```yaml
apiVersion: tekton.dev/v1beta1
kind: Task
metadata:
  name: mytask
spec:
  steps:
    - name: writesomething
      image: ubuntu
      command: ["bash", "-c"]
      args: ["echo 'foo' > /my-cache/bar"]
      volumeMounts:
        - name: my-cache
          mountPath: /my-cache
---
apiVersion: tekton.dev/v1beta1
kind: Pipeline
metadata:
  name: mypipeline
spec:
  tasks:
    - name: task1
      taskRef:
        name: mytask
---
apiVersion: tekton.dev/v1beta1
kind: PipelineRun
metadata:
  name: mypipelinerun
spec:
  pipelineRef:
    name: mypipeline
  podTemplate:
    securityContext:
      runAsNonRoot: true
      runAsUser: 1001
    volumes:
    - name: my-cache
      persistentVolumeClaim:
        claimName: my-volume-claim
```

## Specifying `Workspaces`

If your `Pipeline` specifies one or more `Workspaces`, you must map those `Workspaces` to
the corresponding physical volumes in your `PipelineRun` definition. For example, you
can map a `PersistentVolumeClaim` volume to a `Workspace` as follows:

```yaml
workspaces:
- name: myworkspace # must match workspace name in Task
  persistentVolumeClaim:
    claimName: mypvc # this PVC must already exist
  subPath: my-subdir
```

For more information, see the following topics:
- For information on mapping `Workspaces` to `Volumes`, see [Specifying `Workspaces` in `PipelineRuns`](workspaces.md#specifying-workspaces-in-pipelineruns).
- For a list of supported `Volume` types, see [Specifying `VolumeSources` in `Workspaces`](workspaces.md#specifying-volumesources-in-workspaces).
- For an end-to-end example, see [`Workspaces` in a `PipelineRun`](../examples/v1beta1/pipelineruns/workspaces.yaml).

### Specifying `LimitRange` values

In order to only consume the bare minimum amount of resources needed to execute one `Step` at a
time from the invoked `Task`, Tekton only requests the *maximum* values for CPU, memory, and ephemeral
storage from within each `Step`. This is sufficient as `Steps` only execute one at a time in the `Pod`.
Requests other than the maximum values are set to zero.

When a [`LimitRange`](https://kubernetes.io/docs/concepts/policy/limit-range/) parameter is present in
the namespace in which `PipelineRuns` are executing and *minimum* values are specified for container resource requests,
Tekton searches through all `LimitRange` values present in the namespace and uses the *minimums* instead of 0.

For more information, see the [`LimitRange` code example](../examples/v1beta1/pipelineruns/no-ci/limitrange.yaml).

### Configuring a failure timeout

You can use the `timeout` field to set the `PipelineRun's` desired timeout value in minutes.
If you do not specify this value in the `PipelineRun`, the global default timeout value applies.
If you set the timeout to 0, the `PipelineRun` fails immediately upon encountering an error.

The global default timeout is set to 60 minutes when you first install Tekton. You can set
a different global default timeout value using the `default-timeout-minutes` field in
[`config/config-defaults.yaml`](./../config/config-defaults.yaml).

The `timeout` value is a `duration` conforming to Go's
[`ParseDuration`](https://golang.org/pkg/time/#ParseDuration) format. For example, valid
values are `1h30m`, `1h`, `1m`, and `60s`. If you set the global timeout to 0, all `PipelineRuns`
that do not have an idividual timeout set will fail immediately upon encountering an error.

## Cancelling a `PipelineRun`

To cancel a `PipelineRun` that's currently executing, update its definition
to mark it as cancelled. When you do so, the spanwed `TaskRuns` are also marked
as cancelled and all associated `Pods` are deleted. For example:

```yaml
apiVersion: tekton.dev/v1beta1
kind: PipelineRun
metadata:
  name: go-example-git
spec:
  # […]
  status: "PipelineRunCancelled"
```

---

Except as otherwise noted, the content of this page is licensed under the
[Creative Commons Attribution 4.0 License](https://creativecommons.org/licenses/by/4.0/),
and code samples are licensed under the
[Apache 2.0 License](https://www.apache.org/licenses/LICENSE-2.0).

## Specifying task run specs

Specifies a list of  `PipelineRunTaskSpec` which contains `TaskServiceAccountName`,`TaskPodTemplate` and `TaskName`. Mapping the specs to the corresponding `Task` based upon the `TaskName` a PipelineTask will run with the configured  `TaskServiceAccountName` and `TaskPodTemplate` overwriting the pipeline wide [`ServiceAccountName`](#service-account)  and [`podTemplate`](#pod-template) configuration, for example:

```yaml
spec:
   podTemplate:
    securityContext:
      runAsUser: 1000
      runAsGroup: 2000
      fsGroup: 3000
  taskRunSpecs:
    - taskName: build-task
      taskServiceAccountName: sa-for-build
      taskPodTemplate:
        nodeSelector:
          disktype: ssd
```

If used with this `Pipeline`,  `build-task` will use the task specific pod template (where `nodeSelector` has `disktype` equal to `ssd`). 
