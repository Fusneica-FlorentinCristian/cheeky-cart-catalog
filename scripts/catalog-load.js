// k6 load test for cheeky-cart-catalog REST (L09 HW7).
// Usage (from repo root): k6 run scripts/catalog-load.js
// Prerequisite: go run ./cmd/rest

import http from 'k6/http';
import { check, sleep } from 'k6';

const BASE = __ENV.BASE_URL || 'http://localhost:8080';

export const options = {
  stages: [
    { duration: '10s', target: 20 },
    { duration: '30s', target: 50 },
    { duration: '10s', target: 0 },
  ],
  thresholds: {
    http_req_failed: ['rate<0.01'],
    http_req_duration: ['p(95)<2000'],
  },
};

export default function () {
  const health = http.get(`${BASE}/health`);
  check(health, { 'health status 200': (r) => r.status === 200 });

  const list = http.get(`${BASE}/products`);
  check(list, { 'products status 200': (r) => r.status === 200 });

  const detail = http.get(`${BASE}/products/1`);
  check(detail, { 'product 1 status 200': (r) => r.status === 200 });

  sleep(0.1);
}
