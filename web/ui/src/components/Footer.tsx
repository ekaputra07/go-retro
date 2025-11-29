import type { AppInfo } from '../types'

interface props {
    userCount: number
    appInfo: AppInfo
}

const usersOnlineText = (count: number): string => {
    if(count == 1) {
        return "1 user online"
    } else {
        return `${count} users online`
    }
}

export default function Footer(p: props) {
    return (
        <div className="flex justify-between items-center bg-white px-4 py-2 shadow text-xs text-gray-600 text-center">
            <div className="flex justify-between md:justify-start items-center gap-2 flex-1">
                <div className="flex items-center">
                    <span className="flex w-2 h-2 me-1 bg-green-500 rounded-full"></span> 
                    <span>{usersOnlineText(p.userCount)}</span>
                </div>
            </div>
            <p className="text-xs text-gray-600 text-center hidden md:block">
                <a href="https://github.com/ekaputra07/go-retro" className="underline" target="_blank">{p.appInfo.name} ({p.appInfo.version})</a> - {p.appInfo.tagline}
            </p>
        </div>
    )
}
