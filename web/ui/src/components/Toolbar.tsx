import type { User, UserConnectionsCount } from "../types"

interface props {
    users: User[]
    conn: UserConnectionsCount
    showStandupBtn: boolean
    showTimerBtn: boolean
    onAvatarClick(user: User): void
    onNewStandup(): void
    onNewTimer(): void
    onNewColumn(): void
}

function numConnections(conn: UserConnectionsCount, userId: string): number {
    return conn[userId]
}

export default function Toolbar(p: props) {
    return (
        <div className="flex justify-between items-center pb-2 px-4">
            <div className="flex items-center justify-center gap-2 mb-1">
                {p.users.map(u => (
                    <div key={u.id} title={u.name} onClick={() => p.onAvatarClick(u)} className="flex flex-col relative isolate items-center justify-center cursor-pointer">
                        <img src={import.meta.env.BASE_URL + 'avatar/' + u.avatar_id + '.png'} alt="avatar" className="w-12 h-12 rounded-full border-2 border-white shadow-sm" />

                        {numConnections(p.conn, u.id) > 1 &&
                            <span className="absolute top-0 right-0 z-10 flex items-center justify-around w-4 h-4 bg-white rounded-full text-gray-600 shadow text-xs font-semibold">
                                {numConnections(p.conn, u.id)}
                            </span>
                        }

                        <span className="font-bold text-sm text-gray-700">{u.name}</span>
                    </div>
                ))}

            </div>
            <div className="flex items-center justify-center gap-2">

                {p.showStandupBtn &&
                    <button onClick={p.onNewStandup} title="Start stand-up" className="w-12 h-12 flex items-center justify-center text-white shadow rounded-full bg-sky-600 hover:bg-sky-700 z-9 cursor-pointer border-2 border-white">
                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" width="24" height="24" color="currentColor" fill="none"><path d="M19.5576 4L20.4551 4.97574C20.8561 5.41165 21.0566 5.62961 20.9861 5.81481C20.9155 6 20.632 6 20.0649 6C18.7956 6 17.2771 5.79493 16.1111 6.4733C15.3903 6.89272 14.8883 7.62517 14.0392 9M3 18H4.58082C6.50873 18 7.47269 18 8.2862 17.5267C9.00708 17.1073 9.50904 16.3748 10.3582 15" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"></path><path d="M19.5576 20L20.4551 19.0243C20.8561 18.5883 21.0566 18.3704 20.9861 18.1852C20.9155 18 20.632 18 20.0649 18C18.7956 18 17.2771 18.2051 16.1111 17.5267C15.2976 17.0534 14.7629 16.1815 13.6935 14.4376L10.7038 9.5624C9.63441 7.81853 9.0997 6.9466 8.2862 6.4733C7.47269 6 6.50873 6 4.58082 6H3" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"></path></svg>
                    </button>
                }

                {p.showTimerBtn &&
                    <button onClick={p.onNewTimer} title="Start timer" className="w-12 h-12 flex items-center justify-center text-white shadow rounded-full bg-sky-600 hover:bg-sky-700 z-9 cursor-pointer border-2 border-white">
                        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M9 3V1h6v2zm3 19q-1.85 0-3.488-.712T5.65 19.35t-1.937-2.863T3 13t.713-3.488T5.65 6.65t2.863-1.937T12 4q1.55 0 2.975.5t2.675 1.45l1.4-1.4l1.4 1.4l-1.4 1.4Q20 8.6 20.5 10.025T21 13q0 1.85-.713 3.488T18.35 19.35t-2.863 1.938T12 22m0-2q2.9 0 4.95-2.05T19 13t-2.05-4.95T12 6T7.05 8.05T5 13t2.05 4.95T12 20m-2-3l6-4l-6-4z" /></svg>
                    </button>
                }

                <button onClick={p.onNewColumn} title="New Column" className="w-12 h-12 flex items-center justify-center text-white shadow rounded-full bg-sky-600 hover:bg-sky-700 z-9 cursor-pointer border-2 border-white">
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="none" stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M5 12h14m-7-7v14" /></svg>
                </button>
            </div>
        </div>
    )
}