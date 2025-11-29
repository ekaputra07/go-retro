import type { TimerState } from '../types'

interface props {
    state: TimerState
    sender: (msg: object) => void
}

export default function Timer(p: props) {
    const pause = (): void => {
        p.sender({
            type: 'timer.cmd',
            data: { cmd: 'pause' }
        })
    }
    const resume = (): void => {
        p.sender({
            type: 'timer.cmd',
            data: { cmd: 'start' }
        })
    }
    const stop = (): void => {
        p.sender({
            type: 'timer.cmd',
            data: { cmd: 'stop' }
        })
    }

    if (p.state.status == "stopped") {
        return <></>
    }
    return (
        <div className={"p-3 " + (p.state.status == "done" ? "bg-orange-600" : "bg-green-600")}>
            <div className="flex justify-center items-center text-white gap-2 relative">
                {p.state.status == "done" && <div className="text-xl font-bold leading-none">Time's up<span className="italic">!</span></div>}

                {(p.state.status != "done" && p.state.status != "stopped") &&
                    <div className="flex items-center gap-2">
                        <div className="text-xl font-bold leading-none">{p.state.display}</div>

                        {p.state.status == "running" &&
                            <button type="button" onClick={pause} title="Pause" className="cursor-pointer">
                                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M12 2C6.477 2 2 6.477 2 12s4.477 10 10 10s10-4.477 10-10S17.523 2 12 2m-1 14H9V8h2zm4 0h-2V8h2z" /></svg>
                            </button>
                        }

                        {p.state.status == "paused" &&
                            <button type="button" onClick={resume} title="Resume" className="cursor-pointer">
                                <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M12 2a10 10 0 1 0 10 10A10 10 0 0 0 12 2m-2 14.5v-9l6 4.5z" /></svg>
                            </button>
                        }
                    </div>
                }

                <button type="button" onClick={stop} title="Close" className="absolute top-1 right-0 mr-2 cursor-pointer">
                    <svg xmlns="http://www.w3.org/2000/svg" width="14" height="14" viewBox="0 0 14 14"><path fill="none" stroke="currentColor" strokeLinecap="round" strokeLinejoin="round" d="m13.5.5l-13 13m0-13l13 13" /></svg>
                </button>
            </div>
        </div>
    )
}