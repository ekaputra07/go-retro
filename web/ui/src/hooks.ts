import { useCallback, useEffect, useMemo, useState } from 'react'
import type { Client, UserConnectionsCount, User, Column, Card, ChangeOp, TimerState, Message, MessageList, WSMessage } from './types'

export interface BoardState {
    currentUser: User | null
    users: User[]
    userConnectionsCount: UserConnectionsCount
    columns: Column[]
    cards: Card[]
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

export function useBoardState(
    lastMessage: MessageEvent | null,
    onNotification?: (msg: string) => void,
): BoardState {
    const [currentUser, setCurrentUser] = useState<User | null>(null)
    const [clients, setClients] = useState<Client[]>([])
    const [columns, setColumns] = useState<Column[]>([])
    const [cards, setCards] = useState<Card[]>([])
    const [timerState, setTimerState] = useState<TimerState | null>(null)

    const handleMsg = useCallback((m: WSMessage) => {
        switch (m.type) {
            case "me":
                setCurrentUser((m as Message).user)
                break

            case "columns":
                setColumns(applyChangeOperation(columns, m as ChangeOp<Column>))
                break

            case "cards":
                setCards(applyChangeOperation(cards, m as ChangeOp<Card>))
                break

            case "clients":
                setClients(applyChangeOperation(clients, m as ChangeOp<Client>))
                break

            case "board.notification":
                const msg = m as Message
                // don't show notification to those who triggers it
                if (currentUser && msg.user) {
                    if (currentUser.id != msg.user.id) {
                        if (onNotification) {
                            onNotification(msg.data as string)
                        }
                    }
                }
                break

            case "timer.state":
                {
                    const msg = m as Message
                    const ts = msg.data as TimerState
                    setTimerState(ts)
                    if (ts.status == 'done') {
                        setTimeout(() => {
                            setTimerState({ ...timerState!, status: 'stopped' })
                        }, 5000)
                    }
                }
                break

            default:
                break
        }
    }, [lastMessage])

    useEffect(() => {
        if (lastMessage !== null) {
            const message: WSMessage = JSON.parse(lastMessage.data)
            // if message list, unpack and handle one-by-one
            if (message.type === 'messages') {
                const messages = (message as MessageList).messages
                for (let i = 0; i < messages.length; i++) {
                    handleMsg(messages[i])
                }
            } else {
                handleMsg(message)
            }
        }
    }, [lastMessage])

    const connectionsCount: UserConnectionsCount = useMemo(() => {
        return clients.reduce(
            (acc: UserConnectionsCount, c: Client) => {
                if (!acc[c.user.id]) acc[c.user.id] = 0
                acc[c.user.id]++
                return acc
            },
            {},
        )
    }, [clients])

    const users: User[] = useMemo(() => {
        const uniqueUserIds: string[] = [
            ...new Set(clients.map((c: Client) => c.user.id)),
        ]
        return uniqueUserIds
            .map((id) => clients.find((c) => c.user.id === id) as Client)
            .map(c => c.user)
    }, [clients])

    return {
        currentUser,
        users,
        userConnectionsCount: connectionsCount,
        columns,
        cards,
        timerRunning: timerState !== null && timerState.status !== 'stopped',
        timerState,
    }
}

let timeoutId: number | null
export function useNotification(timeout: number = 3000): [string, (msg: string) => void] {
    const [notification, setNotification] = useState<string>('')

    // only set message when its not empty and not the same as previous one
    const set = (msg: string) => {
        const trimmed = msg.trim()
        if (trimmed === '' || trimmed === notification) return

        setNotification(trimmed)

        // clear timeout before setting new one
        if (timeoutId) {
            clearTimeout(timeoutId)
        }
        timeoutId = setTimeout(() => {
            setNotification('')
        }, timeout)
    }

    return [notification, set]
}