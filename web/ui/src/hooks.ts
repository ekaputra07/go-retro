import { useEffect, useState } from 'react'
import type { Client, UserConnectionsCount, User, Column, Card, ChangeOp, TimerState, WSMessage } from './types'

export interface BoardState {
    currentUser: User | null
    users: User[]
    userConnectionsCount: UserConnectionsCount
    columns: Column[]
    cards: Card[]
    notification: string
    timerRunning: boolean
    timerState: TimerState | null
}

function applyChangeOperation<T>(list: T[], change: ChangeOp<T>): T[] {
    const obj: T = change.obj as T

    // we need to return a new list reference to trigger re-render
    const newList = [...list]
    const idx = newList.findIndex((item: T) => (item as { id: string }).id === change.id)
    if (idx >= 0) {
        if (change.op === "put") {
            newList[idx] = obj
        } else if (change.op === "del") {
            newList.splice(idx, 1)
        }
        return newList.sort((a: T, b: T) => ((a as { created_at: number }).created_at) - ((b as { created_at: number }).created_at))
    } else if (idx === -1 && change.op === "put") {
        newList.push(obj)
        return newList.sort((a: T, b: T) => ((a as { created_at: number }).created_at) - ((b as { created_at: number }).created_at))
    }
    return newList
}

let currentUser: User | null = null
let clients: Client[] = []
let columns: Column[] = []
let cards: Card[] = []
let timerState: TimerState | null = null

export function useBoardState(
    lastMessage: MessageEvent | null,
    onNotification?: (msg: string) => void,
): BoardState {

    // const [currentUser, setCurrentUser] = useState<User | null>(null)
    // const [clients, setClients] = useState<Client[]>([])
    // const [columns, setColumns] = useState<Column[]>([])
    // const [cards, setCards] = useState<Card[]>([])
    const [notification, setNotification] = useState<string>('')
    // const [timerState, setTimerState] = useState<TimerState | null>(null)

    const connectionsCount: UserConnectionsCount = clients.reduce(
        (acc: UserConnectionsCount, c: Client) => {
            if (!acc[c.user.id]) acc[c.user.id] = 0
            acc[c.user.id]++
            return acc
        },
        {},
    )
    const uniqueUserIds: string[] = [
        ...new Set(clients.map((c: Client) => c.user.id)),
    ]
    const users: User[] = uniqueUserIds
        .map((id) => clients.find((c) => c.user.id === id) as Client)
        .map(c => c.user)

    const handleEvent = (event: MessageEvent): void => {
        const e: WSMessage = JSON.parse(event.data)
        switch (e.type) {
            case "me":
                currentUser = e.user
                break

            case "columns":
                columns = applyChangeOperation(columns, e as ChangeOp<Column>)
                break

            case "cards":
                cards = applyChangeOperation(cards, e as ChangeOp<Card>)
                break

            case "clients":
                clients = applyChangeOperation(clients, e as ChangeOp<Client>)
                break

            case "board.notification":
                // don't show notification to those who triggers it
                if (currentUser && e.user) {
                    if (currentUser.id != e.user.id) {
                        setNotification(e.data as string)
                        if (onNotification) {
                            onNotification(e.data as string)
                        }
                    }
                }
                break

            case "timer.state":
                {
                    timerState = e.data as TimerState
                    // if (timerState.status == 'done') {
                    //     setTimeout(() => {
                    //         timerState = { ...timerState, status: 'stopped' }
                    //     }, 5000)

                    // }
                }
                break

            default:
                break
        }
    }
    if (lastMessage) {
        handleEvent(lastMessage)
    }

    // useEffect(() => {
    //     if (lastMessage !== null) {
    //         handleEvent(lastMessage)
    //     }
    //     // eslint-disable-next-line react-hooks/exhaustive-deps
    // }, [lastMessage])

    return {
        currentUser,
        users,
        userConnectionsCount: connectionsCount,
        columns,
        cards,
        notification,
        timerRunning: timerState !== null && timerState.status !== 'stopped',
        timerState,
    }
}

export function useNotification(timeout: number = 3000): [string, (msg: string) => void] {
    const [notification, setNotification] = useState<string>('')

    useEffect(() => {
        if (notification !== '') {
            const timeoutId = setTimeout(() => {
                setNotification('')
            }, timeout)

            return () => clearTimeout(timeoutId)
        }
    }, [notification, timeout])
    return [notification, setNotification]
}