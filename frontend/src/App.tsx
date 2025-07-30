import './App.css'
import { useProjectStore } from './store/projectStore';
import { useDocumentStore } from './store/documentStore';
import { useWebSocketHandler } from './lib/webSocket';
import { useEffect, useState } from 'react';
import type { Project } from './types';


const clientId = "b15377e4-60f1-11f0-9ce3-834692c66f23";


function App() {
  const { createProject, askQuestion, fetchProjects, projects, currentProject, setCurrentProject } = useProjectStore();
  const { documents, fetchDocuments, uploadDocument } = useDocumentStore();
  const { connectionStatus, answer } = useWebSocketHandler(clientId);
  const [question, setQuestion] = useState("");
  // When the this loads, fetch the projects
  useEffect(() => {
    fetchProjects().catch((error) => {
      console.error("Failed to fetch projects:", error);
    });
  }, [fetchProjects]);

  async function uploadFile(event: React.ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0];
    if (file && currentProject) {
      try {
        await uploadDocument(currentProject.id, file);
        alert("File uploaded successfully");
      } catch (error: any) {
        alert(`Error uploading file: ${error.message}`);
      }
    }
  }

  async function sendQuestion() {
    if (!currentProject) {
      alert("Please select a project first.");
      return;
    }
    const questionPayload = {
      question: question,
      projectId: currentProject.id,
    };
    try {
      await askQuestion(questionPayload);
    } catch (error: any) {
      alert(`Error asking question: ${error.message}`);
    }
  }

  async function loadProject(project: Project) {
    try {
      setCurrentProject(project);
      await fetchDocuments(project.id);

    } catch (error: any) {
      alert(`Error loading documents: ${error.message}`);
    }
  }

  return (
    <>
      <div className="parent-container">
        <div>
          {projects.map((project) => (
            <div key={project.id}>
              <h5>{project.title}</h5>
              <button onClick={() => loadProject(project)}>
                Load Documents
              </button>

            </div>
          ))}
        </div>
        <main>
          <h1>Welcome to Pivo</h1>
          <h2>Current Project: {currentProject ? currentProject.title : "None"}</h2>
          <p>WebSocket Status: {connectionStatus}</p>
          <p>{answer}</p>
          <button onClick={() => createProject({ title: "New Project 111" })}>
            Add Project
          </button>
          <input
            type="file"
            onChange={uploadFile}
          />
          <input
            type="text"
            placeholder="Ask a question"
            onChange={(e) => {
              setQuestion(e.target.value);
            }}
          />
          <button onClick={sendQuestion}>Ask question</button>
        </main>
        <div>
          {<ul>
            {documents.map((doc) => (
              <li key={doc.id}>{doc.title}</li>
            ))}
          </ul>}
        </div>
      </div>
    </>
  )
}

export default App
