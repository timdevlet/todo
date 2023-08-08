current_date_time=$(shell date)

minikube:
	minikube start
	minikube addons enable ingress
	minikube addons enable ingress-dns
	minikube addons list

minikube-tunnel:
	minikube tunnel

minikube-ip:
	minikube ip

minikube-hosts:
	sudo echo "127.0.0.1 nginx.mini.app" >> /etc/hosts

k8s-apply:
	kubectl apply -f yaml

ingress-logs:
	kubectl logs -n ingress-nginx -l app.kubernetes.io/name=ingress-nginx

cert:
	helm repo add jetstack https://charts.jetstack.io
	kubectl create namespace cert-manager
	helm install cert-manager jetstack/cert-manager --namespace cert-manager --create-namespace --version v1.12.0 --set installCRDs=true

	kubectl get Issuers,ClusterIssuers,Certificates,CertificateRequests,Orders,Challenges --all-namespaces

// Please note: When you use minikube with docker driver (e.g. on an Apple M1 machine), you won’t be able to access the ingress exposed services via minikube IP. You need to run minikube tunnel to make the service available via local ip 127.0.0.1 instead. You will also need to map minikube.data.gov.au to IP 127.0.0.1 instead in /etc/hosts.

helm-install-front:
	helm install front ./charts/front
	helm get manifest front

	helm install --debug --dry-run front ./charts/front

helm-upgrade-front:
	## or -f new_values.yaml  
	helm upgrade -i --set annotations.time="$(current_date_time)" front ./charts/front

readme:
	## https://cloud.yandex.ru/docs/managed-kubernetes/operations/connect/create-static-conf
	## https://www.unisender.com/

pomerium-install:
	kubectl apply -k github.com/pomerium/ingress-controller/config/default\?ref=v0.22.1

