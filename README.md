# Kubernetes StaticSite Operator

![Project Screenshot](assets/screenshot1.png) ![Project
Screenshot](assets/screenshot2.png)

This project is a **Kubernetes Operator** designed to deploy static
websites (HTML/CSS/JS) on Kubernetes within seconds.

It is developed using **Golang** and **Kubebuilder**, and follows
Kubernetes' core **Reconciliation Loop** principle.

## 🚀 What Does It Do?

The user only needs to create a simple YAML resource (`StaticSite`). The
Operator automatically handles the following in the background:

1.  **Deployment:** Creates Pods running an Nginx container with a
    Git-Sync **Init Container**.
2.  **Service:** Exposes the Pods internally using a ClusterIP Service.
3.  **Ingress:** Creates routing rules to expose the site to the outside
    world (or host machine).
4.  **Self-Healing:** If a Deployment or Service is deleted manually,
    the Operator detects the drift and recreates it within milliseconds.

## 🏗 Architecture

-   **Language:** Go (Golang)
-   **Framework:** Kubebuilder / Controller-Runtime
-   **Pattern:** Init Container Pattern (git-clone → shared-volume →
    nginx)

## 🛠 Installation & Usage

### Prerequisites

-   Kubernetes Cluster (Minikube, Kind, or Cloud)
-   kubectl
-   Go 1.20+

### 1. Install CRDs

``` bash
make install
```

### 2. Run the Operator

``` bash
make run
```

### 3. Create a Sample Static Site

``` yaml
apiVersion: web.mydomain.com/v1
kind: StaticSite
metadata:
  name: my-site
spec:
  gitRepo: "https://github.com/cloudacademy/static-website-example"
  replicas: 1
```

``` bash
kubectl apply -f config/samples/web_v1_staticsite.yaml
```

## 💡 Why It Matters

This project goes beyond deploying a static website. It demonstrates the
ability to **extend Kubernetes itself**, not just consume it.

By implementing a **Custom Resource Definition (CRD)** and a **Custom
Controller**, this Operator shows how real-world, production-grade
platforms are built. Instead of manually managing Kubernetes YAML
manifests, the desired state is declared once and continuously enforced
by the controller.

### What this project proves

-   Deep understanding of **Kubernetes internals** (API Server,
    Controllers, Reconciliation Loop)
-   Ability to design **self-healing, declarative infrastructure**
-   Strong **Golang** skills applied to platform-level problems
-   Practical use of the **Operator Pattern**
-   A **Platform / SRE mindset** focused on reducing operational
    complexity

This is the same architectural approach used by widely adopted
Kubernetes projects such as **Ingress Controllers, Cert-Manager, ArgoCD,
and Prometheus Operator**.

------------------------------------------------------------------------

Developed by **Bilal Yılmaz**

------------------------------------------------------------------------

------------------------------------------------------------------------

# Kubernetes StaticSite Operator (Türkçe)

Bu proje, Kubernetes üzerinde statik web sitelerini (HTML/CSS/JS)
saniyeler içinde yayına almak için geliştirilmiş bir **Kubernetes
Operator** projesidir.

**Golang** ve **Kubebuilder** kullanılarak geliştirilmiştir ve
Kubernetes'in temel çalışma prensibi olan **Reconciliation Loop (Uzlaşma
Döngüsü)** mantığıyla çalışır.

## 🚀 Ne Yapar?

Kullanıcı yalnızca basit bir YAML kaynağı (`StaticSite`) oluşturur.
Operator arka planda aşağıdaki bileşenleri otomatik olarak yönetir:

1.  **Deployment:** Nginx container'ı ve Git-Sync **Init Container**
    içeren Pod'ları oluşturur.
2.  **Service:** Pod'lara erişim için ClusterIP Service tanımlar.
3.  **Ingress:** Siteyi dış dünyaya (veya host makineye) açan
    yönlendirme kurallarını oluşturur.
4.  **Self-Healing:** Deployment veya Service manuel olarak silinirse,
    Operator bunu algılar ve çok kısa sürede yeniden oluşturur.

## 🏗 Mimari

-   **Dil:** Go (Golang)
-   **Framework:** Kubebuilder / Controller-Runtime
-   **Pattern:** Init Container Pattern (git-clone → shared-volume →
    nginx)

## 🛠 Kurulum ve Çalıştırma

### Gereksinimler

-   Kubernetes Cluster (Minikube, Kind veya Cloud)
-   kubectl
-   Go 1.20+

### 1. CRD'leri Yükle

``` bash
make install
```

### 2. Operator'ı Çalıştır

``` bash
make run
```

### 3. Örnek Bir Statik Site Oluştur

``` yaml
apiVersion: web.mydomain.com/v1
kind: StaticSite
metadata:
  name: benim-sitem
spec:
  gitRepo: "https://github.com/cloudacademy/static-website-example"
  replicas: 1
```

``` bash
kubectl apply -f config/samples/web_v1_staticsite.yaml
```

## 💡 Neden Önemli?

Bu proje sadece bir statik site yayına almak için yazılmamıştır. Asıl
değer, **Kubernetes'i kullanan değil, Kubernetes'i genişleten** bir
çözüm olmasıdır.

Bir **Custom Resource Definition (CRD)** ve **Custom Controller**
geliştirilerek, Kubernetes üzerinde gerçek hayatta kullanılan **platform
seviyesinde** bir sistem tasarlanmıştır. Kullanıcı yalnızca "istenen
durumu" tanımlar, Operator bu durumu sürekli olarak korur.

### Bu proje neyi kanıtlar?

-   **Kubernetes iç mimarisi** bilgisi (API Server, Controller'lar,
    Reconciliation Loop)
-   **Self-healing** ve **declarative altyapı** tasarımı
-   **Golang** ile sistem ve altyapı seviyesinde geliştirme becerisi
-   Cloud-native dünyada yaygın olan **Operator Pattern** hakimiyeti
-   Operasyonel karmaşıklığı azaltmayı hedefleyen **Platform / SRE bakış
    açısı**

Bu yaklaşım; **Ingress Controller, Cert-Manager, ArgoCD ve Prometheus
Operator** gibi gerçek dünya Kubernetes projeleriyle aynı mimari
anlayışı paylaşır.

------------------------------------------------------------------------

Geliştiren: **Bilal Yılmaz**
