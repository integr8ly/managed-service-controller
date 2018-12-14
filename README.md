# Managed service controller

Kubernetes Controller to create a managed service namespace.

A managed service namespace will consist of managed services.
A managed service is something that a user of a cluster can make use of without needing to be concerned about the maintenance and overhead of operating that service.

These are heavy weight services that will be single tenant services, e.g. a service shared by a single tenant across multiple consumer namespaces.

## Status
POC

## Dependencies
__NOTE:__ This command should only be run once.
```bash
make setup/dep
```

## Cluster Setup

#### ManagedServiceNamespace CRD
The Managed Service Controller reacts to a `ManagedServiceNamespace` custom resource with the following properties:

```yaml
properties:
    userId:
      type: string
      description: This User will be granted view access to the ManagedServiceNamespace and update access to the Integration resource.
    consumerNamespaces:
      type: array
      description: Names of the namespaces that will be allowed use the services created in ManagedServiceNamespace.
```

#### Managed Service Requirements
The managed-service-controller currently installs the managed services/managed service operators:
- [Fuse](https://github.com/syndesisio/syndesis/tree/master/install/operator)
- [Integration controller](https://github.com/integr8ly/integration-controller)

into a managed service namespace.

#### Deploy the managed services controller

```bash
# login as cluster-admin
$ oc login -u system:admin

# Setup the cluster with the required CRDs, setup, and resources.
make cluster/prepare
```
__NOTE:__ `make cluster/clean` will remove the required CRDs, setup, and resources.

```bash
$ oc new-project <manged-services-controller-namespace>
$ make cluster/deploy PROJECT_NAME=<manged-services-controller-namespace>
```
__NOTE:__ `make cluster/deploy/remove PROJECT_NAME=<manged-services-controller-namespace>` will remove the managed-service-controller deployment.

## Deploy a Custom Resource

There is an example resource in `./deploy/crds/integreatly_v1alpha1_managedservicenamespace_cr.yaml` with the following spec:

```yaml
spec:
  autoEnabledServiceIntegrations: "true"
  userId: "developer"
  consumerNamespaces:
  - "consumer1"
  - "consumer2"
```

The name of the created namespace will be taken from the custom resource name:
```yaml
metadata:
  name: "managed-services-project"
```

__NOTE:__ Ensure the User and any namespaces in the `consumerNamespaces` array exist in the cluster.

```bash
$ oc create user admin
$ oc create namespace consumer1
$ oc create namespace consumer2
```

Create the custom resource:
```bash
$ oc create -f ./deploy/crds/integreatly_v1alpha1_managedservicenamespace_cr.yaml -n <manged-services-controller-namespace>
```

Based on the above custom resource the Managed Service Controller will create a new namespace called `managed-services-project` and populate the namespace with managed services operators/controllers.
A user can then deploy these managed services from any of the namespaces in the `consumerNamespaces` list

Delete the custom resource to remove the managed service namespace: `managed-services-project`
```bash
$ oc delete -f ./deploy/crds/integreatly_v1alpha1_managedservicenamespace_cr.yaml -n <manged-services-controller-namespace>
```