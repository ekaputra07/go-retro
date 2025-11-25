import { EmojiButton } from '@joeattardi/emoji-button';
import { App, Client, ClientConnection, Column, Card } from './types/app';

export default (): App => ({
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
    standup: {
        show: false,
        title: "Stand-up",
        shuffled_user_ids: [],
    },
    initBoard(): void {
        this.username = localStorage.getItem('username') || '';
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
    wsConnect(username: string): void {
        const host = window.location.host;
        const protocol = window.location.protocol;
        const pathname = window.location.pathname;

        const wsProtocol = protocol === 'https:' ? 'wss:' : 'ws:';
        this.socket = new WebSocket(`${wsProtocol}//${host}${pathname}/ws?u=${username}`);
        var app = this;
        this.socket.addEventListener("open", (_event) => {
            console.log(`Connected to WebSocket as ${username}...`);
            app.socket?.send(JSON.stringify({ type: 'me' }));
        });

        this.socket.addEventListener("message", (event) => {
            app.onWsEvent(event);
        });

        this.socket.addEventListener("close", (event) => {
            app.socket = undefined;
            app.dispatchCustomEvents("flash", "You're disconnected. Reconnecting...");
            console.log(`Disconnected from WebSocket: ${event.code} ${event.reason}`);
            setTimeout(() => {
                app.wsConnect(username);
            }, 5000);
        });
    },
    onWsEvent(event: MessageEvent): void {
        const e: any = JSON.parse(event.data);

        switch (e.type) {
            case 'me':
                this.currentUser = e.user;
                break

            case 'columns':
                this.applyOperation(this.columns, e);
                break;

            case 'cards':
                this.applyOperation(this.cards, e);
                break;

            case 'board.users':
                const sortedClients: Client[] = e.data.sort((a: Client, b: Client) => a.joined_at - b.joined_at);
                const uniqueClients: string[] = [...new Set(sortedClients.map((c: Client) => c.user.id))];

                this.numClients = uniqueClients.length;
                this.clientConnections = sortedClients.reduce((acc: ClientConnection, c: Client) => {
                    if (!acc[c.user.id]) acc[c.user.id] = 0;
                    acc[c.user.id]++;
                    return acc;
                }, {});
                this.clients = uniqueClients.map((id: string) => sortedClients.find((c: Client) => c.user.id === id) as Client);

                if (this.standup.show) {
                    this.refreshShuffledUserIds();
                }
                break;

            case 'board.notification':
                // don't show notification to those who triggers it
                if (this.currentUser && e.user) {
                    if (this.currentUser.id != e.user.id) {
                        this.dispatchCustomEvents('flash', e.data);
                    }
                }
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
    applyOperation(list: Column[] | Card[], change: any): void {
        var idx = list.findIndex(c => c.id == change.id)
        if (idx >= 0) {
            if (change.op == "put") {
                list[idx] = change.obj;
            } else if (change.op == "del") {
                list.splice(idx, 1);
            }
            list = list.sort((a, b) => a.created_at! - b.created_at!);
        } else if (idx == -1 && change.op == "put") {
            list.push(change.obj);
            list = list.sort((a, b) => a.created_at! - b.created_at!);
        }
    },
    askUsername(): void {
        this.openUsernameModal = true;
    },
    joinBoard(): void {
        if (this.username == '' || this.username == '') {
            this.askUsername();
        } else {
            localStorage.setItem('username', this.username!);
            this.wsConnect(this.username!);
            this.closeModal('username');
        }
    },
    columnNameById(id: string): string {
        const column = this.columns.find(c => c.id === id);
        return column ? column.name : '';
    },
    editColumn(column: Column): void {
        if (column == undefined) {
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
    saveColumn(): void {
        if (this.tempColumn.name == '') return;
        if (this.tempColumn.id) {
            this.socket?.send(JSON.stringify({ type: 'column.update', data: this.tempColumn }));
        } else {
            this.socket?.send(JSON.stringify({ type: 'column.new', data: { name: this.tempColumn.name } }));
        }
        this.closeModal('column');
    },
    deleteColumn(column: Column): void {
        this.socket?.send(JSON.stringify({ type: 'column.delete', data: { id: column.id } }));
        this.closeModal('column');
    },
    editCard(column: Column, card: Card): void {
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
    saveCard(): void {
        if (this.tempCard.name == '') return;

        if (this.tempCard.id) {
            this.socket?.send(JSON.stringify({ type: 'card.update', data: this.tempCard }));
        } else {
            this.socket?.send(JSON.stringify({ type: 'card.new', data: { name: this.tempCard.name, column_id: this.tempCard.column_id } }));
        }
        this.closeModal('card');
    },
    voteCard(card: Card, vote: number): void {
        if (vote !== 1 && vote !== -1) return;
        this.socket?.send(JSON.stringify({ type: 'card.vote', data: { id: card.id, vote: vote } }));
    },
    deleteCard(card: Card): void {
        this.socket?.send(JSON.stringify({ type: 'card.delete', data: { id: card.id } }));
        this.closeModal('card');
    },
    closeModal(name: string): void {
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
    onDragStart(event: DragEvent, card: Card): void {
        event.dataTransfer?.setData('cardId', card.id);
        event.dataTransfer?.setData('cardColumnId', card.column_id);
        const dragTarget = event.target as HTMLElement | null;
        dragTarget?.classList.add('opacity-10');
    },
    onDragEnd(event: DragEvent): void {
        const dragTarget = event.target as HTMLElement | null;
        dragTarget?.classList.remove('opacity-10');
    },
    onDragOver(event: DragEvent): void {
        event.preventDefault();
    },
    onDragEnter(event: DragEvent): void {
        const dragTarget = event.target as HTMLElement | null;
        if (!dragTarget?.classList.contains('is-dropzone')) return;
        dragTarget?.classList.add('bg-blue-200');
    },
    onDragLeave(event: DragEvent): void {
        const dragTarget = event.target as HTMLElement | null;
        if (!dragTarget?.classList.contains('is-dropzone')) return;
        dragTarget?.classList.remove('bg-blue-200');
    },
    onDrop(event: DragEvent, newColumn: Column): void {
        event.stopPropagation(); // Stops some browsers from redirecting.
        event.preventDefault();
        const dragTarget = event.target as HTMLElement | null;
        dragTarget?.classList.remove('bg-blue-200');

        const cardId = event.dataTransfer?.getData('cardId');
        const cardColumnId = event.dataTransfer?.getData('cardColumnId');

        if (cardColumnId === newColumn.id) {
            event.dataTransfer?.clearData();
            return;
        }

        // element moved automatically by changing in data, so no need to remove it manually
        // const draggableElement = document.getElementById(cardId);
        // const dropzone = event.target;
        // dropzone.removeChild(draggableElement);

        // Update
        let cardIndex = this.cards.findIndex(c => c.id === cardId);
        this.cards[cardIndex].column_id = newColumn.id;
        this.socket?.send(JSON.stringify({ type: 'card.update', data: this.cards[cardIndex] }));
        event.dataTransfer?.clearData();
    },
    startTimer(): void {
        if (this.timer.running) return;
        this.closeModal('timer');
        this.socket?.send(JSON.stringify({ type: 'timer.cmd', data: { cmd: 'start', value: this.timer.duration } }));
    },
    pauseTimer(): void {
        if (!this.timer.running) return;
        this.socket?.send(JSON.stringify({ type: 'timer.cmd', data: { cmd: 'pause' } }));
    },
    resumeTimer(): void {
        if (this.timer.running) return;
        this.socket?.send(JSON.stringify({ type: 'timer.cmd', data: { cmd: 'start' } }));
    },
    stopTimer(isCommand: boolean): void {
        if (isCommand) this.socket?.send(JSON.stringify({ type: 'timer.cmd', data: { cmd: 'stop' } }));

        this.timer.show = false;
        this.timer.running = false;
        this.timer.done = false;
        this.timer.display = "00:00";
    },
    playSound(): void {
        const audio = new Audio('/static/notif.wav');
        audio.play();
    },
    numClientConnections(userId: string): number {
        return this.clientConnections[userId] || 0;
    },
    dispatchCustomEvents(eventName: string, message: string): void {
        let customEvent = new CustomEvent(eventName, { detail: { message: message } });
        window.dispatchEvent(customEvent);
    },
    gridColsClass(): string {
        return {
            1: 'grid-cols-1',
            2: 'grid-cols-2',
            3: 'grid-cols-3',
            4: 'grid-cols-4',
            5: 'grid-cols-5',
            6: 'grid-cols-6',
        }[this.columns.length] || 'grid-cols-4';
    },
    openEmojiPicker(): void {
        this.emojiPicker?.showPicker(this.$event.target);
    },
    getClientById(id: string): Client | undefined {
        return this.clients.find(c => c.user.id === id);
    },
    startStandup(): void {
        if (this.clients.length < 2) {
            this.dispatchCustomEvents('flash', 'Not enough people to start a stand-up!');
            return;
        }
        // Fisher-Yates shuffle algorithm
        const userIds = this.clients.map(c => c.user.id);
        for (let i = userIds.length - 1; i > 0; i--) {
            const j = Math.floor(Math.random() * (i + 1));
            [userIds[i], userIds[j]] = [userIds[j], userIds[i]];
        }
        this.standup.shuffled_user_ids = userIds;
        this.setCurrentStandupUser(this.getClientById(userIds[0])!);
        this.standup.show = true;
    },
    closeStandup(): void {
        this.standup.current_user_id = undefined;
        this.standup.shuffled_user_ids = [];
        this.standup.show = false;
    },
    refreshShuffledUserIds(): void {
        let newShuffledUserIds = [];
        let newUserIds = this.clients.map(c => c.user.id).filter(id => !this.standup.shuffled_user_ids.includes(id));
        let missingUserIds = this.standup.shuffled_user_ids.filter(id => !this.clients.map(c => c.user.id).includes(id));

        if (newUserIds.length > 0) {
            newShuffledUserIds = [...this.standup.shuffled_user_ids, ...newUserIds];
        } else {
            newShuffledUserIds = this.standup.shuffled_user_ids;
        }

        if (missingUserIds.length > 0) {
            newShuffledUserIds = newShuffledUserIds.filter(id => !missingUserIds.includes(id));
        }
        this.standup.shuffled_user_ids = newShuffledUserIds;
    },
    setCurrentStandupUser(client: Client): void {
        this.dispatchCustomEvents('flash', `${client.user.name}'s turn`);
        this.standup.current_user_id = client.user.id;
    }
});
