import { useState } from 'react'
import reactLogo from './assets/react.svg'
import viteLogo from '/vite.svg'
import { Button, Card } from "@radix-ui/themes";
import { Chat } from './components/Chat/Chat'
import './styles/App.css'

function App() {
  const [count, setCount] = useState(0)

  return (
    <>
      <div className='app-container'>
        <div className='top-bar-container'>ivo</div>
        <div className='left-menu-container'>Menu</div>
        <div className='main-container'>
          <div className="left">
            <Button onClick={() => setCount((count) => count + 1)}>
              count is {count}
            </Button>
            <Chat />
          </div>

        </div>
      </div >
      {/* <h1>Vite + React</h1>
      <div className="card">
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p> */}
    </>
  )
}

export default App
