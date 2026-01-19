/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"

	// --- (Kubernetes Objeleri İçin) ---
	appsv1 "k8s.io/api/apps/v1"                   // Deployment oluşturmak için gerekli
	corev1 "k8s.io/api/core/v1"                   // Pod, Container ve Volume tanımları için gerekli
	"k8s.io/apimachinery/pkg/api/errors"          // "NotFound" hatasını yakalamak için gerekli
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1" // ObjectMeta (isim, namespace) için gerekli
	"k8s.io/apimachinery/pkg/types"               // NamespacedName kullanımı için gerekli
	// --------------------------------------------------

	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	webv1 "github.com/imbilalyilmaz/k8s-staticsite-operator/api/v1"
)

// StaticSiteReconciler reconciles a StaticSite object
type StaticSiteReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=web.mydomain.com,resources=staticsites,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=web.mydomain.com,resources=staticsites/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=web.mydomain.com,resources=staticsites/finalizers,verbs=update
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the StaticSite object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.22.4/pkg/reconcile
// Reconcile fonksiyonu
func (r *StaticSiteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// 1. Cluster'daki StaticSite objesini (CR) çekelim
	staticSite := &webv1.StaticSite{}
	err := r.Get(ctx, req.NamespacedName, staticSite)
	if err != nil {
		if errors.IsNotFound(err) {
			// Obje silinmişse bir şey yapmaya gerek yok
			return ctrl.Result{}, nil
		}
		// Başka bir hata varsa tekrar dene
		return ctrl.Result{}, err
	}

	// 2. İstenen Deployment tanımını hazırla (Helper fonksiyonu aşağıda)
	dep := r.deploymentForStaticSite(staticSite)

	// 3. Bu Deployment zaten var mı diye kontrol et
	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: dep.Name, Namespace: dep.Namespace}, found)

	if err != nil && errors.IsNotFound(err) {
		// A) Deployment YOK, o zaman YARAT (Create)
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.Create(ctx, dep)
		if err != nil {
			return ctrl.Result{}, err
		}
		// Başarılı yaratıldı, tekrar sıraya girmeye gerek yok
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// 4. (Opsiyonel ama Önemli) Replicas sayısı değişmiş mi? (Update Mantığı)
	size := staticSite.Spec.Replicas
	if *found.Spec.Replicas != size {
		found.Spec.Replicas = &size
		err = r.Update(ctx, found)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	// --- SERVICE MANTIĞI BAŞLANGICI ---

	// 5. Service tanımını hazırla
	svc := r.serviceForStaticSite(staticSite)

	// 6. Service var mı kontrol et
	foundSvc := &corev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: svc.Name, Namespace: svc.Namespace}, foundSvc)

	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Service", "Service.Namespace", svc.Namespace, "Service.Name", svc.Name)
		err = r.Create(ctx, svc)
		if err != nil {
			return ctrl.Result{}, err
		}
		// Service yaratıldı, tekrar sıraya girmeye gerek yok (Deployment zaten requeue yapmıştı gerekirse)
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// --- SERVICE MANTIĞI BİTİŞİ ---

	// --- INGRESS MANTIĞI BAŞLANGICI ---

	// 7. Ingress tanımını hazırla
	ing := r.ingressForStaticSite(staticSite)

	// 8. Ingress var mı kontrol et
	foundIng := &networkingv1.Ingress{}
	err = r.Get(ctx, types.NamespacedName{Name: ing.Name, Namespace: ing.Namespace}, foundIng)

	if err != nil && errors.IsNotFound(err) {
		log.Info("Creating a new Ingress", "Ingress.Namespace", ing.Namespace, "Ingress.Name", ing.Name)
		err = r.Create(ctx, ing)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	}

	// --- INGRESS MANTIĞI BİTİŞİ ---

	// Her şey yolunda
	return ctrl.Result{}, nil
}

// deploymentForStaticSite, StaticSite objesine bakarak bir Deployment objesi üretir.
func (r *StaticSiteReconciler) deploymentForStaticSite(m *webv1.StaticSite) *appsv1.Deployment {
	ls := labelsForStaticSite(m.Name)
	replicas := m.Spec.Replicas

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					// Volume Tanımı: İki container arasında dosya paylaşmak için
					Volumes: []corev1.Volume{{
						Name: "site-data",
						VolumeSource: corev1.VolumeSource{
							EmptyDir: &corev1.EmptyDirVolumeSource{},
						},
					}},
					// Init Container: Önce klasörü temizler, sonra repoyu çeker (Idempotent)
					InitContainers: []corev1.Container{{
						Name:    "git-cloner",
						Image:   "alpine/git",
						Command: []string{"/bin/sh", "-c"},
						// Önce 'rm -rf' ile eski dosyaları sil, SONRA (&&) 'git clone' yap
						Args: []string{"rm -rf /repo/* && git clone --single-branch -- " + m.Spec.GitRepo + " /repo"},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "site-data",
							MountPath: "/repo",
						}},
					}},
					// Main Container: Nginx çalışır, çekilen dosyaları sunar
					Containers: []corev1.Container{{
						Image: "nginx:alpine",
						Name:  "nginx",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 80,
							Name:          "http",
						}},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "site-data",
							MountPath: "/usr/share/nginx/html", // Nginx'in varsayılan klasörü
						}},
					}},
				},
			},
		},
	}

	// OwnerReference: StaticSite silinirse, Deployment da silinsin diye.
	ctrl.SetControllerReference(m, dep, r.Scheme)
	return dep
}

// Helper: Label üretici
func labelsForStaticSite(name string) map[string]string {
	return map[string]string{"app": "staticsite", "staticsite_cr": name}
}

// SetupWithManager sets up the controller with the Manager.
func (r *StaticSiteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webv1.StaticSite{}).
		Named("staticsite").
		Complete(r)
}

// serviceForStaticSite, Pod'lara erişim sağlayacak Service objesini üretir.
func (r *StaticSiteReconciler) serviceForStaticSite(m *webv1.StaticSite) *corev1.Service {
	ls := labelsForStaticSite(m.Name)

	svc := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
		},
		Spec: corev1.ServiceSpec{
			Selector: ls,
			Ports: []corev1.ServicePort{{
				Port:       80,
				TargetPort: intstr.FromInt(80),
				Protocol:   corev1.ProtocolTCP,
			}},
			// Type: ClusterIP (Varsayılan). Ingress ile dışarı açacağız.
			Type: corev1.ServiceTypeClusterIP,
		},
	}

	// OwnerReference: StaticSite silinirse, Service de silinsin.
	ctrl.SetControllerReference(m, svc, r.Scheme)
	return svc
}

// ingressForStaticSite, siteyi dışarı açacak Ingress kuralını üretir.
func (r *StaticSiteReconciler) ingressForStaticSite(m *webv1.StaticSite) *networkingv1.Ingress {
	ls := labelsForStaticSite(m.Name)

	pathType := networkingv1.PathTypePrefix

	ing := &networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      m.Name,
			Namespace: m.Namespace,
			Labels:    ls,
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{{
				Host: m.Name + ".local",
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{{
							Path:     "/",
							PathType: &pathType,
							Backend: networkingv1.IngressBackend{
								Service: &networkingv1.IngressServiceBackend{
									Name: m.Name,
									Port: networkingv1.ServiceBackendPort{
										Number: 80,
									},
								},
							},
						}},
					},
				},
			}},
		},
	}

	ctrl.SetControllerReference(m, ing, r.Scheme)
	return ing
}
