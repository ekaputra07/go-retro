import { useEffect, useState } from 'react'
import type { User } from '../types'

interface props {
    users: User[]
    onClose: () => void
    onSelectUser: (u: User, first: boolean) => void
}

export function Standup(p: props) {
    const [user, setUser] = useState<User | null>(null)
    const [currentUsers, setCurrentUsers] = useState<User[]>([])

    useEffect(() => {
        if (currentUsers.length == 0) {
            const users = shuffleUsers(p.users)
            setCurrentUsers(users)
            setUser(users[0])
        } else {
            setCurrentUsers(appendNewUsers(currentUsers, p.users))
        }
    }, [p.users])

    useEffect(() => {
        if (currentUsers.length > 0 && user != null) {
            p.onSelectUser(user, user.id == currentUsers[0].id)
        }
    }, [user])

    function shuffleIds(ids: string[]): string[] {
        const shuffledIds: string[] = [...ids]
        for (let i = shuffledIds.length - 1; i > 0; i--) {
            const j = Math.floor(Math.random() * (i + 1));
            [shuffledIds[i], shuffledIds[j]] = [shuffledIds[j], shuffledIds[i]]
        }
        return shuffledIds
    }

    function shuffleUsers(users: User[]): User[] {
        const ids = users.map((u) => u.id)
        const shuffledIds = shuffleIds(ids)
        return shuffledIds.map((id) => users.find((u) => u.id === id) as User)
    }

    function appendNewUsers(users: User[], newUsers: User[]): User[] {
        const ids = users.map((u) => u.id)
        const nusers = newUsers.filter((u) => ids.indexOf(u.id) == -1)
        return [...users, ...nusers]
    }

    return (
        <div className="bg-white pb-4 mr-4 rounded-md shadow overflow-y-auto overflow-x-hidden border-t-8 border-green-600 min-h-[150px] w-[200px]">
            <div className="flex justify-between items-center px-4 py-2">
                <h2 className="font-bold text-gray-800 text-2xl">Stand-up</h2>
                <span onClick={p.onClose} title="Stop stand-up" className="cursor-pointer text-gray-400 hover:text-gray-500">
                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth="1.5" stroke="currentColor" className="size-6">
                        <path strokeLinecap="round" strokeLinejoin="round" d="M6 18 18 6M6 6l12 12" />
                    </svg>
                </span>
            </div>
            <div className="py-2">
                {currentUsers.map((u) => (
                    <div key={u.id} title={u.name} onClick={() => setUser(u)} className={'flex justify-between items-center cursor-pointer px-4 py-2 hover:bg-green-50' + ((u.id == user?.id) ? ' bg-green-100' : '')}>
                        <img src={import.meta.env.BASE_URL + 'avatar/' + u.avatar_id + '.png'} alt="avatar" className="w-6 h-6 rounded-full border-2 border-white shadow-sm mr-2 cursor-pointer" />
                        <span className="flex-1 font-bold text-sm text-gray-700 cursor-pointer">{u.name}</span>
                        {(u.id == user?.id) &&
                            <span className="relative flex size-3 mr-1">
                                <span className="absolute inline-flex h-full w-full animate-ping rounded-full bg-green-700 opacity-100"></span>
                                <span className="relative inline-flex size-3 rounded-full bg-green-600"></span>
                            </span>
                        }
                    </div>
                ))}
            </div >
        </div >
    )
}

export function useStandup(users: User[], notifier: (msg: string) => void): [boolean, () => void, props] {
    const [isOpen, setIsOpen] = useState(false)
    function setOpen() {
        if (users.length < 2) {
            notifier('Not enough users for a stand-up session.')
            return
        }
        setIsOpen(true)
    }
    function onClose() {
        setIsOpen(false)
    }
    function onSelectUser(u: User, first: boolean) {
        if (first) {
            notifier(`${u.name} the lucky first 🥇`)
        } else {
            notifier(`${u.name}'s turn!`)
        }
    }
    return [isOpen, setOpen, {
        users, onClose, onSelectUser
    }]
}