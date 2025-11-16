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
    let username = 'TestUser';

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

    function wsConnect() {
        if (connected) return;
        const url = serverUrlEl.value.trim();
        if (!url) return alert('Informe a URL do WebSocket (ex: ws://localhost:8080/ws)');

        // attach room param to url if not present
        const r = roomInputEl.value.trim() || currentRoom;
        currentRoom = r;
        currentRoomEl.textContent = currentRoom;

        // optional: send token via ?token=... if your server requires auth
        const sep = url.includes('?') ? '&' : '?';
        const full = `${url}${sep}room=${encodeURIComponent(currentRoom)}`;
        try {
            ws = new WebSocket(full);
        } catch (err) {
            console.error(err);
            return alert('Falha ao criar WebSocket: ' + err.message);
        }

        ws.onopen = () => {
            connected = true;
            connectBtn.style.display = 'none';
            disconnectBtn.style.display = 'inline-block';
            setStatus('sim');
            console.log('WS aberto em', full);
            // optionally announce join
            // const joinMsg = { type: 'message', room: currentRoom, content: `${usernameEl.value || username} entrou na sala` };
            const joinMsg = { RoomID: 1, Content: `${usernameEl.value || username} entrou na sala` };
            ws.send(JSON.stringify(joinMsg));
        };

        ws.onmessage = (ev) => {
            try {
                const data = JSON.parse(ev.data);
                console.log('msg recv', data);
                // server uses OutgoingMessage { type, room, user_id, content, timestamp }
                if (data.type === 'message' && data.room === currentRoom) {
                    addMessage({ user_id: data.user_id ?? data.user, content: data.content, timestamp: data.timestamp }, false);
                } else {
                    // generic dump
                    addMessage({ user_id: data.user_id ?? 'srv', content: JSON.stringify(data).slice(0, 200), timestamp: new Date().toISOString() });
                }
            } catch (e) {
                console.warn('msg parse fail', e);
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
        try {
            ws.close();
        } catch (e) { /* ignore */ }
        ws = null;
    }

    function sendMessage() {
        if (!ws || ws.readyState !== WebSocket.OPEN) {
            return alert('Conecte-se primeiro');
        }
        const text = messageInput.value.trim();
        if (!text) return;
        // const payload = {
        //     type: 'message',
        //     room: currentRoom,
        //     content: text
        // };
        const payload = {
            RoomID: 1,
            Content: text
        };
        ws.send(JSON.stringify(payload));
        // optimistically render as "me" (server may echo with user_id)
        addMessage({ user_id: usernameEl.value || username, content: text, timestamp: new Date().toISOString() }, true);
        messageInput.value = '';
    }

    connectBtn.addEventListener('click', wsConnect);
    disconnectBtn.addEventListener('click', wsDisconnect);
    sendBtn.addEventListener('click', sendMessage);
    messageInput.addEventListener('keydown', (e) => { if (e.key === 'Enter') sendMessage(); });

    // room switching via sidebar
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
        // reconnect automatically if already connected
        if (ws && ws.readyState === WebSocket.OPEN) {
            // simple approach: close and reconnect to change ?room param
            ws.close();
            setTimeout(wsConnect, 200);
        }
    });

    // fetch history
    fetchHistoryBtn.addEventListener('click', () => {
        const base = window.location.origin.replace(/^http/, 'http');
        const apiBase = base; // ex: http://localhost:8080
        const room = roomInputEl.value.trim() || currentRoom;
        const limit = 50;
        const url = `${apiBase}/rooms/${encodeURIComponent(room)}/messages?limit=${limit}`;
        fetch(url)
            .then(r => {
                if (!r.ok) throw new Error('failed: ' + r.status);
                return r.json();
            })
            .then(data => {
                clearMessages();
                // expects array of messages [{User, UserID, Content, CreatedAt,...}]
                (data || []).forEach(m => {
                    // adapt to your response shape: try several possibilities
                    const user_id = m.User?.Username ?? m.UserID ?? m.user_id ?? m.user ?? m.username ?? 'unk';
                    const content = m.Content ?? m.content ?? m.content_text ?? JSON.stringify(m);
                    const timestamp = m.CreatedAt ?? m.created_at ?? m.timestamp ?? new Date().toISOString();
                    addMessage({ user_id, content, timestamp }, false);
                });
            })
            .catch(err => {
                console.error(err);
                alert('Falha ao carregar histórico: ' + err.message);
            });
    });

    // quick connect on enter in room input
    roomInputEl.addEventListener('keydown', (e) => { if (e.key === 'Enter') wsConnect(); });

})();