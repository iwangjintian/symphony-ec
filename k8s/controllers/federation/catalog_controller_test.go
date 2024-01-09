package federation

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	federationv1 "gopls-workspace/apis/federation/v1"

	k8smodel "github.com/azure/symphony/k8s/apis/model/v1"
	"github.com/stretchr/testify/assert"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestCatalogReconciler_WithoutUpdate(t *testing.T) {

	//Create a Catalog instance
	catalog := &federationv1.Catalog{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-catalog",
			Namespace: "default",
			DeletionTimestamp: &metav1.Time{
				Time: metav1.Now().Add(1 * time.Hour),
			},
		},
		Spec: k8smodel.CatalogSpec{
			SiteId: "site1",
			Type:   "type1",
			Name:   "name1",
			Metadata: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			ParentName: "parent1",
			Generation: "generation1",
		},
		Status: federationv1.CatalogStatus{
			Properties: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}

	scheme := runtime.NewScheme()

	// Add the clientgo scheme to it (this includes standard kubernetes types)
	_ = clientgoscheme.AddToScheme(scheme)

	// Add your Catalog type to the scheme
	_ = federationv1.AddToScheme(scheme)

	// Create a fake client to mock API calls, and prepopulate it with the Catalog instance
	cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(catalog).Build()
	// Create a Reconciler instance
	r := &CatalogReconciler{
		Client: cl,
	}

	// Create a context and request
	ctx := context.Background()
	req := ctrl.Request{
		NamespacedName: client.ObjectKey{
			Namespace: catalog.ObjectMeta.Namespace,
			Name:      catalog.ObjectMeta.Name,
		},
	}

	// Call Reconcile
	result, err := r.Reconcile(ctx, req)

	// Assert the expected results
	assert.NoError(t, err)
	assert.Equal(t, ctrl.Result{}, result)
}

func TestCatalogReconciler_WithUpdate(t *testing.T) {

	//Create a Catalog instance
	catalog := &federationv1.Catalog{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-catalog",
			Namespace: "default",
		},
		Spec: k8smodel.CatalogSpec{
			SiteId: "site1",
			Type:   "type1",
			Name:   "name1",
			Metadata: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			ParentName: "parent1",
			Generation: "generation1",
		},
		Status: federationv1.CatalogStatus{
			Properties: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}

	scheme := runtime.NewScheme()

	// Add the clientgo scheme to it (this includes standard kubernetes types)
	_ = clientgoscheme.AddToScheme(scheme)

	// Add your Catalog type to the scheme
	_ = federationv1.AddToScheme(scheme)

	// Create a fake client to mock API calls, and prepopulate it with the Catalog instance
	cl := fake.NewClientBuilder().WithScheme(scheme).WithObjects(catalog).Build()
	// Create a Reconciler instance
	r := &CatalogReconciler{
		Client: cl,
	}

	// Create a mock server to mock API calls
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		} else {
			w.Write([]byte("OK"))
		}
	}))
	defer ts.Close()

	// Create a context and request
	ctx := context.Background()
	req := ctrl.Request{
		NamespacedName: client.ObjectKey{
			Namespace: catalog.ObjectMeta.Namespace,
			Name:      catalog.ObjectMeta.Name,
		},
	}

	// Call Reconcile
	result, err := r.Reconcile(ctx, req)

	// Assert the expected results
	assert.NoError(t, err)
	assert.Equal(t, ctrl.Result{}, result)
}
