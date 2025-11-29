import { useRef } from 'react'
import { useDrop, type DropTargetMonitor } from 'react-dnd'
import type { Column, Card } from '../types'
import { ColumnModal, useColumnModal } from './ColumnModal'
import { CardModal, useCardModal } from './CardModal'

interface props extends React.PropsWithChildren {
    column: Column
    sender: (data: object) => void
}

export default function ColumnItem(p: props) {
    const [showColModal, setShowColModal, colModalProps] = useColumnModal(p.sender)
    const [showCardModal, setShowCardModal, cardModalProps] = useCardModal(p.sender)
    const dropZoneRef = useRef<HTMLDivElement>(null)

    const handleCardDrop = (card: Card) => {
        p.sender({
            type: 'card.update',
            data: { ...card, column_id: p.column.id }
        })
    }
    const [{ dropIsOver }, dropConnector] = useDrop(() => ({
        accept: 'card',
        drop: handleCardDrop,
        collect: (monitor: DropTargetMonitor) => ({
            dropIsOver: monitor.isOver(),
        })
    }))
    dropConnector(dropZoneRef)

    return (
        <div className="bg-slate-100 pb-4 rounded-md shadow overflow-y-auto overflow-x-hidden border-t-8 border-sky-600 min-h-[110px]">
            <div className="flex justify-between items-center px-4 py-2 bg-gray-100 mb-2">
                <h2 className="font-bold text-gray-800 text-2xl">{p.column.name}</h2>
                {showCardModal && <CardModal {...cardModalProps} />}
                {showColModal && <ColumnModal {...colModalProps} />}
                <button onClick={() => setShowColModal(p.column)} className="cursor-pointer text-gray-500 hover:text-gray-700" title="Column settings">
                    <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" fill="none" viewBox="0 0 24 24" strokeWidth="1.5" stroke="currentColor" className="size-6">
                        <path strokeLinecap="round" strokeLinejoin="round" d="M9.594 3.94c.09-.542.56-.94 1.11-.94h2.593c.55 0 1.02.398 1.11.94l.213 1.281c.063.374.313.686.645.87.074.04.147.083.22.127.325.196.72.257 1.075.124l1.217-.456a1.125 1.125 0 0 1 1.37.49l1.296 2.247a1.125 1.125 0 0 1-.26 1.431l-1.003.827c-.293.241-.438.613-.43.992a7.723 7.723 0 0 1 0 .255c-.008.378.137.75.43.991l1.004.827c.424.35.534.955.26 1.43l-1.298 2.247a1.125 1.125 0 0 1-1.369.491l-1.217-.456c-.355-.133-.75-.072-1.076.124a6.47 6.47 0 0 1-.22.128c-.331.183-.581.495-.644.869l-.213 1.281c-.09.543-.56.94-1.11.94h-2.594c-.55 0-1.019-.398-1.11-.94l-.213-1.281c-.062-.374-.312-.686-.644-.87a6.52 6.52 0 0 1-.22-.127c-.325-.196-.72-.257-1.076-.124l-1.217.456a1.125 1.125 0 0 1-1.369-.49l-1.297-2.247a1.125 1.125 0 0 1 .26-1.431l1.004-.827c.292-.24.437-.613.43-.991a6.932 6.932 0 0 1 0-.255c.007-.38-.138-.751-.43-.992l-1.004-.827a1.125 1.125 0 0 1-.26-1.43l1.297-2.247a1.125 1.125 0 0 1 1.37-.491l1.216.456c.356.133.751.072 1.076-.124.072-.044.146-.086.22-.128.332-.183.582-.495.644-.869l.214-1.28Z" />
                        <path strokeLinecap="round" strokeLinejoin="round" d="M15 12a3 3 0 1 1-6 0 3 3 0 0 1 6 0Z" />
                    </svg>
                </button>
            </div>

            <div className="px-4">
                <div ref={dropZoneRef} className={"pb-10 rounded-md " + (dropIsOver ? 'bg-blue-200' : '')}>
                    {p.children}
                </div>
                <div className="text-center">
                    <button onClick={() => { setShowCardModal(p.column, null) }} className="inline-flex items-center text-gray-700 text-sm font-medium cursor-pointer">
                        <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 4v16m8-8H4" />
                        </svg>
                        Add Card
                    </button>
                </div>
            </div>
        </div>
    )
}