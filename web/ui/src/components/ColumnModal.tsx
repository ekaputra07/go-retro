import { useState, useRef, useEffect } from 'react'
import type { Column } from '../types'

interface props {
    column: Column
    onCancel(): void
    onDelete(col: Column): void
    onSave(col: Column): void
}

export function ColumnModal(p: props) {
    const inputRef = useRef<HTMLInputElement>(null)
    const [name, setName] = useState(p.column.name)

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
                    {p.column.id && <h2 className="font-semibold text-xl mb-6 text-gray-800">Edit Column</h2>}
                    {!p.column.id && <h2 className="font-semibold text-xl mb-6 text-gray-800">New Column</h2>}
                    <form onSubmit={(e) => { e.preventDefault(); p.onSave({ ...p.column, name: name || p.column.name }) }}>
                        <div className="mb-4">
                            <input
                                ref={inputRef}
                                value={name || p.column.name}
                                onChange={(e) => setName(e.target.value)}
                                placeholder="Column Name"
                                required
                                type="text" className="bg-gray-200 appearance-none border-2 border-gray-200 rounded-md w-full py-2 px-4 text-gray-700 leading-tight focus:outline-none focus:bg-white focus:border-sky-500" />
                        </div>
                        <div className="flex justify-between items-center mt-8 text-right">
                            {p.column.id && <input type="button" value="Delete" onClick={() => p.onDelete(p.column)} className="text-sm text-red-600 font-medium cursor-pointer" />}
                            <div className="flex-1">
                                <input type="button" value="Cancel" onClick={p.onCancel} className="bg-white hover:bg-gray-100 text-gray-700 font-semibold py-1 px-4 border border-gray-300 rounded-md shadow-sm mr-2" />
                                <input value="Save" type="submit" className="text-white font-semibold py-1 px-4 border border-transparent rounded-md shadow-sm bg-sky-600 hover:bg-sky-700" />
                            </div>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    )
}

// provides state and handlers for ColumnModal
export function useColumnModal(sender: (data: object) => void): [boolean, (col: Column | null) => void, props] {
    const [isOpen, setIsOpen] = useState(false)
    const [column, setColumn] = useState<Column>({ name: '' })

    function setOpen(column: Column | null) {
        setIsOpen(true)
        column && setColumn(column)
    }
    function onCancel() {
        setIsOpen(false)
    }
    function onSave(col: Column) {
        if (col.id) {
            sender({
                type: 'column.update',
                data: col
            })
        } else {
            sender({
                type: 'column.new',
                data: col
            })
        }
        setIsOpen(false)
    }
    function onDelete(col: Column) {
        sender({
            type: 'column.delete',
            data: col
        })
        setIsOpen(false)
    }
    return [
        isOpen,
        setOpen,
        { column, onCancel, onSave, onDelete }
    ]
}
