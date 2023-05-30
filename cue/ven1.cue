package main

deployment: {
  apiVersion: "apps/v1"
  kind:       "Deployment"
  metadata: {
    name:   string
    labels: {
      app: string
    }
  }
  spec: {
    replicas: int
    selector: {
      matchLabels: {
        app: string
      }
    }
    template: {
      metadata: {
        labels: {
          app: string
        }
      }
      spec: {
        containers: [...#Container]
      }
    }
  }
}

service: {
  apiVersion: "v1"
  kind:       "Service"
  metadata: {
    name: string
  }
  spec: {
    selector: {
      app: string
    }
    ports: [...#ServicePort]
  }
}

ingress: {
  apiVersion: "networking.k8s.io/v1"
  kind:       "Ingress"
  metadata: {
    annotations: {
      "kubernetes.io/ingress.class": string
    }
    name: string
  }
  spec: {
    rules: [...#IngressRule]
    tls: [...#IngressTLS]
  }
}

#Container: {
  name:  string
  image: string
  env: [...{
    name:  string
    value: string
  }]
  resources: {
    limits:   #Resources
    requests: #Resources
  }
}

#Resources: {
  cpu:    string
  memory: string
}

#ServicePort: {
  protocol:    string
  port:        int
  targetPort:  int
}

#IngressRule: {
  host: string
  http: {
    paths: [...#HTTPPath]
  }
}

#HTTPPath: {
  pathType:    string
  path:        string
  backend: {
    service: {
      name:  string
      port: {
        number: int
      }
    }
  }
}

#IngressTLS: {
  hosts: [...string]
}

QUEUE_SIZE: [7, 10, 15]
CPU: [0.5, 1, 2]
RAM: ["800Mi", "1Gi", "2Gi"]

manifest: {
  deployment: [string]: deployment
  service:    [string]: service
  ingress:    [string]: ingress
}
