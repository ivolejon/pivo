import './App.css'
import { useProjectStore } from './store/projectStore';
import { useWebSocketHandler } from './lib/webSocket';

const clientId = "b15377e4-60f1-11f0-9ce3-834692c66f23";

function App() {
  const { createProject, uploadDocument, askQuestion } = useProjectStore();
  const { connectionStatus, answer } = useWebSocketHandler(clientId);

  async function uploadFile(event: React.ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0];
    if (file) {
      try {
        await uploadDocument("0e1b56a5-71c5-4359-8006-2074c2637d55", file);
        alert("File uploaded successfully");
      } catch (error: any) {
        alert(`Error uploading file: ${error.message}`);
      }
    }
  }

  async function sendQuestion() {
    const question = {
      question: "Vem var ordförande?",
      projectId: "103cc611-e095-4f22-ac8e-87a943e54d23",
    };
    try {
      await askQuestion(question);
    } catch (error: any) {
      alert(`Error asking question: ${error.message}`);
    }
  }

  return (
    <>
      <h1>Welcome to Pivo</h1>
      <div>
        <p>WebSocket Status: {connectionStatus}</p>
        <p>{answer}</p>
        <button onClick={() => createProject({ title: "New Project" })}>
          Add Project
        </button>
        <input
          type="file"
          onChange={uploadFile}
        />
        <button onClick={sendQuestion}>Ask question</button>
      </div>
    </>
  )
}

export default App
