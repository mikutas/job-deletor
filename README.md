# job-deletor

## Prerequisites

https://github.com/kubernetes-sigs/kubebuilder/blob/book-v3/docs/book/src/quick-start.md#prerequisites

## Install CRD

```
make install
```

## Deploy

```
make deploy
```

## Use


Sample CR

`config/samples/mikutas_v1alpha1_jobdeletor_all.yaml`

```
kubectl apply -f config/samples/mikutas_v1alpha1_jobdeletor_all.yaml
```

Sample Job

```
kubectl apply -f config/samples/job-success.yaml
```

## Delete

```
make undeploy
```

## Uninstall CRD

`make undeploy` includes CRD uninstallation

```
make uninstall
```
