import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '30s', target: 500 }, // ramp up to 500 VUs in 30s
        { duration: '1m', target: 500 }, // stay at 500 VUs for 1 minute
    ]
};

const BASE_URL = "http://localhost:5100";

export default function () {
    const userId = 'user1';
    const symbol = 'AAPL';
    const amount = 1000;
    const transactionNum = 1;

    // Test the /buy endpoint
    check(http.post(`${BASE_URL}/buy`, JSON.stringify({
        "user": userId,
        "stock": symbol,
        "amount": amount,
        "transactionNum": transactionNum
    })), { 'buy status was 200': (r) => r.status === 200 });

    sleep(1);
}





