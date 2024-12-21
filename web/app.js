function app() {
    return {
        socket: null,
        openModal: false,
        boardId: '',
        columns: [],
        card: {
            name: '',
            column: '',
        },
        editCard: {},
        cards: [],
        initBoard() {
            const host = window.location.host;
            const pathname = window.location.pathname;
            const boardId = pathname.replace('/b/', '');
            this.boardId = boardId;
            localStorage.setItem('goretro:current-board', boardId);

            this.socket = new WebSocket(`ws://${host}${pathname}/ws`);
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
                    this.columns = e.data.columns;
                    this.cards = (e.data.cards || []).map(c => {
                        c.edit = false; return c;
                    }).sort((a, b) => a.created_at - b.created_at)
                    break;
            }
        },
        showModal(column) {
            this.card.column = column;
            this.openModal = true;
            setTimeout(() => this.$refs.cardName.focus(), 200);
        },
        closModal() {
            // Reset the form
            this.card.name = '';
            this.card.column = '';
            this.card.board = '';

            // close the modal
            this.openModal = false;
        },
        saveCard(card) {
            if (card.name == '') return;
            this.socket.send(JSON.stringify({type: 'card.update', data: card}));
        },
        addCard() {
            if (this.card.name == '') return;

            // data to save
            const newCard = {
                name: this.card.name,
                column: this.card.column.id,
            };

            this.socket.send(JSON.stringify({type: 'card.new', data: newCard}));
            this.closModal();
        },
        deleteCard(card) {
            this.socket.send(JSON.stringify({type: 'card.delete', data: card}));
        },
        onDragStart(event, id) {
            event.dataTransfer.setData('text/plain', id);
            event.target.classList.add('opacity-5');
        },
        onDragOver(event) {
            event.preventDefault();
            return false;
        },
        onDragEnter(event) {
            event.target.classList.add('bg-gray-200');
        },
        onDragLeave(event) {
            event.target.classList.remove('bg-gray-200');
        },
        onDrop(event, column) {
            event.stopPropagation(); // Stops some browsers from redirecting.
            event.preventDefault();
            event.target.classList.remove('bg-gray-200');

            // console.log('Dropped', this);
            const id = event.dataTransfer.getData('text');

            const draggableElement = document.getElementById(id);
            const dropzone = event.target;

            dropzone.appendChild(draggableElement);
            
            // Update
            let cardIndex = this.cards.findIndex(t => t.id === id);

            // Add new data to localStorage Array
            this.cards[cardIndex].column = column;
            this.socket.send(JSON.stringify({type: 'card.update', data: this.cards[cardIndex]}));
            event.dataTransfer.clearData();
        },
    }
}