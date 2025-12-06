export interface AppInfo {
    name: string
    version: string
    tagline: string
}

export interface User {
    id: string
    name: string
    avatar_id: number
}

export interface Client {
    id: string
    user: User
    created_at: number
}

export interface UserConnectionsCount {
    [key: string]: number
}

export interface Column {
    name: string
    id?: string
    created_at?: number
}

export interface Card {
    name: string
    column_id?: string
    id?: string
    created_at?: number
    votes?: number
}

export interface TimerState {
    status: string
    display: string
}

export interface ChangeOp<T> {
    type: "clients" | "columns" | "cards"
    op: "put" | "del"
    id: string
    obj?: T
}

export interface Message {
    type: string
    data: string | TimerState | User
    user: User
}

export interface MessageList {
    type: string
    messages: Message[]
}

export type WSMessage = Message | MessageList | ChangeOp<Client> | ChangeOp<Column> | ChangeOp<Card>