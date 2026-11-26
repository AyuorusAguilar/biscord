'use client';
import { useSearchParams } from 'next/navigation';

export default function Conversacion() {
  const amigo = useSearchParams().get('u');
  return <h1>Chat con {amigo}</h1>;
}