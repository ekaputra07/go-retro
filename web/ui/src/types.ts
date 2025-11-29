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
    joined_at: number
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

export interface ChangeOp {
    op: "put" | "del"
    id: string
    obj?: any
}

export type ChangeableList = Column[] | Card[]
