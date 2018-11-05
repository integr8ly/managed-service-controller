package server

import (
	"encoding/json"
	"fmt"
	"github.com/golang/glog"
	clients "github.com/integr8ly/managed-services-controller/admission-webhook/pkg/clients/v1alpha1"
	utils "github.com/integr8ly/managed-services-controller/admission-webhook/pkg/utils/v1alpha1"
	"io/ioutil"
	"k8s.io/api/admission/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"net/http"
)

func msnValidator(k8sClient kubernetes.Interface) func(http.ResponseWriter, *http.Request) {
	msnc := clients.ManagedServiceNamespaceClient{
		K8sClient: k8sClient,
	}
	return func(w http.ResponseWriter, r *http.Request) {
		glog.V(2).Info("Validating ManagedServiceNamespace.")
		serve(w, r, msnc.Validate)
	}
}

func ListenAndServeTLS(config Config) error {
	//TODO: Better server/router?
	http.HandleFunc("/validate/msn", msnValidator(config.K8sClient))

	server := &http.Server{
		Addr:      fmt.Sprintf(":%s", config.Port),
		TLSConfig: configTLS(config),
	}

	err := server.ListenAndServeTLS("", "")
	if err != nil {
		return err
	}

	return nil
}

type admitFunc func(v1beta1.AdmissionReview) *v1beta1.AdmissionResponse

func serve(w http.ResponseWriter, r *http.Request, admit admitFunc) {
	var body []byte
	if r.Body != nil {
		if data, err := ioutil.ReadAll(r.Body); err == nil {
			body = data
		}
	}

	// verify the content type is accurate
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		glog.Errorf("contentType=%s, expect application/json", contentType)
		return
	}

	var reviewResponse *v1beta1.AdmissionResponse
	ar := v1beta1.AdmissionReview{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(body, nil, &ar); err != nil {
		glog.Error(err)
		reviewResponse = utils.ToAdmissionResponse(err)
	} else {
		reviewResponse = admit(ar)
	}

	response := v1beta1.AdmissionReview{}
	if reviewResponse != nil {
		response.Response = reviewResponse
		response.Response.UID = ar.Request.UID
	}
	// reset the Object and OldObject, they are not needed in a response.
	ar.Request.Object = runtime.RawExtension{}
	ar.Request.OldObject = runtime.RawExtension{}

	resp, err := json.Marshal(response)
	if err != nil {
		glog.Error(err)
	}
	if _, err := w.Write(resp); err != nil {
		glog.Error(err)
	}
}
