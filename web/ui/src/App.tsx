import { useState, useEffect } from 'react'
import useWebSocket from 'react-use-websocket'
import { DndProvider } from 'react-dnd'
import { HTML5Backend } from 'react-dnd-html5-backend'
import Alert from './components/Alert'
import Timer from './components/Timer'
import NameModal from './components/NameModal'
import { TimerModal, useTimerModal } from './components/TimerModal'
import { ColumnModal, useColumnModal } from './components/ColumnModal'
import Toolbar from './components/Toolbar'
import Footer from './components/Footer'
import ColumnItem from './components/ColumnItem'
import CardItem from './components/CardItem'
import { Standup, useStandup } from './components/Standup'
import type { AppInfo, User } from './types'
import { useBoardState, useNotification } from './hooks'

declare global {
  interface Window {
    GORETRO_DATA?: {
      AppName: string
      AppVersion: string
      AppTagline: string
    };
  }
}

const appInfo: AppInfo = {
  name: window.GORETRO_DATA?.AppName || 'GoRetro',
  version: window.GORETRO_DATA?.AppVersion || 'ui-dev',
  tagline: window.GORETRO_DATA?.AppTagline || 'Minimalist retro board for happy teams 😉'
}

const gridColsClass = (length: number): string => {
  return {
    1: 'grid-cols-1',
    2: 'grid-cols-2',
    3: 'grid-cols-3',
    4: 'grid-cols-4',
    5: 'grid-cols-5',
    6: 'grid-cols-6',
  }[length] || 'grid-cols-4'
}

function App() {
  // required to be set
  const nameKey = 'GR_USERNAME'
  const [name, setName] = useState<string>('')
  const [socketUrl, setSocketUrl] = useState<string>('')
  const [connect, setConnect] = useState<boolean>(false)

  // webhook connection
  const { lastMessage, sendJsonMessage } = useWebSocket(socketUrl, {
    onOpen: () => {
      console.log('WebSocket connection opened.')
      sendJsonMessage({ type: 'me' })
    },
    onClose: () => console.log('WebSocket connection closed.'),
    onError: (event) => console.error('WebSocket error observed:', event),
    shouldReconnect: () => true,
  }, connect)

  // board state
  const [notification, setNotification] = useNotification()
  const { users, userConnectionsCount, columns, cards, timerRunning, timerState } = useBoardState(lastMessage, setNotification)
  const [standupOpen, standupSetOpen, standupProps] = useStandup(users, setNotification)
  const [timerModalOpen, timerModalSetOpen, timerModalProps] = useTimerModal(sendJsonMessage)
  const [columnModalOpen, columnModalSetOpen, columnModalProps] = useColumnModal(sendJsonMessage)

  const saveName = (name: string): void => {
    localStorage.setItem(nameKey, name)
    setName(name)
  }

  // read name from local storage
  useEffect(() => {
    const name = localStorage.getItem(nameKey) || ''
    if (name !== '') {
      setName(name)
    }
  }, [])

  // check for connection ready
  useEffect(() => {
    if (name == '') return

    let host = window.location.host
    if (import.meta.env.DEV) {
      host = import.meta.env.VITE_BACKEND_HOST
    }
    const protocol = window.location.protocol
    const pathname = window.location.pathname

    const wsProtocol: string = (protocol === 'https:') ? 'wss:' : 'ws:'
    setSocketUrl(`${wsProtocol}//${host}${pathname}/ws?u=${name}`)
    setConnect(true)
  }, [name])

  return (
    <DndProvider backend={HTML5Backend}>
      <div className="flex flex-col min-h-screen">
        <div className="flex-1">

          {notification !== '' && <Alert text={notification} />}
          {timerRunning && timerState && <Timer state={timerState} sender={sendJsonMessage} />}

          {/* I put a 100ms delay in NameModal so that it won't create a short blip */}
          {!connect && <NameModal onJoin={saveName} />}

          {connect &&
            <div className="py-4 px-6">
              {/* kanban board */}
              <div className="flex justify-items-start mt-4">

                {/* stand-up */}
                <div className="grid grid-cols-1 gap-4 pb-2 items-start">
                  {standupOpen && <Standup {...standupProps} />}
                </div>

                {/* columns */}
                <div className={"flex-1 grid gap-4 pb-2 items-start " + gridColsClass(columns.length)}>
                  {columns.map((col) =>
                    <ColumnItem column={col} sender={sendJsonMessage} key={col.id}>
                      {cards
                        .filter(c => c.column_id === col.id)
                        .map((c) => <CardItem column={col} card={c} sender={sendJsonMessage} key={c.id} />)}
                    </ColumnItem>
                  )}
                </div>
              </div>
            </div>
          }
        </div>

        <div className="flex flex-col">
          {timerModalOpen && <TimerModal {...timerModalProps} />}
          {columnModalOpen && <ColumnModal {...columnModalProps} />}
          <Toolbar
            users={users}
            conn={userConnectionsCount}
            showStandupBtn={!standupOpen}
            showTimerBtn={!timerRunning}
            onAvatarClick={(u: User) => setNotification(u.name)}
            onNewColumn={() => {
              if (columns.length >= 6) {
                setNotification("Can only create maximum 6 columns!")
              } else {
                columnModalSetOpen(null)
              }
            }}
            onNewStandup={standupSetOpen}
            onNewTimer={timerModalSetOpen}
          />
          <Footer
            userCount={users.length}
            appInfo={appInfo}
          />
        </div>
      </div>
    </DndProvider>
  )
}

export default App
