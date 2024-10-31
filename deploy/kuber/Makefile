#!/usr/bin/env make

.PHONY: run_website install_kind install_kubectl create_kind_cluster \
	create_docker_registry connect_registry_to_kind_network \
	connect_registry_to_kind create_kind_cluster_with_registry \
  install_ingress_controller install_app

install_kubectl:
	brew install kubectl || true;

install_kind:
	curl -o ./kind https://github.com/kubernetes-sigs/kind/releases/download/v0.11.1/kind-darwin-arm64

connect_registry_to_kind_network:
	docker network connect kind local-registry || true;

connect_registry_to_kind: connect_registry_to_kind_network
	kubectl apply -f ./kind_configmap.yaml;

create_docker_registry:
	if ! docker ps | grep -q 'local-registry'; \
	then docker run -d -p 5001:5000 --name local-registry --restart=always registry:2; \
	else echo "---> local-registry is already running. There's nothing to do here."; \
	fi

create_kind_cluster: install_kind install_kubectl create_docker_registry
	kind create cluster --image=kindest/node:v1.21.12 --name shortlink.com --config ./kind_config.yaml || true
	kubectl get nodes

delete_docker_registry:
	docker stop local-registry && docker rm local-registry || true;

delete_kind_cluster: delete_docker_registry
	kind delete cluster --name shortlink.com

# deploy rmi go_app
# deploy rmi localhost:5001/shortlink.com
# deploy commit  go_app go_app
# deploy build -t go_app .
# deploy build -f deploy/Dockerfile -t go_app .
# deploy build -f deploy/Dockerfile -t go_app .
# deploy build --no-cache -f deploy/Dockerfile -t go_app .
# deploy build --no-cache -f ../Dockerfile -t go_app .
create_docker_tag:
	docker tag go_app localhost:5001/shortlink.com

push_docker_container:
	docker push localhost:5001/shortlink.com

install_ingress_controller:
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml && \
	sleep 5 && \
	kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s

delete_ingress_nginx:
	kubectl delete -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml

apply_ingress:
	kubectl apply -f ingress.yaml

delete_ingress:
	kubectl delete -f ingress.yaml

log_ingress:
	kubectl get all -n ingress-nginx

build_docker:
	docker build -f Dockerfile -t go_app ../..
	#docker build -f Dockerfile --no-cache -t go_app ../..

install_helm_app:
	helm install  shortlink-website ./chart

uninstall_helm:
	helm uninstall  shortlink-website

delete:
	kubectl delete all -l app=shortlink.com

descr:
	kubectl describe pod -l app=shortlink-website

upgrade_app:
	helm upgrade --atomic --install  shortlink-website ./chart

upgrade_all_app: build_docker create_docker_tag push_docker_container upgrade_app

install_all_cluster: create_kind_cluster_with_registry install_helm_app install_ingress_controller

uninstall_cluster: delete_kind_cluster delete_ingress_nginx uninstall_helm

create_kind_cluster_with_registry: build_docker create_docker_tag
	$(MAKE) create_kind_cluster && $(MAKE) connect_registry_to_kind && $(MAKE) push_docker_container

migration:
	kubectl exec -it $$(kubectl get pod -l app=shortlink-website -o jsonpath='{.items[0].metadata.name}') -- migrate -path database/migration/ -database "postgresql://postgres:postgres@postgres:5432/slink?sslmode=disable" -verbose up

# When you use --atomic with helm upgrade or helm install, Helm will automatically roll back the release if the deployment fails,


# kubectl get all -n ingress-nginx

# kubectl get pods -l app=shortlink-website
# kubectl exec -it shortlink-website-659dcbb44d-fgjr6  -- ls /website
# helm template ./chart
# helm history shortlink-website  // list versions
# helm list // current version
# helm status shortlink-website
# helm history shortlink-website
# helm rollback shortlink-website 1

# deploy commit
#  kubectl get service shortlink-svc

# kubectl get pods
# kubectl get pods -o wide --all-namespaces
# kubectl get pods --all-namespaces
# kubectl describe pod
# kubectl describe service shortlink-svc
# kubectl logs shortlink-website-778d6d5f8c-d4qtt
# kubectl get all -n ingress-nginx

#kubectl delete all -l app.kubernetes.io/name=ingress-nginx

# kind get clusters


# Native Tools

# kubectl get deployments -- list all deployments
# kubectl get deployments -o wide -- list all deployments with details
# kubectl get namespaces --

# I need to run this in new namespace

