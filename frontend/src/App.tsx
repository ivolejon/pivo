import { useState } from 'react'

import './App.css'
import { useProjectStore } from './store/projectStore';

function App() {
  const { createProject } = useProjectStore();

  const [message, setMessage] = useState('Welcome to Pivo')
  return (
    <>
      <h1>Welcome to Pivo</h1>
      <div>
        <button onClick={() => createProject({ title: "New Project" })}>
        </button>
      </div>
    </>
  )
}

export default App
