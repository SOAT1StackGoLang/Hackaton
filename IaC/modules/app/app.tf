terraform {
  required_version = ">= 0.13"

  required_providers {
    kubectl = {
      source  = "gavinbunney/kubectl"
      version = ">= 1.7.0"
    }
  }
}

locals {
  namespace = var.project_name

  svc_hackaton_migs_image = "${var.image_registry}/hackaton:migs-${var.svc_hackaton_image_tag}"
  svc_hackaton_image      = "${var.image_registry}/hackaton:svc-${var.svc_hackaton_image_tag}"

  lb_service_name_hackaton = "lb-hackaton-svc"
  lb_service_port_hackaton = 8080

  svc_hackaton_port = 8080

  svc_hackaton_svc = "svc-hackaton-svc"

  svc_hackaton_uri = "http://${local.svc_hackaton_svc}.${local.namespace}.svc.cluster.local:${local.svc_hackaton_port}"

  kvstore_db_svc_hackaton = 10

  postgres_host     = var.database_host
  postgres_port     = var.database_port
  postgres_user     = var.database_username
  postgres_password = var.database_password
  postgres_db       = var.database_name

  ### deployment annotations time
  redeploy_annotation = var.redeploy_annotation
}


#---------------------------------------------------------------------------------------------------
#  Namespace
#---------------------------------------------------------------------------------------------------

resource "kubectl_manifest" "namespace" {
  yaml_body = <<YAML

apiVersion: v1
kind: Namespace
metadata:
  name: ${local.namespace}
  annotations:
    argocd.argoproj.io/sync-wave: "-1"
    linkerd.io/inject: enabled
  labels:
    app: ${var.project_name}

YAML
}


#---------------------------------------------------------------------------------------------------
#  Secrets
#---------------------------------------------------------------------------------------------------

resource "kubectl_manifest" "secrets" {
  yaml_body = <<YAML

apiVersion: v1
kind: Secret
metadata:
  name: ${var.project_name}-secret
  namespace: ${local.namespace}
type: Opaque
stringData:
  DB_URI: "host=${local.postgres_host} port=${local.postgres_port} user=${local.postgres_user} password=${local.postgres_password} dbname=${local.postgres_db} sslmode=require"

YAML
}

#---------------------------------------------------------------------------------------------------
# Deployment Hackaton
#---------------------------------------------------------------------------------------------------

resource "kubectl_manifest" "svc_hackaton_deployment" {
  yaml_body = <<YAML

apiVersion: apps/v1
kind: Deployment
metadata:
  name: svc-hackaton
  namespace: ${local.namespace}
  labels:
    app: svc-hackaton
  annotations:
    kubectl.kubernetes.io/restartedAt: ${local.redeploy_annotation}
spec:
  selector:
    matchLabels:
      app: svc-hackaton
  replicas: 1
  strategy:
    rollingUpdate:
      maxSurge: 25%
      maxUnavailable: 25%
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: svc-hackaton
      annotations:
        kubectl.kubernetes.io/restartedAt: ${local.redeploy_annotation}
    spec:
      initContainers:
        - name: svc-hackaton-migs
          image: ${local.svc_hackaton_migs_image}
          imagePullPolicy: Always
          securityContext:
            readOnlyRootFilesystem: true
            allowPrivilegeEscalation: false
            runAsNonRoot: true
            runAsUser: 10000
            capabilities:
              drop:
                - ALL
          resources:
            requests:
              cpu: 10m
              memory: 25Mi
            limits:
              cpu: 100m
              memory: 100Mi
          envFrom:
            - secretRef:
                name: ${var.project_name}-secret
      containers:
        - name: svc-hackaton
          image: ${local.svc_hackaton_image}
          imagePullPolicy: Always
          terminationGracePeriodSeconds: 15
          securityContext:
            readOnlyRootFilesystem: true
            allowPrivilegeEscalation: false
            runAsNonRoot: true
            runAsUser: 10000
            capabilities:
              drop:
                - ALL
          resources:
            requests:
              cpu: 100m
              memory: 128Mi
            limits:
              cpu: 500m
              memory: 1000Mi
          livenessProbe:
            tcpSocket:
              port: 8080
            initialDelaySeconds: 5
            timeoutSeconds: 5
            successThreshold: 1
            failureThreshold: 3
            periodSeconds: 10
          readinessProbe:
            tcpSocket:
              port: 8080
            initialDelaySeconds: 5
            timeoutSeconds: 2
            successThreshold: 1
            failureThreshold: 3
            periodSeconds: 10
          envFrom:
            - secretRef:
                name: ${var.project_name}-secret
          ports:
            - containerPort: 8080
              name: web
          
      restartPolicy: Always
YAML
}


resource "kubectl_manifest" "svc_hackaton_service" {
  yaml_body = <<YAML

apiVersion: v1
kind: Service
metadata:
  name: ${local.svc_hackaton_svc}
  namespace: ${local.namespace}
spec:
  selector:
    app: svc-hackaton
  type: ClusterIP
  ports:
  - protocol: TCP
    port: ${local.svc_hackaton_port}
    targetPort: 8080

YAML
}

#---------------------------------------------------------------------------------------------------
# Load Balancers
#---------------------------------------------------------------------------------------------------

resource "kubectl_manifest" "lb-hackaton" {
  yaml_body = <<YAML

apiVersion: v1
kind: Service
metadata:
  name: ${local.lb_service_name_hackaton}
  namespace: ${local.namespace}
  annotations:
    service.beta.kubernetes.io/aws-load-balancer-name: ${local.lb_service_name_hackaton}
    service.beta.kubernetes.io/aws-load-balancer-type: nlb

spec:
  selector:
    app: svc-hackaton
  type: LoadBalancer
  ports:
  - protocol: TCP
    port: ${local.lb_service_port_hackaton}
    targetPort: 8080

YAML
}

#---------------------------------------------------------------------------------------------------
#  HPA
#---------------------------------------------------------------------------------------------------

resource "kubectl_manifest" "svc_hackaton_hpa" {
  yaml_body = <<YAML

apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: svc-hackaton-hpa
  namespace: ${local.namespace}
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: svc-hackaton
  minReplicas: 1
  maxReplicas: 50
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: AverageValue
        averageUtilization: 200
  - type: Resource
    resource:
      name: memory
      target:
        type: AverageValue
        averageUtilization: 200
YAML
}