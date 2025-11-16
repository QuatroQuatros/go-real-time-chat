# Real-Time Chat - Kubernetes Local Setup (Kind/Docker Desktop)

Este guia explica como subir a aplicação localmente usando Kubernetes **Kind** (Kubernetes in Docker), incluindo Deployment, Service, Ingress Controller e configuração de host local.

> ⚠️ A aplicação ainda está em desenvolvimento.
Isso significa que falhas, comportamentos inesperados, interrupções e instabilidades podem ocorrer com frequência.
Novas funcionalidades, ajustes estruturais e mudanças de arquitetura também podem ser introduzidos a qualquer momento.
Use este projeto apenas para fins de estudo, testes locais e experimentação.

---

## 1. Pré-requisitos

Antes de tudo, instale:

- **Docker Desktop** (com Kubernetes ativado)
- **kubectl**
- **Kind** (para criar o cluster)

Verifique o cluster:

    kubectl get nodes

---

## 2. Build e Carregamento da Imagem

A imagem da aplicação deve ser criada e carregada no Kind.

1. **Build da imagem local:**

    ```bash
    docker build -t real-time-chat:v2 .
    ```

2. **Carregue a imagem para o cluster Kind:**

    ```bash
    kind load docker-image real-time-chat:v2
    ```

    Confirme:

    ```bash
    docker images | grep real-time-chat
    ```

---

## 3. Subir o Postgres

A aplicação usa Postgres. Antes suba um Deployment e Service do banco.

Depois, confirme que o pod está funcionando:

    kubectl get pods
    kubectl logs -f ./kubernets/postgres/deploymentsey.yaml

---

## 4. Subir a Aplicação

Aplique os manifestos:

    kubectl apply -f ./kubernets/api/deploymentset.yaml
    kubectl apply -f ./kubernets/api/chat-service.yaml

Verifique:

    kubectl get pods
    kubectl logs -f <pod-da-aplicacao>

---

## 5. Instalar e Configurar o Ingress Controller NGINX (Kind)

O Kind requer a instalação manual do Ingress Controller.

1. **Instale o Ingress Controller:**

    ```bash
    kubectl apply -f [https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml](https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml)
    ```

2. **Verifique se o controller subiu:**

    ```bash
    kubectl get pods -n ingress-nginx
    ```

    Procure por algo como: `ingress-nginx-controller   Running`

---

## 6. Criar o Ingress (com Regras de WebSocket)

Aplique o arquivo:

    kubectl apply -f ingress.yaml

> **Ajuste Importante:** Sua aplicação Go serve o frontend na raiz e já está respondendo corretamente. O Ingress abaixo roteará o tráfego da raiz e do WebSocket.

Para verificar:

    kubectl get ingress

---

## 7. Configurar o Acesso Local (Port Forwarding)

Como você usa o Kind, o `LoadBalancer` não é acessível diretamente pelo IP externo. Usaremos o **Port Forwarding** para fazer a ponte para a porta 80 do seu host.

1. **Encontre o nome do Pod do Ingress Controller:**

    ```bash
    kubectl get pods -n ingress-nginx -l app.kubernetes.io/component=controller
    ```

2. **Abra um terminal NOVO e execute o Port Forwarding (mantenha-o rodando):**
    > **Atenção:** Se não funcionar, tente executar o terminal como Administrador.

    ```bash
    kubectl port-forward -n ingress-nginx <nome-do-pod-do-ingress> 80:80
    ```

3. **Configure o /etc/hosts (Windows/Linux/macOS):**
    Adicione esta linha no seu arquivo de hosts para que `chat.com` aponte para o *localhost*:

    ```
    127.0.0.1   chat.com
    ```

4. **Limpe o cache DNS (Windows):**

    ```bash
    ipconfig /flushdns
    ```

---

## 8. Acessar a Aplicação

Agora basta acessar no navegador ou via `curl`. O Ingress Controller na porta 80 (via `port-forward`) irá rotear para sua aplicação.

- **Acesso ao Frontend (Rota Raiz):**

    ```
    [http://chat.com/](http://chat.com/)
    ```

- **Teste de Health Check:**

    ```bash
    curl [http://chat.com/health](http://chat.com/health)
    ```

---

## 9. Solução de problemas comuns

### A imagem não sobe no pod

Provavelmente o Kind não encontrou a imagem localmente. Garanta que a imagem foi carregada:

    kind load docker-image real-time-chat:v2

E recrie o pod:

    kubectl delete pod -l app=chat-api

---

### Ingress não funciona

- **Verifique o `Port Forwarding`:** Certifique-se de que a janela do `kubectl port-forward 80:80` está aberta e ativa.
- **Verifique o controller:**

      kubectl get pods -n ingress-nginx

- **Verifique o Ingress:**

      kubectl describe ingress chat-ingress

- **Problemas de WebSockets:** O Ingress Controller precisa de configurações específicas para conexões longas de WebSockets (que já estão nas *Annotations* do seu manifesto).

---

Pronto! Sua aplicação agora roda localmente usando Kubernetes com ingress.
