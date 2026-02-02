
# Rulemancer API Endpoints

The server exposes the following API routes under `/api/v1`:

## System Routes

- `GET /api/v1/system/status` - Get system status

## Game Routes  

- `GET /api/v1/game/list` - List available games
- `GET /api/v1/game/{id}` - Get game details

## Room Routes

- `POST /api/v1/room/create` - Create a new game room
- `GET /api/v1/room/list` - List active rooms
- `GET /api/v1/room/{id}` - Get room details
- `DELETE /api/v1/room/{id}` - Delete a room
- `POST /api/v1/room/{id}/assert` - Assert facts to room
- `POST /api/v1/room/{id}/query` - Query facts from room
- `POST /api/v1/room/{id}/join` - Join a room as a client
- `DELETE /api/v1/room/{id}/leave` - Leave a room