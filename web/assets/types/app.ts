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