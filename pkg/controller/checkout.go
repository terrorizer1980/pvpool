package controller

import (
	"context"

	pvpoolv1alpha1 "github.com/puppetlabs/pvpool/pkg/apis/pvpool.puppet.com/v1alpha1"
	"github.com/puppetlabs/pvpool/pkg/app"
	"github.com/puppetlabs/pvpool/pkg/obj"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// +kubebuilder:rbac:groups=pvpool.puppet.com,resources=checkouts,verbs=get;list;watch
// +kubebuilder:rbac:groups=pvpool.puppet.com,resources=checkouts/status,verbs=update
// +kubebuilder:rbac:groups=core,resources=persistentvolumes,verbs=get;list;watch;update
// +kubebuilder:rbac:groups=core,resources=persistentvolumeclaims,verbs=get;list;watch;create;update;delete

type CheckoutReconciler struct {
	cl client.Client
}

var _ reconcile.Reconciler = &CheckoutReconciler{}

func (pr *CheckoutReconciler) Reconcile(ctx context.Context, req reconcile.Request) (r reconcile.Result, err error) {
	klog.InfoS("checkout reconciler: starting reconcile for checkout", "checkout", req.NamespacedName)
	defer klog.InfoS("checkout reconciler: ending reconcile for checkout", "checkout", req.NamespacedName)
	defer func() {
		if err != nil {
			klog.ErrorS(err, "checkout reconciler: failed to reconcile checkout", "checkout", req.NamespacedName)
		}
	}()

	checkout := obj.NewCheckout(req.NamespacedName)
	if ok, err := checkout.Load(ctx, pr.cl); err != nil || !ok {
		return reconcile.Result{}, err
	}

	_, err = app.ApplyCheckoutState(ctx, pr.cl, checkout)
	return
}

func NewCheckoutReconciler(cl client.Client) *CheckoutReconciler {
	return &CheckoutReconciler{
		cl: cl,
	}
}

func AddCheckoutReconcilerToManager(mgr manager.Manager) error {
	r := NewCheckoutReconciler(mgr.GetClient())

	return builder.ControllerManagedBy(mgr).
		For(&pvpoolv1alpha1.Checkout{}).
		Owns(&corev1.PersistentVolumeClaim{}).
		Complete(r)
}
