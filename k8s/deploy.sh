#!/bin/bash

# Ecommerce Kubernetes Deployment Script
# This script automates the deployment of the ecommerce microservices to Kubernetes

set -e

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
NAMESPACE="ecommerce"
REGISTRY="${DOCKER_REGISTRY:-your-registry}"
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

check_dependencies() {
    log_info "Checking dependencies..."
    
    # Check kubectl
    if ! command -v kubectl &> /dev/null; then
        log_error "kubectl not found. Please install kubectl."
        exit 1
    fi
    log_success "kubectl found"
    
    # Check kubectlconnection
    if ! kubectl cluster-info &> /dev/null; then
        log_error "Cannot connect to Kubernetes cluster. Please configure kubectl."
        exit 1
    fi
    log_success "Connected to Kubernetes cluster"
}

update_registry() {
    log_info "Updating Docker registry to: $REGISTRY"
    
    if [ "$REGISTRY" != "your-registry" ]; then
        for file in "$SCRIPT_DIR"/*.yaml; do
            sed -i.bak "s|your-registry|$REGISTRY|g" "$file" && rm "${file}.bak"
        done
        log_success "Registry updated in all manifests"
    else
        log_warning "Using default registry 'your-registry'. Update DOCKER_REGISTRY env var for production."
    fi
}

deploy_namespace_and_secrets() {
    log_info "Creating namespace and secrets..."
    kubectl apply -f "$SCRIPT_DIR/00-namespace-and-secrets.yaml"
    log_success "Namespace and secrets created"
}

deploy_infrastructure() {
    log_info "Deploying infrastructure services (PostgreSQL, Redis, NATS, Elasticsearch, MinIO)..."
    kubectl apply -f "$SCRIPT_DIR/01-infrastructure-services.yaml"
    log_success "Infrastructure services deployed"
    
    log_info "Waiting for infrastructure to be ready (this may take a few minutes)..."
    
    kubectl wait --for=condition=ready pod -l app=postgres -n $NAMESPACE --timeout=300s || log_warning "PostgreSQL not ready yet"
    kubectl wait --for=condition=ready pod -l app=redis -n $NAMESPACE --timeout=300s || log_warning "Redis not ready yet"
    kubectl wait --for=condition=ready pod -l app=nats -n $NAMESPACE --timeout=300s || log_warning "NATS not ready yet"
    
    log_success "Infrastructure services are ready"
}

deploy_configmaps() {
    log_info "Deploying ConfigMaps..."
    kubectl apply -f "$SCRIPT_DIR/02-configmap.yaml"
    log_success "ConfigMaps deployed"
}

deploy_application_services() {
    log_info "Deploying application services..."
    kubectl apply -f "$SCRIPT_DIR/03-application-services-part1.yaml"
    kubectl apply -f "$SCRIPT_DIR/04-application-services-part2.yaml"
    log_success "Application services deployed"
    
    log_info "Waiting for application services to be ready..."
    kubectl wait --for=condition=ready pod -l app=auth -n $NAMESPACE --timeout=300s || log_warning "Auth service not ready yet"
    log_success "Application services are ready"
}

deploy_gateway_and_web() {
    log_info "Deploying API Gateway and Web Frontend..."
    kubectl apply -f "$SCRIPT_DIR/05-api-gateway-and-web.yaml"
    log_success "API Gateway and Web Frontend deployed"
}

deploy_ingress() {
    log_info "Deploying Ingress..."
    kubectl apply -f "$SCRIPT_DIR/06-ingress.yaml"
    log_success "Ingress deployed"
}

deploy_network_policies() {
    log_info "Deploying Network Policies..."
    kubectl apply -f "$SCRIPT_DIR/07-network-policies.yaml"
    log_success "Network Policies deployed"
}

run_migrations() {
    log_info "Running database migrations..."

    # Delete previous migration job if it exists (jobs are immutable)
    kubectl delete job db-migrations -n $NAMESPACE --ignore-not-found=true

    kubectl apply -f "$SCRIPT_DIR/10-migration-job.yaml"
    log_info "Waiting for migration job to complete (timeout: 10 minutes)..."

    if kubectl wait --for=condition=complete job/db-migrations -n $NAMESPACE --timeout=600s; then
        log_success "Database migrations completed successfully"
    else
        log_error "Database migrations failed or timed out"
        log_info "Checking migration job logs..."
        kubectl logs job/db-migrations -n $NAMESPACE || true
        log_error "Aborting deployment due to migration failure"
        exit 1
    fi
}

deploy_hpa() {
    log_info "Deploying Horizontal Pod Autoscalers..."
    kubectl apply -f "$SCRIPT_DIR/08-hpa.yaml"
    log_success "HPA deployed for all services"
}

deploy_pdb() {
    log_info "Deploying Pod Disruption Budgets..."
    kubectl apply -f "$SCRIPT_DIR/09-pdb.yaml"
    log_success "PDB deployed for core services"
}

deploy_monitoring() {
    log_info "Deploying monitoring stack (Prometheus, Grafana, Jaeger, OTel Collector)..."
    kubectl apply -f "$SCRIPT_DIR/11-monitoring.yaml"
    log_success "Monitoring stack deployed"

    log_info "Waiting for monitoring services to be ready..."
    kubectl wait --for=condition=ready pod -l app=prometheus -n $NAMESPACE --timeout=180s || log_warning "Prometheus not ready yet"
    kubectl wait --for=condition=ready pod -l app=grafana -n $NAMESPACE --timeout=180s || log_warning "Grafana not ready yet"
    kubectl wait --for=condition=ready pod -l app=jaeger -n $NAMESPACE --timeout=180s || log_warning "Jaeger not ready yet"
    kubectl wait --for=condition=ready pod -l app=otel-collector -n $NAMESPACE --timeout=180s || log_warning "OTel Collector not ready yet"
    log_success "Monitoring stack is ready"
}

verify_deployment() {
    log_info "Verifying deployment..."
    
    echo ""
    log_info "=== Deployments ==="
    kubectl get deployments -n $NAMESPACE
    
    echo ""
    log_info "=== Pods ==="
    kubectl get pods -n $NAMESPACE
    
    echo ""
    log_info "=== Services ==="
    kubectl get svc -n $NAMESPACE
    
    echo ""
    log_success "Deployment verification complete"
}

show_access_info() {
    echo ""
    log_success "=== Access Information ==="
    echo ""
    echo "Update your /etc/hosts file with the following entries:"
    echo ""
    echo "127.0.0.1 ecommerce.local"
    echo "127.0.0.1 api.ecommerce.local"
    echo "127.0.0.1 kong-admin.ecommerce.local"
    echo "127.0.0.1 mail.ecommerce.local"
    echo "127.0.0.1 minio.ecommerce.local"
    echo ""
    
    if command -v minikube &> /dev/null; then
        echo "For Minikube, use: $(minikube ip) instead of 127.0.0.1"
        echo ""
    fi
    
    echo "Access services at:"
    echo "  Web Frontend: http://ecommerce.local"
    echo "  API Gateway: http://api.ecommerce.local"
    echo "  Kong Admin: http://kong-admin.ecommerce.local"
    echo "  MailHog: http://mail.ecommerce.local"
    echo "  MinIO Console: http://minio.ecommerce.local"
    echo ""
    echo "Monitoring (use kubectl port-forward or add ingress rules):"
    echo "  Prometheus: kubectl port-forward svc/prometheus 9090:9090 -n ecommerce"
    echo "  Grafana:    kubectl port-forward svc/grafana 3000:3000 -n ecommerce  (admin / ecommerce-grafana-admin)"
    echo "  Jaeger UI:  kubectl port-forward svc/jaeger 16686:16686 -n ecommerce"
    echo ""
}

# Main execution
main() {
    log_info "Starting Ecommerce Kubernetes Deployment"
    log_info "Namespace: $NAMESPACE"
    log_info "Registry: $REGISTRY"
    echo ""
    
    check_dependencies
    echo ""
    
    update_registry
    echo ""
    
    deploy_namespace_and_secrets
    echo ""
    
    deploy_infrastructure
    echo ""
    
    deploy_configmaps
    echo ""

    run_migrations
    echo ""

    deploy_application_services
    echo ""

    deploy_gateway_and_web
    echo ""

    deploy_ingress
    echo ""

    deploy_network_policies
    echo ""

    deploy_hpa
    echo ""

    deploy_pdb
    echo ""

    deploy_monitoring
    echo ""

    verify_deployment
    echo ""

    show_access_info
    
    log_success "Deployment completed successfully!"
}

# Run main function
main
