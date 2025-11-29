import { useEffect, useState } from 'react'
import type { Client, UserConnectionsCount, User, Column, Card, ChangeOp, ChangeableList, TimerState } from './types'

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

function applyChangeOperation(list: ChangeableList, change: ChangeOp): ChangeableList {
    // we need to return a new list reference to trigger re-render
    const newList = [...list]
    const idx = newList.findIndex(c => c.id === change.id)
    if (idx >= 0) {
        if (change.op === "put") {
            newList[idx] = change.obj
        } else if (change.op === "del") {
            newList.splice(idx, 1)
        }
        return newList.sort((a, b) => (a.created_at!) - (b.created_at!))
    } else if (idx === -1 && change.op === "put") {
        newList.push(change.obj)
        return newList.sort((a, b) => (a.created_at!) - (b.created_at!))
    }
    return newList
}

export function useBoardState(
    lastMessage: MessageEvent | null,
    onNotification?: (msg: string) => void,
): BoardState {
    const [currentUser, setCurrentUser] = useState<User | null>(null)
    const [clients, setClients] = useState<Client[]>([])
    const [columns, setColumns] = useState<Column[]>([])
    const [cards, setCards] = useState<Card[]>([])
    const [notification, setNotification] = useState<string>('')
    const [timerState, setTimerState] = useState<TimerState | null>(null)

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
        const e: any = JSON.parse(event.data)
        switch (e.type) {
            case "me":
                setCurrentUser(e.user)
                break

            case "columns":
                setColumns(applyChangeOperation(columns, e))
                break

            case "cards":
                setCards(applyChangeOperation(cards, e))
                break

            case "board.users":
                const sortedClients: Client[] = e.data.sort(
                    (a: Client, b: Client) => a.joined_at - b.joined_at,
                )
                setClients(sortedClients)
                break

            case "board.notification":
                // don't show notification to those who triggers it
                if (currentUser && e.user) {
                    if (currentUser.id != e.user.id) {
                        setNotification(e.data)
                        if (onNotification) {
                            onNotification(e.data)
                        }
                    }
                }
                break

            case "timer.state":
                const state: TimerState = e.data
                setTimerState(state)

                if (state.status == 'done') {
                    setTimeout(() => {
                        setTimerState({ ...state, status: 'stopped' })
                    }, 5000)
                }
                break

            default:
                break
        }
    }

    useEffect(() => {
        if (lastMessage !== null) {
            handleEvent(lastMessage)
        }
    }, [lastMessage])

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
    }, [notification])
    return [notification, setNotification]
}