import { EmojiButton } from '@joeattardi/emoji-button';

export interface User {
    id: string;
    name: string;
    avatar_id: number;
}

export interface Client {
    id: string;
    user: User;
    joined_at: number;
}

export interface ClientConnection {
    [key: string]: number;
}

export interface Column {
    id: string;
    name: string;
    created_at?: number;
}

export interface Card {
    id: string;
    name: string;
    column_id: string;
    created_at?: number;
    votes?: number;
}

export interface Timer {
    duration: string;
    show: boolean;
    running: boolean;
    done: boolean;
    display: string;
}

export interface Standup {
    show: boolean;
    title: string;
    current_user_id?: string;
    shuffled_user_ids: string[];
}

export interface EmojiPicker {
    picker?: EmojiButton;
}

export interface App {
    socket?: WebSocket;
    username?: string;
    currentUser?: User;
    openUsernameModal: boolean;
    openCardModal: boolean;
    openColumnModal: boolean;
    openTimerModal: boolean;
    clients: Client[];
    clientConnections: ClientConnection;
    columns: Column[];
    cards: Card[];
    numClients: number;
    tempColumn: Column;
    tempCard: Card;
    timer: Timer;
    standup: Standup;
    emojiPicker?: EmojiButton;

    // Methods
    initBoard(): void;
    wsConnect(username: string): void;
    onWsEvent(event: MessageEvent): void;
    applyOperation(list: any[], change: any): void;
    askUsername(): void;
    joinBoard(): void;
    columnNameById(id: string): string;
    editColumn(column: Column): void;
    saveColumn(): void;
    deleteColumn(column: Column): void;
    editCard(column: Column, card: Card): void;
    saveCard(): void;
    voteCard(card: Card, vote: number): void;
    deleteCard(card: Card): void;
    closeModal(name: string): void;
    onDragStart(event: DragEvent, card: Card): void;
    onDragEnd(event: DragEvent): void;
    onDragOver(event: DragEvent): void;
    onDragEnter(event: DragEvent): void;
    onDragLeave(event: DragEvent): void;
    onDrop(event: DragEvent, newColumn: Column): void;
    startTimer(): void;
    pauseTimer(): void;
    resumeTimer(): void;
    stopTimer(isCommand: boolean): void;
    playSound(): void;
    numClientConnections(userId: string): number;
    dispatchCustomEvents(eventName: string, message: string): void;
    gridColsClass(): string;
    openEmojiPicker(): void;
    getClientById(id: string): Client | undefined;
    startStandup(): void;
    closeStandup(): void;
    refreshShuffledUserIds(): void;
    setCurrentStandupUser(client: Client): void;
}