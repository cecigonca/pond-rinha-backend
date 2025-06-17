import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
  vus: 100,         // 100 usuários simultâneos
  duration: '30s',  // duração total do teste
};

const BASE_URL = 'http://localhost:9999'; // Load balancer

export default function () {
  const clientId = Math.floor(Math.random() * 5) + 1; // IDs de 1 a 5

  const tipo = Math.random() > 0.5 ? 'c' : 'd';
  const valor = Math.floor(Math.random() * 5000) + 1;

  const payload = JSON.stringify({
    valor: valor,
    tipo: tipo,
    descricao: "rinhaTest"
  });

  const headers = { 'Content-Type': 'application/json' };

  const res = http.post(`${BASE_URL}/clientes/${clientId}/transacoes`, payload, { headers });

  check(res, {
    'status is 200 or 422': (r) => r.status === 200 || r.status === 422,
  });

  sleep(0.1); // aguarda 100ms antes de repetir
}