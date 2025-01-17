function app() {
    return {
        socket: null,
        openCardModal: false,
        openColumnModal: false,
        openTimerModal: false,
        boardId: '',
        columns: [],
        cards: [],
        numPeople: 0,
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
        initBoard() {
            const host = window.location.host;
            const protocol = window.location.protocol;
            const pathname = window.location.pathname;
            const boardId = pathname.replace('/b/', '');
            this.boardId = boardId;

            const wsProtocol = protocol === 'https:' ? 'wss:' : 'ws:';
            this.socket = new WebSocket(`${wsProtocol}//${host}${pathname}/ws`);
            var app = this;
            this.socket.addEventListener("open", (event) => {
                app.onWsOpen(event);
            });

            this.socket.addEventListener("message", (event) => {
                app.onWsEvent(event);
            });
        },
        onWsOpen(event) {
            console.log('Connected to WebSocket...');
        },
        onWsEvent(event) {
            const e = JSON.parse(event.data);
            switch (e.type) {
                case 'board.status':
                    this.numPeople = e.data.user_count;
                    this.columns = e.data.columns.sort((a, b) => a.order - b.order);
                    this.cards = (e.data.cards || []).sort((a, b) => a.created_at - b.created_at)
                    break;
                case 'timer.state':
                    this.timer.show = ['running', 'paused', 'done'].indexOf(e.data.status) !== -1;
                    this.timer.done = e.data.status === 'done';
                    this.timer.running = e.data.status === 'running';
                    this.timer.display = e.data.display;
                    if(this.timer.done) {
                        this.playSound();
                        setTimeout(() => this.stopTimer(), 10000);
                    }
            }
        },
        columnNameById(id){
            const column = this.columns.find(c => c.id === id);
            return column ? column.name : '';
        },
        editColumn(column) {
            if(column == null) {
                this.tempColumn.name = '';
                this.tempColumn.id = '';
            } else {
                this.tempColumn.name = column.name;
                this.tempColumn.id = column.id;
            }
            this.openColumnModal = true;
            setTimeout(() => this.$refs.columnName.focus(), 200);
        },
        saveColumn(){
            if (this.tempColumn.name == '') return;
            if(this.tempColumn.id) {
                this.socket.send(JSON.stringify({type: 'column.update', data: this.tempColumn}));
            } else {
                this.socket.send(JSON.stringify({type: 'column.new', data: {name: this.tempColumn.name}}));
            }
            this.closeModal('column');
        },
        deleteColumn(column) {
            this.socket.send(JSON.stringify({type: 'column.delete', data: {id: column.id}}));
            this.closeModal('column');
        },
        editCard(column, card) {
            if(card == null) {
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

            if(this.tempCard.id) {
                this.socket.send(JSON.stringify({type: 'card.update', data: this.tempCard}));
            } else {
                this.socket.send(JSON.stringify({type: 'card.new', data: {name: this.tempCard.name, column_id: this.tempCard.column_id}}));
            }
            this.closeModal('card');
        },
        deleteCard(card) {
            this.socket.send(JSON.stringify({type: 'card.delete', data: {id: card.id}}));
            this.closeModal('card');
        },
        closeModal(name) {
            if(name === 'column') {
                this.tempColumn.name = '';
                this.tempColumn.id = '';
                this.openColumnModal = false;
            }
            if(name === 'card') {
                this.tempCard.id = '';
                this.tempCard.name = '';
                this.tempCard.column_id = '';
                this.openCardModal = false;
            }
            if(name === 'timer') {
                this.openTimerModal = false;
            }
        },
        onDragStart(event, card) {
            event.dataTransfer.setData('cardId', card.id);
            event.dataTransfer.setData('cardColumnId', card.column_id);
            event.target.classList.add('opacity-5');
        },
        onDragOver(event) {
            event.preventDefault();
            return false;
        },
        onDragEnter(event) {
            if(!event.target.classList.contains('is-dropzone')) return;
            event.target.classList.add('bg-blue-200');
        },
        onDragLeave(event) {
            if(!event.target.classList.contains('is-dropzone')) return;
            event.target.classList.remove('bg-blue-200');
        },
        onDrop(event, newColumn) {
            event.stopPropagation(); // Stops some browsers from redirecting.
            event.preventDefault();
            event.target.classList.remove('bg-blue-200');

            const cardId = event.dataTransfer.getData('cardId');
            const cardColumnId = event.dataTransfer.getData('cardColumnId');

            if(cardColumnId === newColumn.id) {
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
            this.socket.send(JSON.stringify({type: 'card.update', data: this.cards[cardIndex]}));
            event.dataTransfer.clearData();
        },
        startTimer() {
            if(this.timer.running) return;
            this.closeModal('timer');
            this.socket.send(JSON.stringify({type: 'timer.cmd', data: {cmd: 'start', value: this.timer.duration}}));
        },
        pauseTimer() {
            if(!this.timer.running) return;
            this.socket.send(JSON.stringify({type: 'timer.cmd', data: {cmd: 'pause'}}));
        },
        resumeTimer() {
            if(this.timer.running) return;
            this.socket.send(JSON.stringify({type: 'timer.cmd', data: {cmd: 'start'}}));
        },
        stopTimer() {
            this.socket.send(JSON.stringify({type: 'timer.cmd', data: {cmd: 'stop'}}));
            this.timer.show = false;
            this.timer.running = false;
            this.timer.done = false;
            this.timer.display = "00:00";
        },
        playSound() {
            const audio = new Audio('/static/notif.wav');
            audio.play();
        }
    }
}