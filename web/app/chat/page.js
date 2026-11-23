import { Suspense } from 'react';
import Conversacion from './conversacion';

export default function ChatPage() {
  return (
    <Suspense fallback={<p>Cargando…</p>}>
      <Conversacion />
    </Suspense>
  );
}