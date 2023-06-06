ven_name : "ven1"
backend_queue_size: "5"
backend_cpus: "0.5"
backend_ram: "512Mi"

db_name_label: ven_name + "-db"
backend_name_label: ven_name + "-backend"
backend_service_name: ven_name + "-backend-service"
db_service_name: ven_name + "-db-service"
ingress_name: ven_name + "-backend-ingress"
ingress_host: "dgvkh-"+ven_name+".nrp-nautilus.io"

backend_image: "39dj29dl2d9l2/vec-ven:6"
db_image: "39dj29dl2d9l2/ven-db:1.0"

container_port: 8080
service_port: 8080

mysql_port_env: "3306"
mysql_port: 3306
mysql_user: "root"
mysql_password: "root_password_vec"
mysql_dbname: "app_db"



#Metadata: {
    name:   string
    labels: { app: string }
}

#Spec: {
    replicas: 1
    selector: matchLabels: {
        app: string
    }
    template: {
        metadata: labels: {
            app: string
        }
        spec: {
            containers: [{
                name:  string
                image: string
                ports: [{ containerPort: int }]
                env: [...{
                    name:  string
                    value: string
                }]
                resources: {
                    limits:   { cpu: string, memory: string }
                    requests: { cpu: string, memory: string }
                }
            }]
        }
    }
}

#Deployment: {
    apiVersion: "apps/v1"
    kind:       "Deployment"
    metadata:   #Metadata
    spec:       #Spec
}

#Service: {
    apiVersion: "v1"
    kind:       "Service"
    metadata: {
        name: string
    }
    spec: {
        selector: {
            app: string
        }
        ports: [{
            protocol:   "TCP"
            port:       int
            targetPort: int
        }]
    }
}

#DBDeployment: #Deployment

#Ingress: {
    apiVersion: "networking.k8s.io/v1"
    kind:       "Ingress"
    metadata: {
        annotations: {
            "kubernetes.io/ingress.class": "haproxy"
        }
        name: string
    }
    spec: {
        rules: [...{
            host: string
            http: {
                paths: [...{
                    pathType: "ImplementationSpecific"
                    path:     "/"
                    backend: {
                        service: {
                            name: string
                            port: {
                                number: int
                            }
                        }
                    }
                }]
            }
        }]
        tls: [...{
            hosts: [...string]
        }]
    }
}

removeEdgeBackend: #Deployment & {
    metadata: name: backend_name_label
    metadata: labels: app: backend_name_label
    spec: selector: matchLabels: app: backend_name_label
    spec: template: metadata: labels: app: backend_name_label
    spec: template: spec: containers: [{
        name:  "backend"
        image: backend_image
        ports: [{ containerPort: container_port }]
        env: [
            { name: "MYSQL_USER",      value: mysql_user },
            { name: "MYSQL_PASSWORD",  value: mysql_password },
            { name: "MYSQL_HOST",      value: db_service_name },
            { name: "MYSQL_PORT",      value: mysql_port_env },
            { name: "MYSQL_DBNAME",    value: mysql_dbname },
            { name: "QUEUE_SIZE",      value: backend_queue_size },
            { name: "CPUS",            value: backend_cpus },
            { name: "RAM",             value: backend_ram }
        ]
        resources: {
            limits:   { cpu: backend_cpus, memory: backend_ram }
            requests: { cpu: backend_cpus, memory: backend_ram }
        }
    }]
}

removeEdgeBackendService: #Service & {
    metadata: name: backend_service_name
    spec: selector: app: backend_name_label
    spec: ports: [{ protocol: "TCP", port: service_port, targetPort: container_port }]
}

removeEdgeDB: #DBDeployment & {
    metadata: name: db_name_label
    metadata: labels: app: db_name_label
    spec: selector: matchLabels: app: db_name_label
    spec: template: metadata: labels: app: db_name_label
    spec: template: spec: containers: [{
        name:  "db"
        image: db_image
        ports: [{ containerPort: mysql_port }]
        env: [
            { name: "MYSQL_ROOT_PASSWORD", value: mysql_password },
            { name: "MYSQL_DATABASE", value: mysql_dbname }
        ]
        resources: {
            limits:   { cpu: "1", memory: "1Gi" }
            requests: { cpu: "1", memory: "1Gi" }
        }
    }]
}

removeEdgeDBService: #Service & {
    metadata: name: db_service_name
    spec: selector: app: db_name_label
    spec: ports: [{ protocol: "TCP", port: mysql_port, targetPort: mysql_port }]
}

removeEdgeIngress: #Ingress & {
    metadata: name: ingress_name
    spec: {
        rules: [{
            host: ingress_host
            http: {
                paths: [{
                    backend: {
                        service: {
                            name: backend_service_name
                            port: { number: service_port }
                        }
                    }
                }]
            }
        }]
        tls: [{
            hosts: [ingress_host]
        }]
    }
}

