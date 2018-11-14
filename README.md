# Managed service controller

Kubernetes Controller to create a managed service namespace.

A managed service namespace will consist of managed services.
A managed service is something that a user of a cluster can make use of without needing to be concerned about the maintenance and overhead of operating that service.

These are usually heavy weight services that will be shared among many users.

## Status
POC

## Cluster Setup

#### Add ManagedServiceNamespace CRD
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

```bash
$ oc create -f ./deploy/crd.yaml
```

#### Managed Service Requirements
The managed-service-controller currently installs the managed services:
- [Fuse](https://github.com/syndesisio/syndesis/tree/master/install/operator)
- [Integration controller](https://github.com/integr8ly/integration-controller)

into the managed service namespace. Setup the required resources for those services.

```bash
# RBACS
$ oc create -f https://raw.githubusercontent.com/integr8ly/integration-controller/master/deploy/enmasse/enmasse-cluster-role.yaml
$ oc create -f https://raw.githubusercontent.com/integr8ly/integration-controller/master/deploy/applications/route-services-viewer-cluster-role.yaml

# Custom resource definitions
$ oc create -f https://raw.githubusercontent.com/syndesisio/syndesis/master/install/operator/deploy/syndesis-crd.yml
$ oc create -f https://raw.githubusercontent.com/integr8ly/integration-controller/master/deploy/crd.yaml

***SHOULD THIS BE MOVED INTO THE CONTROLLER****
# Fuse Image Streams
$ oc create -f https://raw.githubusercontent.com/syndesisio/fuse-online-install/1.4.8/resources/fuse-online-image-streams.yml -n openshift
$ oc create -f ./deploy/fuse-image-stream.yaml -n openshift
```

***GROUPS FIX THIS?****
```bash
# Create enmasse namespace
$ oc create namespace enmasse
```

#### Deploy the manged services controller
```bash
$ oc new-project <manged-services-controller-namespace>
$ oc create -f ./deploy/rbac.yaml
$ oc create -f ./deploy/operator.yaml
```

## Deploy a Custom Resource

There is an example resource in `deploy/cr.yaml` with the following spec:

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

__NOTE:__ Ensure any namespaces in the `consumerNamespaces` array are created.

Create the custom resource:
```bash
$ oc create -f ./deploy/cr.yaml
```

Based on the above custom resource the Managed Service Controller will create a new namespace called `managed-services-project` and populate the namespace with managed services operators/controllers.
A user can then deploy these managed services from any of the namespaces in the `consumerNamespaces` list