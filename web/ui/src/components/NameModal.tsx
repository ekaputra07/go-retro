import { useState, useRef, useEffect } from 'react'

interface props {
    onJoin(name: string): void
}

export default function NameModal(p: props) {
    const inputRef = useRef<HTMLInputElement>(null)
    const [visible, setVisible] = useState(false)
    const [name, setName] = useState('')

    useEffect(() => {
        if (inputRef.current) {
            inputRef.current.focus()
        }
        // show the modal after 100ms delay
        setTimeout(() => setVisible(true), 100)
    }, [])

    return (
        <div className={"fixed inset-0 flex h-screen w-full items-end md:items-center justify-center z-10 " + (visible ? '' : 'invisible')}>
            <div className="absolute inset-0 bg-black opacity-50"></div>
            <div className="md:p-4 md:max-w-lg mx-auto w-full flex-1 relative overflow-hidden">
                <form onSubmit={e => { e.preventDefault(); p.onJoin(name) }}>
                    <div className="w-full rounded-t-lg md:rounded-md bg-white p-8">
                        <h2 className="font-semibold text-xl mb-6 text-gray-800">Your name</h2>
                        <div className="mb-4">
                            <input
                                onChange={e => setName(e.target.value)}
                                value={name}
                                ref={inputRef}
                                required
                                type="text" className="bg-gray-200 appearance-none border-2 border-gray-200 rounded-md w-full py-2 px-4 text-gray-700 leading-tight focus:outline-none focus:bg-white focus:border-sky-500" />
                            <p className="text-gray-500 text-sm mt-2">Name used to show who's joining, cards are anonymous.</p>
                        </div>
                        <div className="flex justify-between items-center mt-8 text-right">
                            <div className="flex-1">
                                <input type="submit" value="Join" className="text-white font-semibold py-1 px-4 border border-transparent rounded-md shadow-sm bg-sky-600 hover:bg-sky-700" />
                            </div>
                        </div>
                    </div>
                </form>
            </div>
        </div>
    )
}