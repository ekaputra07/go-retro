import { useState, useRef, useEffect } from 'react'

interface props {
    onCancel(): void
    onStart(duration: string): void
}

export function TimerModal(p: props) {
    const inputRef = useRef<HTMLInputElement>(null)
    const [duration, setDuration] = useState('5m')

    useEffect(() => {
        if (inputRef.current) {
            inputRef.current.focus()
        }
    }, [])

    return (
        <div className="fixed inset-0 flex h-screen w-full items-end md:items-center justify-center z-10">
            <div className="absolute inset-0 bg-black opacity-50"></div>
            <div className="md:p-4 md:max-w-lg mx-auto w-full flex-1 relative overflow-hidden">
                <div className="w-full rounded-t-lg md:rounded-md bg-white p-8">
                    <h2 className="font-semibold text-xl mb-6 text-gray-800">Timer duration</h2>
                    <form onSubmit={e => { e.preventDefault(); p.onStart(duration) }}>
                        <div className="mb-4">
                            <input
                                ref={inputRef}
                                onChange={e => setDuration(e.target.value)}
                                value={duration}
                                required
                                type="text"
                                className="bg-gray-200 appearance-none border-2 border-gray-200 rounded-md w-full py-2 px-4 text-gray-700 leading-tight focus:outline-none focus:bg-white focus:border-sky-500" />
                            <p className="text-gray-500 text-sm mt-2">Supported formats: <strong>5m</strong> or <strong>30s</strong> or <strong>5m30s</strong></p>
                        </div>
                        <div className="flex justify-between items-center mt-8 text-right">
                            <div className="flex-1">
                                <input onClick={p.onCancel} type="button" value="Cancel" className="bg-white hover:bg-gray-100 text-gray-700 font-semibold py-1 px-4 border border-gray-300 rounded-md shadow-sm mr-2" />
                                <input value="Start" type="submit" className="text-white font-semibold py-1 px-4 border border-transparent rounded-md shadow-sm bg-sky-600 hover:bg-sky-700" />
                            </div>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    )
}

// provides state and handlers for TimerModal
export function useTimerModal(sender: (data: object) => void): [boolean, () => void, props] {
    const [isOpen, setIsOpen] = useState(false)

    function setOpen() {
        setIsOpen(true)
    }

    function onCancel() {
        setIsOpen(false)
    }

    function onStart(duration: string) {
        sender({
            type: 'timer.cmd',
            data: { cmd: 'start', value: duration }
        })
        setIsOpen(false)
    }

    return [
        isOpen,
        setOpen,
        { onCancel, onStart }
    ]
}
