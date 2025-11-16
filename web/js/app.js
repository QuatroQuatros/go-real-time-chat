const loginScreen = document.getElementById('loginScreen');
const loginUsernameEl = document.getElementById('loginUsername');
const usernameEl = document.getElementById('username');
const loginBtn = document.getElementById('loginBtn');
const chatApp = document.getElementById('chatApp');

let username = null;

// Gera um guest aleatório
function generateGuestName() {
    return 'Guest' + Math.floor(Math.random() * 10000);
}

loginBtn.addEventListener('click', () => {
    const name = loginUsernameEl.value.trim();
    username = name || generateGuestName();

    localStorage.setItem('chatUsername', username);

    // Esconde login e mostra chat
    loginScreen.style.display = 'none';
    chatApp.style.display = 'block';

    // Preenche input e status
    usernameEl.value = username;

    // Conecta automaticamente
    wsConnect();
});

(() => {
    const serverUrlEl = document.getElementById('serverUrl');
    const roomInputEl = document.getElementById('roomInput');
    const usernameEl = document.getElementById('username');
    const connectBtn = document.getElementById('connectBtn');
    const disconnectBtn = document.getElementById('disconnectBtn');
    const messagesEl = document.getElementById('messages');
    const sendBtn = document.getElementById('sendBtn');
    const messageInput = document.getElementById('messageInput');
    const statusEl = document.getElementById('status');
    const currentRoomEl = document.getElementById('currentRoom');
    const roomsListEl = document.getElementById('roomsList');
    const fetchHistoryBtn = document.getElementById('fetchHistoryBtn');

    let ws = null;
    let connected = false;
    let currentRoom = 'general';

    const roomsMap = {
        general: 1,
        random: 2,
        games: 3,
        support: 4
    };

    // default values
    serverUrlEl.value = 'ws://localhost:8080/ws';
    roomInputEl.value = currentRoom;
    usernameEl.value = username;

    function setStatus(s) {
        statusEl.textContent = s;
    }

    function addMessage({ user_id, content, timestamp }, me = false) {
        const wrap = document.createElement('div');
        wrap.className = 'message' + (me ? ' me' : '');
        const meta = document.createElement('div');
        meta.className = 'meta';
        const ts = timestamp ? new Date(timestamp).toLocaleString() : (new Date()).toLocaleTimeString();
        meta.textContent = `user:${user_id ?? 'unknown'} • ${ts}`;
        const body = document.createElement('div');
        body.textContent = content;
        wrap.appendChild(meta);
        wrap.appendChild(body);
        messagesEl.appendChild(wrap);
        messagesEl.scrollTop = messagesEl.scrollHeight;
    }

    function clearMessages() {
        messagesEl.innerHTML = '';
    }


    function fetchRoomHistory(roomName) {
        const roomId = roomsMap[roomName];
        fetch(`/rooms/${roomId}/messages?limit=50`)
            .then(r => {
                if (!r.ok) throw new Error('Falha ao carregar histórico: ' + r.status);
                return r.json();
            })
            .then(data => {
                clearMessages();
                (data || []).forEach(m => {
                    const user_id = m.User?.Username ?? m.UserID ?? m.user ?? 'unk';
                    const content = m.Content ?? m.content ?? JSON.stringify(m);
                    const timestamp = m.CreatedAt ?? m.created_at ?? new Date().toISOString();
                    addMessage({ user_id, content, timestamp }, false);
                });
            })
            .catch(err => console.error(err));
    }

    function wsConnect() {
        if (connected) return;
        const url = serverUrlEl.value.trim();
        if (!url) return alert('Informe a URL do WebSocket (ex: ws://localhost:8080/ws)');


        const roomId = roomsMap[currentRoom];
        const sep = url.includes('?') ? '&' : '?';
        const fullUrl = `${url}${sep}room=${roomId}`;

        try {
            ws = new WebSocket(fullUrl);
        } catch (err) {
            console.error(err);
            return alert('Falha ao criar WebSocket: ' + err.message);
        }

        ws.onopen = () => {
            connected = true;
            connectBtn.style.display = 'none';
            disconnectBtn.style.display = 'inline-block';
            setStatus('sim');
            console.log('WS aberto em', fullUrl);

            // Mensagem de join
            ws.send(JSON.stringify({
                RoomID: roomId,
                Content: `${usernameEl.value || username} entrou na sala`
            }));

            //fetchRoomHistory(currentRoom);
        };

        ws.onmessage = (ev) => {
            try {
                const data = JSON.parse(ev.data);
                console.log('msg recv', data);

                // filtrar por sala
                if (data.RoomID === roomsMap[currentRoom]) {
                    addMessage({
                        user_id: data.user_id ?? data.UserID ?? data.user ?? 'srv',
                        content: data.Content ?? data.content ?? '',
                        timestamp: data.CreatedAt ?? data.created_at ?? new Date().toISOString()
                    }, false);
                }
            } catch (e) {
                console.warn('Erro ao processar mensagem WS', e);
            }
        };

        ws.onclose = () => {
            connected = false;
            connectBtn.style.display = 'inline-block';
            disconnectBtn.style.display = 'none';
            setStatus('não');
            console.log('WS fechado');
        };

        ws.onerror = (e) => {
            console.error('WS error', e);
        };
    }

    function wsDisconnect() {
        if (!ws) return;
        const roomId = roomsMap[currentRoom];
        // Mensagem de saída
        ws.send(JSON.stringify({
            RoomID: roomId,
            Content: `${usernameEl.value || username} desconectou-se da sala`
        }));
        ws.close();
        ws = null;
    }

    function sendMessage() {
        if (!ws || ws.readyState !== WebSocket.OPEN) return alert('Conecte-se primeiro');
        const text = messageInput.value.trim();
        if (!text) return;
        const roomId = roomsMap[currentRoom];
        const payload = {
            RoomID: roomId,
            Content: text
        };
        ws.send(JSON.stringify(payload));
        addMessage({ user_id: usernameEl.value || username, content: text, timestamp: new Date().toISOString() }, true);
        messageInput.value = '';
    }

    connectBtn.addEventListener('click', wsConnect);
    disconnectBtn.addEventListener('click', wsDisconnect);
    sendBtn.addEventListener('click', sendMessage);
    messageInput.addEventListener('keydown', (e) => { if (e.key === 'Enter') sendMessage(); });

    // Troca de sala
    roomsListEl.addEventListener('click', (ev) => {
        const node = ev.target.closest('.room');
        if (!node) return;
        const room = node.dataset.room;
        // highlight
        document.querySelectorAll('.room').forEach(el => el.classList.remove('active'));
        node.classList.add('active');
        currentRoom = room;
        roomInputEl.value = room;
        currentRoomEl.textContent = room;
        clearMessages();
        // Se estiver conectado, desconecta e reconecta na nova sala
        if (connected) {
            wsDisconnect();
            setTimeout(wsConnect, 200);
        }
    });

    // Botão histórico
    fetchHistoryBtn.addEventListener('click', () => fetchRoomHistory(currentRoom));
    roomInputEl.addEventListener('keydown', e => { if (e.key === 'Enter') wsConnect(); });

})();