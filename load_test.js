import http from 'k6/http';
import { sleep, check } from 'k6';

export let options = {
  stages: [
    { duration: '30s', target: 10 },
    { duration: '30s', target: 20 },
    { duration: '30s', target: 50 },
    { duration: '30s', target: 100 },
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'],
  },
};

const BASE_URL = 'http://localhost:8080';

function randomUserId() {
  return Math.floor(Math.random() * 10) + 1;  // Random userId between 1 and 10
}

function randomPollId() {
  return Math.floor(Math.random() * 220) + 1;  // Random pollId between 1 and 100
}

export default function () {
  const userId = randomUserId();  
  const pollID = randomPollId(); 

  const action = Math.random();
  
  if (action < 0.5) {
    // Fetch Polls
    const res = http.get(`${BASE_URL}/polls`, {
      headers: { 'userId': String(userId) },
    });
    check(res, { 'fetch polls status 200': (r) => r.status === 200 });
    
  } else if (action < 0.8) {
    // Vote
    const res = http.post(`${BASE_URL}/polls/${pollID}/vote`, JSON.stringify({ optionIndex: 0 }), {
      headers: { 'Content-Type': 'application/json', 'userId': String(userId) },
    });
    check(res, { 'vote status ok or already voted': (r) => r.status === 200 || r.status === 409 });

  } else {
    // Skip
    const res = http.post(`${BASE_URL}/polls/${pollID}/skip`, null, {
      headers: { 'userId': String(userId) },
    });
    check(res, { 'skip status ok or already skipped': (r) => r.status === 200 || r.status === 409 });
  }

  sleep(1);
}
