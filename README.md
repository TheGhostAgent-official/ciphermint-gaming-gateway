# CIPHERMINT™ GAMING GATEWAY – v1.0.0
### Studio Integration Package • Powered by CIPHERMINT STUDIOS™

---

## 1. Overview

The **CIPHERMINT Gaming Gateway™** is a plug-and-play microservice that allows any game studio to instantly integrate real token economies into their existing games without modifying their internal engine.

This system enables:

- Player identity creation and tracking  
- Token earning and spending events  
- Integration registration for each game title  
- Real-time token balance fetch  
- Secure API-key authentication  
- SQLite-backed persistence  
- Lightweight HTTP+JSON communication  

The demo uses **Rackdog™** (ticker: `RACKDOG`) to demonstrate functionality, but the gateway is fully **token-agnostic** and supports any game-specific token a studio chooses to issue.

---

## 2. The Core Concept

CIPHERMINT™ allows studios to:

### • Issue their own in-game tokens  
(Unique currency for each game or franchise)

### • Drive real, persistent digital value  
Tokens hold value **inside AND outside** the game world.

### • Strengthen player engagement  
Tokens can be earned for:
- Wins
- Level-ups
- Daily logins
- Missions  
- Achievements  
- Seasonal events  

### • Expand monetization without predatory tactics  
Players spend tokens on:
- Cosmetics  
- Upgrades  
- Battle Pass items  
- Limited drops  
- Cross-game perks  

### • Preserve engine independence  
Studios don’t need to change ANY engine code.  
The Gateway handles all token logic externally.

---

## 3. How the Gateway Works

1. **Register Integration** – The studio creates a game integration shell.  
2. **Create Player** – The player identity is linked to that integration.  
3. **Reward Tokens** – Tokens are awarded from gameplay events.  
4. **Spend Tokens** – Player spends tokens on in-game content.  
5. **Fetch Balances** – Studio queries the player’s current token state.

Everything communicates through clean HTTP+JSON.

---

## 4. API Endpoints

### GET `/health`
Check if the gateway is active.

### POST `/integration`
Register a new game integration.

### POST `/player/{integration_id}`
Create a player within an integration.

### POST `/player/{integration_id}/{player_id}/earn`
Award tokens.

### POST `/player/{integration_id}/{player_id}/spend`
Spend tokens.

### GET `/player/{integration_id}/{player_id}`
Retrieve token balances.

**Required header:**---
X-API-Key: <your_api_key>

## 5. Running the Demo

Run:
bash scripts/gateway_api_demo.sh

You will see:

- Health check  
- Register Integration  
- Create Player  
- Signup Bonus  
- Earnings  
- Spend events  
- Final balances  

---

## 6. Tech Stack

- Go (Golang)  
- SQLite3  
- `net/http`  
- CIPHERMINT™ plug-and-play architecture  

---

## 7. Studio Benefits

- No blockchain knowledge required  
- Instant cross-game token economies  
- Deeper player retention  
- New revenue channels  
- Lower development overhead  
- Secure, scalable, and engine-agnostic  

---

## 8. Brand Notice

This package is officially produced by:

**CIPHERMINT STUDIOS™**  
Creators of the world’s first universal cross-game token economy gateway.
