import { EmojiButton } from '@joeattardi/emoji-button';

export default () => ({
    socket: null,
    username: null,
    openUsernameModal: false,
    openCardModal: false,
    openColumnModal: false,
    openTimerModal: false,
    clients: [],
    clientConnections: {},
    columns: [],
    cards: [],
    numClients: 0,
    tempColumn: {
        id: '',
        name: '',
    },
    tempCard: {
        id: '',
        name: '',
        column_id: '',
    },
    timer: {
        duration: '5m',
        show: false,
        running: false,
        done: false,
        display: "00:00",
    },
    emojiPicker: null,
    initBoard() {
        this.username = localStorage.getItem('username');
        this.joinBoard();
        this.emojiPicker = new EmojiButton({
            autoHide: false,
            showVariants: false,
            position: "bottom-start",
            zIndex: 99,
        });
        this.emojiPicker.on('emoji', selection => {
            this.tempCard.name += selection.emoji;
        });
    },
    wsConnect(username) {
        const host = window.location.host;
        const protocol = window.location.protocol;
        const pathname = window.location.pathname;

        const wsProtocol = protocol === 'https:' ? 'wss:' : 'ws:';
        this.socket = new WebSocket(`${wsProtocol}//${host}${pathname}/ws?u=${username}`);
        var app = this;
        this.socket.addEventListener("open", (_event) => {
            console.log(`Connected to WebSocket as ${username}...`);
        });

        this.socket.addEventListener("message", (event) => {
            app.onWsEvent(event);
        });
    },
    onWsEvent(event) {
        const e = JSON.parse(event.data);
        switch (e.type) {
            case 'board.users':
                const sortedClients = e.data.sort((a, b) => a.joined_at - b.joined_at);
                const uniqueClients = [...new Set(sortedClients.map(c => c.user.id))];
                this.numClients = uniqueClients.length;
                this.clientConnections = sortedClients.reduce((acc, c) => {
                    if (!acc[c.user.id]) acc[c.user.id] = 0;
                    acc[c.user.id]++;
                    return acc;
                }
                    , {});
                this.clients = uniqueClients.map(id => sortedClients.find(c => c.user.id === id));
                break;

            case 'board.status':
                this.columns = e.data.columns.sort((a, b) => a.order - b.order);
                this.cards = (e.data.cards || []).sort((a, b) => a.created_at - b.created_at)
                break;

            case 'board.notification':
                this.dispatchCustomEvents('flash', e.data);
                break;

            case 'timer.state':
                this.timer.show = ['running', 'paused', 'done'].indexOf(e.data.status) !== -1;
                this.timer.done = e.data.status === 'done';
                this.timer.running = e.data.status === 'running';
                this.timer.display = e.data.display;
                if (this.timer.done) {
                    this.playSound();
                    setTimeout(() => this.stopTimer(false), 5000);
                }
                break;

            default:
                break;
        }
    },
    askUsername() {
        this.openUsernameModal = true;
        setTimeout(() => this.$refs.username.focus(), 200);
    },
    joinBoard() {
        if (this.username == null || this.username == '') {
            this.askUsername();
            return;
        }

        localStorage.setItem('username', this.username);
        this.wsConnect(this.username);
        this.closeModal('username');
    },
    columnNameById(id) {
        const column = this.columns.find(c => c.id === id);
        return column ? column.name : '';
    },
    editColumn(column) {
        if (column == null) {
            if (this.columns.length >= 6) {
                this.dispatchCustomEvents('flash', 'You can only have a maximum of 6 columns');
                return;
            }
            this.tempColumn.name = '';
            this.tempColumn.id = '';
        } else {
            this.tempColumn.name = column.name;
            this.tempColumn.id = column.id;
        }
        this.openColumnModal = true;
        setTimeout(() => this.$refs.columnName.focus(), 200);
    },
    saveColumn() {
        if (this.tempColumn.name == '') return;
        if (this.tempColumn.id) {
            this.socket.send(JSON.stringify({ type: 'column.update', data: this.tempColumn }));
        } else {
            this.socket.send(JSON.stringify({ type: 'column.new', data: { name: this.tempColumn.name } }));
        }
        this.closeModal('column');
    },
    deleteColumn(column) {
        this.socket.send(JSON.stringify({ type: 'column.delete', data: { id: column.id } }));
        this.closeModal('column');
    },
    editCard(column, card) {
        if (card == null) {
            this.tempCard.id = '';
            this.tempCard.name = '';
            this.tempCard.column_id = column.id;
        } else {
            this.tempCard.id = card.id;
            this.tempCard.name = card.name;
            this.tempCard.column_id = column.id;
        }
        this.openCardModal = true;
        setTimeout(() => this.$refs.cardName.focus(), 200);
    },
    saveCard() {
        if (this.tempCard.name == '') return;

        if (this.tempCard.id) {
            this.socket.send(JSON.stringify({ type: 'card.update', data: this.tempCard }));
        } else {
            this.socket.send(JSON.stringify({ type: 'card.new', data: { name: this.tempCard.name, column_id: this.tempCard.column_id } }));
        }
        this.closeModal('card');
    },
    voteCard(card, vote) {
        if (vote !== 1 && vote !== -1) return;
        this.socket.send(JSON.stringify({ type: 'card.vote', data: { id: card.id, vote: vote } }));
    },
    deleteCard(card) {
        this.socket.send(JSON.stringify({ type: 'card.delete', data: { id: card.id } }));
        this.closeModal('card');
    },
    closeModal(name) {
        if (name === 'username') {
            this.openUsernameModal = false;
        }
        if (name === 'column') {
            this.tempColumn.name = '';
            this.tempColumn.id = '';
            this.openColumnModal = false;
        }
        if (name === 'card') {
            this.tempCard.id = '';
            this.tempCard.name = '';
            this.tempCard.column_id = '';
            this.openCardModal = false;
        }
        if (name === 'timer') {
            this.openTimerModal = false;
        }
    },
    onDragStart(event, card) {
        event.dataTransfer.setData('cardId', card.id);
        event.dataTransfer.setData('cardColumnId', card.column_id);
        event.target.classList.add('opacity-10');
    },
    onDragEnd(event) {
        event.target.classList.remove('opacity-10');
    },
    onDragOver(event) {
        event.preventDefault();
        return false;
    },
    onDragEnter(event) {
        if (!event.target.classList.contains('is-dropzone')) return;
        event.target.classList.add('bg-blue-200');
    },
    onDragLeave(event) {
        if (!event.target.classList.contains('is-dropzone')) return;
        event.target.classList.remove('bg-blue-200');
    },
    onDrop(event, newColumn) {
        event.stopPropagation(); // Stops some browsers from redirecting.
        event.preventDefault();
        event.target.classList.remove('bg-blue-200');

        const cardId = event.dataTransfer.getData('cardId');
        const cardColumnId = event.dataTransfer.getData('cardColumnId');

        if (cardColumnId === newColumn.id) {
            event.dataTransfer.clearData();
            return;
        }

        // element moved automatically by changing in data, so no need to remove it manually
        // const draggableElement = document.getElementById(cardId);
        // const dropzone = event.target;
        // dropzone.removeChild(draggableElement);

        // Update
        let cardIndex = this.cards.findIndex(c => c.id === cardId);
        this.cards[cardIndex].column_id = newColumn.id;
        this.socket.send(JSON.stringify({ type: 'card.update', data: this.cards[cardIndex] }));
        event.dataTransfer.clearData();
    },
    startTimer() {
        if (this.timer.running) return;
        this.closeModal('timer');
        this.socket.send(JSON.stringify({ type: 'timer.cmd', data: { cmd: 'start', value: this.timer.duration } }));
    },
    pauseTimer() {
        if (!this.timer.running) return;
        this.socket.send(JSON.stringify({ type: 'timer.cmd', data: { cmd: 'pause' } }));
    },
    resumeTimer() {
        if (this.timer.running) return;
        this.socket.send(JSON.stringify({ type: 'timer.cmd', data: { cmd: 'start' } }));
    },
    stopTimer(isCommand) {
        if (isCommand) this.socket.send(JSON.stringify({ type: 'timer.cmd', data: { cmd: 'stop' } }));

        this.timer.show = false;
        this.timer.running = false;
        this.timer.done = false;
        this.timer.display = "00:00";
    },
    playSound() {
        const audio = new Audio('/static/notif.wav');
        audio.play();
    },
    numClientConnections(userId) {
        return this.clientConnections[userId] || 0;
    },
    dispatchCustomEvents(eventName, message) {
        let customEvent = new CustomEvent(eventName, { detail: { message: message } });
        window.dispatchEvent(customEvent);
    },
    gridColsClass() {
        return {
            1: 'grid-cols-1',
            2: 'grid-cols-2',
            3: 'grid-cols-3',
            4: 'grid-cols-4',
            5: 'grid-cols-5',
            6: 'grid-cols-6',
        }[this.columns.length] || 'grid-cols-4';
    },
    openEmojiPicker() {
        this.emojiPicker.showPicker(this.$event.target);
    }
});
