import { useEffect, useMemo, useState } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';

type AiMessage = {
  chunk: string,
  error: string | null,
  status: 'streaming' | 'completed' | 'error',
}

export const useWebSocketHandler = (clientId: string) => {

  const { lastMessage, readyState } = useWebSocket(`${import.meta.env.VITE_WS_URL}?client_id=${clientId}`);
  const [messageHistory, setMessageHistory] = useState<AiMessage[]>([]);
  const [answer, setAnswer] = useState<string>("");


  useEffect(() => {
    if (lastMessage !== null) {
      setMessageHistory((prev) => prev.concat(lastMessage.data));
      if (lastMessage.data === '') return; // Ignore empty messages
      const parsedMessage = JSON.parse(lastMessage.data) as AiMessage;

      if (parsedMessage.status === 'completed') {
        setAnswer(parsedMessage.chunk);
      }
      if (parsedMessage.status === 'error') {
        alert(`Error: ${parsedMessage.error}`);
      }
      if (parsedMessage.status === 'streaming') {
        setAnswer((prev) => prev + parsedMessage.chunk);
      }
      // Optionally, you can also log the message or handle it in other ways
      console.log('Received message:', parsedMessage);
      // Reset the answer if you want to start fresh for each new question
    }
  }, [lastMessage]);

  const connectionStatus = {
    [ReadyState.CONNECTING]: 'Connecting',
    [ReadyState.OPEN]: 'Open',
    [ReadyState.CLOSING]: 'Closing',
    [ReadyState.CLOSED]: 'Closed',
    [ReadyState.UNINSTANTIATED]: 'Uninstantiated',
  }[readyState];

  return { connectionStatus, messageHistory, answer };

}
