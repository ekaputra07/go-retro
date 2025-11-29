interface props {
    text: string
}

export default function Alert(p: props) {
    return (
        <div className="fixed w-full z-50 flex inset-0 items-start justify-center pointer-events-none md:mt-5">
            <div role="alert" className="w-full px-4 py-4 md:max-w-sm bg-gray-900 md:rounded-md shadow-lg">
                <div className="flex items-center">
                    <div className="shrink-0 mr-3">
                        <svg className="h-6 w-6 text-gray-400" viewBox="0 0 20 20" fill="currentColor">  <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" /></svg>
                    </div>
                    <div className="text-gray-200 text-base">{p.text}</div>
                </div>
            </div>
        </div>
    )
}

