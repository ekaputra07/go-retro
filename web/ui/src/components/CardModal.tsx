import { useState, useRef, useEffect } from 'react'
import { EmojiButton } from "@joeattardi/emoji-button"
import type { Column, Card } from '../types'

interface props {
    column: Column
    card: Card
    onCancel(): void
    onDelete(col: Card): void
    onSave(col: Card): void
}

const emojiPicker = new EmojiButton({
    autoHide: true,
    showVariants: false,
    position: "bottom-start",
    zIndex: 99,
})

export function CardModal(p: props) {
    const inputRef = useRef<HTMLInputElement>(null)
    const emojiBtnRef = useRef<HTMLButtonElement>(null)
    const [name, setName] = useState(p.card.name)

    const openEmojiPicker = (): void => {
        emojiPicker.showPicker(emojiBtnRef.current as HTMLElement)
    }

    useEffect(() => {
        inputRef.current && inputRef.current.focus()

        emojiPicker.on("emoji", selection => {
            const newName = name + selection.emoji
            setName(newName)
        })
    }, [name])

    return (
        <div className="fixed inset-0 flex h-screen w-full items-end md:items-center justify-center z-10">
            <div className="absolute inset-0 bg-black opacity-50"></div>
            <div className="md:p-4 md:max-w-lg mx-auto w-full flex-1 relative overflow-hidden">
                <div className="w-full rounded-t-lg md:rounded-md bg-white p-8">

                    {p.card.id && <h2 className="font-semibold text-xl mb-6 text-gray-800">Edit card</h2>}
                    {!p.card.id && <h2 className="font-semibold text-xl mb-6 text-gray-800">New card for <span className="leading-normal text-sky-600">{p.column.name}</span></h2>}

                    <form onSubmit={(e) => { e.preventDefault(); p.onSave({ ...p.card, name: name }) }}>
                        <div className="mb-4 relative">
                            <button type="button" onClick={openEmojiPicker} ref={emojiBtnRef} className="absolute top-0 right-0 mt-2 mr-2 cursor-pointer" title="Add emoji">
                                <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth="1.5" stroke="currentColor" className="size-6 text-gray-500 hover:text-gray-700">
                                    <path strokeLinecap="round" strokeLinejoin="round" d="M15.182 15.182a4.5 4.5 0 0 1-6.364 0M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0ZM9.75 9.75c0 .414-.168.75-.375.75S9 10.164 9 9.75 9.168 9 9.375 9s.375.336.375.75Zm-.375 0h.008v.015h-.008V9.75Zm5.625 0c0 .414-.168.75-.375.75s-.375-.336-.375-.75.168-.75.375-.75.375.336.375.75Zm-.375 0h.008v.015h-.008V9.75Z" />
                                </svg>
                            </button>
                            <input
                                ref={inputRef}
                                value={name}
                                onChange={(e) => setName(e.target.value)}
                                required
                                placeholder="What's on your mind?"
                                type="text" className="bg-gray-200 appearance-none border-2 border-gray-200 rounded-md w-full py-2 px-4 pr-8 text-gray-700 leading-tight focus:outline-none focus:bg-white focus:border-sky-500" />
                        </div>

                        <div className="flex justify-between items-center mt-8 text-right">
                            {p.card.id && <input type="button" value="Delete" onClick={() => p.onDelete(p.card)} className="text-sm text-red-600 font-medium cursor-pointer" />}
                            <div className="flex-1">
                                <input onClick={p.onCancel} type="button" value="Cancel" className="bg-white hover:bg-gray-100 text-gray-700 font-semibold py-1 px-4 border border-gray-300 rounded-md shadow-sm mr-2" />
                                <input type="submit" value="Save" className="text-white font-semibold py-1 px-4 border border-transparent rounded-md shadow-sm bg-sky-600 hover:bg-sky-700" />
                            </div>
                        </div>
                    </form>
                </div>
            </div>
        </div>
    )
}

export function useCardModal(sender: (data: object) => void): [boolean, (col: Column, card: Card | null) => void, props] {
    const [isOpen, setIsOpen] = useState(false)
    const [column, setColumn] = useState<Column>({ name: '' })
    const [card, setCard] = useState<Card>({ name: '' })

    function setOpen(column: Column, card: Card | null) {
        setIsOpen(true)
        setColumn(column)
        if (card) {
            setCard(card)
        } else {
            setCard({ name: '', column_id: column.id })
        }
    }
    function onCancel() {
        setIsOpen(false)
    }
    function onSave(c: Card) {
        if (c.id) {
            sender({
                type: 'card.update',
                data: c
            })
        } else {
            sender({
                type: 'card.new',
                data: { ...c, column_id: column.id }
            })
        }
        setIsOpen(false)
    }
    function onDelete(c: Card) {
        sender({
            type: 'card.delete',
            data: c
        })
        setIsOpen(false)
    }
    return [
        isOpen,
        setOpen,
        { column, card, onCancel, onSave, onDelete }
    ]
}