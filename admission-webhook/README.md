# Managed Service Namespace Admission Webhook

## Cluster Setup

Mutating and Validating webhooks are not enabled by default on Openshift 3.10.0. To enable complete the following setup.

### Minishift

```bash
$ minishift start --openshift-version v3.10.0
$ minishift ssh

docker@minishift $ sed -i -e 's/"admissionConfig":{"pluginConfig":null}/"admissionConfig": {\
    "pluginConfig": {\
        "ValidatingAdmissionWebhook": {\
            "configuration": {\
                "apiVersion": "v1",\
                "kind": "DefaultAdmissionConfig",\
                "disable": false\
            }\
        },\
        "MutatingAdmissionWebhook": {\
            "configuration": {\
                "apiVersion": "v1",\
                "kind": "DefaultAdmissionConfig",\
                "disable": false\
            }\
        }\
    }\
}/' /var/lib/minishift/base/kube-apiserver/master-config.yaml

$ exit

$ minishift stop
$ minishift start
```


### oc cluster up

```bash
$ oc cluster up --write-config

$ sed -i -e 's/"admissionConfig":{"pluginConfig":null}/"admissionConfig": {\
    "pluginConfig": {\
        "ValidatingAdmissionWebhook": {\
            "configuration": {\
                "apiVersion": "v1",\
                "kind": "DefaultAdmissionConfig",\
                "disable": false\
            }\
        },\
        "MutatingAdmissionWebhook": {\
            "configuration": {\
                "apiVersion": "v1",\
                "kind": "DefaultAdmissionConfig",\
                "disable": false\
            }\
        }\
    }\
}/' openshift.local.clusterup/kube-apiserver/master-config.yaml

$ oc cluster up --server-loglevel=5
```

## Deploying to a Cluster
```bash
# first create a new project to deploy to
$ oc new-project <namespace>

# User should be able to create cluster level resources
$ make deploy_template NAMESPACE=<namespace>

# To clean up the deployment, run
$ make delete_template NAMESPACE=<namespace>
```

The kube-api server needs to authenticate the webhook server so a server cert and server key for the webhook server and a CA cert for the config are needed. `make deploy_template` sets up [the template](https://github.com/integr8ly/managed-service-controller/tree/master/admission-webhook/deploy/webhook.template.yaml) for you but you can do these steps manually too.

### Manually Create Credentials
```bash
# first create a new project to deploy to
$ oc new-project <namespace>>

# Then create the service credentials
$ make service_context NAMESPACE=<namespace>

Base64 encoded CA cert
LS0tLS1CRUdJTiBDRVJUSUZJQ...


Base64 encoded Server key
LS0tLS1CRUdJTiBSU0EgUFJJV...


Base64 encoded Server cert
LS0tLS1CRUdJTiBDRVJUSUZJQ...
```

The output of `make service_context` can be plugged straight into the template.

```bash
oc process -f ./deploy/webhook.template.yaml --param=NAMESPACE=<namespace> --param=SERVER_KEY=LS0tLS1CRUdJTiBSU0EgUFJJV... --param=SERVER_CERT=LS0tLS1CRUdJTiBDRVJUSUZJQ... --param=CA_BUNDLE=LS0tLS1CRUdJTiBDRVJUSUZJQ... | oc create -f - -n <namespace>
```

## Develop Locally
__NOTE__: Make sure `KUBERNETES_CONFIG` is set

```bash
$ export KUBERNETES_CONFIG=/$HOME/.kube/config
```

```bash
# first create a new project to deploy to
$ oc new-project <namespace>

# User should be able to create cluster level resources
$ make deploy_local_config IP_ADDRESS=<your_host_ip_address>

# Running locally needs a physical key and cert.
validatingwebhookconfiguration.admissionregistration.k8s.io/managed-service-namespaces-webhook-config created
Server key is at ./build/tmp/server-credentials/server.key.pem
Server cert is at ./build/tmp/server-credentials/server.cert.pem

# make run is set to look for those credentials in ./build/tmp/server-credentials/ by default
# run the server
$ make run

# To clean up the local deployment, run
$ make delete_local_config
```

Again `make deploy_local_config` sets up [the WebhookConfiguration](https://github.com/integr8ly/managed-service-controller/tree/master/admission-webhook/deploy/webhook.external.yaml) for you but you can do these steps manually too. The next steps will create some credentials for a locally running server.

```bash
$ make local_context IP_ADDRESS=<your_host_ip_address>

Base64 encoded CA cert
LS0tLS1CRUdJTiBDRVJUSUZJQ0...


Server key is at ./build/tmp/server-credentials/server.key.pem
Server cert is at ./build/tmp/server-credentials/server.cert.pem
```

Then create the ValidatingWebhookConfiguration:
```bash
# Plug in the host ip and the Base64 encoded CA cert from above.
$ cat deploy/webhook.external.yaml | sed -e 's/<DOMAIN>/<your_host_ip_address>/g' -e 's/<CA_BUNDLE>/LS0tLS1CRUdJTiBDRVJUSUZJQ0.../g' | oc create -f -
```