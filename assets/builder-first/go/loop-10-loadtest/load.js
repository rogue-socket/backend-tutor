// k6 load test for the links service.
//
// Run:
//   k6 run load.js
//
// Or against a deployed instance:
//   BASE_URL=https://my-service.fly.dev k6 run load.js
//
// Open-loop arrival-rate executor produces requests at the target rate
// regardless of how slow the service is — this is what tells you when the
// service can't keep up (vs closed-loop, which would just queue and hide
// tail latency).

import http from 'k6/http';
import { check } from 'k6';

export const options = {
  scenarios: {
    ramp_then_soak: {
      executor: 'ramping-arrival-rate',
      startRate: 10,
      timeUnit: '1s',
      preAllocatedVUs: 100,
      maxVUs: 500,
      stages: [
        { duration: '5m', target: 1000 },
        { duration: '5m', target: 1000 },
        { duration: '1m', target: 0 },
      ],
    },
  },
  thresholds: {
    http_req_failed:    ['rate<0.01'],
    http_req_duration:  ['p(95)<200', 'p(99)<500'],
  },
};

const BASE = __ENV.BASE_URL || 'http://localhost:8080';

export default function () {
  if (Math.random() < 0.9) {
    const r = http.get(`${BASE}/links/popular`);
    check(r, { 'popular 200': (r) => r.status === 200 });
  } else {
    const r = http.post(`${BASE}/links`, JSON.stringify({
      url: 'https://x', title: 'x',
    }), { headers: { 'Content-Type': 'application/json' } });
    check(r, { 'create 201': (r) => r.status === 201 });
  }
}
