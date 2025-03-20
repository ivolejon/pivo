import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import App from './App.tsx'
import { Theme } from "@radix-ui/themes";
import "@radix-ui/themes/styles.css";

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <Theme radius="small" scaling="95%">
      <App />
    </Theme>
  </StrictMode>,
)
