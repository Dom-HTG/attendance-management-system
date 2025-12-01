# Workspace Setup — attendance-management + frontend

This workspace file and VS Code configuration allow you to manage both the Go backend and the frontend app from the same VS Code window.

Files added:
- `attendance-management.code-workspace` — opens both backend and frontend folders in the same workspace
- `.vscode/tasks.json` — tasks to run backend and frontend, and a compound task
- `.vscode/extensions.json` — recommended extensions

Quick start

1. Open the workspace file in VS Code:

```bash
code attendance-management.code-workspace
```

2. Install recommended extensions if prompted.

3. Open the Command Palette (Ctrl+Shift+P) → `Tasks: Run Task` → choose `Run: Backend + Frontend`.
   - This runs `docker compose up` in the backend and `npm install && npm start` in the frontend concurrently.

4. Alternatively run tasks separately:
   - `Backend: Docker Compose Up` (runs in backend folder)
   - `Backend: Go Run (dev)` (if you prefer running Go directly)
   - `Frontend: npm start` (runs in frontend folder)

Notes

- The workspace references the frontend at `C:/Users/HP/Desktop/dev/macro-logger-frontend/attend-ai-live`.
  If your frontend is at a different path, edit `attendance-management.code-workspace` and update the `path` for the `attend-ai-live` folder.

- The tasks use `${workspaceFolder:attendance-management}` and `${workspaceFolder:attend-ai-live}` so they run in the right directories.

- The default integrated terminal profile is set to WSL in the workspace file. If you prefer another terminal, change the setting in `attendance-management.code-workspace` or VS Code settings.

- If your frontend uses `yarn` or another script to start, edit `.vscode/tasks.json` and replace the frontend command.

Troubleshooting

- "docker compose up" failing: Ensure Docker Desktop/Engine is running. You can also run the task `Backend: Go Run (dev)` to run the Go server locally without Docker.

- Frontend failing to start: Check `package.json` in the frontend folder. If `npm start` isn't available, run the proper start script or adapt the task.

If you want, I can:
- Detect and adapt the frontend task to `yarn` if `yarn.lock` exists
- Add `launch.json` entries for Go and Chrome debugging
- Add a `tasks` to seed or clear the dev database

Tell me which of these you'd like next.