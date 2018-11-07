# Managed service controller

Kubernetes Controller to create a managed service namespace.

A managed service namespace will consist of managed services.
A managed service is something that a user of a cluster can make use of without needing to be concerned about the maintenance and overhead of operating that service.

These are usually heavy weight services that will be shared among many users.

## Status
POC

## Deploy to a cluster

#### Add CRD
```bash
$ oc create -f deploy/crd.yaml
```

#### Add Managed services RBAC
The managed-service-controller currently installs [Fuse](https://github.com/syndesisio/syndesis/tree/master/install/operator) and the [Integration controller](https://github.com/integr8ly/integration-controller) into the managed service namespace.
Set the required RBAC for those services.

```bash
$ oc create -f https://raw.githubusercontent.com/integr8ly/integration-controller/master/deploy/enmasse/enmasse-cluster-role.yaml
$ oc create -f https://raw.githubusercontent.com/integr8ly/integration-controller/master/deploy/applications/route-services-viewer-cluster-role.yaml
```

#### Deploy the manged services controller
```bash
$ oc create -f deploy/rbac.yaml
$ oc create -f deploy/operator.yaml
```

## Deploy a Custom Resource

The Managed Service Controller reacts to a `ManagedServiceNamespace` custom resource.

```yaml
properties:
  managedNamespace:
    type: string
  consumerNamespaces:
    type: array
```

There is an example resource in `deploy/cr.yaml`

```yaml
apiVersion: "integreatly.org/v1alpha1"
kind: "ManagedServiceNamespace"
metadata:
  name: "managed-services-project"
spec:
  managedNamespace: "managed-services-project"
  consumerNamespaces:
  - "myproject"
```

```bash
$ oc create -f deploy/cr.yaml -n <namespace_of_managed_service_controller>
```

Based on the above custom resource the Managed Service Controller will create a new namespace called `managed-services-project` and populate the namespace with managed services operators/controllers.
A user can then deploy these managed services from any of the namespaces in the `consumerNamespaces` list