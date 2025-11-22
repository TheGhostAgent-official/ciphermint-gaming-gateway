# CipherMint Gaming Gateway – RACKDOG™ Integration Contract

This document defines the HTTP contract for any game or integration that wants
to plug into the CipherMint Gaming Gateway and award RACKDOG™ (and other tokens)
to players.

Base URL (dev):
    http://localhost:8080

---

## 1. Health Check

**GET /health**

Response:
{
  "service": "CipherMint Gaming Gateway",
  "status": "ok"
}

---

## 2. Register Integration  
**POST /v1/game**

Request:
{
  "id": "ghostops_cod",
  "name": "Ghost Ops -- CoD Integration",
  "company_id": ""
}

Response:
{
  "id": "ghostops_cod",
  "name": "Ghost Ops -- CoD Integration",
  "company_id": ""
}

---

## 3. Create / Attach Player  
**POST /v1/game/{integration_id}/player**

Request:
{
  "player_id": "player123",
  "alias": "GhostPlayer"
}

Response:
{
  "player_id": "player123",
  "alias": "GhostPlayer",
  "integration_id": "ghostops_cod",
  "balances": {}
}

---

## 4. Earn Tokens  
**POST /v1/game/{integration_id}/player/{player_id}/earn**

Request:
{
  "token": "RACKDOG",
  "amount": 100,
  "source": "signup_bonus"
}

Response:
{
  "status": "ok",
  "token": "RACKDOG"
}

---

## 5. Fetch Player  
**GET /v1/game/{integration_id}/player/{player_id}**

Response:
{
  "player_id": "player123",
  "alias": "GhostPlayer",
  "integration_id": "ghostops_cod",
  "balances": {
    "RACKDOG": 100
  }
}

---

## 6. Usage Pattern

1. Register integration  
2. Register/attach player  
3. Earn tokens (login, win, action, purchase)  
4. Fetch balances  

---

## 7. Demo Script

scripts/rackdog_demo.sh  
Performs: health → integration → player → bonus → fetch.

