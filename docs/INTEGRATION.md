# Frontend Integration Guide

This document gives the essential details a frontend developer needs to integrate with the backend API and support QR code workflows.

Authentication
- Login endpoints (student and lecturer) return a JWT token that must be sent in subsequent requests in the Authorization header as `Bearer <token>`.
- Store tokens securely in the frontend (in-memory or secure HTTP-only cookie for SPAs).

Important headers
- `Content-Type: application/json` for JSON bodies
- `Authorization: Bearer <token>` for protected endpoints

QR code handling (frontend responsibilities)
1. Lecturer flow:
   - Call POST `/api/lecturer/qrcode/generate` with course/event details.
   - The response contains `qr_code_data` (base64 PNG) and `qr_token` (UUID). Display the PNG by setting `src` to `data:image/png;base64,<qr_code_data>` in an <img> tag.
   - Keep event_id and qr_token for optional later operations.

2. Student flow:
   - Student scans the QR image with a device or the frontend reads the `qr_token` value (if the scanner provides token).
   - Frontend sends POST `/api/attendance/check-in` with `{ "qr_token": "<token>" }` and Authorization header with student token.
   - Handle error responses for time-window violations or duplicate check-ins.

Client-side validation
- Validate required fields (non-empty strings) before sending requests.
- Use client timezone awareness to display event start/end times; backend validates check-in times using server time.

Error handling
- Map backend status codes to UI messages:
  - 400: show validation message or QR invalid
  - 401: redirect to login
  - 403: show "access denied"
  - 409: show "already checked in"

Sample fetch usage (pseudo-code)
```js
// login
const res = await fetch(`${BASE_URL}/api/auth/login-student`, { method: 'POST', body: JSON.stringify({email, password}), headers: {'Content-Type':'application/json'} });
const json = await res.json();
localStorage.setItem('token', json.token);

// generate QR (lecturer)
await fetch(`${BASE_URL}/api/lecturer/qrcode/generate`, { method: 'POST', headers: { 'Content-Type':'application/json', 'Authorization': `Bearer ${token}` }, body: JSON.stringify(payload) });

// check-in (student)
await fetch(`${BASE_URL}/api/attendance/check-in`, { method: 'POST', headers: { 'Content-Type':'application/json', 'Authorization': `Bearer ${token}` }, body: JSON.stringify({ qr_token }) });
```

CORS and environment
- Ensure the backend has CORS configured to allow your frontend origin in development.
- Use environment variables to configure BASE_URL in frontend.

Testing helpers
- Use `docs/API.md` for exact request/response shapes.
- On QR generation, the `qr_token` is what students need; `qr_code_data` is for displaying QR image.

Notes
- Time window enforcement is by server: do not rely on client time for enforcing check-in validity.
- Avoid persisting long-lived JWTs in localStorage for security; prefer HTTP-only cookies when possible.
