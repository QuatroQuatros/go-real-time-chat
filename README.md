# Real-Time Chat - Kubernetes Local Setup

Este guia explica como subir a aplicação localmente usando Kubernetes, incluindo
Deployment, Service, Ingress Controller e configuração de host local.

> ⚠️ A aplicação ainda está em desenvolvimento.
Isso significa que falhas, comportamentos inesperados, interrupções e instabilidades podem ocorrer com frequência.
Novas funcionalidades, ajustes estruturais e mudanças de arquitetura também podem ser introduzidos a qualquer momento.
Use este projeto apenas para fins de estudo, testes locais e experimentação.

---

## 1. Pré-requisitos

Antes de tudo, instale:

- Docker
- kubectl
- Minikube (ou outro cluster local como Kind ou K3d)
- Ingress Controller NGINX (instalado no cluster)

Verifique o cluster:

    kubectl get nodes

---

## 2. Build da imagem local

A imagem da aplicação deve ser criada localmente e carregada no Minikube:

    eval $(minikube docker-env)
    docker build -t real-time-chat:v2 .

Confirme:

    docker images | grep real-time-chat

---

## 3. Subir o Postgres

A aplicação usa Postgres. Antes suba um Deployment e Service do banco.

Depois, confirme que o pod está funcionando:

    kubectl get pods
    kubectl logs -f <pod-do-postgres>

---

## 4. Subir a aplicação

Aplique os manifestos:

    kubectl apply -f deployment.yaml
    kubectl apply -f service.yaml

Verifique:

    kubectl get pods
    kubectl logs -f <pod-da-aplicacao>

---

## 5. Instalar o Ingress Controller NGINX

Execute:

    minikube addons enable ingress

Verifique se o controller subiu:

    kubectl get pods -n ingress-nginx

Procure por algo como:

    ingress-nginx-controller   Running

---

## 6. Criar o Ingress

Aplique o arquivo:

    kubectl apply -f ingress.yaml

Para verificar:

    kubectl get ingress

---

## 7. Configurar o /etc/hosts

Adicione esta linha no arquivo de hosts:

    127.0.0.1   chat.com

Depois limpe o cache DNS (Windows):

    ipconfig /flushdns

No Linux/macOS:

    sudo dscacheutil -flushcache
    sudo systemctl restart systemd-resolved

---

## 8. Acessar a aplicação

Agora basta acessar no navegador:

    http://chat.com/

Se quiser testar via curl:

    curl http://chat.com/

---

## 9. Solução de problemas comuns

### A imagem não sobe no pod

Provavelmente o Kubernetes está tentando puxar do Docker Hub.
Garanta que a flag do Minikube está ativa:

    eval $(minikube docker-env)

E recrie o pod:

    kubectl delete pod -l app=chat-api

---

### Ingress não funciona

- Verifique se o controller está rodando:

      kubectl get pods -n ingress-nginx

- Verifique se o Ingress pegou o endereço IP:

      kubectl describe ingress chat-ingress

---

Pronto! Sua aplicação agora roda localmente usando Kubernetes com ingress.
