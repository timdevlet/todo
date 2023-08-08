local tls =  [
  {
    "hosts": [
      "nginx.mini.ru"
    ],
    "secretName": "nginx-app-cert-tls"
  }
];

{
  "apiVersion": "networking.k8s.io/v1",
  "kind": "Ingress",
  "metadata": {
    "name": "hello-world",
    "namespace": "default",
    "annotations": {
      "nginx.ingress.kubernetes.io/rewrite-target": "/$1",
      "cert-manager.io/cluster-issuer": "selfsigned-issuer"
    }
  },
  "spec": {
    "rules": [
      {
        "host": "nginx.mini.ru",
        "http": {
          "paths": [
            {
              "pathType": "Prefix",
              "path": "/",
              "backend": {
                "service": {
                  "name": "nginx",
                  "port": {
                    "number": 80
                  }
                }
              }
            }
          ]
        }
      }
    ],
    "tls": tls
  }
}