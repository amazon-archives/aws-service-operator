# AWS Service Operator

This is the Helm chart for the [AWS Service Operator](https://github.com/awslabs/aws-service-operator)

## Prerequisites

- Kubernetes 1.9+

## Installing the chart
Create AWS resources with Kubernetes
The chart can be installed by running:

```bash
$ helm repo add incubator http://storage.googleapis.com/kubernetes-charts-incubator
$ helm install incubator/aws-service-operator
```

## Configuration

The following table lists the configurable parameters of the aws-service-operator chart and their default values.

| Parameter                 | Description                            | Default                                            |
| ------------------------- | -------------------------------------- | -------------------------------------------------- |
| `image.repository`        | Container image repository             | `awsserviceoperator/aws-service-operator`          |
| `image.tag`               | Container image tag                    | `v0.0.1-alpha4`                                    |
| `image.pullPolicy`        | Container pull policy                  | `IfNotPresent`                                     |
| `image.pullSecret`        | Container pull secret (secret created not by this chart) | ``                               |
| `operator.accountId`      | AWS Account ID to operator on          | `""`                                               |
| `operator.bucket`         | Base bucket to store resources in      | `aws-operator`                                     |
| `operator.clusterName`    | Used to label generated CF templates   | `aws-operator`                                     |
| `operator.config`         | Config file                            | `$HOME/.aws-operator.yaml`                         |
| `operator.kubeconfig`     | Path to local kubeconfig file          | `""`                                               |
| `operator.logfile`        | Log file                               | `""`                                               |
| `operator.loglevel`       | Log level                              | `Info`                                             |
| `operator.region`         | AWS Region for created resources       | `us-west-2`                                        |
| `operator.resources`      | Comma Delimited list of CRDS to deploy | `s3bucket,dynamodb`                                |
| `affinity`                | affinity settings for pod assignment   | `{}`                                               |
| `extraArgs`               | Optional CLI argument                  | `[]`                                               |
| `extraEnv`                | Optional environment variables         | `[]`                                               |
| `extraVolumes`            | Custom Volumes                         | `[]`                                               |
| `extraVolumeMounts`       | Custom VolumeMounts                    | `[]`                                               |
| `nodeSelector`            | Node labels for pod assignment         | `{}`                                               |
| `podAnnotations`          | Annotations to attach to pod           | `{}`                                               |
| `rbac.create`             | Create RBAC roles                      | `true`                                             |
| `rbac.serviceAccountName` | Existing ServiceAccount to use         | `default`                                          |
| `replicas`                | Deployment replicas                    | `1`                                                |
| `resources`               | container resource requests and limits | `{}`                                               |
| `tolerations`             | Toleration labels for pod assignment   | `[]`                                               |
