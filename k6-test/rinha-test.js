import http from 'k6/http';
import { check, sleep } from 'k6';
import { Trend, Rate, Counter } from 'k6/metrics';

const BASE_URL = 'http://localhost:9999';
const CLIENT_IDS = [1, 2, 3, 4, 5];

// --- Configuração ---
export const options = {
  stages: [
     { duration: '20s', target: 200 },
    { duration: '3m', target: 200 },
  ],
  thresholds: {
    http_req_failed: ['rate<0.001'],
    http_req_duration: ['p(99)<600'],
    checks: ['rate>0.999'],
    transacao_req_duration: ['p(99)<500'],
    extrato_req_duration: ['p(99)<400'],
    transacao_errors_422_count: ['count<100000'],
    error_rate_422_transacao: ['rate<0.50'],
    error_rate_404: ['rate<0.0001'],
    error_rate_generic: ['rate<0.0001'],
  },
  summaryTrendStats: ['avg', 'min', 'med', 'max', 'p(90)', 'p(95)', 'p(99)', 'count'],
  discardResponseBodies: false,
};

// --- Métricas personalizadas ---
const extratoReqDuration = new Trend('extrato_req_duration');
const transacaoReqDuration = new Trend('transacao_req_duration');
const transacao422Counter = new Counter('transacao_errors_422_count');
const errorRate422Transacao = new Rate('error_rate_422_transacao');
const errorRate404 = new Rate('error_rate_404');
const errorRateGeneric = new Rate('error_rate_generic');
const successfulTransacoes = new Counter('successful_transacoes');
const successfulExtratos = new Counter('successful_extratos');
const descriptions = [
  "padaria", "aluguel", "mercado", "pix", "invest", "ifood", "uber",
  "farmacia", "luz", "agua", "salario", "bonus", "reembolso",
  "transfer", "servico", "netflix", "spotify", "cafe", "lanche", "jogo"
];

function getRandomClientId() {
  return CLIENT_IDS[Math.floor(Math.random() * CLIENT_IDS.length)];
}

function getRandomDescription() {
  return descriptions[Math.floor(Math.random() * descriptions.length)];
}

function getRandomTipo() {
  return Math.random() < 0.75 ? 'd' : 'c';
}

function getRandomValor(tipo) {
  if (tipo === 'c') {
    return [100, 500, 1000, 5000, 10000][Math.floor(Math.random() * 5)];
  } else {
    const p = Math.random();
    if (p < 0.4) return Math.floor(Math.random() * 5000) + 1;
    if (p < 0.8) return Math.floor(Math.random() * 40000) + 5001;
    return Math.floor(Math.random() * 200000) + 50000;
  }
}

// --- Main logic ---
export default function () {
  const clientId = getRandomClientId();
  // 40% POST /transacoes
  if (Math.random() < 0.4) {
    const tipo = getRandomTipo();
    const valor = getRandomValor(tipo);
    const descricao = getRandomDescription();
    const payload = JSON.stringify({ valor, tipo, descricao });
    const res = http.post(`${BASE_URL}/clientes/${clientId}/transacoes`, payload, {
      headers: { 'Content-Type': 'application/json' },
      tags: { name: 'TransacaoPOST' },
    });

    transacaoReqDuration.add(res.timings.duration);
    const is200 = res.status === 200;
    const is422 = res.status === 422;
    check(res, {
      '[Transacao] status 200 ou 422': is200 || is422,
    });

    if (is200) {
      successfulTransacoes.add(1);
      errorRate422Transacao.add(0);
    } else if (is422) {
      transacao422Counter.add(1);
      errorRate422Transacao.add(1);
    } else if (res.status === 404) {
      errorRate404.add(1);
    } else {
      errorRateGeneric.add(1);
    }
  }
  // 60% GET /extrato
  else {
    const res = http.get(`${BASE_URL}/clientes/${clientId}/extrato`, {
      tags: { name: 'ExtratoGET' },
    });
    extratoReqDuration.add(res.timings.duration);
    check(res, {
      '[Extrato] status 200 ou 404': res.status === 200 || res.status === 404,
    });
    if (res.status === 200) {
      successfulExtratos.add(1);
    } else if (res.status === 404) {
      errorRate404.add(1);
    } else {
      errorRateGeneric.add(1);
    }
  }
  //sleep(0.1);
}