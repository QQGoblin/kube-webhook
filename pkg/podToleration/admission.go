package podToleration

import (
	"context"
	"encoding/json"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type podToleration struct {
	Client  client.Client
	decoder *admission.Decoder
}

func NewHandler(c client.Client) *podToleration {
	return &podToleration{
		Client: c,
	}
}

func (a *podToleration) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}

	klog.Infof("pod toleration hook: UID<%v> %s", req.UID, req.Kind.String())
	err := a.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if pod.Spec.Tolerations == nil {
		pod.Spec.Tolerations = make([]corev1.Toleration, 0)
	}

	pod.Spec.Tolerations = append(pod.Spec.Tolerations, corev1.Toleration{
		Key:      "lqingcloud.cn/not-ready",
		Operator: "Exists",
		Effect:   corev1.TaintEffectNoSchedule,
	})

	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

func (a *podToleration) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
	return nil
}
